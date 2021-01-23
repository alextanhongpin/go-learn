# Naming Convention

From [effective go](https://golang.org/doc/effective_go.html#interface-names):

```
By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.

There are a number of such names and it's productive to honor them and the function names they capture. Read, Write, Close, Flush, String and so on have canonical signatures and meanings. To avoid confusion, don't give your method one of those names unless it has the same signature and meaning. Conversely, if your type implements a method with the same meaning as a method on a well-known type, give it the same name and signature; call your string-converter method String not ToString.
```

## Agentive suffix (`-er`)

`-er` in the sense of writ-er, bak-er, is the **agentive suffix**. It turns a verb into a noun that refers to the agent that performs that verb.


# Basic Overwrite

```go
package main

import (
	"fmt"
)

type A interface {
	Spew() string
}

type aImpl struct {
}

func (a *aImpl) Spew() string {
	return "hello"
}

type bImpl struct {
	A
}

// This will overwrite the A's Spew(), but A's Spew() can still be accessed by b.A.Spew()
func (b *bImpl) Spew() string {
	return "world"
}

func main() {
	a := new(aImpl)
	b := bImpl{a}

  // Use the latter (shorthand) version, it keeps the code shorter and more concise.
	fmt.Println(b.A.Spew(), b.Spew())
}
```

## Type-Casting

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Event interface {
	isEvent()
}

func (p PersonCreated) isEvent() {}
func (p PersonUpdated) isEvent() {}
func (p PersonRemoved) isEvent() {}

type PersonCreated struct{}
type PersonUpdated struct{}
type PersonRemoved struct{}

// Why this approach is not ideal.
type event interface {
	EventName() string
}

func getName(i interface{}) string {
	if t := reflect.TypeOf(i); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
func (t TodoCreated) EventName() string { return getName(t) }
func (t TodoUpdated) EventName() string { return getName(t) }
func (t TodoRemoved) EventName() string { return getName(t) }

type TodoCreated struct{}
type TodoUpdated struct{}
type TodoRemoved struct{}

func main() {
	personEvents := []Event{PersonCreated{}, PersonUpdated{}, PersonRemoved{}}
	for _, p := range personEvents {
		switch t := p.(type) {
		case PersonCreated:
			fmt.Println("creating person", t)
		case PersonUpdated:
			fmt.Println("updating person", t)
		case PersonRemoved:
			fmt.Println("removing person", t)
		default:
			log.Fatal(errors.New("not implemented"))
		}
	}

	// This does not require type-casting.
	// However, just knowing the type does not allow us to query the fields.
	// We still need to use the method above to get the type-casted fields.
	todoEvents := []event{TodoCreated{}, TodoUpdated{}, TodoRemoved{}}
	for _, t := range todoEvents {
		switch t.EventName() {
		case getName(TodoCreated{}):
			fmt.Println("creating todo", t)
		case getName(TodoUpdated{}):
			fmt.Println("updating todo", t)
		case getName(TodoRemoved{}):
			fmt.Println("removing todo", t)
		default:
			log.Fatal(errors.New("not implemented"))
		}
	}
}
```


## References
- http://objology.blogspot.com/2011/09/one-of-best-bits-of-programming-advice.html
- https://softwareengineering.stackexchange.com/questions/131667/interface-naming-prefix-can-vs-suffix-able


## Choices - Using adapter for converting between interface

The pragmatic way is to define a partial interface at the caller side to select only methods that are in used.

However, when programming, this makes it hard to trace back to the original struct that implements this interface. 

Hence, it could be easier to expose an adapter that shows the original interface (or struct) to partial interface. The original interface should be placed at the same folder that implements it.
```go
package main

import (
	"fmt"
)

type UserRepository interface {
	Create()
	Update()
	Delete()
}

type LoginRepository interface {
	Create()
}

type userRepository struct {
}

// Prefer "r" over "l", because changing the name of the repository requires major change.
func (r *userRepository) Create() {
	fmt.Println("create")
}

func (r *userRepository) Update() {
	fmt.Println("update")
}

func (r *userRepository) Delete() {
	fmt.Println("delete")
}

func adapter(r UserRepository) LoginRepository {
	return r
}

func main() {
	r := new(userRepository)
	r.Create()
	r.Update()
	r.Delete()

	l := adapter(r)
	l.Create()
}
```
