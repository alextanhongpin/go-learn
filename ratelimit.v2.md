## Basic
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

	upper := time.Now()
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


## With result

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
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
	}
	time.Sleep(1 * time.Second)
	fmt.Println()
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

type result struct {
	allow     bool
	limit     int
	remaining int
	resetAt   time.Time
	retryAt   time.Time
}

func (r *result) String() string {
	return fmt.Sprintf("allow: %t (%d/%d), next: %s, reset: %s",
		r.allow,
		r.remaining,
		r.limit,
		r.retryAt,
		r.resetAt,
	)
}

func (g *gcra) Allow() *result {
	now := time.Now()
	interval := g.period / time.Duration(g.limit)
	burst := time.Duration(g.burst) * interval
	lower := now.Add(-burst)
	if gte(lower, g.ts) {
		// reset.
		g.ts = now.Add(-burst)
	}

	resetAt := now.Truncate(g.period).Add(g.period)
	remaining := int(math.Floor(float64(resetAt.Sub(now)) / float64(interval)))
	if gte(g.ts, now) {
		return &result{
			allow:     false,
			limit:     g.limit + g.burst,
			remaining: remaining,
			resetAt:   resetAt,
			retryAt:   g.ts,
		}
	}

	g.ts = g.ts.Add(interval)

	return &result{
		allow:     true,
		limit:     g.limit + g.burst,
		remaining: remaining,
		resetAt:   resetAt,
		retryAt:   g.ts,
	}
}

func gte(a, b time.Time) bool {
	return !a.Before(b)
}
```
