# Basic Overwrite

```go
package main

import (
	"fmt"
)

type A interface {
	Spew() string
}

type aImpl struct {
}

func (a *aImpl) Spew() string {
	return "hello"
}

type bImpl struct {
	A
}

// This will overwrite the A's Spew(), but A's Spew() can still be accessed by b.A.Spew()
func (b *bImpl) Spew() string {
	return "world"
}

func main() {
	a := new(aImpl)
	b := bImpl{a}

  // Use the latter (shorthand) version, it keeps the code shorter and more concise.
	fmt.Println(b.A.Spew(), b.Spew())
}
```
