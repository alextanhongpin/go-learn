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


## Generic

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sort"
)

type Greet struct {
	Message string
}

func HandleGreet(ctx context.Context, greet Greet) error {
	fmt.Println("greetings,", greet.Message)

	return nil
}

func HandleFormalizedGreet(ctx context.Context, greet Greet) error {
	fmt.Println("ahem, greetings,", greet.Message)

	return nil
}

func main() {
	obs := New[Greet]()
	if err := obs.On(HandleGreet, HandleFormalizedGreet); err != nil {
		log.Fatalf("failed to register: %v", err)
	}

	if err := obs.Emit(context.Background(), Greet{Message: "hello world"}); err != nil {
		log.Fatalf("failed to emit: %v", err)
	}

	mgr := &ObserverManager{greet: obs}
	if err := mgr.Greet(context.Background(), Greet{Message: "hello world"}); err != nil {
		log.Fatalf("failed to emit: %v", err)
	}
}

type ObserverManager struct {
	greet *Observer[Greet]
}

func (mgr *ObserverManager) Greet(ctx context.Context, greet Greet) error {
	return mgr.greet.Emit(ctx, greet)
}

type GenericHandler[T any] func(ctx context.Context, req T) error

type Observer[T any] struct {
	handlers map[string]GenericHandler[T]
}

func New[T any]() *Observer[T] {
	return &Observer[T]{
		handlers: make(map[string]GenericHandler[T]),
	}
}

func (o *Observer[T]) On(handlers ...GenericHandler[T]) error {
	for _, handler := range handlers {
		funcName := GetFunctionName(handler)
		_, ok := o.handlers[funcName]
		if ok {
			return errors.New("handler exists")
		}

		o.handlers[funcName] = handler
		fmt.Printf("registered: handler=%s\n", funcName)
	}

	return nil
}

func (o *Observer[T]) Emit(ctx context.Context, evt T) error {
	keys := GetKeys(o.handlers)
	sort.Strings(keys)

	for _, key := range keys {
		handler := o.handlers[key]
		if err := handler(ctx, evt); err != nil {
			return err
		}
	}

	return nil
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetTypeName[T any](unk T) string {
	t := reflect.TypeOf(unk)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t.Name()
}

func GetKeys[T comparable, U any](m map[T]U) []T {
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
```
