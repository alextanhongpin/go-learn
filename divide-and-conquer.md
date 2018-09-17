```go
package main

import (
	"fmt"
	"math"
)

func main() {
	var ar []int
	for i := 0; i < 10000; i++ {
		ar = append(ar, i)
	}

	fmt.Println("Expected big 0:", math.Log2(10000.0))
	steps := divideAndConquer(69, 0, ar)

	fmt.Printf("took %d steps to complete", steps)
}

func divideAndConquer(target, steps int, arr []int) int {
	if len(arr) == 1 {
		return steps
	}
	midpoint := (len(arr) / 2) - 1
	midvalue := arr[midpoint]
	if target > midvalue {
		return divideAndConquer(target, steps+1, arr[midpoint:len(arr)-1])
	} else if target < midvalue {
		return divideAndConquer(target, steps+1, arr[0:midpoint])
	}
	return steps
}
```
