```go
package tools

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"log"
	"maps"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
	"golang.org/x/tools/go/packages"
)

var (
	ctxType = reflect.TypeFor[context.Context]()
	errType = reflect.TypeFor[error]()
)

func IsFunc(a any) bool {
	t := reflect.TypeOf(a)
	return t.Kind() == reflect.Func
}

type Tools struct {
	tools api.Tools
	funcs map[string]func(map[string]any) ([]byte, error)
}

func NewTools() *Tools {
	return &Tools{
		funcs: make(map[string]func(map[string]any) ([]byte, error)),
	}
}

func (t *Tools) Tools() api.Tools {
	return t.tools
}

func (t *Tools) AddMCP(ctx context.Context, session *mcp.ClientSession, exclude ...string) error {
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return err
	}
	log.Println("registering MCP tools")
	for _, tool := range toolsResult.Tools {
		if slices.Contains(exclude, tool.Name) {
			log.Println("skipping", tool.Name)
			continue
		}
		log.Printf(" - %s: %s %v\n", tool.Name, tool.Description, tool.InputSchema)
		if _, ok := t.funcs[tool.Name]; ok {
			return fmt.Errorf("tools: mcp tool exists: %q", tool.Name)
		}
		b, err := json.Marshal(tool.InputSchema)
		if err != nil {
			return err
		}
		var parameters api.ToolFunctionParameters
		if err := json.Unmarshal(b, &parameters); err != nil {
			return err
		}
		log.Println("schema is", tool.Name, string(b), parameters)
		t.tools = append(t.tools, api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  parameters,
			},
		})
		var schema jsonschema.Schema
		if err := json.Unmarshal(b, &schema); err != nil {
			return err
		}
		rs, err := schema.Resolve(nil)
		if err != nil {
			return err
		}
		t.funcs[tool.Name] = func(args map[string]any) ([]byte, error) {
			if err := rs.Validate(args); err != nil {
				return nil, err
			}

			fmt.Printf("\nCalling Tool %s(%v)...\n", tool.Name, args)
			res, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name:      tool.Name,
				Arguments: args,
			})
			if err != nil {
				return nil, err
			}
			if res.StructuredContent != nil {
				return json.Marshal(res.StructuredContent)
			}
			fmt.Printf("Output:\n%v\n\n", res.StructuredContent)

			return json.Marshal(res.Content[0])
		}
	}

	return nil
}
func (t *Tools) Add(fn any, opts *api.ToolFunction) error {
	opts = cmp.Or(opts, new(api.ToolFunction))
	name := cmp.Or(opts.Name, GetShortFunctionName(fn))
	name = strings.ToLower(name)
	if _, ok := t.funcs[name]; ok {
		return fmt.Errorf("function: already exists: %q", name)
	}
	log.Printf("registered tool: %s\n", name)

	description := cmp.Or(opts.Description, GetFunctionDescription(fn))

	fv := reflect.ValueOf(fn)
	ft := fv.Type()

	var arg reflect.Type
	nargs := ft.NumIn()
	// Must be of type func(T) (V, error) or func(context.Context, T) (V, error)
	switch {
	case nargs == 1 && ft.In(0) != ctxType:
		arg = ft.In(0)
	case nargs == 2 && ft.In(0) == ctxType:
		arg = ft.In(1)
	default:
		return fmt.Errorf("function: invalid input signature: %s", ft)
	}
	if arg.Kind() == reflect.Pointer {
		return fmt.Errorf("function: arg cannot be pointer: %s", ft)
	}
	if arg.Kind() != reflect.Struct {
		return fmt.Errorf("function: arg must be struct: %s", ft)
	}

	// Must be of type func(...) (T, error)
	nrets := ft.NumOut()
	switch {
	case nrets == 2 && ft.Out(1).Implements(errType):
	default:
		return fmt.Errorf("function: invalid output signature: %s", ft)
	}

	// Get JSON Schema for input.
	schema, err := jsonschema.ForType(arg, nil)
	if err != nil {
		return err
	}
	rs, err := schema.Resolve(nil)
	if err != nil {
		return err
	}

	b, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	var parameters api.ToolFunctionParameters
	if err := json.Unmarshal(b, &parameters); err != nil {
		return err
	}
	log.Println("schema is", name, string(b), parameters)

	// Append function definition.
	t.tools = append(t.tools, api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	})

	// Register function.
	t.funcs[name] = func(args map[string]any) ([]byte, error) {
		if err := rs.Validate(args); err != nil {
			return nil, err
		}

		data, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		// Create new instance of that type and unmarshal.
		argPtr := reflect.New(arg)
		if err := json.Unmarshal(data, argPtr.Interface()); err != nil {
			return nil, err
		}

		in := make([]reflect.Value, nargs)
		if nargs == 2 {
			in[0] = reflect.Indirect(reflect.New(ctxType))
			in[1] = argPtr.Elem()
		} else {
			in[0] = argPtr.Elem()
		}

		// Call function with unmarshaled value (dereference pointer).
		out := fv.Call(in)
		errValue := out[nrets-1]
		if errValue.IsValid() && !errValue.IsNil() {
			err, ok := errValue.Interface().(error)
			if ok && err != nil {
				return nil, err
			}
		}
		res := out[0].Interface()
		return json.Marshal(res)
	}

	return nil
}

func (t *Tools) Exec(name string, args map[string]any) ([]byte, error) {
	log.Printf("tools.Exec: %s %v. Available: %v\n", name, args, slices.Collect(maps.Keys(t.funcs)))

	fn, ok := t.funcs[name]
	if !ok {
		return nil, fmt.Errorf("function not found: %q", name)
	}
	return fn(args)
}

func (t *Tools) Load(name string) (func(map[string]any) ([]byte, error), bool) {
	fn, ok := t.funcs[name]
	return fn, ok
}

// GetFunctionName returns the full name of the function passed to it.
func GetFunctionName(a any) string {
	// Use reflect.ValueOf to get the function's value
	pc := reflect.ValueOf(a).Pointer()

	// Use runtime.FuncForPC to get function information
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}

	// f.Name() returns the full path, e.g., "main.main" or "main.GetFunctionName"
	fullName := f.Name()
	return fullName
}

// GetShortFunctionName returns just the name without the package prefix.
func GetShortFunctionName(a any) string {
	fullName := GetFunctionName(a)
	// Split by '/' and take the last part, then split by '.' and take the last part
	strs := strings.Split(fullName, ".")

	// Struct method calls, e.g. (*Struct).Method will add a suffix, e.g. -fm
	// behind because the compiler creates a closure, and the suffix is to avoid
	// collision in name.
	name := strs[len(strs)-1]
	name, _, _ = strings.Cut(name, "-")
	return name
}

func GetFunctionDescription(fn any) string {
	description := "no description available"

	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:  ".",
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		panic(err)
	}

	name := GetShortFunctionName(fn)

	for _, pkg := range pkgs {
		for _, fileAST := range pkg.Syntax {
			ast.Inspect(fileAST, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if fn.Doc != nil && fn.Name.Name == name {
						description = strings.TrimSpace(fn.Doc.Text())
						return false
					}
				}

				return true
			})
		}
	}

	return description
}
```

## Alternative 

```
package tools

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"log"
	"log/slog"
	"reflect"
	"runtime"
	"strings"

	"github.com/alextanhongpin/go-ollama/types"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ollama/ollama/api"
	"golang.org/x/tools/go/packages"
)

var (
	ctxType = reflect.TypeFor[context.Context]()
	errType = reflect.TypeFor[error]()

	ErrToolExists   = errors.New("tool: function exists")
	ErrToolNotFound = errors.New("tool: not found")
)

type ToolFunc = func(args map[string]any) ([]byte, error)

type Tools struct {
	tools api.Tools
	funcs map[string]ToolFunc
}

func NewTools() *Tools {
	return &Tools{
		funcs: make(map[string]ToolFunc),
	}
}

func (t *Tools) Tools() api.Tools {
	return t.tools
}

func (t *Tools) Has(name string) bool {
	_, ok := t.funcs[name]
	return ok
}

func (t *Tools) Set(tool api.Tool, fn ToolFunc) error {
	name := tool.Function.Name
	t.tools = append(t.tools, tool)
	t.funcs[name] = fn
	slog.Default().Info("added tool", "name", tool.Function.Name, "description", tool.Function.Description)
	return nil
}

func (t *Tools) Add(tool api.Tool, fn ToolFunc) error {
	name := tool.Function.Name
	if _, ok := t.funcs[name]; ok {
		return fmt.Errorf("%w: %s", ErrToolExists, name)
	}

	return t.Set(tool, fn)
}

func (t *Tools) Exec(name string, args map[string]any) ([]byte, error) {
	log.Printf("\n\n[ToolCalling] %s(%v)\n", name, args)

	fn, ok := t.funcs[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrToolNotFound, name)
	}
	res, err := fn(args)
	if err != nil {
		log.Printf("\n\n[ToolCallingFailed] %s(%v) = Error(%v)\n", name, args, err)
		return nil, err
	}
	log.Printf("\n\n[ToolCallingSuccess] %s(%v) = Result(%s)\n", name, args, res)

	return res, nil
}

func AddAnyTool(ctx context.Context, t *Tools, fn any, opts *api.ToolFunction) error {
	a, err := AnyToToolAPI(fn, opts)
	if err != nil {
		return err
	}
	afn, err := AnyToToolFunc(fn)
	if err != nil {
		return err
	}

	return t.Add(a, afn)
}

func AddGoTool[T, V any](ctx context.Context, t *Tools, fn func(context.Context, T) (V, error), opts *api.ToolFunction) error {
	a, err := GoFuncToToolAPI(fn, opts)
	if err != nil {
		return err
	}
	afn := GoFuncToToolFunc(fn)

	return t.Add(a, afn)
}

func GoFuncToToolAPI[T, V any](fn func(context.Context, T) (V, error), opts *api.ToolFunction) (api.Tool, error) {
	opts = cmp.Or(opts, new(api.ToolFunction))

	name := cmp.Or(opts.Name, GetShortFunctionName(fn))
	description := cmp.Or(opts.Description, GetFunctionDescription(fn))
	parameters := opts.Parameters

	if reflect.DeepEqual(parameters, api.ToolFunctionParameters{}) {
		schema, err := jsonschema.For[T](nil)
		if err != nil {
			return api.Tool{}, err
		}
		b, err := json.Marshal(schema)
		if err != nil {
			return api.Tool{}, err
		}
		parameters, err = types.As[api.ToolFunctionParameters](b)
		if err != nil {
			return api.Tool{}, err
		}
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}, nil
}

func GoFuncToToolFunc[K, V any](fn func(context.Context, K) (V, error)) ToolFunc {
	return func(args map[string]any) ([]byte, error) {
		b, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		in, err := types.As[K](b)
		if err != nil {
			return nil, err
		}
		out, err := fn(context.Background(), in)
		if err != nil {
			return nil, err
		}

		return json.Marshal(out)
	}
}

func AnyToToolAPI(fn any, opts *api.ToolFunction) (api.Tool, error) {
	var zero api.Tool
	opts = cmp.Or(opts, new(api.ToolFunction))
	name := cmp.Or(opts.Name, GetShortFunctionName(fn))
	description := cmp.Or(opts.Description, GetFunctionDescription(fn))

	p, err := newParser(fn)
	if err != nil {
		return zero, err
	}

	var parameters api.ToolFunctionParameters
	if p.hasIn {
		b, err := json.Marshal(p.schema)
		if err != nil {
			return zero, err
		}
		parameters, err = types.As[api.ToolFunctionParameters](b)
		if err != nil {
			return zero, err
		}
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}, nil
}

func AnyToToolFunc(fn any) (ToolFunc, error) {
	p, err := newParser(fn)
	if err != nil {
		return nil, err
	}
	return p.fn, nil
}

type parser struct {
	hasIn    bool
	inType   reflect.Type
	hasCtx   bool
	hasOut   bool
	hasErr   bool
	schema   *jsonschema.Schema
	resolved *jsonschema.Resolved
	fn       func(args map[string]any) ([]byte, error)
}

func newParser(fn any) (*parser, error) {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()

	p := &parser{}
	return p, cmp.Or(
		p.validateInSignature(ft),
		p.validateOutSignature(ft),
		p.buildInSchema(),
		p.buildFunc(fv),
	)
}

func (p *parser) buildFunc(fv reflect.Value) error {
	p.fn = func(args map[string]any) ([]byte, error) {
		if p.hasIn {
			if err := p.resolved.Validate(args); err != nil {
				return nil, err
			}
		}

		var in []reflect.Value
		if p.hasCtx {
			in = append(in, reflect.ValueOf(context.Background()))
		}
		if p.hasIn {
			b, err := json.Marshal(args)
			if err != nil {
				return nil, err
			}
			inVal, err := TypeUnmarshal(b, p.inType)
			if err != nil {
				return nil, err
			}
			in = append(in, inVal)
		}

		vals := fv.Call(in)
		if p.hasErr {
			errVal := vals[len(vals)-1]
			if errVal.IsValid() && !errVal.IsNil() {
				err, ok := errVal.Interface().(error)
				if ok && err != nil {
					return nil, err
				}
			}
		}
		if p.hasOut {
			return json.Marshal(vals[0].Interface())
		}

		return nil, nil
	}

	return nil
}

func (p *parser) buildInSchema() error {
	if !p.hasIn {
		return nil
	}

	// Get JSON Schema for input.
	schema, err := jsonschema.ForType(p.inType, nil)
	if err != nil {
		return err
	}

	rs, err := schema.Resolve(nil)
	if err != nil {
		return err
	}

	p.schema = schema
	p.resolved = rs
	return nil
}

func (p *parser) validateInSignature(ft reflect.Type) error {
	in := ft.NumIn()
	switch {
	case in == 0:
	case in == 1 && isCtx(ft.In(0)):
		p.hasCtx = true

	case in == 1 && !isCtx(ft.In(0)):
		p.hasIn = true
		p.inType = ft.In(0)

	case in == 2 && isCtx(ft.In(0)) && !isCtx(ft.In(1)):
		p.inType = ft.In(1)
		p.hasIn = true
		p.hasCtx = true

	default:
		return fmt.Errorf("function: invalid input signature: %s", ft)
	}

	return nil
}

func (p *parser) validateOutSignature(ft reflect.Type) error {
	out := ft.NumOut()
	switch {
	case out == 0:
	case out == 1 && isErr(ft.Out(0)):
		p.hasErr = true

	case out == 1 && !isErr(ft.Out(0)):
		p.hasOut = true

	case out == 2 && !isErr(ft.Out(0)) && isErr(ft.Out(1)):
		p.hasOut = true
		p.hasErr = true

	default:
		return fmt.Errorf("function: invalid output signature: %s", ft)
	}

	return nil
}

// GetFunctionName returns the full name of the function passed to it.
func GetFunctionName(a any) string {
	// Use reflect.ValueOf to get the function's value
	pc := reflect.ValueOf(a).Pointer()

	// Use runtime.FuncForPC to get function information
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}

	// f.Name() returns the full path, e.g., "main.main" or "main.GetFunctionName"
	fullName := f.Name()
	return fullName
}

// GetShortFunctionName returns just the name without the package prefix.
func GetShortFunctionName(a any) string {
	fullName := GetFunctionName(a)
	// Split by '/' and take the last part, then split by '.' and take the last part
	strs := strings.Split(fullName, ".")

	// Struct method calls, e.g. (*Struct).Method will add a suffix, e.g. -fm
	// behind because the compiler creates a closure, and the suffix is to avoid
	// collision in name.
	name := strs[len(strs)-1]
	name, _, _ = strings.Cut(name, "-")
	return name
}

func GetFunctionDescription(fn any) string {
	description := "no description available"

	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:  ".",
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		panic(err)
	}

	name := GetShortFunctionName(fn)

	for _, pkg := range pkgs {
		for _, fileAST := range pkg.Syntax {
			ast.Inspect(fileAST, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if fn.Doc != nil && fn.Name.Name == name {
						description = strings.TrimSpace(fn.Doc.Text())
						return false
					}
				}

				return true
			})
		}
	}

	return description
}

func TypeUnmarshal(data []byte, t reflect.Type) (reflect.Value, error) {
	p := reflect.New(t)
	if err := json.Unmarshal(data, p.Interface()); err != nil {
		return reflect.Value{}, err
	}
	return p.Elem(), nil
}
func isCtx(t reflect.Type) bool {
	return t == ctxType
}

func isErr(t reflect.Type) bool {
	return t == errType
}

```
