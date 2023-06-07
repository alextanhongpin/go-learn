## Functional Optional pattern when you need to share

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	HTTP(IgnoreHeaders(), IgnoreFields())
	fmt.Println()
	JSON(IgnoreFields())
}

func HTTP(opts ...HTTPOption) {
	fmt.Println("HTTP")
	for _, opt := range opts {
		switch o := opt.(type) {
		case *IgnoreHeadersOption:
			fmt.Println("ignore headers", o)
		case *IgnoreFieldsOption:
			fmt.Println("ignore fields", o)
		}
	}
}

func JSON(opts ...JSONOption) {
	fmt.Println("JSON")
	for _, opt := range opts {
		switch o := opt.(type) {
		case *IgnoreFieldsOption:
			fmt.Println("ignore fields", o)
		}
	}
}

type HTTPOption interface {
	isHTTP()
}
type JSONOption interface {
	isJSON()
}

type IgnoreHeadersOption struct {
}

func IgnoreHeaders() *IgnoreHeadersOption {
	return &IgnoreHeadersOption{}
}
func (o *IgnoreHeadersOption) isHTTP() {}

func IgnoreFields() *IgnoreFieldsOption {
	return &IgnoreFieldsOption{}
}

type IgnoreFieldsOption struct{}

func (o *IgnoreFieldsOption) isHTTP() {}
func (o *IgnoreFieldsOption) isJSON() {}
```
