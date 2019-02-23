## Min, max

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(min(1, -2))
	fmt.Println(max(1, -2))
}

func min(hd int, rest ...int) int {
	o := hd
	for _, r := range rest {
		if r < o {
			o = r
		}
	}
	return o
}

func max(hd int, rest ...int) int {
	o := hd
	for _, r := range rest {
		if r > o {
			o = r
		}
	}
	return o
}
```

## Concurrent Set

```go
package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Name string
	Age  int
}

func main() {
	set := NewSet()
	set.Add(1)
	set.Add(20)
	fmt.Println(set.Has(1)) // true
	set.Remove(1)
	fmt.Println(set.Has(1))   // false
	fmt.Println(set.Has(20))  // true
	fmt.Println(set.Has(100)) // false

	var r1 = Result{Name: "John"}
	set.Add(r1)

	fmt.Println(set.Has(r1))                            // true
	fmt.Println(set.Has(Result{Name: "John"}))          // true
	fmt.Println(set.Has(Result{Name: "John", Age: 10})) // false
	fmt.Println(set.Has(Result{Name: "Doe"}))           // false
	set.Remove(r1)
	fmt.Println(set.Has(r1)) // false
}

type Set struct {
	mu    *sync.RWMutex
	value map[interface{}]struct{}
}

func NewSet() *Set {
	return &Set{
		mu:    new(sync.RWMutex),
		value: make(map[interface{}]struct{}),
	}
}

func (s *Set) Add(val interface{}) {
	s.mu.Lock()
	s.value[val] = struct{}{}
	s.mu.Unlock()
}

func (s *Set) Has(val interface{}) bool {
	s.mu.RLock()
	_, exist := s.value[val]
	s.mu.RUnlock()
	return exist
}

func (s *Set) Remove(val interface{}) {
	s.mu.Lock()
	delete(s.value, val)
	s.mu.Unlock()
}
```

## Enums

```go
package main

import (
	"fmt"
	"strings"
)

const (
	West  = "west"
	South = "south"
	East  = "east"
	North = "north"
)

func main() {
	west := NewEnum("west")
	fmt.Println(west.Is("West"))
	fmt.Println(west.Is("west"))
	fmt.Println(west.IsStrict("West"))

	directions := NewEnums(West, South, East, North)
	fmt.Println(directions.IsValid("west"))
	fmt.Println(directions.IsValid("northeast"))
}

type Enum struct {
	value string
}

func NewEnum(e string) Enum {
	return Enum{e}
}

func (e *Enum) OneOf(values ...string) bool {
	for _, val := range values {
		if eq := strings.EqualFold(e.value, val); eq {
			return eq
		}
	}
	return false
}

func (e *Enum) OneOfStrict(values ...string) bool {
	for _, val := range values {
		if eq := e.value == val; eq {
			return eq
		}
	}
	return false
}

func (e *Enum) Is(value string) bool {
	return strings.EqualFold(e.value, value)
}

func (e *Enum) IsStrict(value string) bool {
	return e.value == value
}

type Enums struct {
	values []string
}

func NewEnums(values ...string) Enums {
	return Enums{values}
}

func (e *Enums) IsValid(val string) bool {
	for _, e := range e.values {
		if eq := strings.EqualFold(e, val); eq {
			return true
		}
	}
	return false
}
```

## Strings comparison with levenshtein

This is particularly useful when we do not know what type of string to expect (mapping csv columns to struct), and where equality is not important (Uppercase, lowercase etc). We are only looking for the probability that the two strings might be equal.
```go
package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func main() {
	fmt.Println(score("IC number", "ic number"))
}

func score(a, b string) float64 {
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	parse := func(a string) string {
		a = strings.ToLower(a)
		a = reg.ReplaceAllString(a, "")
		res := strings.Split(a, " ")
		sort.Strings(res)
		return strings.Join(res, " ")
	}
	a = parse(a)
	b = parse(b)
	l := max(len(a), len(b))
	return 1 - (float64(levenshtein([]rune(a), []rune(b))) / float64(l))
}

func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = min(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func min(hd int, rest ...int) int {
	o := hd
	for _, v := range rest {
		if v < o {
			o = v
		}
	}
	return o
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```
