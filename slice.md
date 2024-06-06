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

## Append item in slice

```go
package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3}

	// Remove from an index.
	i := 1
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)

	// Append at index.
	a = append(a[:i+1], a[i:]...)
	a[i] = 100
	fmt.Println(a)
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

## Stack and Queue

```go
package main

import (
	"fmt"
)

// First in, first out (FIFO)
type Queue struct {
	list []int
}

func NewQueue() *Queue {
	return &Queue{
		list: make([]int, 0),
	}
}

func (q *Queue) Push(i ...int) {
	q.list = append(q.list, i...)
}

func (q *Queue) Size() int {
	return len(q.list)
}
func (q *Queue) Pop() int {
	var head int
	head, q.list = q.list[0], q.list[1:]
	return head
}

type Stack struct {
	list []int
}

// Last-in, first-out (LIFO).
func NewStack() *Stack {
	return &Stack{
		list: make([]int, 0),
	}
}
func (s *Stack) Push(i ...int) {
	s.list = append(s.list, i...)
}
func (s *Stack) Size() int {
	return len(s.list)
}
func (s *Stack) Pop() int {
	var tail int
	s.list, tail = s.list[:len(s.list)-1], s.list[len(s.list)-1]
	return tail
}

func main() {
	fmt.Println("QUEUE")
	q := NewQueue()
	q.Push(1, 2, 3, 4, 5)
	for q.Size() > 0 {
		fmt.Println(q.Pop())
	}

	fmt.Println("STACK")
	s := NewStack()
	s.Push(1, 2, 3, 4, 5)
	for s.Size() > 0 {
		fmt.Println(s.Pop())
	}
}
```

## Diff

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}
	removed := diff(b, a)
	added := diff(a, b)
	fmt.Println(removed, added)
	fmt.Println("Hello, 世界")
}

func diff[T ~[]V, V comparable](a, b T) T {
	m := make(map[V]bool)
	for _, v := range a {
		m[v] = true
	}
	var res T
	for _, v := range b {
		if !m[v] {
			res = append(res, v)
		}
	}
	return res
}
```
