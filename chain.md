# A promise like chain

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	c := Chain{}.
		Then(foo).
		Then(bar)
	if err := c.Err(); err != nil {
		log.Fatal(err)
	}
}

func foo() error {
	fmt.Println("calling foo")
	return errors.New("foo")
}
func bar() error {
	fmt.Println("calling bar")
	return errors.New("bar")
}

type Chain struct {
	err error
}

func (c Chain) Then(fn func() error) Chain {
	if c.err != nil {
		return c
	}
	c.err = fn()
	return c
}

func (c Chain) Err() error {
	return c.err
}
```
