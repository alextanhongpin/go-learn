## Unwrapping multiple errors

We unwrap the errors from left to right:

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

var NotFoundErr = NewError("not found")
var NotFoundFactory = NewFactory[int]("not found")

func main() {
	var err error = NotFoundErr
	err = fmt.Errorf("%w: not found", err)
	fmt.Println(errors.Is(err, NotFoundErr))
	err = fmt.Errorf("%w: %w, %w", err, NotFoundFactory(32), NotFoundFactory(42))

	for {
		fmt.Println()
		fmt.Println("msg", err)
		var berr *BaseError
		if errors.As(err, &berr) {
			fmt.Println("data", berr.data)
		}
		type multiError interface {
			Unwrap() []error
		}
		v, ok := err.(multiError)
		if ok {
			errs := v.Unwrap()
			fmt.Println("yes", errs)
			err = errors.Join(errs[1:]...)
		} else {
			err = errors.Unwrap(err)
		}

		if err == nil {
			break
		}
	}

	fmt.Println("Hello, 世界")
}

type BaseError struct {
	error
	data any
	code string
}

func NewFactory[T any](msg string) func(t T) *BaseError {
	return func(t T) *BaseError {
		return &BaseError{
			error: errors.New(msg),
			data:  t,
		}
	}
}

func NewError(msg string) *BaseError {
	return &BaseError{
		error: errors.New(msg),
	}
}

func (b *BaseError) Is(err error) bool {
	return errors.Is(b.error, err)
}

func (b *BaseError) Unwrap() error {
	return b.error
}

func (b *BaseError) Data() any {
	return b.data
}

func (b *BaseError) Code() string {
	return b.code
}
```
