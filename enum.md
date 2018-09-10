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
