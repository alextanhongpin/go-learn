# Generic Set
```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	s := New[string]()
	s.Add("hello")
	fmt.Println(s.Has("hello"))

	si := New[int]()
	si.Add(1, 2, 3)
	fmt.Println(si.Has(1))

	s2 := New[int]()
	s2.Add(1, 2)
	fmt.Println(si.Intersect(*s2))
	fmt.Println(si.List())
}

type Set[T comparable] map[T]bool

func New[T comparable]() *Set[T] {
	return &Set[T]{}
}

func (s Set[T]) Add(t ...T) {
	for _, k := range t {
		s[k] = true
	}
}

func (s Set[T]) Has(t T) bool {
	return s[t]
}

func (s Set[T]) Intersect(other Set[T]) []T {
	if len(s) > len(other) {
		return other.Intersect(s)
	}
	var result []T
	for k := range s {
		if other.Has(k) {
			result = append(result, k)
		}
	}
	return result
}

func (s Set[T]) List() []T {
	var result []T
	for k := range s {
		result = append(result, k)
	}
	return result
}
```

## Generic Map

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	users := []User{
		{"John", 10},
		{"Jane", 20},
	}
	ages := mapR(users, func(u User) int {
		return u.Age
	})
	fmt.Println(ages)
}

func mapR[T any, R any](in []T, fn func(T) R) []R {
	res := make([]R, len(in))
	for i, k := range in {
		res[i] = fn(k)
	}
	return res
}
```

## Pointer and value receiver


```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	msg := "hello world"
	msgp := Reference(msg)
	fmt.Println(msgp)
	fmt.Println(Value(msgp))
	fmt.Println(Value((*string)(nil)))
}

func Reference[T any](t T) *T {
	return &t
}

func Value[T any](t *T) (T, bool) {
	if t == nil {
		var tt T
		return tt, false
	}
	return *t, true
}
```

## Generic Middleware

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	createUser := LogTime(AddPronoun("Mr.", CreateUser))
	createUser(context.Background(), "john")
	fmt.Println("Hello, 世界")
}

type Decorator[Req any, Res any] func(ctx context.Context, req Req) (Res, error)

func CreateUser(ctx context.Context, name string) (id int, err error) {
	fmt.Println("creating user:", name)
	time.Sleep(1 * time.Second)
	return 0, errors.New("not implemented")
}

func LogTime[Req any, Res any](fn Decorator[Req, Res]) Decorator[Req, Res] {
	return func(ctx context.Context, req Req) (Res, error) {
		start := time.Now()
		defer func() {
			fmt.Println(time.Since(start))
		}()

		return fn(ctx, req)
	}
}

func AddPronoun(pronoun string, fn Decorator[string, int]) Decorator[string, int] {
	return func(ctx context.Context, req string) (int, error) {
		return fn(ctx, fmt.Sprintf("%s %s", pronoun, req))
	}
}
```


## Generic Decorators

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	res, err := Retry(CallAPI, 3)(context.Background(), "john")
	fmt.Println(res, err)
}

func CallAPI(ctx context.Context, name string) (id int, err error) {
	fmt.Println("calling", name)
	return 0, errors.New("bad request")
}

type AnyFunc[R any, W any] func(ctx context.Context, req R) (W, error)

func Retry[R any, W any](fn AnyFunc[R, W], n int) AnyFunc[R, W] {
	return func(ctx context.Context, req R) (res W, err error) {
		for i := 0; i < n; i++ {
			res, err = fn(ctx, req)
			if err == nil {
				return res, nil
			}
			seconds := (i + 1) * 1000
			duration := time.Millisecond * time.Duration((rand.Intn(seconds) + seconds/2))
			fmt.Println("retrying in", duration)
			time.Sleep(duration)
		}
		fmt.Println("retry failed")
		return res, err
	}
}
```

### Generic hint

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
	printInt := TypedFunc[int](print)
	printInt(42)
}

// If external lib still havent use generics, we can always assign type hints using generics.
func print(v any) error {
	fmt.Println(v)
	return nil
}

func TypedFunc[T any](fn func(any) error) func(T) error {
	return func(t T) error {
		return fn(t)
	}
}
```

### Generic Filter

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	list := List[int]([]int{1, 2, 3})
	result := list.Filter(func(i int) bool {
		return i > 2
	})
	fmt.Println(result)
}

type List[T any] []T

func (list List[T]) Filter(fn func(T) bool) List[T] {
	result := make([]T, 0, len(list))
	for _, item := range list {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}
```

## Generic Slice

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type Person struct {
	_    struct{}
	name string
	age  int
}

func main() {
	fmt.Println("Hello, 世界")
	numbers := []int{1, 2, 3, 4, 5}
	greaterThree := Filter(numbers, func(i int) bool {
		return i > 3
	})
	fmt.Println(greaterThree)

	people := []Person{
		{name: "john", age: 10},
		{name: "jane", age: 20},
	}
	personByName := ToMap(people, func(p Person) string {
		return p.name
	})
	fmt.Println(personByName)
}

func Filter[T any](list []T, fn func(T) bool) []T {
	result := make([]T, 0, len(list))
	for _, item := range list {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

func ToMap[K comparable, V any](list []V, getKeyFn func(V) K) map[K]V {
	result := make(map[K]V)
	for _, item := range list {
		key := getKeyFn(item)
		_, ok := result[key]
		if ok {
			panic(fmt.Errorf("key exists: %s", key))
		}
		result[key] = item
	}
	return result
}
```

## Setter getter again

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotSet             = errors.New("not set")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type User struct {
	// NOTE: Don't use private field.
	// name *Field[string]
	// age  *Field[int]
	Name  Getter[string] // Interface will panic.
	Age   Getter[int]
	Email Getter[string]
}

func (u *User) Valid() bool {
	return ValidGetter(u.Name) &&
		ValidGetter(u.Age)
}

func main() {
	u := User{
		Name:  NewField("john"),
		Age:   NewField(10),
		Email: NewEmail("john.doe@mail.com"), // Can use value object too.
	}
	if valid := u.Valid(); !valid {
		// Should panic here.
	}
	fmt.Println(u.Valid())
	fmt.Println(u.Name.Get()) // Suddenly it becomes a bad idea, should only use setter/getter with public fields.
	fmt.Println(u.Age.Get())  // Will panic.
	fmt.Println(u.Email.Get())
}

func ValidGetter[T any](getter Getter[T]) bool {
	return getter != nil && getter.Valid()
}

type SetterGetter[T any] interface {
	Setter[T]
	Getter[T]
}

type Setter[T any] interface {
	Set(T)
}

type Getter[T any] interface {
	Get() (T, bool)
	MustGet() T
	Valid() bool
}

type Field[T any] struct {
	value       T
	dirty       bool
	constructed bool
}

func NewField[T any](t T) *Field[T] {
	return &Field[T]{
		value:       t,
		dirty:       true,
		constructed: true,
	}
}

func (f Field[T]) Valid() bool {
	return f.Validate() == nil && f.dirty
}

func (f *Field[T]) Validate() error {
	if f == nil || !f.constructed {
		return ErrNotSet
	}
	return nil
}

func (f *Field[T]) Set(t T) {
	f.value = t
	f.dirty = true
}

func (f Field[T]) Get() (t T, valid bool) {
	if err := f.Validate(); err != nil {
		return
	}
	return f.value, f.dirty
}

func (f Field[T]) MustGet() T {
	if err := f.Validate(); err != nil {
		panic(err)
	}
	return f.value
}

type Email struct {
	value       string
	constructed bool
}

func NewEmail(email string) *Email {
	return &Email{
		value:       email,
		constructed: true,
	}
}

func (e *Email) Validate() error {
	if e == nil || !e.constructed {
		return ErrNotSet
	}
	if !strings.Contains(e.value, "@") {
		return ErrInvalidEmailFormat
	}
	return nil
}

func (e *Email) Valid() bool {
	return e.Validate() == nil
}
func (e *Email) Get() (string, bool) {
	if !e.Valid() {
		return "", false
	}
	return e.value, true
}
func (e *Email) MustGet() string {
	if !e.Valid() {
		panic(ErrInvalidEmailFormat)
	}
	return e.value
}
```