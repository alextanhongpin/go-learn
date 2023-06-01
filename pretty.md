## Pretty print struct


```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strings"
	"unicode"
)

type User struct {
	Name    string
	age     int
	hobbies []string
	data    map[string]any
}

func main() {
	u := User{
		Name:    "John",
		age:     10,
		hobbies: []string{"hello"},
		data: map[string]any{
			"profile":   "https://img.com",
			"isMarried": false,
		},
	}
	s := fmt.Sprintf("%#v", u)
	fmt.Println(s)
	fmt.Println(format(u))
	fmt.Println(format([]int{1, 2, 3}))
}

func countLeadingSpace(s string) int {
	return len(s) - len(strings.TrimSpace(s))
}

func format(v any) string {
	s := fmt.Sprintf("%#v", v)
	var res []string
	push := func(s string) {
		res = append(res, strings.TrimSuffix(s, " "))
	}

	sb := new(strings.Builder)

	for i, r := range s {
		h := i - 1
		j := i + 1
		k := i + 2
		if i == len(s)-1 {
			j = i
			k = i
		} else if i == 0 {
			h = i
		}
		nextIsDigit := unicode.IsDigit(rune(s[j]))
		nextIsLetter := unicode.IsLetter(rune(s[j]))
		nextIsQuote := s[j] == '"'
		nextIsArray := s[j] == '[' && s[k] == ']'

		if r == '{' && s[j] != '}' { // Skip if map[string]interface{}
			sb.WriteRune(r)
			push(sb.String())
			sb.Reset()

			var n int
			if len(res) > 0 {
				last := res[len(res)-1]
				n = countLeadingSpace(last)
			}
			n += 2
			sb.WriteString(strings.Repeat(" ", n))
		} else if r == '}' && s[h] != '{' {
			push(sb.String())
			sb.Reset()

			var n int
			if len(res) > 0 {
				last := res[len(res)-1]
				n = countLeadingSpace(last)
			}
			n -= 2
			sb.WriteString(strings.Repeat(" ", n))
			sb.WriteRune(r)
		} else if r == ',' && s[j] == ' ' {
			// New line.
			sb.WriteRune(r)
			push(sb.String())
			sb.Reset()

			var n int
			if len(res) > 0 {
				last := res[len(res)-1]
				n = countLeadingSpace(last)
			}
			n -= 1 // Because there will be a leading space, we deduct 1.
			sb.WriteString(strings.Repeat(" ", n))
		} else if r == ':' && (nextIsDigit || nextIsLetter || nextIsArray || nextIsQuote) { // To avoid matching http://
			sb.WriteRune(r)
			sb.WriteRune(' ') // Add a space after colon.
		} else {
			sb.WriteRune(r)
		}
	}
	push(sb.String())
	return strings.Join(res, "\n")
}
```

Output:

```go
main.User{Name:"John", age:10, hobbies:[]string{"hello"}, data:map[string]interface {}{"isMarried":false, "profile":"https://img.com"}}
main.User{
  Name: "John",
  age: 10,
  hobbies: []string{
    "hello"
  },
  data: map[string]interface {}{
    "isMarried": false,
    "profile": "https://img.com"
  }
}
[]int{
  1,
  2,
  3
}
```
