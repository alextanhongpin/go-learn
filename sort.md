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
