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


## Memento, Facade and Command

```go
package main

import (
	"errors"
	"fmt"
)

type Command interface {
	GetValue() interface{}
}

type Volume byte

func (v Volume) GetValue() interface{} {
	return v
}

type Mute bool

func (m Mute) GetValue() interface{} {
	return m
}

type Memento struct {
	memento Command
}

type originator struct {
	Command Command
}

func (o *originator) NewMemento() Memento {
	return Memento{memento: o.Command}
}

func (o *originator) ExtractAndStoreCommand(m Memento) {
	o.Command = m.memento
}

type careTaker struct {
	mementoStack []Memento
}

func (c *careTaker) Add(m Memento) {
	c.mementoStack = append(c.mementoStack, m)
}
func (c *careTaker) Memento(i int) (Memento, error) {
	if len(c.mementoStack) < i || i < 0 {
		return Memento{}, errors.New("index not found")
	}
	return c.mementoStack[i], nil
}
func (c *careTaker) Pop() Memento {
	if len(c.mementoStack) > 0 {
		tempMemento := c.mementoStack[len(c.mementoStack)-1]
		c.mementoStack = c.mementoStack[0 : len(c.mementoStack)-1]
		return tempMemento
	}
	return Memento{}
}

type MementoFacade struct {
	originator originator
	careTaker  careTaker
}

func (m *MementoFacade) SaveSettings(c Command) {
	m.originator.Command = c
	m.careTaker.Add(m.originator.NewMemento())
}

func (m *MementoFacade) RestoreSettings(i int) Command {
	mem, _ := m.careTaker.Memento(i)
	m.originator.ExtractAndStoreCommand(mem)
	return m.originator.Command
}

func main() {
	m := MementoFacade{}

	m.SaveSettings(Volume(4))
	m.SaveSettings(Mute(true))

	assertCommand(m.RestoreSettings(0))
	assertCommand(m.RestoreSettings(1))
}

func assertCommand(c Command) {
	switch v := c.(type) {
	case Volume:
		fmt.Println("Volume:", v)
	case Mute:
		fmt.Println("Mute:", v)
	}
}
```


## Visitor

```go
package main

import (
	"fmt"
	"io"
	"os"
)

type MessageA struct {
	Msg    string
	Output io.Writer
}

func (m *MessageA) Accept(v Visitor) {
	v.VisitA(m)
}

func (m *MessageA) Print() {
	if m.Output == nil {
		m.Output = os.Stdout
	}
	fmt.Fprintf(m.Output, "A: %s", m.Msg)
}

type MessageB struct {
	Msg    string
	Output io.Writer
}

func (m *MessageB) Accept(v Visitor) {
	v.VisitB(m)
}

func (m *MessageB) Print() {
	if m.Output == nil {
		m.Output = os.Stdout
	}
	fmt.Fprintf(m.Output, "B: %s", m.Msg)
}

type Visitor interface {
	VisitA(*MessageA)
	VisitB(*MessageB)
}

type Visitable interface {
	Accept(Visitor)
}

type MessageVisitor struct {
}

func (mf *MessageVisitor) VisitA(m *MessageA) {
	m.Msg = fmt.Sprintf("%s %s\n", m.Msg, "(Visited A)")
}

func (mf *MessageVisitor) VisitB(m *MessageB) {
	m.Msg = fmt.Sprintf("%s %s\n", m.Msg, "(Visited B)")
}

func main() {
	//	testHelper := &TestHelper{}
	visitor := &MessageVisitor{}
	{
		msg := MessageA{Msg: "hello world"}

		msg.Accept(visitor)
		msg.Print()
	}
	{
		msg := MessageB{Msg: "hallo welt"}

		msg.Accept(visitor)
		msg.Print()
	}
}
```

## State

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type GameState interface {
	executeState(*GameContext) bool
}
type GameContext struct {
	SecretNumber int
	Retries      int
	Won          bool
	Next         GameState
}

type StartState struct {
}

func (s *StartState) executeState(c *GameContext) bool {
	c.Next = &AskState{}

	rand.Seed(time.Now().UnixNano())
	c.SecretNumber = rand.Intn(10)
	fmt.Println("Introduce a number of retries to set the difficulty:")
	fmt.Fscanf(os.Stdin, "%d", &c.Retries)
	return true
}

type AskState struct{}

func (a *AskState) executeState(c *GameContext) bool {
	fmt.Printf("Introduce a number between 0 and 10, you have %d tries left\n", c.Retries)
	var n int
	fmt.Fscanf(os.Stdin, "%d", &n)
	c.Retries = c.Retries - 1
	if n == c.SecretNumber {
		c.Won = true
		c.Next = &FinishState{}
	}
	if c.Retries == 0 {
		c.Next = &FinishState{}
	}
	return true
}

type FinishState struct{}

func (f *FinishState) executeState(c *GameContext) bool {
	if c.Won {
		c.Next = &WinState{}
	} else {
		c.Next = &LoseState{}
	}
	return true
}

type WinState struct{}

func (w *WinState) executeState(c *GameContext) bool {
	fmt.Println("congrats, you won")
	return false
}

type LoseState struct{}

func (l *LoseState) executeState(c *GameContext) bool {
	fmt.Println("you lose")
	return false
}
func main() {

	start := StartState{}
	game := GameContext{
		Next: &start,
	}

	for game.Next.executeState(&game) {
	}
}
```

## Mediator

TODO: Find examples

## Observer

```go
package main

import (
	"fmt"
)

type Observer interface {
	Notify(string)
}

type Publisher struct {
	ObserverList []Observer
}

func (p *Publisher) AddObserver(o Observer) {
	p.ObserverList = append(p.ObserverList, o)

}

func (p *Publisher) RemoveObserver(o Observer) {
	var indexToRemove int

	for i, observer := range p.ObserverList {
		if observer == o {
			indexToRemove = i
			break
		}
	}
	p.ObserverList = append(p.ObserverList[:indexToRemove], p.ObserverList[indexToRemove+1:]...)
}

func (p *Publisher) NotifyObservers(m string) {
	fmt.Printf("Publisher received message '%s' to notify observers\n", m)
	for _, o := range p.ObserverList {
		o.Notify(m)
	}
}

type TestObserver struct {
	ID      int
	Message string
}

func (p *TestObserver) Notify(m string) {
	fmt.Printf("Observer %d: message '%s' received\n", p.ID, m)
	p.Message = m
}

func main() {

	testObserver1 := &TestObserver{1, ""}
	testObserver2 := &TestObserver{2, ""}
	testObserver3 := &TestObserver{3, ""}
	publisher := Publisher{}
	publisher.AddObserver(testObserver1)
	publisher.AddObserver(testObserver2)
	publisher.AddObserver(testObserver3)

	publisher.NotifyObservers("new message")
}
```
