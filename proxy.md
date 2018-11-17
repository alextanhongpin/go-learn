```go
package main

import (
	"fmt"
)

type SecretSigner interface {
	Sign() string
}

type Signer struct {
	secret string
}

func (s *Signer) Sign() string {
	return "A"
}

type signerProxy struct {
	proxy SecretSigner
}

func (s *signerProxy) Sign() string {
	// s.proxy.Sign() // returns A
	return "B"
}

func SignMethod(s SecretSigner) string {
	return s.Sign()
}

func main() {
	s := new(Signer)
	fmt.Println("sign:", SignMethod(s))
	fmt.Println("sign:", SignMethod(&signerProxy{s}))
}
```
