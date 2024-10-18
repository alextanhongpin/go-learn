```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type User struct {
	name string
}

func main() {
	Set("hello")
	Set(User{"john"})
	fmt.Println("Hello, 世界", Get[string](), Get[User]())
}

var registry = map[any]any{}

func Get[T any]() T {
	var v T
	return registry[reflect.TypeOf(v)].(T)
}

func Set(v any) {
	registry[reflect.TypeOf(v)] = v
}
```
