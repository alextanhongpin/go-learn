## Step-Driven

Maybe we can snapshot the steps, so that we have better visibility on the logic?

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"github.com/davecgh/go-spew/spew"

	"context"
	"fmt"
)

type contextKey string

var stepHistoryContextKey = contextKey("step_history")

type StepHistory struct {
	steps []Step
}

func (s *StepHistory) Dump() {
	t := len(s.steps)
	for i, s := range s.steps {
		n := i + 1
		fmt.Printf("step %d: %s\n", n, s.Name)
		spew.Dump(s.Args...)

		if n != t {
			fmt.Println("---")
		}
	}
}

type Step struct {
	Name string
	Args []any
}

func main() {
	hist := &StepHistory{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, stepHistoryContextKey, hist)
	Foo(ctx, FooDto{"foo"})
	hist.Dump()
}

func step(ctx context.Context, name string, args ...any) {
	hist := ctx.Value(stepHistoryContextKey).(*StepHistory)
	hist.steps = append(hist.steps, Step{
		Name: name,
		Args: args,
	})
}

type FooDto struct {
	msg string
}

func Foo(ctx context.Context, dto FooDto) {
	step(ctx, "exec foo", dto)
	Bar(ctx, dto.msg)
}

func Bar(ctx context.Context, msg string) {
	step(ctx, "exec bar", msg)
}
```
