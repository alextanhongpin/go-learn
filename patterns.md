# Factory pattern 

For creating user with different scenario. Delegating business logic to factory.

```go
package main

import (
	"log"
)

func main() {
	factory := &UserFactory{}
	factory.WithEmail("john.doe@mail.com").
		WithEncryptedPassword("hello world")
	user, err := factory.Build()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}

type User struct {
	EncryptedPassword string
	Email             string
}

func encrypt(plaintext string) (string, error) {
	return "encrypted:" + plaintext, nil
}

type UserFactory struct {
	user User
	err  error
}

func (u *UserFactory) WithEmail(email string) *UserFactory {
	if u.err != nil {
		return u
	}
	u.user.Email = email
	return u
}

func (u *UserFactory) WithEncryptedPassword(password string) *UserFactory {
	if u.err != nil {
		return u
	}
	u.user.EncryptedPassword, u.err = encrypt(password)
	return u
}

func (u *UserFactory) Build() (User, error) {
	return u.user, u.err
}
```

## Strategy vs Factory?


One of the most common misconception is that we need to provide different strategy for a given functionality (hashing algorithm etc). Strategy pattern is meant to change the algorithm’s behaviour during runtime. If we only need a single implementation, factory pattern is the way to go when initialising a class/function with a given behaviour.

```go
package main

import (
	"fmt"
)

type UserService struct {
	pwdMgr PasswordManager
}

type Strategy int

const (
	Argon2 Strategy = iota
	Bcrypt
)

func main() {
	fmt.Println("Hello, playground")
}

// PasswordManager represents the operation for the password.
type PasswordManager interface {
	Encrypt(password string) (string, error)
	Compare(password, encryted string) (bool, error)
}

type Argon2Strategy struct{}
type BcryptStrategy struct{}

// Usage:
// mgr := NewPasswordManager(Argon2Strategy{})
// encrypted, err := mgr.Encrypt(plainText)
// match, err := mgr.Compare(plainText, encrypted)

func NewUserService(strategy Strategy) *UserService {
	switch Strategy {
	case Argon2:
		return &UserService{Argon2Strategy{}}
	case Bcrypt:
		return &UserService{BcryptStrategy{}}
	default:
	}
}
```

## Delegating to Models

```go
package main

import (
	"fmt"
)

type ChangePassword func(email, password string) (bool, error)

type User struct {
	Password string
}

func (u *User) UpdatePassword(oldPassword, newPassword string) (bool, error) {
	encrypted, err := encrypt(oldPassword)
	if err != nil {
		return nil, err
	}
	if !encrypted.Match(u.Password) {
		return false, nil
	}
	u.Password = encrypt(newPassword).Value()
	return true, nil
}

func (u *User) ResetPassword(password string) {

}

func Save(repo Repository, u *User) (bool, error) {
	return repo.Save(u)
}

type encryptedPassword struct {
	password string
}

func (e *encryptedPassword) Match(password string) bool {
	return constantTimeCompare(e.password, password)
}

func (e *encryptedPassword) Value() string {
	return e.password
}

func encrypt(password string) (*encryptedPassword, error) {
	// Some encryption logic.
	return &encryptedPassword(password), nil
}

func main() {
	fmt.Println("Hello, playground")
}
```

## Alternative

A better approach is to delegate the password to a model:

```go
package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type EncryptedPassword interface {
	Compare(password string) error
}

type User struct {
	Password EncryptedPassword
}

type BcryptPassword string

func (b BcryptPassword) Compare(plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(b), []byte(plainText))
}

// BcryptPasswordFactory
func NewBcryptPassword(plainText string) (BcryptPassword, error) {
	// NOTE: Use higher MinCost.
	cipher, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.MinCost)
	return BcryptPassword(string(cipher)), err
}

func main() {
	password := "hello world"
	encrypted, err := NewBcryptPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encrypted)
	if err := encrypted.Compare(password); err != nil {
		log.Fatal(err)
	}
	if err := encrypted.Compare("wrong password"); err != nil {
		log.Fatal(err)
	}
}
```

## Using Types vs Primitives Values

In the example above, we use a type to define encrypted password (or an interface, because there could be several polymorphic types of password encryption, e.g. using bcrypt or argon2id). But this presents some additional issue, such as storing the values to the database. We now need to implement `values` and `scanner` for the types in order to store them. An alternative is to create a type, but not using the type directly. Instead, create a method for that type that will return the primitive value directly.

```go
package main

type Argon2Password struct {
}

type BcryptPassword struct {
	value string
}

func NewBcryptPassword(password string) *BcryptPassword {
	// Your bcrypt encryption method here.
	encryptBcrypt := func(value string) string {
		return value
	}
	return &BcryptPassword{value: encryptBcrypt(password)}
}

func (b *BcryptPassword) Value() string {
	return b.value
}

func (b *BcryptPassword) Compare(password string) bool {
	// Compare your password here.
	return false
}

type User struct {
	encryptedPassword string
}

func main() {
	var u User
	u.encryptedPassword = NewBcryptPassword("hello world").Value()
}
```
