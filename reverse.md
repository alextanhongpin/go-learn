# Reversing a slice

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(Reverse([]int{1, 2, 3, 4, 5, 6, 7}))
}

func Reverse(nums []int) []int {
	// Make a copy of the slice.
	tmp := make([]int, len(nums))
	copy(tmp, nums)

	n := len(tmp) - 1
	for i := 0; i < len(tmp)/2; i++ {
		tmp[i], tmp[n-i] = tmp[n-i], tmp[i]
	}
	return tmp
}
```

Test:

```go
package main

import (
	"log"
	"testing"
	"testing/quick"
)

func TestReverse(t *testing.T) {
	f := func(nums []int) bool {
		n := len(nums) - 1
		val := Reverse(nums)
		out := Reverse(val)
		for i, v := range nums {
			if v != val[n-i] {
				return false
			}
			if v != out[i] {
				return false
			}
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		log.Fatal(err)
	}
}
```

## Reversing Integer

Can be used to check if an integer is a palindrone or not. This way of reversing is done without converting the integer to string first. Note that reversing `100` will result in `1`, since the zeros will be removed.

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(reverseInt(987654321))
}

func reverseInt(n int) int {
	var reversed int
	for n > 0 {
		reversed = (reversed * 10) + (n % 10)
		n /= 10
	}
	return reversed
}
```
