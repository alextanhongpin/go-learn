```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)

func main() {
	rl := &gcra{
		burst:  1,
		limit:  5,
		period: time.Second,
	}
	for i := range 10 {
		fmt.Println(i, rl.Allow())
		time.Sleep(100 * time.Millisecond)
		if i == 8 {
			time.Sleep(200 * time.Millisecond)
		}
	}
	time.Sleep(1 * time.Second)
	for i := range 10 {
		fmt.Println(i, rl.Allow())
	}
	fmt.Println("Hello, 世界")
}

type gcra struct {
	// State.
	ts    time.Time
	burst int

	// Option.
	period time.Duration
	limit  int
}

func (g *gcra) Allow() bool {
	interval := g.period / time.Duration(g.limit)
	burst := time.Duration(g.burst) * interval
	lower := time.Now().Add(-burst)
	if gte(lower, g.ts) {
		// reset.
		g.ts = time.Now().Add(-burst)
	}

	upper := time.Now().Truncate(g.period).Add(g.period)
	if gte(g.ts, upper) {
		return false
	}

	g.ts = g.ts.Add(interval)

	return true
}

func gte(a, b time.Time) bool {
	return !a.Before(b)
}
```
