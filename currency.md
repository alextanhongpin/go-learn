## Formatting local currency

```go
package main

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer *message.Printer

func init() {
	base := language.MustParse("en-SG")
	printer = message.NewPrinter(base)
}

func SGD(amount float64) string {
	return printer.Sprintf("%s %.2f", "SGD", amount)
}

func main() {
	fmt.Printf("%.2f\n", 12345678.90)
	fmt.Println(SGD(12345678.90))
}

```
