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

## With Termination, Single State

**NOTE**: This might be a little over-engineered, see the solution below:
```go
package main

import (
	"fmt"
)

type State string

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	if !exist {
		return Invalid, false
	}
	return next, next == Ended
}

const (
	Started = State("^")
	Ended   = State("$") // We need this to differentiate between non-existing state and completed state.
	Invalid = State("invalid")

	Submitted = State("submitted")
	Approved  = State("approved")
	Published = State("published")
)

func main() {
	var states = StateMachine{
		Started:   Submitted,
		Submitted: Approved,
		Approved:  Published,
		Published: Ended,
		Ended:     Ended,
	}
	var completed bool
	var initialState = Started
	state := initialState

	for i := 0; i < 5; i++ {
		fmt.Printf("prev: %s", state)
		state, completed = states.Next(state)
		fmt.Printf(" next: %s, completed: %t\n", state, completed)
		// Break to avoid infinite loop.
		if completed {
			break
		}
	}
}
```

Simplified solution:

```go
package main

import (
	"fmt"
	"strings"
)

type State string

func (s State) String() string {
	return string(s)
}

func (s State) Eq(in string) bool {
	return strings.EqualFold(string(s), in)
}

func (s State) EqStrict(in string) bool {
	return string(s) == in
}

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	return next, !exist
}

const (
	Invalid   = State("invalid")
	Started   = State("started")
	Submitted = State("submitted")
	Approved  = State("approved")
	Published = State("published")
)

func main() {
	var states = StateMachine{
		Started:   Submitted,
		Submitted: Approved,
		Approved:  Published,
	}
	var completed bool
	var initialState = Started
	state := initialState

	for i := 0; i < 5; i++ {
		fmt.Printf("prev: %s", state)
		state, completed = states.Next(state)
		fmt.Printf(" next: %s, completed: %t\n", state, completed)
		// Break to avoid infinite loop.
		if completed {
			break
		}
	}
}
```

## With Termination, Multiple States

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
	Initialized = State("initialized") // Start
	Submitted   = State("submitted")
	Rejected    = State("rejected")
	Approved    = State("approved") // End
)

func main() {
	// Assuming the states are describing payment.
	var states = States{
		Initialized: []State{Submitted},
		Submitted:   []State{Approved, Rejected},
		Rejected:    []State{Submitted},
		Approved:    []State{}, // After approved, there are no other continuation.
	}
	var initialState = Initialized
	state := initialState
	for i := 0; i < 3; i++ {
		choices, completed := states.Next(state)
		if completed {
			fmt.Println("completed")
			break
		}
		state = choices[0]
		fmt.Printf("next state is: %s, completed: %t\n", state, completed)
	}
}
```

## With switch, single state

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


## With switch, multiple states
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
