## Sequence

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
