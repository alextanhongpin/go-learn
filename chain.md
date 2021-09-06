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


## Passing context

```go
package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	p := new(Prepper)

	bar := func(ctx context.Context) error {
		return p.Bar(ctx, "bar")
	}
	ctx := context.Background()
	c := NewChain(ctx).
		Then(p.Foo("hello foo")).
		Then(bar)
	if err := c.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}

type Prepper struct {
	foo string
	bar string
}

func (p *Prepper) Foo(fooParams string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		p.foo = fooParams
		fmt.Println("calling foo")
		return nil
	}
}

func (p *Prepper) Bar(ctx context.Context, barParams string) error {
	fmt.Println("calling bar")
	p.bar = p.foo + barParams
	return nil
}

type Chain struct {
	ctx context.Context
	err error
}

func NewChain(ctx context.Context) *Chain {
	return &Chain{ctx: ctx}
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


## Map

```go
package main

import (
	"errors"
	"fmt"
)

func main() {

	incrby1 := func(p *Processor) error {
		p.value += 1
		return nil
	}

	incrby10 := func(p *Processor) error {
		p.value += 10
		return nil
	}
	
	badChain := func(p *Processor) error {
		p.value += 5
		return errors.New("bad chain")
	}
	
	newproc := Processor{}.
		Map(incrby1).
		Map(incrby10).
		Map(badChain)
	fmt.Println(newproc.value)
}

type Processor struct {
	value int
	err   error
}

func (p Processor) Map(fn func(*Processor) error) Processor {
	if p.err != nil {
		return p
	}
	cp := p
	if err := fn(&cp); err != nil {
		p.err = err
		return p
	}
	return cp
}
```
