```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func main() {
	fmt.Println("Hello, 世界")
	sp := NewStringParser()
	v := Compile(sp, "optional,min=10,max=20")
	fmt.Println(v.Validate(""))
}

var ErrSkip = errors.New("skip")

type Validator[T any] interface {
	Validate(T) error
}

type validator[T any] struct {
	funcs []func(T) error
}

func (v *validator[T]) Validate(t T) error {
	for _, fn := range v.funcs {
		if err := fn(t); err != nil {
			if errors.Is(err, ErrSkip) {
				return nil
			}
			return err
		}
	}
	return nil
}

type ParserMap[T any] interface {
	Get(key string) Parser[T]
}

func Compile[T any](m ParserMap[T], exprs string) Validator[T] {
	vals := &validator[T]{}
	for _, expr := range strings.Split(exprs, ",") {
		k, v, _ := strings.Cut(expr, "=")
		fn := m.Get(k)(v)
		vals.funcs = append(vals.funcs, fn)
	}
	return vals
}

type Parser[T any] func(string) func(T) error

type parser[T any] struct {
	parsers map[string]Parser[T]
}

func (p *parser[T]) Get(key string) Parser[T] {
	return p.parsers[key]
}

func (p *parser[T]) Set(key string, fn Parser[T]) {
	p.parsers[key] = fn
}

var (
	alphanum = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	numeric  = regexp.MustCompile(`^[0-9]+$`)
	email    = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
)

func NewStringParser() ParserMap[string] {
	return &parser[string]{
		parsers: map[string]Parser[string]{
			"required": func(string) func(string) error {
				return func(val string) error {
					if len(val) == 0 {
						return errors.New("must not be empty")
					}
					return nil
				}
			},
			"optional": func(string) func(string) error {
				return func(val string) error {
					if len(val) == 0 {
						return ErrSkip
					}
					return nil
				}
			},
			"min": func(params string) func(string) error {
				n := toInt(params)
				return func(val string) error {
					if len([]rune(val)) < n {
						return fmt.Errorf("min %d characters", n)
					}
					return nil
				}
			},
			"max": func(params string) func(string) error {
				n := toInt(params)
				return func(val string) error {
					if len([]rune(val)) > n {
						return fmt.Errorf("max %d characters", n)
					}
					return nil
				}
			},
			"len": func(params string) func(string) error {
				n := toInt(params)
				return func(val string) error {
					if len(val) != n {
						return fmt.Errorf("must have %d characters", n)
					}
					return nil
				}
			},
			"eq": func(params string) func(string) error {
				return func(val string) error {
					if val != params {
						return fmt.Errorf("must be %q", params)
					}
					return nil
				}
			},
			"neq": func(params string) func(string) error {
				return func(val string) error {
					if val == params {
						return fmt.Errorf("must not be %q", params)
					}
					return nil
				}
			},
			"oneof": func(params string) func(string) error {
				vals := strings.Fields(params)
				return func(val string) error {
					for _, v := range vals {
						if v == val {
							return nil
						}
					}
					return fmt.Errorf("must be one of %s", strings.Join(vals, ", "))
				}
			},
			"alphanum": func(params string) func(string) error {
				return func(val string) error {
					if !alphanum.MatchString(val) {
						return errors.New("must be alphanumeric")
					}
					return nil
				}
			},
			"numeric": func(string) func(string) error {
				return func(val string) error {
					if !numeric.MatchString(val) {
						return errors.New("must be numeric")
					}
					return nil
				}
			},
			"email": func(string) func(string) error {
				return func(val string) error {
					if !email.MatchString(val) {
						return errors.New("invalid email format")
					}
					return nil
				}
			},
			"url": func(string) func(string) error {
				return func(val string) error {
					_, err := url.Parse(val)
					if err != nil {
						return errors.New("invalid url format")
					}
					return nil
				}
			},
			"ip": func(string) func(string) error {
				return func(val string) error {
					ip := net.ParseIP(val)
					if ip == nil {
						return errors.New("invalid ip format")
					}
					return nil
				}
			},
		},
	}
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func toFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

type Number interface {
	constraints.Integer | constraints.Float
}

func NewFloat64Parser() ParserMap[float64] {
	return &parser[float64]{
		parsers: map[string]Parser[float64]{
			"required": func(params string) func(float64) error {
				return func(val float64) error {
					if val == 0.0 {
						return errors.New("required")
					}
					return nil
				}
			},
			"optional": func(params string) func(float64) error {
				return func(val float64) error {
					if val == 0.0 {
						return ErrSkip
					}
					return nil
				}
			},
			"min": func(params string) func(float64) error {
				f := toFloat64(params)
				return func(val float64) error {
					if val < f {
						return fmt.Errorf("min %f", f)
					}
					return nil
				}
			},
			"max": func(params string) func(float64) error {
				f := toFloat64(params)
				return func(val float64) error {
					if val > f {
						return fmt.Errorf("max %f", f)
					}
					return nil
				}
			},
			"eq": func(params string) func(float64) error {
				f := toFloat64(params)
				return func(val float64) error {
					if val != f {
						return fmt.Errorf("must be %f", f)
					}
					return nil
				}
			},
			"neq": func(params string) func(float64) error {
				f := toFloat64(params)
				return func(val float64) error {
					if val == f {
						return fmt.Errorf("must not be %f", f)
					}
					return nil
				}
			},
			"oneof": func(params string) func(float64) error {
				vals := strings.Fields(params)
				fs := make([]float64, len(vals))
				for i, v := range vals {
					fs[i] = toFloat64(v)
				}

				return func(val float64) error {
					for _, v := range fs {
						if v == val {
							return nil
						}
					}
					return fmt.Errorf("must be one of %s", strings.Join(vals, ", "))
				}
			},
			"positive": func(params string) func(float64) error {
				return func(val float64) error {
					if val < 0 {
						return errors.New("must be positive")
					}

					return nil
				}
			},
			"negative": func(params string) func(float64) error {
				return func(val float64) error {
					if val > 0 {
						return errors.New("must be negative")
					}

					return nil
				}
			},
		},
	}
}

```

### Validator

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type Num int
type Str string

func main() {
	s := "hello world"
	fmt.Println(Validate(s, Required, MinLength(10)))
	fmt.Println(Validate(s, Required, MinLength(100)))
	fmt.Println(Validate(s, StringExpr("required,min=100")...))

	// Unable to handle type definition.
	str := Str("hello world")
	fmt.Println(Validate(string(str), Required, MinLength(10)))
	fmt.Println(Validate(string(str), Required, MinLength(100)))

	num := 10.4
	fmt.Println(Validate(num, Required, Min(10.0)))
	fmt.Println(Validate(num, Required, Min(100.0)))

	arr := []Num{1, 2, 3}
	fmt.Println(Validate(arr, Required, MinItem[Num](3), Each[Num](Min(Num(10)), Max(Num(100)))))
	fmt.Println(Validate(arr, Required, MinItem[Num](3), Each[Num](Max(Num(2)))))

	var ns *string
	ns = &s
	fmt.Println(errors.Join(
		Validate(ns, Required),
		Validate(Value(ns), MinLength(100)),
	))
}

func Value[T any](v *T) T {
	if v == nil {
		var t T
		return t
	}
	return *v
}

func StringExpr(expr string) []func(string) error {
	var res []func(string) error
	for _, e := range strings.Split(expr, ",") {
		k, v, _ := strings.Cut(e, "=")
		switch k {
		case "required":
			res = append(res, Required[string])
		case "optional":
			res = append(res, Optional[string])
		case "min":
			res = append(res, MinLength(toInt(v)))
		case "max":
			res = append(res, MaxLength(toInt(v)))
		default:
			panic(fmt.Sprintf("unknown expr %q", e))
		}
	}
	return res
}
func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

var ErrSkip = errors.New("skip")

func Validate[T any](v T, fns ...func(T) error) error {
	for _, fn := range fns {
		if err := fn(v); err != nil {
			if errors.Is(err, ErrSkip) {
				return nil
			}
			return err
		}
	}
	return nil
}

func Required[T any](v T) error {
	if reflect.ValueOf(v).IsZero() {
		return errors.New("required")
	}
	return nil
}

func Optional[T any](v T) error {
	if reflect.ValueOf(v).IsZero() {
		return ErrSkip
	}
	return nil
}

func MinLength(n int) func(string) error {
	return func(s string) error {
		if len([]rune(s)) < n {
			return fmt.Errorf("min length of %d", n)
		}
		return nil
	}
}

func MaxLength(n int) func(string) error {
	return func(s string) error {
		if len([]rune(s)) > n {
			return fmt.Errorf("max length of %d", n)
		}
		return nil
	}
}

func Item[S ~[]E, E any](n int) func(S) error {
	return func(s S) error {
		if len(s) != n {
			return fmt.Errorf("must have %d item(s)", n)
		}
		return nil
	}
}

func MinItem[T any](n int) func([]T) error {
	return func(s []T) error {
		if len(s) < n {
			return fmt.Errorf("min %d item(s)", n)
		}
		return nil
	}
}

func MaxItem[T any](n int) func([]T) error {
	return func(s []T) error {
		if len(s) > n {
			return fmt.Errorf("max %d item(s)", n)
		}
		return nil
	}
}

func Length(n int) func(string) error {
	return func(s string) error {
		if len([]rune(s)) == n {
			return fmt.Errorf("length must be %d", n)
		}
		return nil
	}
}

func Equal[T comparable](v T) func(T) error {
	return func(vv T) error {
		if vv != v {
			return fmt.Errorf("not equal %v", v)
		}
		return nil
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

func Min[T Number](n T) func(T) error {
	return func(v T) error {
		if v < n {
			return fmt.Errorf("min %v", n)
		}
		return nil
	}
}

func Max[T Number](n T) func(T) error {
	return func(v T) error {
		if v > n {
			return fmt.Errorf("min %v", n)
		}
		return nil
	}
}

func Between[T Number](lo, hi T) func(T) error {
	return func(v T) error {
		if v < lo || v > hi {
			return fmt.Errorf("must be between %v and %v", lo, hi)
		}
		return nil
	}
}

func Each[T any](fns ...func(T) error) func([]T) error {
	return func(vs []T) error {
		for _, v := range vs {
			for _, fn := range fns {
				if err := fn(v); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
```
