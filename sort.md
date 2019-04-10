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
