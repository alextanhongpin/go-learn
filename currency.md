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


## Calculations

Don't use float, use dedicated library:

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {
	money := 45.1235

	fmt.Println(isPreciseStd(money, 2))
	fmt.Println(isPreciseStd(money, 4)) // Note that we already have precision error here ...
	fmt.Println(isPrecise(money, 2))
	fmt.Println(isPrecise(money, 4))

	v := decimal.NewFromFloat(money)
	fmt.Println(v.RoundCash(10))
	fmt.Println(v.RoundBank(4))
}

func isPrecise(f float64, decimals int) bool {
	v := decimal.NewFromFloat(f)
	return v == v.Round(int32(decimals))
}

// Don't use this, it's not accurate.
func isPreciseStd(f float64, decimals int) bool {
	p := float64(tenPower(decimals))
	// With 45.1234, we have precision error here ...
	// 45.1234 == 45.1233
	return f == float64(int(f*p))/p
}

func tenPower(n int) int {
	res := 1
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}
```
