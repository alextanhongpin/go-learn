```golang
package main

import (
	"fmt"
)

// Skipping the first value.
type Code int

const (
	_ Code = iota
	A
	B
	C
)

// Starts from zero.
type Status int

const (
	X Status = iota
	Y
	Z
)

// Bitwise enum.
const (
	M int = 1 << iota
	N
	O
)

func main() {
	fmt.Println(A, B, C)
	fmt.Println(X, Y, Z)
	fmt.Println(M, N, O)
}
```


## Bitwise Operation

```go
package main

import (
	"fmt"
)

type Code int

const (
	A Code = 1 << iota
	B
	C
	D
)

func (c Code) Is(cc Code) bool {
	return c&cc == c|c
}

func (c Code) Has(cc Code) bool {
	return c&cc != 0
}
func main() {
	fmt.Println(A, B, C, D)
	var c Code
	c = 3
	fmt.Println(c.Is(A), c.Is(B), c.Is(C), c.Is(D))
	fmt.Println(c.Has(A), c.Has(B), c.Has(C), c.Has(D))
	fmt.Println(A.Is(c), B.Is(c), C.Is(c), D.Is(c))
}
```


## Enumerable with Comparable



Comparable strings enums in golang. Use case is to check for statuses, or if Tuesday is greater than Monday. For statuses (state machine), we want to know if the next transition is allowable by comparing if the next state is larger than the previous ones (this is assuming that the next state is always greater than the previous one and the state only goes in one direction). Other examples includes tier pricing, we want to know if a plan is upgradable or downgradable by checking the enum type of each tier plan. In order for the products to be upgraded or downgraded, the tier values should not be the same.

```go
package main

import (
	"fmt"
)

type Enumerable int

const (
	Less    Enumerable = -1
	Equal   Enumerable = 0
	Greater Enumerable = 1
)

type Tier string

var TierEnumerable = map[Tier]int{
	Bronze: 1,
	Silver: 2,
	Gold:   4,
}

const (
	Bronze Tier = "bronze"
	Silver Tier = "silver"
	Gold   Tier = "gold"
)

func (t Tier) Cmp(tt Tier) Enumerable {
	l, r := TierEnumerable[t], TierEnumerable[tt]
	if l > r {
		return Greater
	}
	if l < r {
		return Less
	}
	return Equal
}

func main() {
	switch Gold.Cmp(Silver) {
	case Less:
		fmt.Println("is less")
	case Equal:
		fmt.Println("is equal")
	case Greater:
		fmt.Println("is greater")
	}
}
```



## Loading Enums in runtime

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strings"
)

var Direction = NewEnum("direction", `Up Down Left Right`)

func main() {
	fmt.Println(Direction.Is("Up"), Direction.Is(0), Direction.Is(1))
}

type Enum int

type Enums struct {
	name     string
	min, max int
	value    map[Enum]string
}

func (e Enums) Name() string {
	return e.name
}

func (e Enums) Is(v interface{}) bool {
	switch i := v.(type) {
	case string:
		return e.isString(i)
	case int:
		return e.isInt(i)
	default:
		panic("invalid type")
	}
}

func (e Enums) isInt(n int) bool {
	return n >= e.min && n <= e.max
}

func (e Enums) isString(s string) bool {
	for k := range e.value {
		if e.value[k] == s {
			return true
		}
	}
	return false
}

func NewEnum(name, in string) Enums {
	enums := strings.Fields(strings.TrimSpace(in))
	min, max := 1, 0
	value := make(map[Enum]string)
	for i, e := range enums {
		max = min + i
		value[Enum(max)] = strings.TrimSpace(e)
	}
	return Enums{
		name:  name,
		min:   min,
		max:   max,
		value: value,
	}
}
```
