# Shows how to copy the value while maintaining the pointer address
```go
package main

import (
	"fmt"
)

func main() {
	req := new(Request)
	fmt.Printf("req: %p %v\n", req, req)
	a := Request{Name: "Car"}
	req.Replace(&a)
	fmt.Printf("a: %p %v\n", &a, a)
	fmt.Printf("req: %p %v\n", req, req)

	b := &Request{Name: "Paper"}
	req.Replace(b)
	fmt.Printf("b: %p %v\n", b, b)
	fmt.Printf("req: %p %v\n", req, req)
	req.Name = "hello"

	fmt.Printf("a: %p %v\n", &a, a)
	fmt.Printf("b: %p %v\n", b, b)
	fmt.Printf("req: %p %v\n", req, req)
}

type Request struct {
	Name string
}

func (r *Request) Replace(rr *Request) {
	*r = *rr
}
```
