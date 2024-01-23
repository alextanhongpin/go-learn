```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Result struct {
	data  interface{}
	error error
}

type Work interface {
	Source(interface{}) *Result
	Sink(*Result) *Result
}

type split struct {
	Work
}

func (s *split) Source(in interface{}) *Result {
	if in == nil {
		return &Result{
			error: fmt.Errorf("split input is empty"),
		}
	}
	t, ok := in.(string)
	if !ok {
		return &Result{
			error: errors.New("type assertion error"),
		}
	}
	o := strings.Split(t, " ")
	return &Result{
		data:   o,
	}
}

type prefix struct {
	key string
	Work
}

func (p *prefix) Sink(in *Result) *Result {
	if in == nil || in.error != nil || in.data == nil{
		return in
	}
	t, ok := in.data.([]string)
	if !ok {
		return &Result{
			error: errors.New("type assertion error"),
		}
	}
	for i, v := range t {
		t[i] = p.key + v
	}
	return &Result{
		data: t,
	}
}

type suffix struct {
	Work
	key string
}

func (s *suffix) Sink(in *Result) *Result {
	if in == nil || in.error != nil || in.data == nil {
		return nil
	}

	v, ok := in.data.([]string)
	if !ok {
		return nil
	}
	for i, str := range v {
		v[i] = str + s.key
	}
	return &Result{
		data: v,
	}
}

func main() {
	i := new(split)
	works := []Work{&prefix{key: "start-"}, &suffix{key: "-end"}}

	o := i.Source(nil)
	for _, w := range works {
		v := w.Sink(o)
		if v == nil || v.error != nil || v.data == nil {
			fmt.Println(v)
			return
		}
		if v.error == nil {
			o = v
		}
	}
	if o.error != nil {
		fmt.Println(o.error)
	}
	fmt.Println("out", o.data)
}
```

## Pipeline pattern

To reduce the number of error checks

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	pipe := New()
	pipe.Add(func() error { fmt.Println("hello"); return nil })
	pipe.Add(func() error { fmt.Println("world"); return nil })
	fmt.Println(pipe.Exec())
}

type Pipeline struct {
	funcs []func() error
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) Add(fn func() error) {
	p.funcs = append(p.funcs, fn)
}

func (p *Pipeline) Exec() error {
	for _, fn := range p.funcs {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

```

## Generics

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
)

func main() {
	uc := NewCreateAccount()
	fmt.Println(uc.Do(context.Background()))
	fmt.Println("Hello, 世界")
}

type CreateAccountState struct {
	Events []string
}

type CreateAccount struct {
}

func NewCreateAccount() *CreateAccount {
	return &CreateAccount{}
}

func (c *CreateAccount) Do(ctx context.Context) (*CreateAccountState, error) {
	req := new(CreateAccountState)
	p := &Pipeline[*CreateAccountState]{
		Steps: []PipelineFunc[*CreateAccountState]{
			c.Validate,
			c.Create,
			c.Notify,
		},
	}
	return req, p.Exec(ctx, req)
}

func (c *CreateAccount) Validate(ctx context.Context, req *CreateAccountState) error {
	req.Events = append(req.Events, "validate")
	return nil
}

func (c *CreateAccount) Create(ctx context.Context, req *CreateAccountState) error {
	req.Events = append(req.Events, "create")
	return nil
}
func (c *CreateAccount) Notify(ctx context.Context, req *CreateAccountState) error {
	req.Events = append(req.Events, "notify")
	return nil
}

type PipelineFunc[T any] func(ctx context.Context, t T) error

type Pipeline[T any] struct {
	Steps []PipelineFunc[T]
}

func (p *Pipeline[T]) Exec(ctx context.Context, t T) error {
	for _, step := range p.Steps {
		if err := step(ctx, t); err != nil {
			return err
		}
	}
	return nil
}
```
