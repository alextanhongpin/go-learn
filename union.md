# Union Type

The only way to emulate union type is by using interface at the moment.


```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	events := []PersonEvent{
		PersonCreated{Name: "John", Age: 10},
		PersonNameUpdated{Name: "Jane"},
	}
	for _, evt := range events {
		switch e := evt.(type) {
		case PersonCreated:
			fmt.Println("PersonCreated", e)
		case PersonNameUpdated:
			fmt.Println("PersonNameUpdated", e)
		}
	}
}

// union PersonEvent = PersonCreated | PersonNameUpdated
type PersonEvent interface {
	IsPersonEvent()
}

type PersonCreated struct {
	Name string
	Age  int
}

func (PersonCreated) IsPersonEvent() {}

type PersonNameUpdated struct {
	Name string
}

func (PersonNameUpdated) IsPersonEvent() {}
```
