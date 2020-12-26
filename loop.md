# Breaking loop

A cleaner way to break from loop after the conditions inside match.

```go
package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}

	// Print items in "a" that is not in "b".
	for _, i := range a {
		var match bool
		for _, j := range b {
			if i == j {
				match = true
				break
			}
		}
		if match {
			continue
		}
		fmt.Println(i)
	}

mainloop:
	for _, i := range a {
		for _, j := range b {
			if i == j {
				continue mainloop
			}
		}
		fmt.Println(i)
	}
}

```
