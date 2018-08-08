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


## Slice with different len and capacity

```
package main

import "fmt"

func main() {
	vals := make([]int, 5)
	fmt.Println("[vals] start", vals, len(vals), cap(vals))
	for i := 0; i < 5; i++ {
		vals = append(vals, i)
	}
	fmt.Println("[vals] end", vals, len(vals), cap(vals))

	var ints []int
	fmt.Println("[ints] start", ints, len(ints), cap(ints))
	for i := 0; i < 5; i++ {
		ints = append(ints, i)
	}
	fmt.Println("[ints] end", ints, len(ints), cap(ints))

	anotherInts := make([]int, 0, 2)
	fmt.Println("[anotherInts] start", anotherInts, len(anotherInts), cap(anotherInts))

	for i := 0; i < 5; i++ {
		anotherInts = append(anotherInts, i)
	}
	fmt.Println("[anotherInts] end", anotherInts, len(anotherInts), cap(anotherInts))
}
```

Output:

```
[vals] start [0 0 0 0 0] 5 5
[vals] end [0 0 0 0 0 0 1 2 3 4] 10 12
[ints] start [] 0 0
[ints] end [0 1 2 3 4] 5 8
[anotherInts] start [] 0 2
[anotherInts] end [0 1 2 3 4] 5 8
```
