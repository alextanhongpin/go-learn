```go
package main

import (
	"context"
	"log"
	"sort"
)

type Cancel func(ctx context.Context)

func main() {

	split()
	complex()
	bitwiseTrick()
}

func split() {
	// In this case, original refers to the state at the beginning of the transaction, not the initial state (unless it is the only transaction, but this doesn't make it a saga anymore).
	// If we have multiple steps, the original state should be the previous successfull state.
	var original int
	i := original
	tx := func() {
		i += 10
	}
	comp := func() {
		// Can this be reversed? Only if the state has changed...
		// If the compensation is not idempotent, invoking it might cause irreversible side-effects.
		if i != original {
			i -= 10
		}
	}

	tx()
	log.Println("i is now", i, original)
	comp()
	comp()
	log.Println("i is now", i, original)
}

func complex() {
	txs := map[string]func(initial int) (next int){
		"add": func(i int) int {
			return i + 100
		},
		"multiply": func(i int) int {
			return i * 2
		},
		"divide": func(i int) int {
			return i / 4
		},
	}
	// Note: The compensation steps must be in the reverse order.
	// This is like a stack, first-in, first-out.
	cps := map[string]func(prev int) (next int){
		"divide": func(i int) int {
			return i * 4
		},
		"multiply": func(i int) int {
			return i / 2
		},
		"add": func(i int) int {
			return i - 100
		},
	}
	_ = cps
	var initialState int
	state := initialState

	// NOTE: Ranging in maps will not guarante the sequence. Use Slice instead.
	steps := []string{"add", "multiply", "divide"}
	var stepsWithErr string
	var completed []string
	rollback := func() {
		for _, act := range completed {
			state = cps[act](state)
			log.Println("rollback", act, state)
		}
		log.Println("state is now", state, initialState)
	}
	stepsWithErr = "divide"
	for _, act := range steps {
		if act == stepsWithErr {
			// Find all the previous steps, and rollback.
			rollback()
			return
		}
		state = txs[act](state)
		// Note: The compensation steps must be in the reverse order.
		completed = append([]string{act}, completed...)
		log.Println("executing", act, state)
	}
}

type Step uint
type Steps []Step

func (s Steps) Len() int           { return len(s) }
func (s Steps) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Steps) Less(i, j int) bool { return s[i] < s[j] }

func (s *Step) Set(step Step) {
	*s |= step
}

func (s *Step) Unset(step Step) {
	*s &= ^step
}

func (s Step) Has(step Step) bool {
	return s&step > 0
}

const (
	Add Step = 1 << iota
	Multiply
	Divide
)

func bitwiseTrick() {
	var i int
	transactions := map[Step]func(){
		Add: func() {
			i += 100
		},
		Multiply: func() {
			i *= 2
		},
		Divide: func() {
			i /= 4
		},
	}
	compensations := map[Step]func(){
		Add: func() {
			i -= 100
		},
		Multiply: func() {
			i /= 2
		},
		Divide: func() {
			i *= 4
		},
	}
	steps := Steps{Add, Multiply, Divide}
	var step Step

	// We intentionally skip the last state here.
	for _, s := range steps[:len(steps)-1] {
		if step.Has(s) {
			log.Println("skipping", step)
			continue
		}
		transactions[s]()
		// Apply the steps after it has been executed.
		step.Set(s)
		log.Println("applied", step, i)
	}
	// Idempotent...
	for _, s := range steps {
		if step.Has(s) {
			log.Println("skipping", s)
			continue
		}
		transactions[s]()
		// Apply the steps after it has been executed.
		step.Set(s)
		log.Println("applied", step, i)
	}
	// Reverse the steps.
	sort.Sort(sort.Reverse(steps))
	for _, s := range steps {
		compensations[s]()
		// Unset the steps once it have been executed.
		step.Unset(s)
		log.Println("compensated", step, i)
	}
}
```
