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
	next, exist := s[prev]
	return next, exist
}

func (s StateMachine) IsCompleted(prev State) bool {
	next, exist := s[prev]
	return exist && len(next) == 0
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
		Paid:     []State{}, // There are no more states after paid.
	}

	var initialState = NotPaid
	next, _ := states.Next(initialState)
	completed := states.IsCompleted(initialState)
	fmt.Printf("%#v, completed: %t\n", next, completed)

	next, _ = states.Next(next[0])
	completed = states.IsCompleted(next[0])
	fmt.Printf("%#v, completed: %t\n", next, completed)
	{
		// After paid, there are no other states. It should be completed.
		completed := states.IsCompleted(next[0])
		fmt.Printf("completed: %t\n", completed)
	}
	{
		// If the payment is rejected, we go back to payment_pending.
		next, _ := states.Next(next[1])
		completed := states.IsCompleted(initialState)
		fmt.Printf("%#v, completed: %t\n", next, completed)
	}
}
```
