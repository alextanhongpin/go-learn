Simplify domain model validation:

```go
package main

import (
	"fmt"
	"log"
	"strings"
)

type User struct {
	FirstName string
	LastName  string
}

func New(firstName, lastName string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
	}
}

func saveUser(user User) error {
	if err := validate(user.FirstName, "firstName"); err != nil {
		return err
	}
	if err := validate(user.LastName, "lastName"); err != nil {
		return err
	}
	fmt.Println("save user")
	return nil
}

func validate(val, field string) error {
	if strings.TrimSpace(val) == "" {
		return fmt.Errorf("%q is required", field)
	}
	return nil
}

func main() {
	u := New("", "doe")
	if err := saveUser(*u); err != nil {
		log.Fatal(err)
	}
}
```
