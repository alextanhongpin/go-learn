```golang
package main

import (
	"fmt"
)

// Skipping the first value.
type Code int

const (
	_ Code = iota
	A
	B
	C
)

// Starts from zero.
type Status int

const (
	X Status = iota
	Y
	Z
)

// Bitwise enum.
const (
	M int = 1 << iota
	N
	O
)

func main() {
	fmt.Println(A, B, C)
	fmt.Println(X, Y, Z)
	fmt.Println(M, N, O)
}
```


## Bitwise Operation

```
package main

import (
	"fmt"
)

type Code int

const (
	A Code = 1 << iota
	B
	C
	D
)

func (c Code) Is(cc Code) bool {
	return c&cc == c|c
}

func (c Code) Has(cc Code) bool {
	return c&cc != 0
}
func main() {
	fmt.Println(A, B, C, D)
	var c Code
	c = 3
	fmt.Println(c.Is(A), c.Is(B), c.Is(C), c.Is(D))
	fmt.Println(c.Has(A), c.Has(B), c.Has(C), c.Has(D))
	fmt.Println(A.Is(c), B.Is(c), C.Is(c), D.Is(c))
}
```
