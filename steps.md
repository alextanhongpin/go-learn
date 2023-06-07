## Step-Driven

Maybe we can snapshot the steps, so that we have better visibility on the logic?

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
)

type contextKey string

var stepHistoryContextKey = contextKey("step_history")

type StepHistory struct {
	steps []string
}

func main() {
	hist := &StepHistory{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, stepHistoryContextKey, hist)
	Foo(ctx)
	fmt.Println(hist.steps)
}

func step(ctx context.Context, name string, args ...any) {
	hist := ctx.Value(stepHistoryContextKey).(*StepHistory)
	hist.steps = append(hist.steps, name)
}

func Foo(ctx context.Context) {
	step(ctx, "foo")
	Bar(ctx)
}

func Bar(ctx context.Context) {
	step(ctx, "bar")
}
```
