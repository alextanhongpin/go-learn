```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

var ErrFoo = errors.New("foo")
var ErrBar = errors.New("bar")

func main() {
	var err error = fmt.Errorf("%w: %w", ErrFoo, ErrBar)
	fmt.Println(err)
	fmt.Println(errors.Unwrap(err))
	fmt.Println(errors.Is(err, ErrFoo))
	fmt.Println(errors.Is(err, ErrBar))
	fmt.Println()
	err = &wrapErr{
		err: ErrFoo,
		ori: ErrBar,
	}
	fmt.Println(err)
	fmt.Println(errors.Unwrap(err))
	fmt.Println(errors.Is(err, ErrFoo))
	fmt.Println(errors.Is(err, ErrBar))
	/*
foo: bar
<nil>
true
true

foo: bar
bar
true
true
	*/
}

type wrapErr struct {
	err error
	ori error
}

func (t *wrapErr) Error() string {
	return fmt.Sprintf("%s: %s", t.err, t.ori)
}

func (t *wrapErr) Unwrap() error {
	return t.ori
}

func (t *wrapErr) Is(err error) bool {
	return errors.Is(t.err, err) || errors.Is(t.ori, err)
}
```
