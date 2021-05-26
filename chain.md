# A promise like chain

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func main() {
	c := Chain{ctx: context.Background()}.
		Then(foo).
		Then(bar)
	if err := c.Err(); err != nil {
		log.Fatal(err)
	}
}

func foo(ctx context.Context) error {
	fmt.Println("calling foo")
	return errors.New("foo")
}
func bar(ctx context.Context) error {
	fmt.Println("calling bar")
	return errors.New("bar")
}

type Chain struct {
	ctx context.Context
	err error
}

func (c Chain) Then(fn func(ctx context.Context) error) Chain {
	if c.err != nil {
		return c
	}
	c.err = fn(c.ctx)
	return c
}

func (c Chain) Err() error {
	return c.err
}
```
