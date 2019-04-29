```go
package main

import (
	"fmt"
	"unsafe"
)

type A struct {
	a bool    // 1 byte
	b float64 // 8 bytes
	c int32   // 4 bytes
}

type B struct {
	a bool    // 1 byte
	c int32   // 4 bytes
	b float64 // 8 bytes
}

func main() {
	fmt.Println(unsafe.Sizeof(A{})) // 24 bytes
	fmt.Println(unsafe.Sizeof(B{})) // 16 bytes
	{
		var i int8
		fmt.Println("int8", unsafe.Sizeof(i))
	}
	{
		var i int16
		fmt.Println("int16", unsafe.Sizeof(i))
	}
	{
		var i int32
		fmt.Println("int32", unsafe.Sizeof(i))
	}
	{
		var i int64
		fmt.Println("int64", unsafe.Sizeof(i))
	}
}
```

References:
- https://medium.com/@felipedutratine/how-to-organize-the-go-struct-in-order-to-save-memory-c78afcf59ec2
