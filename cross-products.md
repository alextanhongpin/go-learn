```go
package main

import (
	"fmt"
)

func main() {
	out := crossProducts(
		[]string{"a", "b"},
		[]string{"left", "right"},
		[]string{"x", "y", "z"},
	)
	fmt.Println(out, len(out))
}

func crossProducts(a []string, b ...[]string) []string {
	if len(b) == 0 {
		return a
	}
	result := make([]string, len(a)*len(b[0]))
	for i, m := range a {
		for j, n := range b[0] {
			result[i*len(b[0])+j] = m + ":" + n
		}
	}
	return crossProducts(result, b[1:]...)
}
```
