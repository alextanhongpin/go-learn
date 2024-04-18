```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(0)
	cap := time.Minute
	base := 100 * time.Millisecond
	for i := 0; i < 10; i++ {
		attempt := i
		temp := min(int(cap), int(base)*pow(2, attempt))
		sleep := temp/2 + rand.Intn(temp/2)
		fmt.Println(time.Duration(sleep).Round(1 * time.Millisecond))
	}
	/*
		80ms
		126ms
		209ms
		517ms
		1.589s
		2.222s
		4.887s
		10.677s
		22.179s
		36.143s
	*/

}
func pow(a, b int) int {
	t := 1
	for _ = range b {
		t *= a
	}
	return t
}
```
