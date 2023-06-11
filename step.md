# Define steps as interface

Allows steps to be mocked.

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	uc := &UserUsecase{}
	uc.steps = uc
	uc.Do()

	uc.steps = &mockSteps{}
	uc.Do()
	fmt.Println("Hello, 世界")
}

type mockSteps struct{}

func (m *mockSteps) Process() {
	fmt.Println("mocked")
}

type UserUsecase struct {
	steps interface {
		Process()
	}
}

func (uc *UserUsecase) Do() {
	fmt.Println("do")
	uc.steps.Process()
}

func (uc *UserUsecase) Process() {
	fmt.Println("process")
}
```
