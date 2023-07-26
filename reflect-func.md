# How to reflect the function type and cast any to the function arg

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	v := reflect.ValueOf(hello)
	fmt.Println(v.Type().NumIn())
	v.Call([]reflect.Value{
		reflect.ValueOf("ssup"),
	})
	arg0 := v.Type().In(0)
	s := "world"
	sv := reflect.TypeOf(s)
	fmt.Println(sv.AssignableTo(arg0))
	fmt.Println(arg0)
	fmt.Println("Hello, 世界")

	b, err := json.Marshal("hello")
	if err != nil {
		panic(err)
	}
	nt := reflect.New(arg0).Elem().Interface()
	fmt.Println(json.Unmarshal(b, &nt))
	fmt.Println(reflect.TypeOf(nt).AssignableTo(arg0))
	v.Call([]reflect.Value{reflect.ValueOf(nt)})
}

func hello(s string) {
	fmt.Println("hi", s)
}
```
