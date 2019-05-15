
## Sequence

NOTE: Also take a look at `bitwise.md`.

Application:
- saga
- sequential workflow
- state machine

Goal is to ensure the sequence of execution is respected, that is to say:
- the next state can only be executed if the previous state is done
- the sequence should not be skippable, they must be executed in order
- once the sequence is executed, it should not be executable again (once-only, idempotent)

```go
package main

import (
	"fmt"
	"sort"
)

type Sequences []uint

func (s Sequences) Len() int           { return len(s) }
func (s Sequences) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sequences) Less(i, j int) bool { return s[i] < s[j] }

type Sequence struct {
	state     byte
	sequences Sequences // The list of sequences that are allowed.
}

func NewSequence(sequences Sequences) *Sequence {
	// Ensure the sequences are sorted.
	sort.Sort(sequences)
	return &Sequence{sequences: sequences}
}

func (s *Sequence) Completed() bool {
	var state byte
	for _, seq := range s.sequences {
		state |= (1 << seq)
	}
	return state == s.state
}

func (s *Sequence) Set(nextState uint) bool {
	curr := s.state
	// The expected state up till now.
	var exp byte
	for _, seq := range s.sequences {
		notSet := (curr & (1 << seq)) == 0
		isCurr := seq == nextState
		allSet := exp == curr
		if isCurr && notSet && allSet {
			// Set the value before returning.
			s.state |= (1 << seq)
			return true
		}
		exp |= (1 << seq)
	}
	return false
}

func main() {
	seq := NewSequence(Sequences{1, 2, 3, 4, 5, 6, 7, 8})
	fmt.Println(seq.Set(1))
	fmt.Println(seq.Set(1))
	fmt.Println(seq.Set(2))
	fmt.Println(seq.Set(1))
	fmt.Println(seq.Set(8))
	fmt.Println(seq.Set(4))
	fmt.Println(seq.Set(3))
	fmt.Println(seq.Set(4))
	fmt.Println(seq.Set(5))
	fmt.Println(seq.Completed())
	for i := 0; i < 10; i++ {
		fmt.Println(seq.Set(uint(i)))
	}
	fmt.Println(seq.Completed())
}
```

## State Machine 
```go
package main

import (
	"fmt"
	"sort"
)

type Sequence []uint

func (s Sequence) Len() int           { return len(s) }
func (s Sequence) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sequence) Less(i, j int) bool { return s[i] < s[j] }

type StateMachine struct {
	state    byte
	sequence Sequence
}

func NewStateMachine(seq Sequence) *StateMachine {
	sort.Sort(seq)
	return &StateMachine{
		sequence: seq,
	}
}

func (s *StateMachine) CanSet(tgt uint) bool {
	// If the state already exists, it should be false.
	if s.state&(1<<tgt) != 0 {
		return false
	}

	var curr byte
	for _, seq := range s.sequence {
		noSkip := curr == s.state
		isCurr := seq == tgt
		if noSkip && isCurr {
			return true
		}
		curr |= (1 << seq)
	}
	return false
}

func (s *StateMachine) CanUnset(tgt uint) bool {
	// If the state is not yet set, it should be false.
	if s.state&(1<<tgt) == 0 {
		return false
	}
	var curr byte
	for _, seq := range s.sequence {
		curr |= (1 << seq)
		noSkip := curr == s.state
		isCurr := seq == tgt
		if noSkip && isCurr {
			return true
		}
	}
	return false
}

func (s *StateMachine) Set(tgt uint) bool {
	if s.CanSet(tgt) {
		s.state |= (1 << tgt)
		return (s.state & (1 << tgt)) != 0
	}
	return false
}

func (s *StateMachine) Unset(tgt uint) bool {
	if s.CanUnset(tgt) {
		s.state &^= (1 << tgt)
		return (s.state & (1 << tgt)) == 0
	}
	return false
}

func (s *StateMachine) IsZero() bool {
	return s.state == 0
}
func (s *StateMachine) IsCompleted() bool {
	var curr byte
	for _, seq := range s.sequence {
		curr |= (1 << seq)
	}
	return curr == s.state
}
func main() {
	sm := NewStateMachine(Sequence{5, 4, 2, 1, 3})
	fmt.Println(sm.Set(1))
	fmt.Println(sm.Set(2))
	fmt.Println(sm.Set(3))
	fmt.Println(sm.Unset(2))
	fmt.Println(sm.Unset(3))
	fmt.Println(sm.Unset(4))
	var i uint = 5
	for !sm.IsZero() {
		sm.Unset(i)
		i--
	}
	fmt.Println(sm.state)
	for !sm.IsCompleted() {
		sm.Set(i)
		i++
	}
	fmt.Println(sm.state)
}
```
