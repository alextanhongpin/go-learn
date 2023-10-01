# Calculating error rate

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	var total int
	var failed int

	s := rate.Sometimes{
		Interval: 100 * time.Millisecond,
	}
	now := time.Now()
	for i := 0; i < 200; i++ {
		total += 1
		if rand.Intn(10) < 4 {
			failed++
		}
		sleep := time.Duration(rand.Intn(25)) * time.Millisecond
		time.Sleep(sleep)
		s.Do(func() {
			fmt.Println("rate", float64(failed)/float64(total), failed, total, time.Since(now))
			failed = 0
			total = 0
		})
	}
}
```
