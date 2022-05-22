
Let's explore the usecases for generics. Not all implementation here are idiomatic, so take it with a grain of salt. We will not explore implementation of generic Set, Map and Slice since they are pretty common.

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

## Generic Pointer/Value method

```go
func Pointer[T comparable](t T) *T {
	return &t
}

func Value[T comparable](t *T) (res T, ok bool) {
	if t == nil {
		return
	}
	
	return *t, true
}
```

## Field-level hook

There are some limitation when using generics

1) generic type cannot be applied on struct methods

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func StructFields(in interface{}) func() *Required {
	t := reflect.Indirect(reflect.ValueOf(in)).Type()

	fields := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields[i] = t.Field(i).Name
	}

	return func() *Required {
		return NewRequired(fields[0], fields[1:]...)
	}
}

var UserFields = StructFields(&User{})

func main() {
	required := UserFields()

	var u User
	u.Name = Hook("Name", "john", required.Set)
	u.Age = Hook("Age", 10, required.Set)
	// u.Married = Hook("Married", true, required.Set)

	fmt.Println(required.Error()) // missing fields: Married
}

type User struct {
	Name    string
	Age     int
	Married bool
}

func Hook[T any](name string, t T, fn func(name string)) T {
	fn(name)
	return t
}

var ErrMissingFields = errors.New("missing fields")

type Required struct {
	fields []string
	value  map[string]bool
}

func NewRequired(field string, fields ...string) *Required {
	fields = append(fields, field)
	value := make(map[string]bool)
	for _, field := range fields {
		value[field] = false
	}
	return &Required{
		fields: fields,
		value:  value,
	}
}

func (r *Required) Set(name string) {
	r.value[name] = true
}

func (r *Required) Valid() bool {
	for _, field := range r.fields {
		if !r.value[field] {
			return false
		}
	}
	return len(r.fields) == len(r.value)
}

func (r *Required) Error() error {
	missing := make([]string, 0, len(r.fields))
	for _, field := range r.fields {
		if !r.value[field] {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingFields, strings.Join(missing, ", "))
	}

	return nil
}
```

## Generic Builder

This example below is not idiomatic go, use it only if it fits your usecase. 

There are some limitations for this approach:
- private fields cannot be inferred through reflection
- no type safety

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name    string
	Age     int
	Married bool
}

var UserBuilder = BuilderFactory(&User{})

func main() {
	b := UserBuilder()
	user, err := b.Set("Name", "john").
		Set("Age", 10).
		Set("Married", true).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

var ErrMissingFields = errors.New("missing fields")

func BuilderFactory[T any](in T) func() *Builder[T] {
	t := reflect.Indirect(reflect.ValueOf(in)).Type()

	fields := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields[i] = t.Field(i).Name
	}

	return func() *Builder[T] {
		return NewBuilder[T](fields...)
	}
}

type Builder[T any] struct {
	fields map[string]bool
	values map[string]interface{}
}

func NewBuilder[T any](requiredFields ...string) *Builder[T] {
	fields := make(map[string]bool)
	for _, f := range requiredFields {
		fields[f] = false
	}
	return &Builder[T]{
		values: make(map[string]interface{}),
		fields: fields,
	}
}

func (b *Builder[T]) Set(key string, value interface{}) *Builder[T] {
	if _, exists := b.fields[key]; exists {
		b.fields[key] = true
	}
	b.values[key] = value
	return b
}

func (b *Builder[T]) Build() (t T, err error) {
	var missing []string
	for field, set := range b.fields {
		if !set {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return t, fmt.Errorf("%w: %s", ErrMissingFields, strings.Join(missing, ", "))
	}

	bt, err := json.Marshal(b.values)
	if err != nil {
		return t, err
	}

	dec := json.NewDecoder(bytes.NewReader(bt))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		return t, err
	}

	return
}
```


## Builder v2


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

### Generic Slice

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

import (
	"fmt"

	"play.ground/slice"
)

func main() {
	numbers := []int{1, 2, 3}
	ids := slice.Map(numbers, func(n int, _ int) string {
		return fmt.Sprint(n)
	})

	fmt.Println(ids)
	isPositive := slice.All(numbers, func(i int) bool {
		return i > 0
	})
	fmt.Println(isPositive)

	hasTwo := slice.Any(ids, func(id string) bool {
		return id == "2"
	})
	fmt.Println(hasTwo)

	idx, found := slice.FindIndex(numbers, func(i int) bool {
		return i == 42
	})
	fmt.Println(idx, found)
}
-- go.mod --
module play.ground
-- slice/slice.go --
package slice

func Map[T any, R any](list []T, fn func(T, int) R) []R {
	result := make([]R, len(list))
	for i, t := range list {
		result[i] = fn(t, i)
	}
	return result
}

func All[T any](list []T, fn func(T) bool) bool {
	for _, item := range list {
		if !fn(item) {
			return false
		}
	}
	return true
}

func Any[T any](list []T, fn func(T) bool) bool {
	for _, item := range list {
		if fn(item) {
			return true
		}
	}
	return false
}

func Find[T any](list []T, fn func(T) bool) (t T, found bool) {
	for _, item := range list {
		if fn(item) {
			t = item
			found = true
			return
		}
	}
	return
}

func FindIndex[T any](list []T, fn func(T) bool) (index int, found bool) {
	for i, item := range list {
		if fn(item) {
			index = i
			found = true
			return
		}
	}
	return
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

## Some more generic slice
```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	fmt.Println(Slice[int]{}.FillZero(3))
	fmt.Println(Slice[bool]{}.FillZero(3))
}

func Echo[T any](fn func() T) T {
	return fn()
}

type Slice[T any] []T

func (s Slice[T]) FillFunc(n int, fill func() T) []T {
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = fill()
	}
	return result
}

func (s Slice[T]) FillZero(n int) []T {
	result := make([]T, n)
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

## Getter v2
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	Name *Getter[string] `json:"name,omitempty"`
	Age  *Getter[int]    `json:"age"`
}

func (u *User) Valid() bool {
	return u.Name.Valid() && u.Age.Valid()
}

var ErrNotSet = errors.New("not set")

func main() {
	var u User
	if err := json.Unmarshal([]byte(`{
		"name": "alice",
		"age": 13
	}`), &u); err != nil {
		panic(err)
	}
	fmt.Println(u, u.Name.Valid(), u.Age.Valid(), u.Valid())

	u2 := User{
		Name: NewGetter("bob"),
		Age:  NewGetter(42),
	}
	fmt.Println(u2, u2.Valid())
	b, err := json.Marshal(u2)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

type Getter[T any] struct {
	value T
	dirty bool
}

func NewGetter[T any](t T) *Getter[T] {
	return &Getter[T]{value: t, dirty: true}
}

func (g *Getter[T]) Validate() error {
	if g == nil || !g.dirty {
		return ErrNotSet
	}
	return nil
}

func (g *Getter[T]) Valid() bool {
	return g.Validate() == nil
}

func (g *Getter[T]) Get() (t T, valid bool) {
	if !g.Valid() {
		return
	}
	return g.value, true
}
func (g *Getter[T]) MustGet() (t T) {
	if !g.Valid() {
		panic(ErrNotSet)
	}
	return g.value
}
func (g Getter[T]) String() string {
	return fmt.Sprint(g.value)
}

func (g Getter[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.value)
}

func (g *Getter[T]) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}
	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}
	g.value = t
	g.dirty = true
	return nil
}
```

## Type Converter

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
)

type CreateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	req := CreateUserRequest{
		Name: "john",
		Age:  13,
	}
	user, err := TypeConverter[User](req)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

func TypeConverter[T any](s any) (T, error) {
	var t T
	b, err := json.Marshal(s)
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return t, err
	}
	return t, nil
}
```

## Typed template

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"html/template"
	"io"
	"os"
)

var greet = NewTypedTemplate[*Greet](template.Must(template.New("").Parse(`hello {{.Name}}`)))
var greet2 = TypedTemplateFunc[*Greet](template.Must(template.New("").Parse(`hello {{.Name}}`)))

type Greet struct {
	Name string
}

type Person struct {
	Name string
	Age  int
}

type Template[T any] interface {
	Execute(wr io.Writer, data T) error
}

type TypedTemplate[T any] struct {
	tpl Template[any]
}

func NewTypedTemplate[T any](tpl Template[any]) *TypedTemplate[T] {
	return &TypedTemplate[T]{
		tpl: tpl,
	}
}

func (t *TypedTemplate[T]) Execute(wr io.Writer, data T) error {
	return t.tpl.Execute(wr, data)
}

func TypedTemplateFunc[T any](t Template[any]) func(io.Writer, T) error {
	return func(wr io.Writer, data T) error {
		return t.Execute(wr, data)
	}
}

func main() {
	if err := greet.Execute(os.Stdout, &Greet{"world"}); err != nil {
		panic(err)
	}
	if err := greet2(os.Stdout, &Greet{"world"}); err != nil {
		panic(err)
	}
}

```

## Service Hooks

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"play.ground/service"
)

type Person struct {
	Salutation string
	Name       string
}

func main() {
	hook := service.New[string, *Person](GreetError)
	// hook := service.New[string, *Person](Greet)
	hook.Prepend(Precondition)
	hook.Append(Postcondition)
	hook.Decorate(
		Log[string, *Person], // Order matters. 
		Retry[string, *Person](3),
	)
	res, err := hook.Handle(context.Background(), "john")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func Precondition(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("required")
	}
	return name, nil
}

func Postcondition(ctx context.Context, person *Person) (*Person, error) {
	if !strings.EqualFold(person.Salutation, "Mr.") {
		return nil, errors.New("missing prefix")
	}
	return person, nil
}

func Greet(ctx context.Context, name string) (*Person, error) {
	return &Person{
		Salutation: "Mr.",
		Name:       name,
	}, nil
}

func GreetError(ctx context.Context, name string) (*Person, error) {
	return nil, errors.New("bad person")
}

func Log[Req any, Res any](fn service.Handler[Req, Res]) service.Handler[Req, Res] {
	return func(ctx context.Context, req Req) (Res, error) {
		start := time.Now()
		defer func() {
			fmt.Println(time.Since(start))
		}()
		return fn(ctx, req)
	}
}

func Retry[Req any, Res any](n int) service.Decorator[Req, Res] {
	return func(fn service.Handler[Req, Res]) service.Handler[Req, Res] {
		return func(ctx context.Context, req Req) (res Res, err error) {
			for i := 0; i < n; i++ {
				res, err = fn(ctx, req)
				if err == nil {
					return res, err
				}

				ms := (i + 1) * 1000
				ms = rand.Intn(ms) + ms/2
				fmt.Printf("sleep for %d ms\n", ms)
				time.Sleep(time.Duration(ms) * time.Millisecond)
			}
			err = fmt.Errorf("%w: too many retries")
			return
		}
	}
}
-- go.mod --
module play.ground
-- service/service.go --
package service

import (
	"context"
)

type Handler[Req any, Res any] func(ctx context.Context, req Req) (Res, error)

type Decorator[Req any, Res any] func(h Handler[Req, Res]) Handler[Req, Res]

type Func[T any] func(ctx context.Context, t T) (T, error)

type Hook[Req any, Res any] struct {
	handle     Handler[Req, Res]
	before     []Func[Req]
	after      []Func[Res]
	decorators []Decorator[Req, Res]
}

func New[Req any, Res any](fn Handler[Req, Res]) *Hook[Req, Res] {
	return &Hook[Req, Res]{
		handle: fn,
	}
}

func (hook *Hook[Req, Res]) Decorate(fns ...Decorator[Req, Res]) {
	hook.decorators = append(hook.decorators, fns...)
}

func (hook *Hook[Req, Res]) Append(fns ...Func[Res]) {
	hook.after = append(hook.after, fns...)
}

func (hook *Hook[Req, Res]) Prepend(fns ...Func[Req]) {
	hook.before = append(hook.before, fns...)
}

func (hook *Hook[Req, Res]) Handle(ctx context.Context, req Req) (res Res, err error) {
	for _, fn := range hook.before {
		req, err = fn(ctx, req)
		if err != nil {
			return
		}
	}

	decorated := hook.handle
	for _, decorator := range Reverse(hook.decorators) {
		decorated = decorator(decorated)
	}

	res, err = decorated(ctx, req)
	if err != nil {
		return
	}

	for _, fn := range hook.after {
		res, err = fn(ctx, res)
		if err != nil {
			return
		}
	}
	return
}

func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
```




### Base Value Object

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func main() {
	age, err := NewAge(10)
	if err != nil {
		panic(err)
	}
	fmt.Println(age.Valid())
	fmt.Println(age.IsZero())
	fmt.Println(age.Validate())
	fmt.Println(age.Set(-10))
	fmt.Println(age.Get())
	fmt.Println(age.MustGet())
	fmt.Println(age)

	b, err := json.Marshal(age)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshall", string(b))

	u := User{Age: age}
	b, err = json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("user", string(b))

	u = User{}
	b, err = json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("user", string(b))

	var john User
	if err := json.Unmarshal([]byte(`{"age": -1}`), &john); err != nil {
		if errors.Is(err, ErrInvalidValue) {
			// Check if this is error due to value object, which
			// will be a validation error.
			fmt.Println("unmarshal error", err)
		} else {
			panic(err)
		}
	}
	fmt.Println("john", john)
	fmt.Println("age valid", john.Age.Valid())

	var v *Value[string]
	fmt.Println(v.IsZero())
	fmt.Println(v.Valid())
	fmt.Println(v.Validate())
	v = new(Value[string])
	v.Set("hello")

	fmt.Println(v.Valid())
	fmt.Println(v.Validate())
	v.SetValidator(func(s string) error {
		if len(s) < 10 {
			return errors.New("too short")
		}
		return nil
	})
	fmt.Println(v.Validate())
	vv, err := v.With("hello world")
	fmt.Println(vv, err)
}

type User struct {
	Age *Age `json:"age"`
}

// Age value object.
type Age struct {
	// *Value[int] Don't use pointer, there will be issue with unmarshalling
	Value[int]
}

func (a *Age) UnmarshalJSON(raw []byte) error {
	if a == nil {
		// TODO
	}
	var v Value[int]
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	// Set back the age validator manually here.
	v.SetValidator(ValidateAge)
	a.Value = v

	// Additionally perform validation here
	// return nil
	return a.Validate()
}

var ErrInvalidAgeRange = fmt.Errorf("%w: invalid age", ErrInvalidValue)

func ValidateAge(age int) error {
	if age < 0 {
		return ErrInvalidAgeRange
	}
	return nil

}
func NewAge(age int) (*Age, error) {
	value, err := NewValue(age, WithValidator(ValidateAge))
	if err != nil {
		return nil, err
	}
	return &Age{*value}, nil
}

var ErrNotSet = errors.New("not set")
var ErrInvalidValue = errors.New("invalid value")

// Value represents a generic value object.
type Value[T any] struct {
	value     T
	dirty     bool
	validator func(T) error
}

type ValueOption[T any] func(*Value[T]) *Value[T]

func WithValidator[T any](validator func(T) error) ValueOption[T] {
	return func(v *Value[T]) *Value[T] {
		v.validator = validator
		return v
	}

}

func Must[T any](v *Value[T], err error) *Value[T] {
	if err != nil {
		panic(err)
	}
	return v
}

func NewValue[T any](t T, options ...ValueOption[T]) (*Value[T], error) {
	v := &Value[T]{
		value:     t,
		dirty:     true,
		validator: nil,
	}
	for _, opt := range options {
		opt(v)
	}
	return v, v.Validate()
}

func (v *Value[T]) IsZero() bool {
	return v == nil || !v.dirty
}
func (v *Value[T]) IsSet() bool {
	return !v.IsZero() && v.dirty
}

func (v *Value[T]) SetValidator(fn func(T) error) *Value[T] {
	v.validator = fn
	return v
}

func (v *Value[T]) With(t T) (*Value[T], error) {
	if err := v.validate(t); err != nil {
		return v, err
	}
	return NewValue[T](t, WithValidator(v.validator))
}

func (v *Value[T]) Set(t T) error {
	if err := v.validate(t); err != nil {
		return err
	}
	v.value = t
	v.dirty = true
	return nil
}

func (v *Value[T]) Get() (t T, isSet bool) {
	if !v.IsSet() {
		return
	}
	return v.value, v.dirty
}

func (v *Value[T]) MustGet() (t T) {
	if !v.IsSet() {
		panic(ErrNotSet)
	}
	return v.value
}

func (v *Value[T]) validate(t T) error {
	if validate := v.validator; validate != nil {
		return validate(t)
	}
	return nil
}

func (v *Value[T]) Validate() error {
	if v.IsZero() {
		return ErrNotSet
	}
	return v.validate(v.value)
}

func (v *Value[T]) Valid() bool {
	return v.Validate() == nil
}

func (v *Value[T]) String() string {
	if v.IsZero() {
		return "NOT SET"
	}
	return fmt.Sprint(v.value)
}

func (v Value[T]) MarshalJSON() ([]byte, error) {
	if v.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(v.value)
}

// UnmarshalJSON does not add back the validator - figure out how to add it back through reflection (NOTE: manually add in value object, see Age example).
func (v *Value[T]) UnmarshalJSON(raw []byte) error {
	if v == nil {
		return errors.New("unmarshal to nil Value[T]")
	}
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}
	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}
	*v = *Must(NewValue[T](t))
	return nil
}

```


## Generic HTTP Fetch

`server.go`:
```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var ErrNotFound = errors.New("not found")

type User struct {
	_    struct{}
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Book struct {
	_      struct{}
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error,omitempty"`
}

func NewResponse[T any](data T) *Response[T] {
	return &Response[T]{
		Data: data,
	}
}

func NewErrorResponse(err error) *Response[any] {
	return &Response[any]{
		Error: err.Error(),
	}
}

var users = []User{
	{Name: "Alice", Age: 10},
	{Name: "Bob", Age: 13},
}

var books = []Book{
	{Title: "Thinking Fast & Slow Summary", Author: "Daniel Kahneman"},
	{Title: "Influence: Science and Practice", Author: "Robert Cialdini"},
}

func main() {
	r := chi.NewRouter()
	r.Get("/users", getUsers)
	r.Get("/users/{user_id}", getUserByID)
	r.Get("/books", getBooks)
	r.Get("/books/{book_id}", getBookByID)

	fmt.Println("listening to port *:3333. press ctrl + c to cancel")
	http.ListenAndServe(":3333", r)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(NewResponse(users))
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	id, _ := strconv.Atoi(userID)
	id-- // Id starts from 1

	if id < 0 || id > len(users)-1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewErrorResponse(ErrNotFound))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(users[id]))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(NewResponse(books))
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "book_id")
	id, _ := strconv.Atoi(bookID)
	id-- // Id starts from 1

	if id < 0 || id > len(books)-1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(NewErrorResponse(ErrNotFound))
		return
	}

	json.NewEncoder(w).Encode(NewResponse(books[id]))
}
```

`client.go`:
```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var baseURL = "http://localhost:3333"

type User struct {
	_    struct{}
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Book struct {
	_      struct{}
	Title  string `json:"title"`
	Author string `json:"book"`
}

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error,omitempty"`
}

func main() {
	// Fetch users.
	fetchUsers := NewFetcher[Response[[]User], Response[any]](baseURL + "/users")
	usersResult, err := fetchUsers.Fetch(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("users data: %+v\n", usersResult.Data)
	fmt.Printf("users error: %+v\n", usersResult.Error)

	// Fetch one user.
	fetchUser := NewFetcher[Response[User], Response[any]](baseURL + "/users/%d")
	userResult, err := fetchUser.Fetch(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("users data: %+v\n", userResult.Data)
	fmt.Printf("users error: %+v\n", userResult.Error)

	// Fetch non-existing user.
	userResult, err = fetchUser.Fetch(context.Background(), -1)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("users data: %+v\n", userResult.Data)
	fmt.Printf("users error: %+v\n", userResult.Error)

	// Fetch books.
	fetchBooks := NewFetcher[Response[[]Book], Response[any]](baseURL + "/books")
	booksResult, err := fetchBooks.Fetch(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("books data: %+v\n", booksResult.Data)
	fmt.Printf("books error: %+v\n", booksResult.Error)

	// Fetch one book.
	fetchBook := NewFetcher[Response[Book], Response[any]](baseURL + "/books/%d")
	bookResult, err := fetchBook.Fetch(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("books data: %+v\n", bookResult.Data)
	fmt.Printf("books error: %+v\n", bookResult.Error)

	// Fetch non-existing book.
	bookResult, err = fetchBook.Fetch(context.Background(), -1)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("books data: %+v\n", bookResult.Data)
	fmt.Printf("books error: %+v\n", bookResult.Error)
}

type Result[T any, E any] struct {
	Data  T
	Error E
}

type Fetcher[T any, E any] struct {
	url string
}

func NewFetcher[T any, E any](url string) *Fetcher[T, E] {
	return &Fetcher[T, E]{
		url: url,
	}
}

func (f *Fetcher[T, E]) Fetch(ctx context.Context, args ...any) (*Result[T, E], error) {
	resp, err := http.Get(fmt.Sprintf(f.url, args...))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result := new(Result[T, E])

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var t T
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			return nil, err
		}
		result.Data = t
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var e E
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return nil, err
		}
		result.Error = e
	} else {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(body))
	}

	return result, nil
}
```
