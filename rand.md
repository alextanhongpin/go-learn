# Generating unique random number with normal distribution (?)

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(generateUniqueRandomNumber(0, 20, 10))
}

func generateUniqueRandomNumber(lo, hi int, k int) []int {
	if lo > hi {
		panic("lo cannot be higher than hi")
	}
	delta := hi - lo

	if delta < k {
		panic("unable to generate unique number with limited range")
	}

	// There will be a lot of probability of collision, since all numbers are used.
	// Use naive fisher-yates (?) shuffle.
	res := make([]int, delta)
	for i := 0; i < delta; i++ {
		res[i] = lo + i
	}
	for i := delta - 1; i > -1; i-- {
		n := rand.Intn(i + 1)
		res[i], res[n] = res[n], res[i]
	}
	return res[:k]
}
```
