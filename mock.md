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

// TokenFactory takes an input/request/params that is needed to produce an output.
type TokenFactory interface {
	Build(...TokenModifier) (*Token, error)
}

type tokenFactory struct {
	// Takes a default, and populate the random values.
	// Can also take a pointer to an existing data, and RESETS the values to defaults.
	defaults Token

	// A list of initial modifiers that does not need arguments.
	modifiers []TokenModifier

	// The final modifier that overwrites everything. Useful for mocking.
	override TokenModifier
}

func (t *tokenFactory) SetOverride(override TokenModifier) {
	t.override = override
}

func NewTokenFactory(defaults Token, modifiers ...TokenModifier) *tokenFactory {
	return &tokenFactory{
		defaults:  defaults,
		modifiers: modifiers,
	}
}

func (t *tokenFactory) Build(extras ...TokenModifier) (*Token, error) {
	var err error
	// Make a copy.
	token := t.defaults
	// For each of them, apply the modifier. This works if you do not need to set the values in the correct order.
	// If you need that, defined each of them as a pipeline steps.
	for _, mod := range append(t.modifiers, extras...) {
		err = mod(&token)
		if err != nil {
			return nil, err
		}
	}
	// Overrides everything, this makes testing easier.
	if t.override != nil {
		t.override(&token)
	}
	return &token, err
}

type Service interface {
	// No matter what, there are mutation values that needs to be passed from outside the function to mock.
	CreateToken(now time.Time) (*Token, error)
}

type service struct {
	tokenFactory *tokenFactory
}

// A unit of work.
func makeExpireAtModifier(tm time.Time) TokenModifier {
	return func(t *Token) error {
		t.ExpireAt = tm.Add(time.Duration(t.ExpiresIn) * time.Second).Unix()
		return nil
	}
}

func valueModifier(t *Token) error {
	var err error
	t.Value, err = randomString(32)
	return err
}

func NewService(tf *tokenFactory) *service {
	return &service{
		tokenFactory: tf,
	}
}

func (s *service) CreateToken(now time.Time) (*Token, error) {
	// If we need to pass down modifiers with argument, do it in the function to mock.
	return s.tokenFactory.Build(makeExpireAtModifier(now))
}

func main() {
	// Create a factory for a token factory.
	defaults := Token{ExpiresIn: 3600}
	modifiers := []TokenModifier{valueModifier}
	tf := NewTokenFactory(defaults, modifiers...)
	token, err := tf.Build(makeExpireAtModifier(time.Now()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ori: %+v\n", token)
	// What if I need to use the input for building the response?

	// Create a mock token builder - which will just return back anything that we set.
	// We explicitly say return this value that I set.
	// Somehow, in order to mock the data, we have to pass in the data we want to mock to the function.
	// If we want to pass the pass the data, we need to mock the data externally.
	defaults = Token{
		Value:     "xyz",
		ExpiresIn: 3600,
		ExpireAt:  -100,
	}
	tf = NewTokenFactory(defaults)
	token, err = tf.Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("no modifiers: %+v\n", token)

	// Create new service
	svc := NewService(tf)
	token, err = svc.CreateToken(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("service: %+v\n", token)

	// Set an override for the token factory.
	tf.SetOverride(func(t *Token) error {
		t.ExpiresIn = 0
		t.ExpireAt = 0
		t.Value = "mocked"
		return nil
	})

	token, err = svc.CreateToken(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("override: %+v\n", token)

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
