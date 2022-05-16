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
