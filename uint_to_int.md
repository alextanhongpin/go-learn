# "Shifting" uint to int
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("int64", math.MinInt64, math.MaxInt64)
	fmt.Println("int32", math.MinInt32, math.MaxInt32)
	fmt.Println("uint32", 0, uint32(math.MaxUint32))
	fmt.Println("uint64", 0, uint64(math.MaxUint64))

	fmt.Println(uint64ToInt64(uint64(math.MaxUint64)))
	fmt.Println(uint64ToInt64(uint64(math.MaxUint64)) - 1)
	fmt.Println(uint64ToInt64(0))
	fmt.Println(uint32ToInt32(uint32(math.MaxUint32)))
	fmt.Println(uint32ToInt32(uint32(math.MaxUint32) / 2))
	fmt.Println(uint32ToInt32(uint32(math.MaxUint32) - uint32(-math.MinInt32)))
	fmt.Println(uint32ToInt32(uint32(math.MaxUint32) - uint32(math.MaxInt32)))
	fmt.Println(uint32ToInt32(0))
}

func uint64ToInt64(u64 uint64) int64 {
	if u64 > uint64(math.MaxInt64) {
		return int64(u64 - uint64(math.MaxInt64) - 1)
	}
	return int64(u64) - math.MaxInt64 - 1
}

func uint32ToInt32(u32 uint32) int32 {
	if u32 > uint32(math.MaxInt32) {
		return int32(u32 - uint32(math.MaxInt32) - 1)
	}
	return int32(u32) - math.MaxInt32 - 1
}
```
