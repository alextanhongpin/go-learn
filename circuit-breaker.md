```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type State string

type States map[State]State

func (s States) Next(prev State) (State, bool) {
	next, exist := s[prev]
	return next, exist
}

const (
	// When everything is normal, the circuit breaker remains in the closed
	// state and all calls pass through to the services. When the number of
	// failures exceeds a predetermined threshold the breaker trips, and it
	// goes into the Open state
	Closed = State("closed")

	// The circuit breaker returns an error for calls without executing the
	// function.
	Opened = State("opened")

	// After a timeout period, the circuit switches to a half-open state to
	// test if the underlying problem still exists. If a single call fails
	// in this half-open state, the breaker is once again tripped.
	// If it succeeds, the circuit breaker resets back to the normal,
	// closed state.
	HalfOpened = State("half-opened")
)

var states = States{
	Opened:     HalfOpened,
	Closed:     Opened,
	HalfOpened: Closed,
}

type CircuitBreaker struct {
	sync.RWMutex
	threshold int
	counter   int

	state State

	resetDuration time.Duration
	breakTime     time.Time
}

func NewCircuitBreaker(threshold int, resetDuration time.Duration) *CircuitBreaker {
	if threshold <= 0 {
		threshold = 5
	}
	return &CircuitBreaker{
		threshold:     threshold,
		state:         Closed,
		resetDuration: resetDuration,
	}
}

func (c *CircuitBreaker) Exec(fn func() error) error {
	err := fn()
	if err != nil {
		c.Lock()
		c.counter++
		if c.counter >= c.threshold {
			c.state = Opened
			c.breakTime = time.Now()
		}
		c.Unlock()
		return err
	}
	
	c.RLock()
	state := c.state
	c.RUnlock()
	
	switch state {
	case Closed:
		return nil
	case HalfOpened, Opened:
		c.Lock()
		if time.Since(c.breakTime) > c.resetDuration {
			c.counter = 0
			c.state, _ = states[c.state]
		}
		c.Unlock()
		return nil
	default:
		return nil
	}
}

func max(hd int, rest ...int) int {
	for _, i := range rest {
		if i > hd {
			hd = i
		}
	}
	return hd
}

func main() {
	threshold := 3
	cb := NewCircuitBreaker(threshold, 1*time.Second)

	for i := 0; i < 10; i++ {
		err := cb.Exec(func() error {
			if i < threshold {
				return errors.New("bad error")
			}
			if i > threshold && i < threshold+3 {
				time.Sleep(1 * time.Second)
			}
			if i > threshold+3 {
				return errors.New("bad error")
			}
			return nil
		})
		fmt.Println(err, cb.state)
	}

	fmt.Printf("%#v", cb.state)
}
```
