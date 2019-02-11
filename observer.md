# Simple event emitter

```go
package main

import (
	"fmt"
)

type ObserverFunc func(...interface{}) error
type Observer interface {
	On(event string, fn ObserverFunc)
	Emit(event string, params ...interface{}) error
}

type ObserverImpl struct {
	events map[string][]ObserverFunc
}

func (o *ObserverImpl) On(event string, fn ObserverFunc) {
	_, exist := o.events[event]
	if !exist {
		o.events[event] = make([]ObserverFunc, 0)
	}
	o.events[event] = append(o.events[event], fn)
}

func (o *ObserverImpl) Emit(event string, params ...interface{}) error {
	fns, exist := o.events[event]
	if !exist {
		return fmt.Errorf(`event "%s" does not exist`, event)
	}
	for _, fn := range fns {
		err := fn(params...)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewObserver() *ObserverImpl {
	events := make(map[string][]ObserverFunc)
	return &ObserverImpl{events}
}

type User struct {
	Observer
}

func (u *User) Greet(msg string) {
	fmt.Println("greeting", msg)
	u.Emit("greet", msg)
	u.Emit("something", msg, "something")
	u.Emit("car", msg)
}

func main() {
	user := &User{NewObserver()}
	user.On("greet", func(msg ...interface{}) error {
		fmt.Println("got", msg)
		return nil
	})
	user.On("something", func(msg ...interface{}) error {
		fmt.Println("got", msg)
		return nil
	})
	user.Greet("hello!")
}
```
