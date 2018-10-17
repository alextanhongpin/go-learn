## Integer limit values

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(1<<7 - 1)
	fmt.Println(-1 << 7)
	fmt.Println(1<<15 - 1)
	fmt.Println(-1 << 15)
	fmt.Println(1<<31 - 1)
	fmt.Println(-1 << 31)
	fmt.Println(1<<63 - 1)
	fmt.Println(-1 << 63)

	fmt.Println(math.MaxInt8)
	fmt.Println(math.MinInt8)
	fmt.Println(math.MaxInt16)
	fmt.Println(math.MinInt16)
	fmt.Println(math.MaxInt32)
	fmt.Println(math.MinInt32)
	fmt.Println(math.MaxInt64)
	fmt.Println(math.MinInt64)

	fmt.Println(1<<8 - 1)
	fmt.Println(1<<16 - 1)
	fmt.Println(1<<32 - 1)
	// fmt.Println(1<<64 - 1)

	fmt.Println(math.MaxUint8)
	fmt.Println(math.MaxUint16)
	fmt.Println(math.MaxUint32)
	// fmt.Println(math.MaxUint64)
}
```

Output:

```
127
-128
32767
-32768
2147483647
-2147483648
9223372036854775807
-9223372036854775808
127
-128
32767
-32768
2147483647
-2147483648
9223372036854775807
-9223372036854775808
255
65535
4294967295
255
65535
4294967295
```
