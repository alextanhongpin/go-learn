## Create fixed bytes from string

```go
var b [64]byte
copy(b[:], "hello world")

// b is now [64]byte

// This converts array b to slice b
fmt.Println(b[:])
```

## Remove item from slice

Where a is the slice, and i is the index of the element you want to delete:

```go
package main

import "fmt"

func main() {
	var a []string
	a = []string{"a", "b", "c", "d"}

	// Index to remove
	i := 1
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
	// Returns [a c d]
}
```
