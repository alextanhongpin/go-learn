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

type States map[State][]State

func (s States) Next(prev State) ([]State, bool) {
	next, exist := s[prev]
	if !exist {
		return nil, false
	}
	return next, len(next) == 0
}

const (
	Initialized = State("payment_initialized") // Start
	Submitted   = State("payment_submitted")
	Rejected    = State("payment_rejected")
	Approved    = State("payment_approved") // End
)

func main() {
	var states = States{
		Initialized: []State{Submitted},
		Submitted:   []State{Approved, Rejected},
		Rejected:    []State{Submitted},
		Approved:    []State{}, // After approved, there are no other continuation.
	}
	var initialState = Initialized
	state := initialState
	for i := 0; i < 3; i++ {
		next, completed := states.Next(state)
		if completed {
			fmt.Println("completed")
			break
		}
		state = next[0]
		fmt.Printf("next state is: %s, completed: %t\n", state, completed)
	}
}
```

## Using switch 

```go
package main

import (
	"fmt"
)

func main() {
	var initialState = Red
	state := initialState
	for i := 0; i < 3; i++ {
		state = Next(state)
		fmt.Println("next state is:", state)
	}
}

type State string

const (
	Invalid = State("")
	Red     = State("red")
	Green   = State("green")
	Yellow  = State("yellow")
)

func Next(prev State) State {
	switch prev {
	case Red:
		return Green
	case Green:
		return Yellow
	case Yellow:
		return Red
	default:
		return Invalid
	}
}
```


## With multiple states
```go
package main

import (
	"fmt"
)

func main() {
	var initialState = Initialized
	state := initialState
	for i := 0; i < 3; i++ {
		states, completed := Next(state)
		if completed {
			fmt.Println("completed")
			break
		}
		state = states[0]
		fmt.Printf("next state is: %s, completed: %t\n", state, completed)
	}
}

type State string

const (
	Initialized = State("payment_initialized") // Start
	Submitted   = State("payment_submitted")
	Rejected    = State("payment_rejected")
	Approved    = State("payment_approved") // End
)

func Next(prev State) ([]State, bool) {
	switch prev {
	case Initialized:
		return []State{Submitted}, false
	case Submitted:
		return []State{Approved, Rejected}, false
	case Rejected:
		return []State{Submitted}, false
	case Approved:
		return []State{}, true
	default:
		return []State{}, false
	}
}
```
