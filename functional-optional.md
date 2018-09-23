## Functional optional in golang

```go
package main

import (
	"fmt"
)

type API struct {
	Name string
	Age  int
}

type Option func(*API)

func New(options ...Option) *API {
	opts := new(API)
	for _, o := range options {
		o(opts)
	}
	return opts
}

func Name(n string) Option {
	return func(o *API) {
		o.Name = n
	}
}

func Age(a int) Option {
	return func(o *API) {
		o.Age = a
	}
}

func main() {
	opts := New(Name("john"), Age(100))
	fmt.Println(opts)
}
```
