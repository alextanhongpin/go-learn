## State Machine with no termination

```go
package main

import (
	"fmt"
)

type State string

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	return next, exist
}

const (
	GreenLight  = State("green_light")
	OrangeLight = State("orange_light")
	RedLight    = State("red_light")
)

func main() {
	var states = StateMachine{
		GreenLight:  OrangeLight,
		OrangeLight: RedLight,
		RedLight:    OrangeLight,
	}
	var initialState = GreenLight
	state := initialState
	for i := 0; i < 3; i++ {
		fmt.Printf("prev: %s", state)
		state, _ = states.Next(state)
		fmt.Printf(" next: %s\n", state)
	}
}
```


## With Termination

```go
package main

import (
	"fmt"
)

type State string

type StateMachine map[State][]State

func (s StateMachine) Next(prev State) ([]State, bool) {
	next, completed := s[prev]
	return next, !completed
}

const (
	NotPaid  = State("")
	Rejected = State("payment_rejected")
	Pending  = State("payment_pending")
	Paid     = State("paid")
)

func main() {
	var states = StateMachine{
		NotPaid:  []State{Pending},
		Pending:  []State{Paid, Rejected},
		Rejected: []State{Pending},
		// Paid:            []State{Paid},
	}

	var initialState = NotPaid
	next, _ := states.Next(initialState)
	fmt.Printf("%#v\n", next) // []main.State{"payment_pending"}

	next, _ = states.Next(next[0])
	fmt.Printf("%#v\n", next) // []main.State{"paid", "payment_rejected"}
	{
		// After paid, there are no other states. It should be completed.
		_, completed := states.Next(next[0])
		fmt.Printf("completed: %t\n", completed) // completed: true
	}
	{
		// If the payment is rejected, we go back to payment_pending.
		next, _ := states.Next(next[1])
		fmt.Printf("%#v\n", next) // []main.State{"payment_pending"}
	}
}
```
