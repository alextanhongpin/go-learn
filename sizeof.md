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

	// From https://blog.golang.org/strings:
	// In short, Go source code is UTF-8, so the source code for the string literal is UTF-8 text.
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
Output:
```
empty struct struct{}{} is 0 bytes
struct is 12 bytes

default int is 4 bytes
int8 is 1 byte
int16 is 2 bytes
int32 is 4 bytes
int64 is 8 bytes

rune 'a' is 4 bytes
string "a" is 8 bytes
string "ab" is 8 bytes

bool "true" is 1 byte
bool "false" is 1 byte

slice []int{1,2,3} is 12 bytes
string "hello" is 8 bytes
```
