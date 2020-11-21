# Get struct/interface name

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	p := new(Person)
	fmt.Println(p.Name())
}

type Person struct{}

func (p *Person) Name() string {
	return className(p)
}

func className(unknown interface{}) string {
	t := reflect.TypeOf(unknown)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
```
