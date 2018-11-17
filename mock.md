## Simple builder

```go
package main

import (
	"fmt"
)

type Token struct {
	A, B, C string
}

type TokenModifier func(t *Token)

type TokenBuilder struct {
	defaults Token
	override TokenModifier
}

func (t *TokenBuilder) Build(modifiers ...TokenModifier) *Token {
	result := t.defaults
	for _, mod := range modifiers {
		mod(&result)
	}
	if t.override != nil {
		t.override(&result)
	}
	return &result
}

func main() {
	tb := new(TokenBuilder)
	tb.override = func(t *Token) {
		// Test the original implementation.
		fmt.Println(t.A, t.B)
		// Then mock them for testing.
		t.A = "override a"
		t.B = "override b"
	}
	token := tb.Build(
		// Update single params. Can be extracted as a function.
		func(t *Token) { t.A = "a" },
		func(t *Token) { t.B = "b" },
		func(t *Token) {
			// Update multiple params at the same time.
			t.A = "aa"
			t.B = "bb"
			t.C = "cc"
		},
	)
	// How about errors? We can simply return the error in the modifier if needed.
	fmt.Println("Hello, playground", token)
}
```
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

## Alternative with Builder 

Sample builder with readOnly options.

```go
package main

import (
	"fmt"
	"time"
)

type Token struct {
	CreatedAt time.Time
	Secret    string
}

func NewToken() *Token {
	return &Token{
		CreatedAt: time.Now(),
		Secret:    "randomly_generated_value",
	}
}

type TokenBuilder struct {
	defaults Token
	override func(t *Token)
	readOnly bool // writable false
}

type TokenBuilderOptions func(t *TokenBuilder)

func ReadOnly(readOnly bool) TokenBuilderOptions {
	return func(t *TokenBuilder) {
		t.readOnly = readOnly
	}
}

func Defaults(defaults Token) TokenBuilderOptions {
	return func(t *TokenBuilder) {
		t.defaults = defaults
	}
}

func NewTokenBuilder(opts ...TokenBuilderOptions) *TokenBuilder {
	tb := TokenBuilder{}
	for _, o := range opts {
		o(&tb)
	}
	return &tb
}

func (t *TokenBuilder) SetSecret(secret string) {
	if t.readOnly {
		return
	}
	t.defaults.Secret = secret
}

func (t *TokenBuilder) SetCreatedAt(createdAt time.Time) {
	if t.readOnly {
		return
	}
	t.defaults.CreatedAt = createdAt
}

func (t *TokenBuilder) Build() *Token {
	result := t.defaults
	if t.override != nil && !t.readOnly {
		t.override(&result)
	}
	return &result
}

func (t *TokenBuilder) SetOverride(override func(t *Token)) {
	t.override = override
}

func GenerateToken(tb *TokenBuilder) *Token {
	tb.SetSecret("random_secret")
	tb.SetCreatedAt(time.Now())
	return tb.Build()
}

func main() {
	tb := NewTokenBuilder()
	token := tb.Build()
	fmt.Printf("build empty: %#v\n", token)

	tb.SetSecret("secret")
	tb.SetCreatedAt(time.Now())
	token = tb.Build()
	fmt.Printf("build after setting: %#v\n", token)

	tb.SetOverride(func(t *Token) {
		t.Secret = "overwrite secret"
		t.CreatedAt = time.Unix(0, 0)
	})
	token = tb.Build()
	fmt.Printf("build override: %#v\n", token)

	tb2 := NewTokenBuilder(Defaults(Token{CreatedAt: time.Now()}))
	token = GenerateToken(tb2)
	fmt.Printf("build override: %#v\n", token)
	tb2.SetOverride(func(t *Token) {
		// This is useful, because you can test the randomly generated values, as well as mocking them to your desired result.
		t.Secret = "overwritten: " + t.Secret
		t.CreatedAt = time.Unix(0, 0)
	})
	token = GenerateToken(tb2)
	fmt.Printf("build override: %#v\n", token)

	tb3 := NewTokenBuilder(
		Defaults(Token{Secret: "immutable_secret", CreatedAt: time.Now()}),
		ReadOnly(true),
	)
	tb3.SetSecret("can't overwrite immutable_secret")
	token = tb3.Build()
	fmt.Printf("build readOnly: %#v\n", token)
}
```
