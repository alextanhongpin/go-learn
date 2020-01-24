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
