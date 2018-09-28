## Type Alias in Golang

Surprisingly, both `A` and `B` now share the same struct methods!:

```golang
package main

import (
	"fmt"
	"log"
)

type A struct {
	ID string
}

type B = A

func (b *B) Print() {
	log.Println("my id", b.ID)
}

func main() {
	a := A{"1"}
	a.Print()
	
	b := B{"2"}
	b.Print()
}
```

Output:

```
2009/11/10 23:00:00 my id 1
2009/11/10 23:00:00 my id 2
```
