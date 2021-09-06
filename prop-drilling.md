# Using interface to whitelist fields across different layers
```go
package main

import (
	"fmt"
)

func main() {
	api := UserAPI{
		svc: &UserService{
			repo: &UserRepository{},
		},
	}
	api.Update("john")
	fmt.Println("Hello, playground")
}

type User interface {
	Name() string
}

type UserAPI struct {
	svc *UserService
}

type userDto struct {
	name string
}

func (u userDto) Name() string {
	return u.name
}

func (a *UserAPI) Update(name string) {
	var u userDto
	u.name = name
	fmt.Println("api", u)
	if err := a.svc.Update(u); err != nil {
		panic(err)
	}
}

type UserService struct {
	repo *UserRepository
}

func (svc *UserService) Update(u User) error {
	fmt.Println("service", u.Name())
	return svc.repo.Update(u)
}

type UserRepository struct{}

func (r *UserRepository) Update(u User) error {
	fmt.Println("repo", u.Name())
	return nil
}
```
