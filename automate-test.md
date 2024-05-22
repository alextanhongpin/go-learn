```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

func main() {
	t := newCounter() // We need to know the number of errors in advance.
	f := &flow{t: t}
	t.rng(func() {
		fmt.Println()
		fmt.Println("iteration 1", len(t.seen))
		fmt.Println(f.exec())
	})

}

type counter struct {
	seen map[string]int
	i    int
}

func newCounter() *counter {
	return &counter{seen: make(map[string]int)}
}
func (c *counter) rng(fn func()) {
	for {
		fn()
		m := len(c.seen)
		for _, v := range c.seen {
			m = min(m, v)
		}
		// If the least item is seen twice, it means all the items in the map has been traversed.
		if m == 2 {
			break
		}
	}
}

func (c *counter) error(s string) error {
	defer func() {
		c.seen[s]++
	}()
	if _, ok := c.seen[s]; ok {
		return nil
	}
	return errors.New(s)
}

type flow struct {
	t *counter
}

func (f *flow) foo() error {
	return f.t.error("foo error")
}

func (f *flow) bar() error {
	return f.t.error("bar error")
}

func (f *flow) bizz() error {
	return f.t.error("bizz error")
}

func (f *flow) zee() error {
	return f.t.error("zee error")
}

func (f *flow) exec() error {
	if err := f.foo(); err != nil {
		return fmt.Errorf("foo: %w", err)
	}

	if err := f.bar(); err != nil {
		return fmt.Errorf("bar: %w", err)
	}

	if err := f.bizz(); err != nil {
		return fmt.Errorf("bizz: %w", err)
	}

	if err := f.zee(); err != nil {
		return fmt.Errorf("zee: %w", err)
	}

	return nil
}
```
