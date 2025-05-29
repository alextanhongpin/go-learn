# Setting `any` value
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var i int
	set(&i, 100)
	fmt.Println(i)

	j := reflect.New(reflect.TypeOf(i)).Elem().Interface()
	set(&j, 200)
	fmt.Println(j)

	handle(i)
}

func set(value any, v any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("value must be a non-nil pointer, got %T", value)
	}

	rv.Elem().Set(reflect.ValueOf(v))
	return nil
}

func handle(int) {}
```
