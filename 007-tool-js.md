```go
package main

import (
	"fmt"

	"github.com/dop251/goja"
)

func main() {
	vm := goja.New()
	vm.Set("Foo", map[string]any{
		"bar": func(args string) (string, error) {
			fmt.Println(args)
			return "foo", nil
		},
		"json": func(data any) (any, error) {
			fmt.Println("FN:json", data)
			return data, nil
		},
	})
	v, err := vm.RunString(`
const result = Foo.bar(1)
Foo.json({result})
`)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Export().(map[string]any))
}
```
