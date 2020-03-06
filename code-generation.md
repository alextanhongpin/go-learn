```go
package main

import (
	"flag"
	"html/template"
	"os"
)

//go:generate go run main.go -name=Int -type=int
func main() {
	type data struct {
		Name string
		Type string
	}

	var d data
	flag.StringVar(&d.Name, "name", "", "the name of the struct")
	flag.StringVar(&d.Type, "type", "", "the type of the struct")
	flag.Parse()

	t := template.Must(template.New("").Parse(tpl))
	t.Execute(os.Stdout, d)
}

var tpl = `func New{{.Name}}(in {{.Type}}) {{.Type}} {
	return in
}
`
```
