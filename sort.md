```go
package main

import (
	"log"
	"sort"
)

type Steps []uint

func (a Steps) Len() int           { return len(a) }
func (a Steps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Steps) Less(i, j int) bool { return a[i] < a[j] }

func main() {
	res := Steps{1, 4, 2, 10}
	sort.Sort(res)
	log.Println(res)

	sort.Sort(sort.Reverse(res))
	log.Println(res)
}
```


## Sort by multiple keys

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
	"strings"
)

type User struct {
	Name string
	Age  int
}

func main() {
	users := []User{
		{"A", 10},
		{"B", 10},
		{"C", 20},
		{"D", 5},
		{"A", 100},
	}
	sort.Slice(users, func(i, j int) bool {
		lhs, rhs := users[i], users[j]
		byAge := lhs.Age - rhs.Age
		byName := strings.Compare(lhs.Name, rhs.Name) // Returns 0 if equal, -1 if lhs is less than rhs, and 1 if lhs is greater than rhs
		
		// The + sign is not necessary, but adds clarity as it means increasing in value, aka ascending.
		// sortBy(+byAge, +byName) // Sort by age asc, by name asc
		// sortBy(-byAge, +byName) // Sort by age desc, by name asc
		// sortBy(+byName, +byAge) // Sort by name asc, by age asc
		return sortBy(-byAge, -byName) // Sort by age desc, by name desc
	})
	fmt.Println(users)
}

func sortBy(sc ...int) bool {
	for _, c := range sc {
		if c != 0 {
			return c < 0
		}
	}
	return sc[len(sc)-1] < 0
}
```

## Reverse sort

Reverse the slice. Note that this is different than sorting in ascending or descending order (see below for descending sort).

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
)

func main() {
	n := []int64{1, 2, 5, 6, 4}
	sort.Slice(n, func(i, j int) bool {
		return true
	})
	fmt.Println(n)

	s := []string{"a", "b", "c", "k", "d"}
	sort.Slice(s, func(i, j int) bool {
		return true
	})
	fmt.Println(s)
}
```


Sort descending:
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
)

func main() {
	keys := []int{3, 2, 8, 1}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	fmt.Println(keys)
}
```

## Partitioning with sort


Say if you have a list of items with id and stock counts, and some of them are _out of stock_. You want to show the out of stock items first (or last, depending on your requirements). There are many ways to do it, such as filtering them and combining them later etc, but most of the time, you want to preserve the order of the existing stock counts. 

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
)

func main() {
	items := []item{
		{1, 0},
		{2, 10},
		{3, 30},
		{4, 30},
		{5, 0},
		{6, 5},
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].count == 0 && items[j].count != 0
	})
	fmt.Println(items) // [{1 0} {5 0} {2 10} {3 30} {4 30} {6 5}]

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].count != 0 && items[j].count == 0
	})
	fmt.Println(items) // [{2 10} {3 30} {4 30} {6 5} {1 0} {5 0}]
}

type item struct {
	id    int64
	count int64
}
```


## Reversing generic slice

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int64{10, 5, 15, 20, 1, 100, -1}
	ReverseSlice(nums)
	fmt.Println(nums)

	strs := []string{"hello", "world"}
	ReverseSlice(strs)
	fmt.Println(strs)

	runes := []rune{'h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd'}
	ReverseSlice(runes)
	for _, r := range runes {
		fmt.Print(string(r), " ")
	}
	fmt.Println()
	
	// For maps, implement it differently
	m1 := map[string]string{"a": "1"}
	m2 := map[string]string{"b": "2"}
	m3 := map[string]string{"c": "3"}
	maps := []map[string]string{m1, m2, m3}
	reverse(maps)
	fmt.Println(maps)
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
```

Output:
```
[-1 100 1 20 15 5 10]
[world hello]
d l r o w o l l e h 
[map[c:3] map[b:2] map[a:1]]
Program exited.
```
