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


## Replace

```go
package main

import (
	"log"
	"strings"
)

func main() {

	// Replace A with T, and vice-versa.
	// Replace C with G, and vice-versa.
	var (
		dna      = "ATTTACGATGC"
		expected = "TAAATGCTACG"
	)

	var dnaReplacer = strings.NewReplacer(
		"A", "T",
		"T", "A",
		"C", "G",
		"G", "C",
	)

	if result := dnaReplacer.Replace(dna); result != expected {
		log.Fatalf("expected %s, got %s", expected, result)
	}

	replaceFn := func(r rune) rune {
		switch r {
		case 'A':
			return 'T'
		case 'T':
			return 'A'
		case 'C':
			return 'G'
		case 'G':
			return 'C'
		}
		return r
	}
	if result := strings.Map(replaceFn, dna); result != expected {
		log.Fatalf("expected %s, got %s", expected, result)
	}
}
```


## Matching string the right way (exact match)

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// Matches exact.
	fmt.Println(match("food", "foo"))
	fmt.Println(match("(foo=bar)", "foo"))

	// False positive. Matches partial.
	fmt.Println(strings.Contains("food", "foo"))
	fmt.Println(strings.Contains("(foo=bar)", "foo"))
}

func match(src, tgt string) bool {
	// A simpler way is to just match src == tgt without regex.
	// However, using regex allows us to replace string more easily.
	ok, _ := regexp.MatchString(fmt.Sprintf("\\b%s\\b", tgt), src)
	return ok
}
```

## Runes length
Demonstrates that runes length are different than bytes.
```go
package main

import (
	"fmt"
)

func main() {
	// len(s) vs len([]rune(s))
	// ‘£’ takes two bytes as per UTF-8
	fmt.Println(len("£"))
	fmt.Println(len([]rune("£")))
}
```

## Unicode Normalization

For password, it is important to normalize them by using NFKD, so that they are treated the same regardless of what machine the user is on.
```go
package main

import (
	"fmt"

	"golang.org/x/text/unicode/norm"
)

func main() {
	// latin small letter e with acute
	e1 := "\u00e9"
	fmt.Println(e1)
	fmt.Println([]byte(e1))
	fmt.Println(norm.NFKC.Bytes([]byte(e1)))
	fmt.Println(norm.NFKD.Bytes([]byte(e1)))
	fmt.Println(len(e1))         // 2
	fmt.Println(len([]rune(e1))) // 1

	// latin small letter e followed by combining acute accent
	e2 := "\u0065\u0301"
	fmt.Println(e2)
	fmt.Println([]byte(e2))
	fmt.Println(norm.NFKC.Bytes([]byte(e2)))
	fmt.Println(norm.NFKD.Bytes([]byte(e2)))
	fmt.Println(len(e2))         // 3, a character can span multiple runes, len() return the number of bytes in a string
	fmt.Println(len([]rune(e2))) // 2, this counts the actual length of the character
}
```
