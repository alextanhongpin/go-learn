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

// https://medium.com/@felipedutratine/how-to-organize-the-go-struct-in-order-to-save-memory-c78afcf59ec2
// Each of them hold an array of size 8
type BadOrder struct {
	Bool1 bool // Takes 1 byte of storage.
	Name string // Takes 8 bytes of storage.
	Bool2 bool // Takes 1 bytes of storage.
}

// Age and IsBool can be placed in the same bucket.
type GoodOrder struct {
	Name string // Takes 8 bytes of storage.
	// This would be grouped together
	Bool1 bool // Takes 1 byte of storage.
	Bool2 bool // Takes 1 byte of storage.
}

func main() {
	// All measurements are in bytes.
	print(`empty struct struct{}{}`, unsafe.Sizeof(struct{}{}))
	print("struct", unsafe.Sizeof(Sample{}))
	print("GoodOrder", unsafe.Sizeof(GoodOrder{}))
	print("BadOrder", unsafe.Sizeof(BadOrder{}))

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

	print(`slice []int{1,2,3}`, unsafe.Sizeof([]int{1, 2, 3}))
	print(`string "hello"`, unsafe.Sizeof("hello"))

	fmt.Println("")
	
	// To compare size
	print("string", unsafe.Sizeof("abcdefghijklmnopqrstuvwxyz"))
	print("bytes", unsafe.Sizeof([]byte("a")))
	print("bytes", unsafe.Sizeof([]byte("abcdefghijklmnopqrstuvwxz")))
	print("byte", unsafe.Sizeof(byte('a')))
	print("rune", unsafe.Sizeof(rune('a')))
	
		fmt.Println("")
		
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
GoodOrder is 12 bytes
BadOrder is 16 bytes

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

string is 8 bytes
bytes is 12 bytes
bytes is 12 bytes
byte is 1 byte
rune is 4 bytes
```

## Verifying memory usage

We can verify the memory usage in golang using this simple test. This is how the `main_test.go` looks like. For each types we are interested to test the size, we allocate 1million of those types.
```go
// main_test.go
package main_test

import (
	"strings"
	"testing"
)

const Million = 1000000

func TestSliceInt32(t *testing.T) {
	arr := make([]int32, Million)
	for i := 0; i < Million; i++ {
		arr[i] = int32(i)
	}
}

func TestSliceInt64(t *testing.T) {
	arr := make([]int64, Million)
	for i := 0; i < Million; i++ {
		arr[i] = int64(i)
	}
}

func TestString(t *testing.T) {
	arr := make([]string, Million)
	for i := 0; i < Million; i++ {
		arr[i] = "a"
	}
	strings.Join(arr, "")
}

func TestByte(t *testing.T) {
	r := make([]byte, Million)
	for i := 0; i < Million; i++ {
		r[i] = 'a'
	}
}

func TestRune(t *testing.T) {
	r := make([]rune, Million)
	for i := 0; i < Million; i++ {
		r[i] = 'a'
	}
}

func TestBool(t *testing.T) {
	l := make([]bool, Million)
	r := make([]bool, Million)
	for i := 0; i < Million; i++ {
		l[i] = true
		r[i] = false
	}
}
```

To run the test with memory profiling:
```bash
$ go test -memprofile=mem.out
```

We can then run `go tool pprof mem.out`. Here's the output:

```bash
$ go tool pprof mem.out
Type: alloc_space
Time: Oct 16, 2018 at 11:13am (+08)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 36.17MB, 100% of 36.17MB total
      flat  flat%   sum%        cum   cum%
   15.27MB 42.20% 42.20%    17.52MB 48.43%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestString
    7.63MB 21.10% 63.31%     7.63MB 21.10%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestSliceInt64
    3.82MB 10.57% 73.87%     3.82MB 10.57%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestRune
    3.82MB 10.57% 84.44%     3.82MB 10.57%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestSliceInt32
    2.25MB  6.22% 90.66%     2.25MB  6.22%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestBool
    2.25MB  6.22% 96.89%     2.25MB  6.22%  strings.Join
    1.13MB  3.11%   100%     1.13MB  3.11%  _/Users/alextanhongpin/Documents/all/storage-test_test.TestByte
         0     0%   100%    36.17MB   100%  testing.tRunner
```

We can dig into individual test to see where the memory is allocated using `list <FunctionName>`:

```bash
(pprof) list TestSliceInt64
Total: 36.17MB
ROUTINE ======================== _/Users/alextanhongpin/Documents/all/storage-test_test.TestSliceInt64 in /Users/alextanhongpin/Documents/all/storage-test/main_test.go
    7.63MB     7.63MB (flat, cum) 21.10% of Total
         .          .     13:		arr[i] = int32(i)
         .          .     14:	}
         .          .     15:}
         .          .     16:
         .          .     17:func TestSliceInt64(t *testing.T) {
    7.63MB     7.63MB     18:	arr := make([]int64, Million)
         .          .     19:	for i := 0; i < Million; i++ {
         .          .     20:		arr[i] = int64(i)
         .          .     21:	}
         .          .     22:}
         .          .     23:
(pprof)
```

Here we see that 1 million int64 takes up 7.63MB. Remember that int64 takes up 8 bytes each, and to convert it to MB, we need to divide it by 1024 * 1024.

```
1e6 * 8 bytes / (1024 * 1024) = 7.63MB
```
