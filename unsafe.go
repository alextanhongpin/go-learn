# Converting one struct to another using unsafe

Works for
- private fields can be set
- sequence of types must be correct

```go
package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

func main() {
	pdb := personDB{
		Age:     10,
		Name:    "john",
		Married: true,
	}

	p := (*Person)(unsafe.Pointer(&pdb))
	fmt.Println(p)
}

type Person struct {
	// The sequence of types must match.
	name    string
	age     int
	married bool
}

type personDB struct {
	Name    string
	Age     int
	Married bool
	Extra   json.RawMessage // This field is ignored.
}
```
