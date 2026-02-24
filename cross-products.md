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


## Generics

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	a := []int{1, 2}
	b := []int{4, 4}
	c := []int{8, 10}
	fmt.Println(combinations(a, b, c))
	fmt.Println("Hello, 世界")
}

func combinations[T any](vs ...[]T) [][]T {
	var res [][]T
	for len(vs) > 0 {
		var h []T
		h, vs = vs[0], vs[1:]
		if len(res) == 0 {
			for _, v := range h {
				res = append(res, []T{v})
			}
			continue
		}
		var tmp [][]T
		for _, v := range h {
			for _, w := range res {
				tmp = append(tmp, append(w, v))
			}
		}
		res = tmp
	}
	return res
}
```

Another alternative:
```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	a := []int{1, 2}
	b := []int{4, 5}
	c := []int{7, 8}
	fmt.Println(combinations(a, b, c))

}

func combinations[T any](vs ...[]T) [][]T {
	var res [][]T
	var total = 1
	for i := range vs {
		total *= len(vs[i])
	}
	for i := range total {
		var tmp []T
		var num = 1
		for j := range len(vs) {
			den := num
			num *= len(vs[j])
			tmp = append(tmp, vs[j][i%num/den])
		}
		res = append(res, tmp)
	}
	return res
}
```
