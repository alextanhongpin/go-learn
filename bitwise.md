```golang
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

## Another example

```go
package main

import "fmt"

type Code int

const (
	None Code = 1 << iota
	A
	B
	C
)

func (c Code) Has(codes Code) bool {
	return c & codes != 0
}

func (c Code) Is(codes Code) bool {
	// c | codes == codes
	return c&codes == c|codes
}

func main() {
	var abc Code
	abc |= A
	abc |= B
	abc |= C
	abc |= A // Pass A twice, but there won't be duplicate.

	fmt.Println("has A", abc.Has(A))
	fmt.Println("has B", abc.Has(B))
	fmt.Println("has C", abc.Has(C))
	fmt.Println("has A, B", abc.Has(A|B))
	fmt.Println("has B, C", abc.Has(B|C))
	fmt.Println("has A, C", abc.Has(A|C))
	fmt.Println("has A, B, C", abc.Has(A|B|C))

	fmt.Println("is A", abc.Is(A))
	fmt.Println("is B", abc.Is(B))
	fmt.Println("is C", abc.Is(C))
	fmt.Println("is A, B", abc.Is(A|B))
	fmt.Println("is B, C", abc.Is(B|C))
	fmt.Println("is A, C", abc.Is(A|C))
	fmt.Println("is A, B, C", abc.Is(A|B|C))

}
```

## Use Case: Ensuring Steps are completed in order for state machine

```go
package main

import (
	"fmt"
)

type State uint

const (
	Initialised State = 1 << iota
	Checkout
	Submitted
	Verified
	Completed
)

var Steps = []State{Initialised, Checkout, Submitted, Verified, Completed}

func main() {
	validateSteps := func(state State, i int) bool {
		for j, step := range Steps[:i] {
			if step&state == 0 {
				fmt.Println("skipped", j)
				return false
			}
		}
		return true
	}
	{
		state := Initialised | Checkout
		valid := validateSteps(state, 2)
		fmt.Println(valid)
	}

	{
		state := Initialised | Completed
		valid := validateSteps(state, 2)
		fmt.Println(valid)
	}

	{
		state := Initialised | Checkout | Submitted | Verified | Completed
		valid := validateSteps(state, 4)
		fmt.Println(valid)
	}

	{
		state := Initialised | Checkout | Submitted | Verified | Completed
		valid := validateSteps(state, 5)
		fmt.Println(valid)
	}
}
```
