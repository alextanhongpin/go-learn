
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
