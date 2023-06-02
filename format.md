# Formating code

Should be as simple as:
```go
content, err := format.Source(content)
// check error
file.Write(content)
```

Example:
```go
// Other: How to format Go code programmatically
//
// Code can be formatted programmatically in the same way like running go fmt,
// using the go/format package
package main

import (
	"fmt"
	"go/format"
	"log"
)

func Example() {
	unformatted := `
package main
       import "fmt"

func  main(   )  {
    x :=    12
fmt.Printf(   "%d",   x  )
	}


`
	formatted, err := format.Source([]byte(unformatted))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(formatted))
	// Output:
	// package main
	//
	// import "fmt"
	//
	// func main() {
	//	x := 12
	//	fmt.Printf("%d", x)
	// }
}
```

References:
https://golang.org/pkg/go/format/

## Format with import

```go
package main

import "golang.org/x/tools/imports"

// FormatSource is gofmt with addition of removing any unused imports.
func FormatSource(source []byte) ([]byte, error) {
	return imports.Process("", source, &imports.Options{
		AllErrors: true, Comments: true, TabIndent: true, TabWidth: 8,
	})
}
```
