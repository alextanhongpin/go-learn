## Base62 

```go
package main

import (
	"fmt"
)

func main() {
	var base62Alphabets = [...]rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	var hashDigits []int
	dividend := 100
	remainder := 0
	for dividend > 0 {
		remainder = dividend % 62
		dividend = dividend / 62
		// Prepend
		hashDigits = append([]int{remainder}, hashDigits...)
	}

	fmt.Println(hashDigits, dividend, remainder)

	var hashString string
	for _, v := range hashDigits {
		hashString += string(base62Alphabets[v])
	}
	fmt.Println(hashString)
}
```

## Getting rune value as int

```go
func main() {
	// Getting rune as int.
	var total int
  // Note that "hello world" and "helol world" will return the same total.
	for _, v := range []rune("helol world") {
		total += int(v - '0')
	}
	fmt.Println("total", total)
}
```
