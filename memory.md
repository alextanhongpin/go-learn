```go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GC()
	_ = make([]int64, 1024*1024)
	PrintMemoryUsage()
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/1024/1024)
}
```
