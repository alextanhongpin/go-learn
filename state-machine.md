## State Machine with no termination

```go
package main

import (
	"fmt"
)

type State string

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	return next, exist
}

const (
	GreenLight  = State("green_light")
	OrangeLight = State("orange_light")
	RedLight    = State("red_light")
)

func main() {
	var states = StateMachine{
		GreenLight:  OrangeLight,
		OrangeLight: RedLight,
		RedLight:    OrangeLight,
	}
	var initialState = GreenLight
	state := initialState
	for i := 0; i < 3; i++ {
		fmt.Printf("prev: %s", state)
		state, _ = states.Next(state)
		fmt.Printf(" next: %s\n", state)
	}
}
```

## With Termination, Single State

**NOTE**: This might be a little over-engineered, see the solution below:
```go
package main

import (
	"fmt"
)

type State string

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	if !exist {
		return Invalid, false
	}
	return next, next == Ended
}

const (
	Started = State("^")
	Ended   = State("$") // We need this to differentiate between non-existing state and completed state.
	Invalid = State("invalid")

	Submitted = State("submitted")
	Approved  = State("approved")
	Published = State("published")
)

func main() {
	var states = StateMachine{
		Started:   Submitted,
		Submitted: Approved,
		Approved:  Published,
		Published: Ended,
		Ended:     Ended,
	}
	var completed bool
	var initialState = Started
	state := initialState

	for i := 0; i < 5; i++ {
		fmt.Printf("prev: %s", state)
		state, completed = states.Next(state)
		fmt.Printf(" next: %s, completed: %t\n", state, completed)
		// Break to avoid infinite loop.
		if completed {
			break
		}
	}
}
```

Simplified solution:

```go
package main

import (
	"fmt"
	"strings"
)

type State string

func (s State) String() string {
	return string(s)
}

func (s State) Eq(in string) bool {
	return strings.EqualFold(string(s), in)
}

func (s State) EqStrict(in string) bool {
	return string(s) == in
}

type StateMachine map[State]State

func (s StateMachine) Next(prev State) (State, bool) {
	next, exist := s[prev]
	return next, !exist
}

const (
	Invalid   = State("invalid")
	Started   = State("started")
	Submitted = State("submitted")
	Approved  = State("approved")
	Published = State("published")
)

func main() {
	var states = StateMachine{
		Started:   Submitted,
		Submitted: Approved,
		Approved:  Published,
	}
	var completed bool
	var initialState = Started
	state := initialState

	for i := 0; i < 5; i++ {
		fmt.Printf("prev: %s", state)
		state, completed = states.Next(state)
		fmt.Printf(" next: %s, completed: %t\n", state, completed)
		// Break to avoid infinite loop.
		if completed {
			break
		}
	}
}
```

## With Termination, Multiple States

```go
package main

import (
	"fmt"
)

type State string

type States map[State][]State

func (s States) Next(prev State) ([]State, bool) {
	next, exist := s[prev]
	if !exist {
		return nil, false
	}
	return next, len(next) == 0
}

const (
	Initialized = State("initialized") // Start
	Submitted   = State("submitted")
	Rejected    = State("rejected")
	Approved    = State("approved") // End
)

func main() {
	// Assuming the states are describing payment.
	var states = States{
		Initialized: []State{Submitted},
		Submitted:   []State{Approved, Rejected},
		Rejected:    []State{Submitted},
		Approved:    []State{}, // After approved, there are no other continuation.
	}
	var initialState = Initialized
	state := initialState
	for i := 0; i < 3; i++ {
		choices, completed := states.Next(state)
		if completed {
			fmt.Println("completed")
			break
		}
		state = choices[0]
		fmt.Printf("next state is: %s, completed: %t\n", state, completed)
	}
}
```

## With switch, single state

```go
package main

import (
	"fmt"
)

func main() {
	var initialState = Red
	state := initialState
	for i := 0; i < 3; i++ {
		state = Next(state)
		fmt.Println("next state is:", state)
	}
}

type State string

const (
	Invalid = State("")
	Red     = State("red")
	Green   = State("green")
	Yellow  = State("yellow")
)

func Next(prev State) State {
	switch prev {
	case Red:
		return Green
	case Green:
		return Yellow
	case Yellow:
		return Red
	default:
		return Invalid
	}
}
```


## With switch, multiple states
```go
package main

import (
	"fmt"
)

func main() {
	var initialState = Initialized
	state := initialState
	for i := 0; i < 3; i++ {
		states, completed := Next(state)
		if completed {
			fmt.Println("completed")
			break
		}
		state = states[0]
		fmt.Printf("next state is: %s, completed: %t\n", state, completed)
	}
}

type State string

const (
	Initialized = State("payment_initialized") // Start
	Submitted   = State("payment_submitted")
	Rejected    = State("payment_rejected")
	Approved    = State("payment_approved") // End
)

func Next(prev State) ([]State, bool) {
	switch prev {
	case Initialized:
		return []State{Submitted}, false
	case Submitted:
		return []State{Approved, Rejected}, false
	case Rejected:
		return []State{Submitted}, false
	case Approved:
		return []State{}, true
	default:
		return []State{}, false
	}
}
```

## A more complex example

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

type CommandHandler interface {
	Handle() (string, error)
}

type EndCommand struct {
}

func (c *EndCommand) Handle() (string, error) {
	return "end", nil
}

type State struct {
	Commands map[string]CommandHandler `json:"commands"`
	invoke   func() error
}

func (s *State) On(cmd string, fn CommandHandler) *State {
	s.Commands[cmd] = fn
	return s
}

// TODO: Change to entry/exit.
func (s *State) Invoke(fn func() error) *State {
	s.invoke = fn
	return s
}

type StateMachine struct {
	ID      string            `json:"id"`
	Initial string            `json:"initial"`
	States  map[string]*State `json:"states"`
}

func NewStateMachine(initial string) *StateMachine {
	return &StateMachine{
		Initial: initial,
		States:  make(map[string]*State),
	}
}

func (sm *StateMachine) State(state string) *State {
	sm.States[state] = &State{
		Commands: make(map[string]CommandHandler),
	}
	return sm.States[state]
}

func (sm *StateMachine) Send(cmd string) error {
	state := sm.States[sm.Initial]
	if state.invoke != nil {
		if err := state.invoke(); err != nil {
			return err
		}
	}

	handler, ok := state.Commands[cmd]
	if handler == nil || !ok {
		// Invoke-only state machine.
		return nil
	}

	evt, err := handler.Handle()
	sm.Initial = evt
	if err != nil {
		return err
	}
	return nil
}

type OnCommand struct{}

func (o *OnCommand) Handle() (string, error) {
	return "on", nil
}

type OffCommand struct{}

func (o *OffCommand) Handle() (string, error) {
	return "off", nil
}

func main() {
	lightbulb()
	saga()
	log.Println(getFunctionName(lightbulb))
	name := func() {}
	log.Println(getFunctionName(name))
}

func lightbulb() {
	sm := NewStateMachine("off")
	sm.State("off").
		Invoke(func() error {
			log.Println("init off")
			return nil
		}).
		On("on", new(OnCommand))

	sm.State("on").
		Invoke(func() error {
			log.Println("init on")
			return nil
		}).
		On("off", new(OffCommand))

	if err := sm.Send("off"); err != nil {
		log.Println(err)
	}
	if err := sm.Send("on"); err != nil {
		log.Println(err)
	}
	if err := sm.Send("off"); err != nil {
		log.Println(err)
	}
	b, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

func saga() {
	saga := &Saga{
		StateMachine: NewStateMachine("0"),
	}
	createOrder := func() error {
		// ORDER_CREATED
		log.Println("create order")
		return nil
	}
	createPayment := func() error {
		// ORDER_CREATED
		log.Println("create payment")
		return errors.New("payment failed")
	}
	cancelOrder := func() error {
		// PAYMENT_CANCELLED
		log.Println("cancel order")
		return nil
	}
	createDelivery := func() error {
		// PAYMENT_CREATED
		log.Println("create delivery")
		// return errors.New("create delivery failed")
		return nil
	}
	cancelPayment := func() error {
		// DELIVERY_CANCELLED
		log.Println("cancel payment")
		return nil
	}
	confirmOrder := func() error {
		// DELIVERY_CREATED
		log.Println("confirm order")
		// return errors.New("order not confirmed")
		return nil
	}
	cancelDelivery := func() error {
		// ORDER REJECTED
		log.Println("cancel delivery")
		return nil
	}
	done := func() error {
		// ORDER_CONFIRMED
		log.Println("done")
		return nil
	}
	noop := func() error {
		return nil
	}

	saga.Add(createOrder, cancelOrder)
	saga.Add(createPayment, cancelPayment)
	saga.Add(createDelivery, cancelDelivery)
	saga.Add(confirmOrder, noop)
	saga.Add(done, noop)

	for !saga.Done() {
		if err := saga.Do(); err != nil {
			log.Println(err)
		}
	}

	b, err := json.Marshal(saga)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

type Saga struct {
	*StateMachine
	Step     int
	Progress int
	Logs     []string
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (s *Saga) Add(tx func() error, c func() error) {
	s.State(fmt.Sprint(s.Step)).Invoke(tx)
	s.State(fmt.Sprint(s.Step + 1)).Invoke(c)
	s.Step += 2
}
func (s *Saga) Done() bool {
	return s.Progress < 0 || s.Progress > s.Step-1
}

func (s *Saga) Do() error {
	if s.Done() {
		return nil
	}
	if s.Progress&1 == 0 {
		if err := s.Send(fmt.Sprint(s.Progress)); err != nil {
			s.Progress -= 1
			s.StateMachine.Initial = fmt.Sprint(s.Progress)
			return err
		}
		s.Progress += 2
		s.StateMachine.Initial = fmt.Sprint(s.Progress)
		return nil
	} else {
		if err := s.Send(fmt.Sprint(s.Progress)); err != nil {
			return err
		}
		s.Progress -= 2
		s.StateMachine.Initial = fmt.Sprint(s.Progress)
		return nil
	}
}
```


## Lightbulb

```go
package main

import (
	"fmt"
)

type State string

const (
	On  State = "on"
	Off State = "off"
)

func main() {
	m := &Machine{Initial: Off}
	m.Send(Off)
	m.Send(On)
	m.Send(On)
	m.Send(Off)
}

type Machine struct {
	Initial State
}

func (m *Machine) Send(next State) {
	switch m.Initial {
	case On:
		switch next {
		case Off:
			fmt.Println("turning off")
			m.Initial = Off
		default:
			fmt.Println("invalid next state")
		}

	case Off:
		switch next {
		case On:
			fmt.Println("turning on")
			m.Initial = On
		default:
			fmt.Println("invalid next state")
		}
	default:
		fmt.Println("invalid initial state")
	}
}
```
## Another implementation

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type (
	EventType string
	StateType string

	Action interface {
		Handle() (StateType, error)
	}

	State struct {
		// Entry func()
		// Exit func()
		Events map[EventType]Action
	}
)

type Machine struct {
	ID      string
	Initial StateType
	States  map[StateType]State
}

func (m *Machine) getAction(e EventType) Action {
	if state, ok := m.States[m.Initial]; ok {
		if action, ok := state.Events[e]; ok {
			return action
		}
		return nil
	}
	return nil
}

// TODO: Add context, and dynamic args.
func (m *Machine) Send(e EventType) error {
	if m.Initial == End {
		fmt.Println("done")
		return nil
	}
	action := m.getAction(e)
	if action == nil {
		fmt.Println("invalid event")
		return nil
	}

	// TODO:
	//if state.Entry == nil {
	//	return nil
	//}

	next, err := action.Handle()
	if err != nil {
		return err
	}
	if _, ok := m.States[next]; !ok {
		return errors.New("invalid transition")
	}
	m.Initial = next
	return nil
}

const (
	OrderCreated    StateType = "OrderCreated"
	OrderCancelled  StateType = "OrderCancelled"
	PaymentCreated  StateType = "PaymentCreated"
	PaymentRejected StateType = "PaymentRejected"
	PaymentRefunded StateType = "PaymentRefunded"
	Start           StateType = ""
	End             StateType = "end"

	CreateOrder   EventType = "CreateOrder"
	CancelOrder   EventType = "CancelOrder"
	CreatePayment EventType = "CreatePayment"
	RefundPayment EventType = "RefundPayment"
)

type ActionCreateOrder struct{}

func (a *ActionCreateOrder) Handle() (StateType, error) {
	fmt.Println("creating order")
	return OrderCreated, nil
}

type ActionCancelOrder struct{}

func (a *ActionCancelOrder) Handle() (StateType, error) {
	fmt.Println("cancelling order")
	return End, nil
}

type ActionCreatePayment struct{}

func (a *ActionCreatePayment) Handle() (StateType, error) {
	fmt.Println("creating payment")
	// return PaymentCreated, nil
	return PaymentRejected, errors.New("payment rejected: insufficient credit")
}

type ActionRefundPayment struct{}

func (a *ActionRefundPayment) Handle() (StateType, error) {
	fmt.Println("refunding payment")
	return PaymentRefunded, nil
}

// TODO: Use builder pattern to construct individual states.
func main() {
	m := &Machine{
		Initial: Start,
		States: map[StateType]State{
			Start: State{
				Events: map[EventType]Action{
					CreateOrder: new(ActionCreateOrder),
				},
			},
			// The resulting state.
			OrderCreated: State{
				Events: map[EventType]Action{
					// The create payment may have two possible states, success or failure.
					CreatePayment: new(ActionCreatePayment),
					CancelOrder:   new(ActionCancelOrder),
				},
			},
			PaymentCreated: State{
				Events: map[EventType]Action{
					RefundPayment: new(ActionRefundPayment),
				},
			},
			PaymentRefunded: State{
				Events: map[EventType]Action{
					CancelOrder: new(ActionCancelOrder),
				},
			},
			PaymentRejected: State{
				Events: map[EventType]Action{
					CancelOrder: new(ActionCancelOrder),
				},
			},
			OrderCancelled: State{
				Events: map[EventType]Action{},
			},
			End: State{},
		},
	}
	handleError(m.Send(CreateOrder))
	handleError(m.Send(CreatePayment))
	handleError(m.Send(CancelOrder))
	handleError(m.Send(CreateOrder))
	handleError(m.Send(RefundPayment))
	handleError(m.Send(CancelOrder))
	handleError(m.Send(CreateOrder))
	handleError(m.Send(CancelOrder))

	b, _ := json.MarshalIndent(m, "", "  ")
	log.Println(string(b))

}

func handleError(err error) {
	if err != nil {
		log.Println("error:", err)
	}
}

```
