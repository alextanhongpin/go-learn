```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Option func(interface{}) error

type Builder func(ctx context.Context, optFn Option) (interface{}, error)

type Factory struct {
	builder  map[string]Builder
	entities map[string]interface{}
}

func NewFactory() *Factory {
	return &Factory{
		builder:  make(map[string]Builder),
		entities: make(map[string]interface{}),
	}
}

func (f *Factory) Register(entity string, builder Builder) bool {
	_, ok := f.builder[entity]
	if ok {
		return false
	}
	f.builder[entity] = builder
	return true
}
func (f *Factory) Get(entity string) (interface{}, bool) {
	e, ok := f.entities[entity]
	return e, ok
}

func (f *Factory) Build(ctx context.Context, entity string, optFn Option) (interface{}, error) {
	ent, ok := f.Get(entity)
	if ok {
		return ent, nil
	}
	builder, ok := f.builder[entity]
	if !ok {
		return nil, errors.New("entity not registered")
	}
	obj, err := builder(ctx, optFn)
	if err != nil {
		return nil, err
	}
	f.entities[entity] = obj
	return obj, nil
}

type User struct {
	ID   string
	Name string
}

type Photo struct {
	ID     string
	UserID string
}

func main() {
	f := NewFactory()
	f.Register("user", func(ctx context.Context, fnOpt Option) (interface{}, error) {
		u := User{
			ID:   "user-id",
			Name: "john",
		}
		err := fnOpt(&u)
		if err != nil {
			return nil, err
		}
		return u, nil
	})
	f.Register("photo", func(ctx context.Context, fnOpt Option) (interface{}, error) {
		u, err := f.Build(ctx, "user", fnOpt)
		if err != nil {
			return nil, err
		}
		user, ok := u.(User)
		if !ok {
			return nil, errors.New("invalid user")
		}
		p := Photo{
			ID:     "photo-id",
			UserID: user.ID,
		}
		if err := fnOpt(&p); err != nil {
			return nil, err
		}
		// Save to DB.
		return p, nil
	})
	ctx := context.Background()
	p, err := f.Build(ctx, "photo", func(in interface{}) error {
		switch u := in.(type) {
		case *User:
			u.ID = "user-2"
		case *Photo:
			u.ID = "photo-2"
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	photo := p.(Photo)
	fmt.Println("got photo", photo)
	fmt.Println(f.Get("user"))
	fmt.Println(f.Get("photo"))
}
```
