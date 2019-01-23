## Strategy

```go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"
)

const (
	TEXT_STRATEGY  = "text"
	IMAGE_STRATEGY = "image"
)

type PrintStrategy interface {
	Print() error
	SetLog(io.Writer)
	SetWriter(io.Writer)
}

type PrintOutput struct {
	Writer    io.Writer
	LogWriter io.Writer
}

func (d *PrintOutput) SetLog(w io.Writer) {
	d.LogWriter = w
}
func (d *PrintOutput) SetWriter(w io.Writer) {
	d.Writer = w
}

type TextSquare struct {
	PrintOutput
}

func (c *TextSquare) Print() error {
	c.Writer.Write([]byte("Square"))
	return nil
}

type ImageSquare struct {
	PrintOutput
	DestinationFilePath string
}

func (t *ImageSquare) Print() error {
	width := 800
	height := 600
	origin := image.Point{0, 0}
	bgImage := image.NewRGBA(image.Rectangle{
		Min: origin,
		Max: image.Point{X: width, Y: height},
	})
	bgColor := image.Uniform{color.RGBA{R: 70, G: 70, B: 70, A: 0}}
	quality := &jpeg.Options{Quality: 75}
	draw.Draw(bgImage, bgImage.Bounds(), &bgColor, origin, draw.Src)

	if t.Writer == nil {
		return errors.New("no writer stored on ImageSquare")
	}

	if err := jpeg.Encode(t.Writer, bgImage, quality); err != nil {
		return errors.New("error writing image to disk")
	}

	if t.LogWriter != nil {
		r := bytes.NewReader([]byte("image writen in provided writer"))
		io.Copy(t.LogWriter, r)
	}
	return nil
}

func NewPrinter(s string) (PrintStrategy, error) {
	switch s {
	case TEXT_STRATEGY:
		return &TextSquare{
			PrintOutput: PrintOutput{
				LogWriter: os.Stdout,
			},
		}, nil
	case IMAGE_STRATEGY:
		return &ImageSquare{
			PrintOutput: PrintOutput{
				LogWriter: os.Stdout,
			},
		}, nil
	default:
		return nil, fmt.Errorf("strategy '%s' not found", s)
	}
}

func main() {
	output := "text"
	activeStrategy, err := NewPrinter(output)
	if err != nil {
		log.Fatal(err)
	}
	switch output {
	case TEXT_STRATEGY:
		activeStrategy.SetWriter(os.Stdout)
	case IMAGE_STRATEGY:
		w, err := os.Create("/tmp/image.jpg")
		if err != nil {
			log.Fatal("error opening image")
		}
		defer w.Close()
		activeStrategy.SetWriter(w)
	}
	err = activeStrategy.Print()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Chain of Responsibility

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

type ChainLogger interface {
	Next(string)
}

type FirstLogger struct {
	NextChain ChainLogger
}

func (f *FirstLogger) Next(s string) {
	fmt.Printf("First logger: %s\n", s)
	if f.NextChain != nil {
		f.NextChain.Next(s)
	}
}

type SecondLogger struct {
	NextChain ChainLogger
}

func (l *SecondLogger) Next(s string) {
	if strings.Contains(strings.ToLower(s), "hello") {
		fmt.Printf("Second logger: %s\n", s)

		if l.NextChain != nil {
			l.NextChain.Next(s)
		}
		return
	}
	fmt.Println("finishing in second logging")
}

type WriterLogger struct {
	NextChain ChainLogger
	Writer    io.Writer
}

func (w *WriterLogger) Next(s string) {
	if w.Writer != nil {
		w.Writer.Write([]byte("WriterLogger: " + s))
	}
	if w.NextChain != nil {
		w.NextChain.Next(s)
	}
}

type myTestWriter struct {
	receivedMessage *string
}

func (m *myTestWriter) Write(p []byte) (int, error) {
	if m.receivedMessage == nil {
		m.receivedMessage = new(string)
	}
	tmpMessage := fmt.Sprintf("%s%s", *m.receivedMessage, p)
	m.receivedMessage = &tmpMessage
	return len(p), nil
}

func (m *myTestWriter) Next(s string) {
	m.Write([]byte(s))
}

type ClosureChain struct {
	NextChain ChainLogger
	Closure   func(string)
}

func (c *ClosureChain) Next(s string) {
	if c.Closure != nil {
		c.Closure(s)
	}

	if c.NextChain != nil {
		c.Next(s)
	}
}
func main() {
	myWriter := myTestWriter{}

	closureLogger := ClosureChain{
		Closure: func(s string) {
			fmt.Printf("My Closure Logger: %s\n", s)
			myWriter.receivedMessage = &s
		},
	}

	writerLogger := WriterLogger{Writer: &myWriter}
	writerLogger.NextChain = &closureLogger
	second := SecondLogger{NextChain: &writerLogger}
	chain := FirstLogger{NextChain: &second}
	chain.Next("message that breaks the chain")
	fmt.Println(myWriter.receivedMessage)

	chain.Next("hello\n")
	fmt.Println(myWriter.receivedMessage)
}
```

## Command

```go
package main

import (
	"fmt"
)

type Command interface {
	Execute()
}

type ConsoleOutput struct {
	message string
}

func (c *ConsoleOutput) Execute() {
	fmt.Println(c.message)
}

func CreateCommand(s string) Command {
	fmt.Println("creating command")
	return &ConsoleOutput{
		message: s,
	}
}

type CommandQueue struct {
	queue []Command
}

func (p *CommandQueue) AddCommand(c Command) {
	p.queue = append(p.queue, c)
	if len(p.queue) == 3 {
		for _, command := range p.queue {
			command.Execute()
		}
		p.queue = make([]Command, 3)
	}

}
func main() {
	queue := CommandQueue{}
	queue.AddCommand(CreateCommand("first message"))
	queue.AddCommand(CreateCommand("second message"))
	queue.AddCommand(CreateCommand("third message"))
	queue.AddCommand(CreateCommand("fourth message"))
	queue.AddCommand(CreateCommand("fifth message"))
}
```


## Template

```go
package main

import (
	"fmt"
	"strings"
)

type MessageRetriever interface {
	Message() string
}

type Template interface {
	first() string
	third() string
	ExecuteAlgorithm(MessageRetriever) string
}

type TemplateImpl struct{}

func (t *TemplateImpl) first() string {
	return "first"
}
func (t *TemplateImpl) third() string {
	return "third"
}
func (t *TemplateImpl) ExecuteAlgorithm(m MessageRetriever) string {
	return strings.Join([]string{t.first(), m.Message(), t.third()}, " ")
}

type TestStruct struct {
	Template
}

func (m *TestStruct) Message() string {
	return "world"
}

type AnonymousTemplate struct{}

func (a *AnonymousTemplate) first() string {
	return "a:first"
}
func (a *AnonymousTemplate) third() string {
	return "a:third"
}
func (a *AnonymousTemplate) ExecuteAlgorithm(m MessageRetriever) string {
	return strings.Join([]string{a.first(), m.Message(), a.third()}, " ")
}

func main() {
	tpl := new(TemplateImpl)
	s := &TestStruct{}
	res := tpl.ExecuteAlgorithm(s)
	fmt.Println(res)

	atpl := new(AnonymousTemplate)
	res = atpl.ExecuteAlgorithm(s)
	fmt.Println(res)
}
```

## Memento

```go
package main

import (
	"errors"
	"fmt"
)

type State struct {
	Description string
}

type memento struct {
	state State
}

type originator struct {
	state State
}

func (o *originator) NewMemento() memento {
	return memento{state: o.state}
}

func (o *originator) ExtractAndStoreState(m memento) {
	o.state = m.state
}

type careTaker struct {
	mementoList []memento
}

func (c *careTaker) Add(m memento) {
	c.mementoList = append(c.mementoList, m)
}

func (c *careTaker) Memento(i int) (memento, error) {
	if len(c.mementoList) < i || i < 0 {
		return memento{}, errors.New("index not found")
	}
	return c.mementoList[i], nil
}

func main() {
	o := originator{}
	o.state = State{Description: "Idle"}
	ct := careTaker{}
	ct.Add(o.NewMemento())
	fmt.Println(ct)
}
```
