## Shadowing fields

This program demonstrates how to exclude fields in the json output
```go
package main

import (
	"encoding/json"
	"log"
)

type UserPrivate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPublic struct {
	*UserPrivate
	Password bool `json:"password,omitempty"`
}

func main() {

	usrPriv := UserPrivate{"john.doe@mail.com", "123456"}
	usrPub := UserPublic{
		UserPrivate: &usrPriv,
	}
	// Convert it to bytes
	out, err := json.Marshal(usrPub)
	if err != nil {
		log.Println(err)
	}
	log.Printf("with shadowing: %s\n", string(out))
}
```

## Shadowing Dates

If the dates are not set, golang will return `0001-01-01` to the client side as the json response. To avoid that:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var MaxDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.Local)

type User struct {
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (u *User) IsDateOfBirthValid() bool {
	return !u.DateOfBirth.Equal(MaxDate)
}

type Response struct {
	User
	DateOfBirth *time.Time `json:"date_of_birth"`
}

func main() {
	user := User{
		DateOfBirth: time.Date(9999, 12, 31, 0, 0, 0, 0, time.Local),
	}
	res := Response{
		User: user,
	}
	if res.User.IsDateOfBirthValid() {
		res.DateOfBirth = &res.User.DateOfBirth
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
```
