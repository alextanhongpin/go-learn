## With Jennifer
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


## With Jennifer

```go
package main

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

type Gen struct {
	Constructor             string
	PrivateName, PublicName string
	ShortName               string
}

func main() {
	f := NewFile("main")

	entity := Gen{
		Constructor: "NewPerson",
		PrivateName: "person",
		PublicName:  "Person",
		ShortName:   "p",
	}
	builder := Gen{
		Constructor: entity.Constructor + "Builder",
		PrivateName: entity.PrivateName + "Builder",
		PublicName:  entity.PublicName + "Builder",
		ShortName:   entity.ShortName + "b",
	}
	builderOption := Gen{
		Constructor: builder.Constructor + "Option",
		PrivateName: builder.PrivateName + "Option",
		PublicName:  builder.PublicName + "Option",
		ShortName:   builder.ShortName + "o",
	}
	f.Type().Id(builder.PublicName).Struct(
		Id(entity.PrivateName).Id(entity.PublicName),
	)
	f.Type().Id(builderOption.PublicName).Func().Params(Id(entity.ShortName).Op("*").Id(builder.PublicName)).Error()
	f.Func().Id(builder.Constructor).Params(Id("opts").Op("...").Id(builderOption.PublicName)).Parens(Op("*").List(Id(builder.PublicName), Error())).Block(
		Var().Id(builder.ShortName).Id(builder.PublicName),
		For(List(Id("_"), Id("o")).Op(":=").Range().Id("opts")).Block(
			If(
				Err().Op(":=").Id("o").Call(Op("&").Id(builder.ShortName)),
				Err().Op("!=").Nil(),
			).Block(
				Return(Nil(), Err()),
			),
		),
		Return(Op("&").Id(builder.ShortName), Nil()),
	)
	f.Func().Id("WithPersonName").Params(Id("name").Id("string")).Id(builderOption.PublicName).Block(
		Return(Func().Params(Id(entity.ShortName).Op("*").Id(builder.PublicName)).Error().Block(
			Id(entity.ShortName).Dot(entity.PrivateName).Dot("name").Op("=").Id("name"),
			Return(Nil()),
		)),
	)
	fmt.Printf("%#v", f)
}
```

## Using environment variable 


```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// $ world=$(cat Makefile) go run test.go
	fmt.Println(os.ExpandEnv("hello $world"))
}
```
