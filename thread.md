# Thread-like clojure pattern? 
```go
package main

import (
	"fmt"
)

type User struct {
	Name string
}

func main() {
	fetch := func() (*User, error) {
		return &User{Name: "john"}, nil
	}
	print := func(u *User) error {
		fmt.Println(u)
		return nil
	}
	setName := func(u *User) error {
		u.Name = "jane"
		return nil
	}
	t := &Thread{}
	u, err := t.Do(fetch, print, setName)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}

type ThreadArg0 func() (*User, error)
type ThreadArgN func(u *User) error
type Thread struct{}

func (t *Thread) Do(arg ThreadArg0, args ...ThreadArgN) (*User, error) {
	u, err := arg()
	if err != nil {
		return nil, err
	}
	for _, a := range args {
		if err := a(u); err != nil {
			return nil, err
		}
	}
	return u, nil
}
```

## Mutation with rollback

```go
package main

import (
	"errors"
	"fmt"
)

type User struct {
	Name    string
	Hobbies []string
}

func (u *User) Mutate(fn func(u *User) error) error {
	c := *u
	if err := fn(&c); err != nil {
		return err
	}
	*u = c
	return nil
}

func main() {
	u := &User{}
	err := u.Mutate(func(u *User) error {
		u.Name = "john"
		u.Hobbies = append(u.Hobbies, "cycling")
		return errors.New("hello")
	})
	if err != nil {
		fmt.Println(u)
		panic(err)
	}
	fmt.Println(u)
}
```
