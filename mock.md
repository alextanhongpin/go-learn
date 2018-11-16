## Mock

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

type Token struct {
	ExpiresIn int64
	Value     string
	ExpireAt  int64
}

type TokenModifier func(t *Token) error

// Concrete, will create a new token with *RANDOM* value. Hard to mock.
func NewToken(expiresIn int64) (*Token, error) {
	v, err := randomString(32)
	if err != nil {
		return nil, err
	}
	return &Token{
		Value:     v,
		ExpiresIn: expiresIn,
		ExpireAt:  time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}, nil
}

type TokenFactory interface {
	Build() (*Token, error)
}

type tokenFactory struct {
	// Takes a default, and populate the random values.
	// Can also take a pointer to an existing data, and RESETS the values to defaults.
	defaults  Token
	modifiers []TokenModifier
}

func NewTokenFactory(defaults Token, modifiers ...TokenModifier) *tokenFactory {
	return &tokenFactory{defaults, modifiers}
}

func (t *tokenFactory) Build(extras ...TokenModifier) (*Token, error) {
	var err error
	token := t.defaults
	for _, mod := range append(t.modifiers, extras...) {
		err = mod(&token)
		if err != nil {
			return nil, err
		}
	}
	return &token, err
}

type Service interface {
	CreateToken() (*Token, error)
}

type service struct {
	tokenFactory        *tokenFactory
	higherOrderModifier func(time.Time, int64) TokenModifier
}

func setExpiredAt(tm time.Time, duration int64) TokenModifier {
	return func(t *Token) error {
		t.ExpireAt = tm.Add(time.Duration(duration) * time.Second).Unix()
		return nil
	}
}

func setTokenValue() TokenModifier {
	return func(t *Token) error {
		var err error
		t.Value, err = randomString(32)
		return err
	}
}

func NewService(tf *tokenFactory, setExpiredAt func(tm time.Time, duration int64) TokenModifier) *service {
	return &service{
		higherOrderModifier: setExpiredAt,
		tokenFactory:        tf,
	}
}

func (s *service) CreateToken(now time.Time) (*Token, error) {
	token, err := s.tokenFactory.Build(s.higherOrderModifier(now, 3600))
	return token, err
}

func main() {
	// Create a factory for a token factory.
	defaults := Token{ExpiresIn: 3600}
	modifiers := []TokenModifier{setTokenValue()}
	tb := NewTokenFactory(defaults, modifiers...)
	token, err := tb.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ori:", token)
	// What if I need to use the input for building the response?

	// Create a mock token builder - which will just return back anything that we set.
	// We explicitly say return this value that I set.
	// Somehow, in order to mock the data, we have to pass in the data we want to mock to the function.
	// If we want to pass the pass the data, we need to mock the data externally.
	defaults = Token{
		Value:     "xyz",
		ExpiresIn: 3600,
	}
	mtb := NewTokenFactory(defaults)
	token, err = mtb.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mock:", token)

	// Create new service
	svc := NewService(mtb, setExpiredAt)
	token, err = svc.CreateToken(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	// mock the set expired at
	mockSetExpiredAt := func(tm time.Time, duration int64) TokenModifier {
		return func(t *Token) error {
			t.ExpireAt = 0
			return nil
		}
	}
	// Create new service
	svc = NewService(mtb, mockSetExpiredAt)
	token, err = svc.CreateToken(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func randomString(n int) (string, error) {
	b, err := randomBytes(n)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
```
