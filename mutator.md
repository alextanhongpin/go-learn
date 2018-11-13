# Mutator pattern

Looking for better ways to handle dependencies aside from interface.

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// Default function.
	fmt.Println(NewRandom(10))
	
	// Function with overwrite capabilities.
	fmt.Println(NewRandom(10, func() int { return 100 }))
}

type RandomModifier func() int

func NewRandom(n int, fns ...RandomModifier) int {
	if len(fns) == 0 {
		return rand.Intn(n)
	}
	return fns[0]()
}
```
