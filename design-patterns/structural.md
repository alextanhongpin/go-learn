## Composite

```go
package main

import (
	"fmt"
)

type Athlete struct{}

func (a *Athlete) Train() {
	fmt.Println("training")
}

type CompositeSwimmerA struct {
	MyAthlete *Athlete
	MySwim    func()
}

func Swim() {
	fmt.Println("swimming")
}

func main() {
	swimmer := CompositeSwimmerA{
		MySwim: Swim,
	}
	swimmer.MyAthlete.Train()
	swimmer.MySwim()

	trainAndSwimStruct(swimmer)
	trainAndSwim(swimmer.MyAthlete)
}

type Trainable interface {
	Train()
}

func trainAndSwimStruct(swimmer struct {
	MyAthlete *Athlete
	MySwim    func()
}) {
	swimmer.MyAthlete.Train()
	swimmer.MySwim()
}

func trainAndSwim(swimmer interface{ Trainable }) {
	swimmer.Train()
}
```

## Adapter 

```go
package main

import (
	"fmt"
)

type LegacyPrinter interface {
	Print(s string) string
}

type MyLegacyPrinter struct{}

func (l *MyLegacyPrinter) Print(s string) (newMsg string) {
	newMsg = fmt.Sprintf("legacy printer: %s\n", s)
	return
}

type ModernPrinter interface {
	PrintStored() string
}

type PrinterAdapter struct {
	OldPrinter LegacyPrinter
	Msg        string
}

func (p *PrinterAdapter) PrintStored() (newMsg string) {
	if p.OldPrinter != nil {
		newMsg = fmt.Sprintf("Adapter: %s", p.Msg)
		newMsg = p.OldPrinter.Print(newMsg)
		return
	}
	newMsg = p.Msg
	return
}

func main() {
	msg := "hello world"
	adapter := PrinterAdapter{
		OldPrinter: &MyLegacyPrinter{},
		Msg:        msg,
	}
	returnedMsg := adapter.PrintStored()
	fmt.Println(returnedMsg)

	newAdapter := PrinterAdapter{
		Msg: msg,
	}
	returnedMsg = newAdapter.PrintStored()
	fmt.Println(returnedMsg)
}
```
