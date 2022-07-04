```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type User struct {
	Name      string
	Age       int
	Hobbies   []string
	Married   *bool
	CreatedAt time.Time
}

func main() {
	alice := &User{
		Name:      "alice",
		Age:       10,
		Hobbies:   []string{"programming", "reading"},
		Married:   nil,
		CreatedAt: time.Now().Add(1 * time.Second),
	}

	married := true
	bob := &User{
		Name:      "bob",
		Age:       23,
		Hobbies:   []string{"sleeping", "reading", "programming"}, // Order doesn't matter, only the different data will be printed.
		Married:   &married,
		CreatedAt: time.Now(),
	}
	// Other ignore option available.
	// https://pkg.go.dev/github.com/google/go-cmp/cmp/cmpopts#example-IgnoreFields-Testing
	opts := []cmp.Option{
		cmpopts.IgnoreFields(User{}, "CreatedAt"), // Ignore fields that doesn't require testing.
		cmpopts.SortSlices(func(l, r any) bool {
			// Sort slices so that they order doesn't matter. Only sort for string slices.
			lhs, lok := l.(string)
			rhs, rok := r.(string)

			return lok && rok && lhs < rhs
		}),
	}
	
	exp, got := alice, bob
	if diff := cmp.Diff(exp, got, opts...); diff != "" {
		fmt.Printf("-exp, +got:\n%s", diff)
	}
}
```
