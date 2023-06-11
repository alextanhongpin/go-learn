# Define steps as interface

Allows steps to be mocked.

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	// Actually why do this? Just pass a repo and implement all the steps.
	uc := NewUserUsecase(nil)
	uc.Do()

	uc = NewUserUsecase(&mockSteps{})
	uc.Do()
}

func NewUserUsecase(steps steps) *UserUsecase {
	uc := &UserUsecase{steps: steps}
	if steps == nil {
		uc.steps = uc
	}
	return uc
}

type mockSteps struct{}

func (m *mockSteps) Process() {
	fmt.Println("mocked")
}

type steps interface {
	Process()
}
type UserUsecase struct {
	steps
}

func (uc *UserUsecase) Do() {
	fmt.Println("do")
	uc.steps.Process()
}

func (uc *UserUsecase) Process() {
	fmt.Println("process")
}

```
