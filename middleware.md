Simple chainable handler:

```go
package handle

import (
	"context"
)

/*
	Usage:

	func main() {
		handler := Chain[string](&StringHandler{}, UppercaseHandler(), SplitHandler())
		fmt.Println(handler.Handle(context.Background(), "hello world"))
	}

type StringHandler struct{}

	func (h *StringHandler) Handle(ctx context.Context, s string) error {
		fmt.Println("string handler", s)
		return nil
	}

	func UppercaseHandler() Adapter[string] {
		return func(next Handler[string]) Handler[string] {
			return HandlerFunc[string](func(ctx context.Context, s string) error {
				fmt.Println("before:uppercase", s)
				s = strings.ToUpper(s)
				if err := next.Handle(ctx, s); err != nil {
					return err
				}
				fmt.Println("after:uppercase", s)
				return nil
			})
		}
	}

	func SplitHandler() Adapter[string] {
		return func(next Handler[string]) Handler[string] {
			return HandlerFunc[string](func(ctx context.Context, s string) error {
				fmt.Println("before:split", s)
				s = strings.Split(s, " ")[0]
				if err := next.Handle(ctx, s); err != nil {
					return err
				}
				fmt.Println("after:split", s)
				return nil
			})
		}
	}
*/
type Handler[T any] interface {
	Handle(ctx context.Context, t T) error
}

type HandlerFunc[T any] func(ctx context.Context, t T) error

func (fn HandlerFunc[T]) Handle(ctx context.Context, t T) error {
	return fn(ctx, t)
}

type Adapter[T any] func(next Handler[T]) Handler[T]

func Chain[T any](handler Handler[T], handlers ...Adapter[T]) Handler[T] {
	head := handler
	for i := len(handlers) - 1; i > -1; i-- {
		h := handlers[i]
		head = h(head)
	}
	return head
}
```

## HTTP Middleware

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func main() {
	ts := httptest.NewServer(Guard("two", Guard("one", http.HandlerFunc(do))))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	fmt.Println("Hello, 世界")
}

func do(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "do")
}

func Guard(msg string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before", msg)
		next.ServeHTTP(w, r)
		fmt.Println("after")
	})
}
```

Output:
```
before two
before one
after
after
do
Hello, 世界
```
