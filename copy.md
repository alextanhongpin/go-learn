
## Copying struct pointer in slice
```go
package main

import (
	"fmt"
)

type Template struct {
	Name string
}

func copy(templates []*Template) []*Template {
	result := make([]*Template, len(templates))
	for i, tpl := range templates {
		result[i] = &Template{}
		*result[i] = *tpl
	}
	return result
}

func main() {
	var templates []*Template
	templates = append(templates, &Template{"john"})
	cptemplates := copy(templates)

	templates[0].Name = "doe"

	fmt.Println(templates[0], cptemplates[0])
}
```

## Avoiding copy of struct

See [here](https://github.com/golang/go/issues/8005#issuecomment-190753527).

```go
type hello struct {
	noCopy noCopy
}
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

func greet(h hello) {
	fmt.Println(h)
}

func greet2(h *hello) {
	fmt.Println(h)
}

func main() {
	h := hello{}
	greet(h) // Error
	greet2(&h) // Valid
}
```
