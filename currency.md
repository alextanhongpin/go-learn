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

## Another example

```go
package main

import (
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func main() {
	n := display.Tags(language.English)
	for _, lcode := range []string{"en_US", "pt_BR", "de", "ja", "hi"} {
		lang := language.MustParse(lcode)
		cur, _ := currency.FromTag(lang)
		scale, _ := currency.Cash.Rounding(cur) // fractional digits
		dec := number.Decimal(100000.00, number.Scale(scale))
		p := message.NewPrinter(lang)
		p.Printf("%24v (%v): %v%v\n", n.Name(lang), cur, currency.Symbol(cur), dec)
	}
}
```
