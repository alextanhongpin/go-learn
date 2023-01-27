```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

func main() {
	onState := NewLightOn()
	fmt.Println("onState:", onState)

	offState := onState.Off()
	fmt.Println("offState:", offState)

	b, err := json.Marshal(offState)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	var state LightState
	if err := json.Unmarshal(b, &state); err != nil {
		panic(err)
	}
	fmt.Printf("state: %#v\n", state)
	fmt.Println(state.Type())
	fmt.Println(state.Off.On())
	// fmt.Println(state.On.Off()) // This will panic.
}

type LightStateType int

const (
	LightStateTypeZero = iota
	LightStateTypeOn
	LightStateTypeOff
)

func (t LightStateType) String() string {
	switch t {
	case LightStateTypeOn:
		return "light_on"
	case LightStateTypeOff:
		return "light_off"
	default:
		return ""
	}
}

// Structs are useful to represent union types.
// There can only be one valid state at a time.
// Each state is represented by separate structs.
// Each struct has their own method for state transition.
// The LightState struct is simply a factory that unmarshals the
// json representation of the state to the current state.
type LightState struct {
	On  *LightOn
	Off *LightOff
}

func (s *LightState) Type() LightStateType {
	switch {
	case s.On != nil:
		return LightStateTypeOn
	case s.Off != nil:
		return LightStateTypeOff
	default:
		return LightStateTypeZero
	}
}

func (s *LightState) UnmarshalJSON(data []byte) error {
	type state struct {
		Type LightStateType `json:"type"`
	}
	var p state
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	switch p.Type {
	case LightStateTypeOn:
		s.On = NewLightOn()
	case LightStateTypeOff:
		s.Off = NewLightOff()
	default:
		return errors.New("unknown state")
	}
	return nil
}

type LightOn struct {
	Type LightStateType `json:"type"`
}

func NewLightOn() *LightOn {
	return &LightOn{
		Type: LightStateTypeOn,
	}
}
func (l *LightOn) Off() *LightOff {
	// NOTE: This check is important.
	// Methods can still be called on nil structs without panicking.
	if l == nil {
		panic("state not initialized")
	}

	return NewLightOff()
}

type LightOff struct {
	Type LightStateType `json:"type"`
}

func NewLightOff() *LightOff {
	return &LightOff{
		Type: LightStateTypeOff,
	}
}
func (l *LightOff) On() *LightOn {
	if l == nil {
		panic("state not initialized")
	}

	return NewLightOn()
}
```
