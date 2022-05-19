```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
	a := Address{State: "s"}
	fmt.Println(a.Valid())
}

// All address must be filled or remains empty.
type Address struct {
	Street1    string
	Street2    string
	PostalCode string
	State      string
	City       string
}

func (a *Address) Valid() bool {
	var count, total int

	boolInt := func(b bool) int {
		total++
		if b {
			return 1
		}
		return 0
	}

	count += boolInt(a.Street1 == "")
	count += boolInt(a.Street2 == "")
	count += boolInt(a.PostalCode == "")
	count += boolInt(a.State == "")
	count += boolInt(a.City == "")

	return count%total == 0
}
```
