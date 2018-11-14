## Functional optional in golang

```go
package main

import (
	"fmt"
)

type API struct {
	Name string
	Age  int
}

type Option func(*API)

func New(options ...Option) *API {
	opts := new(API)
	for _, o := range options {
		o(opts)
	}
	return opts
}

func Name(n string) Option {
	return func(o *API) {
		o.Name = n
	}
}

func Age(a int) Option {
	return func(o *API) {
		o.Age = a
	}
}

func main() {
	opts := New(Name("john"), Age(100))
	fmt.Println(opts)
}
```


## Using functional optional to decouple dependencies
```go
package main

import (
	"fmt"
)

func main() {
	newClient := new(Client)
	service := new(ClientService)
	service.generator = func(s string) (string, error) {
		return "hello", nil
	}
	service.Register(newClient)
	fmt.Println(newClient)
}

type Client struct {
	ClientID     string
	ClientSecret string
}

type ClientMutator func(c *Client) error

func ClientID(generator func(id string) (string, error)) ClientMutator {
	return func(c *Client) error {
		clientID, err := generator(c.ClientID)
		if err != nil {
			return err
		}
		c.ClientID = clientID
		return nil
	}
}

func ClientSecret(secret string) ClientMutator {
	return func(c *Client) error {
		c.ClientSecret = secret
		return nil
	}
}

type ClientService struct {
	generator func(string) (string, error)
}

func (c *ClientService) Register(client *Client) *Client {
	apply(client, ClientID(c.generator), ClientSecret("car"))
	return client
}

func apply(client *Client, ops ...ClientMutator) error {
	var err error
	for _, o := range ops {
		err = o(client)
		if err != nil {
			return err
		}
	}
	return err
}
```

## Another alternative
```go
package main

import (
	"fmt"
)

type Modifier func(u *User) error

type User struct {
	A, B, C string
}

type UserModifier struct {
	User *User
	
}

func (u *UserModifier) Apply(modifiers ...Modifier) error {
	for _, m := range modifiers {
		if err := m(u.User); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	aSetter := func(u *User) error {
		u.A = "hello a"
		return nil
	}
	u := new(User)
	m := UserModifier{u}
	m.Apply(aSetter)
	fmt.Println("Hello, playground", m.User)
}
```
