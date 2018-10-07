# Size

Finding the size of each type in Golang using the `unsafe.Sizeof` package:

```golang
package main

import (
	"fmt"
	"unsafe"
)

type Sample struct {
	Name string // Takes 8 bytes of storage.
	Age  int32  // Takes 4 bytes of storage.
}

func main() {
	// All measurements are in bytes.
	print(`empty struct struct{}{}`, unsafe.Sizeof(struct{}{}))
	print("struct", unsafe.Sizeof(Sample{}))

	fmt.Println("")

	print("default int", unsafe.Sizeof(int(1)))
	print("int8", unsafe.Sizeof(int8(1)))
	print("int16", unsafe.Sizeof(int16(1)))
	print("int32", unsafe.Sizeof(int32(1)))
	print("int64", unsafe.Sizeof(int64(1)))

	fmt.Println("")

	print(`rune 'a'`, unsafe.Sizeof('a'))
	print(`string "a"`, unsafe.Sizeof("a"))
	print(`string "ab"`, unsafe.Sizeof("ab"))

	fmt.Println("")

	print(`bool "true"`, unsafe.Sizeof(true))
	print(`bool "false"`, unsafe.Sizeof(false))
	
	fmt.Println("")

	print(`slice []int{1,2,3}`, unsafe.Sizeof([]int{1,2,3}))
	print(`string "hello"`, unsafe.Sizeof("hello"))
	
}

func print(name string, size uintptr) {
	label := "bytes"
	if size == 1 {
		label = "byte"
	}
	fmt.Println(name, "is", size, label)
}
```
