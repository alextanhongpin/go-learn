# Generic Set
```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	s := New[string]()
	s.Add("hello")
	fmt.Println(s.Has("hello"))

	si := New[int]()
	si.Add(1, 2, 3)
	fmt.Println(si.Has(1))

	s2 := New[int]()
	s2.Add(1, 2)
	fmt.Println(si.Intersect(*s2))
	fmt.Println(si.List())
}

type Set[T comparable] map[T]bool

func New[T comparable]() *Set[T] {
	return &Set[T]{}
}

func (s Set[T]) Add(t ...T) {
	for _, k := range t {
		s[k] = true
	}
}

func (s Set[T]) Has(t T) bool {
	return s[t]
}

func (s Set[T]) Intersect(other Set[T]) []T {
	if len(s) > len(other) {
		return other.Intersect(s)
	}
	var result []T
	for k := range s {
		if other.Has(k) {
			result = append(result, k)
		}
	}
	return result
}

func (s Set[T]) List() []T {
	var result []T
	for k := range s {
		result = append(result, k)
	}
	return result
}
```

## Generic Map

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	users := []User{
		{"John", 10},
		{"Jane", 20},
	}
	ages := mapR(users, func(u User) int {
		return u.Age
	})
	fmt.Println(ages)
}

func mapR[T any, R any](in []T, fn func(T) R) []R {
	res := make([]R, len(in))
	for i, k := range in {
		res[i] = fn(k)
	}
	return res
}
```
