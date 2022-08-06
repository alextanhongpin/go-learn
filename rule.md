# Rule-based lazy exec

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

func main() {
	r := RegisterRule{
		Email:    "john.doe@mail.com",
		Password: "12345678",
	}
	r.ValidateRequest()
	r.EncryptPassword()
	fmt.Println(r.Exec(), r)
}

type Rule struct {
	execs []func() error
}

func (r *Rule) Add(exec func() error) {
	r.execs = append(r.execs, exec)
}

func (r *Rule) Exec() error {
	for _, exec := range r.execs {
		if err := exec(); err != nil {
			return err
		}
	}

	return nil
}

type RegisterRule struct {
	Rule

	Email    string
	Password string

	EncryptedPassword string
}

func (r *RegisterRule) ValidateRequest() {
	r.Add(func() error {
		if r.Email == "" {
			return errors.New("email is required")
		}
		if r.Password == "" {
			return errors.New("password is required")
		}

		return nil
	})
}

func (r *RegisterRule) EncryptPassword() {
	r.Add(func() error {
		r.EncryptedPassword = fmt.Sprintf("enc:%s", r.Password)

		return nil
	})
}
```
