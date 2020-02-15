```go
package main

import (
	"errors"
	"fmt"
)

var ErrOne = errors.New("one")

func main() {
	e1 := ErrOne
	e2 := fmt.Errorf("two: %w", e1)
	e3 := fmt.Errorf("three: %w", e2)

	fmt.Println(e1)
	fmt.Println(e2)
	fmt.Println(e3)
	fmt.Println(errors.Unwrap(e2))
	fmt.Println(errors.Unwrap(e3))
	fmt.Println(errors.Unwrap(errors.Unwrap(e3)))
	fmt.Println(errors.Is(e1, ErrOne))
	fmt.Println(errors.Is(e2, ErrOne))	
	fmt.Println(errors.Is(e3, ErrOne))	
}
```
