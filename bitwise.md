```
package main

import (
	"fmt"
)

const (
	code1 uint32 = 1 << iota
	code2
	code3
	code4
)

func main() {
	var i uint32
	i |= code1
	i |= code4
	fmt.Println("list: [code1, code2]")

	fmt.Println("codes:", code1, code2, code3, code4)
	fmt.Println("")

	// Check individual codes
	fmt.Println("has code 1:", i&code1 != 0)
	fmt.Println("has code 2:", i&code2 != 0)
	fmt.Println("has code 3:", i&code3 != 0)
	fmt.Println("has code 4:", i&code4 != 0)
	fmt.Println("")

	// Check if either codes exist
	fmt.Println("code 2 or 3 exist:", i&(code2|code3) != 0)
	fmt.Println("code 2 or 4 exist:", i&(code2|code4) != 0) // This will return true, because 4 is present even if 2 is not
	fmt.Println("code 1 or 4 exist:", i&(code1|code4) != 0) // To check if 1 or 4 is present
	fmt.Println("")

	// Check exact match
	fmt.Println("code 2 or 3 exist", i == (code2|code3))
	fmt.Println("code 2 or 4 exist", i == (code2|code4))
	fmt.Println("code 1 or 4 exist", i == (code1|code4)) // To check if 1 and 4 is present
}
```

Output:

```
list: [code1, code2]
codes: 1 2 4 8
has code 1: true
has code 2: false
has code 3: false
has code 4: true

code 2 or 3 exist: false
code 2 or 4 exist: true
code 1 or 4 exist: true

code 2 or 3 exist false
code 2 or 4 exist false
code 1 or 4 exist true
```
