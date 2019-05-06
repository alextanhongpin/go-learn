## Comparing strings effectively

tl;dr, use `strings.EqualFold` for comparison:

```go
// Good
if ok := strings.ToLower(a) == strings.ToLower(b); ok {}

// Better
if ok := strings.EqualFold(a, b); ok {}
```

References:
- https://www.digitalocean.com/community/questions/how-to-efficiently-compare-strings-in-go

## Remove duplicate

```go
package main

import (
	"fmt"
	"log"
	"strings"
	"testing/quick"
)

func main() {
	fmt.Println(removeDups("hello aaa aaa aaaaaplayground"))
	f := func(s string) bool {
		o := removeDups(s)
		return checkDups(o)
	}
	if err := quick.Check(f, nil); err != nil {
		log.Fatal(err)
	}
}

func removeDups(str string) string {
	var out strings.Builder
	m := make(map[rune]bool)
	for _, s := range str {
		if m[s] {
			continue
		}
		m[s] = true
		out.WriteRune(s)
	}
	return out.String()
}

func checkDups(str string) bool {
	m := make(map[rune]int, 0)
	for _, s := range str {
		m[s]++
		if m[s] > 1 {
			return false
		}
	}
	return true
}
```
