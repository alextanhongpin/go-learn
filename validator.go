# Struct validator using golang reflect

The `Validate` method is unnecessary - just implement your own `struct.Validate()` method to give better control.

```go
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Request struct {
	Foo    *Foo `json:"foo"`
	Bar    *Bar
	Nested *Nested `json:"nested"`
	Slice  []Foo
}

type Nested struct {
	Foo Foo
}

func (n *Nested) Validate() error {
	if n == nil {
		return errors.New("nested required")
	}

	return n.Foo.Validate()
}

type Foo string

func (f *Foo) Validate() error {
	if f == nil {
		return errors.New("foo required")
	}

	if len(*f) == 0 {
		return errors.New("foo is required")
	}

	return nil
}

type Bar int

func (b *Bar) Validate() error {
	if b == nil {
		return errors.New("bar required") // Mandatory.
	}

	if *b < 0 {
		return errors.New("bar cannot be negative value")
	}

	return nil
}

func main() {
	var r Request
	fmt.Println(Validate(r))

	foo := Foo("hello")
	r.Foo = &foo

	fmt.Println(Validate(r))

	bar := Bar(0)
	r.Bar = &bar
	fmt.Println(Validate(r))

	r.Nested = &Nested{}
	fmt.Println(Validate(r))

	r.Nested.Foo = "world"
	fmt.Println(Validate(r))

	r.Slice = append(r.Slice, Foo(""))
	fmt.Println(Validate(r))

	r.Slice[0] = Foo("hello")
	fmt.Println(Validate(r))

	fzero := Foo("")
	fmt.Println("zero pointer foo:", Validate(&fzero))
	fmt.Println("zero foo:", Validate(Foo("")))
	fmt.Println("non-zero foo", Validate(Foo("hello")))
	fmt.Println(Validate([]Foo{Foo("")}))
	fmt.Println(Validate([]Foo{Foo("Foo")}))
}

func Validate(r any) error {
	v := reflect.ValueOf(r)

	return validate("", v, false)
}

func validate(prefix string, v reflect.Value, omitempty bool) error {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i)
			ti := t.Field(i)

			name := ti.Tag.Get("json")
			omitempty := strings.Contains(name, ",omitempty")

			if strings.Contains(name, ",") {
				name = strings.Split(name, ",")[0]
			}
			if name == "" {
				name = ti.Name
			}
			if prefix != "" {
				name = prefix + "." + name
			}

			switch vi.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array:
				if err := validate(name, vi, omitempty); err != nil {
					return err
				}
			default:
				if err := validateImplementsInterface(vi, omitempty); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			vi := v.Index(i)
			if !isValidatable(vi) {
				return nil
			}

			if err := validate(prefix, vi, omitempty); err != nil {
				return err
			}
		}
	default:
		return validateImplementsInterface(v, omitempty)
	}

	return nil
}

type validator interface {
	Validate() error
}

func isValidatable(v reflect.Value) bool {
	in := reflect.TypeOf((*validator)(nil)).Elem()

	if v.CanAddr() {
		return v.Addr().Type().Implements(in)
	}

	return v.Type().Implements(in)
}

func validateImplementsInterface(v reflect.Value, omitempty bool) error {
	if !isValidatable(v) {
		return nil
	}

	if omitempty && ((v.Kind() == reflect.Pointer && v.IsNil()) || (v.IsZero())) {
		return nil
	}

	in := v.Interface()
	if v.CanAddr() {
		in = v.Addr().Interface()
	}

	vi, ok := in.(validator)
	if !ok {
		return nil
	}

	return vi.Validate()
}
```
