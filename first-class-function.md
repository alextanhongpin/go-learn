```go
package main

import (
	"fmt"
	"sync"
)

type User struct {
	Name string
}

type Hub struct {
	users      map[string]User
	register   chan User
	unregister chan User
}

func NewHub() (*Hub, func()) {
	hub := &Hub{
		users:      make(map[string]User),
		register:   make(chan User),
		unregister: make(chan User),
	}
	stop := hub.loop()
	return hub, stop
}

func (h *Hub) Register(u User) {
	h.register <- u
}

func (h *Hub) Unregister(u User) {
	h.unregister <- u
}

func (h *Hub) loop() func() {
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan interface{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case u, ok := <-h.register:
				if !ok {
					return
				}
				h.users[u.Name] = u
			case u, ok := <-h.unregister:
				if !ok {
					return
				}
				delete(h.users, u.Name)
			}
		}
	}()
	return func() {
		close(done)
		wg.Wait()
	}
}

type Hub2 struct {
	users map[string]User
	ch    chan func(map[string]User)
}

func NewHub2() (*Hub2, func()) {
	hub := &Hub2{
		users: make(map[string]User),
		ch:    make(chan func(map[string]User)),
	}
	stop := hub.loop()
	return hub, stop
}

func (h *Hub2) Register(u User) {
	h.ch <- func(users map[string]User) {
		users[u.Name] = u
	}
}

func (h *Hub2) Unregister(u User) {
	h.ch <- func(users map[string]User) {
		delete(users, u.Name)
	}
}

func (h *Hub2) loop() func() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for o := range h.ch {
			o(h.users)
		}
	}()
	return func() {
		close(h.ch)
		wg.Wait()
	}
}

func main() {
	u := User{
		Name: "john",
	}
	{
		h, stop := NewHub()
		defer func() {
			stop()
			fmt.Println(h.users)
		}()
		h.Register(u)
		fmt.Println(h.users)
		h.Unregister(u)
	}

	{
		h, stop := NewHub2()
		defer func() {
			stop()
			fmt.Println(h.users)
		}()
		h.Register(u)
		fmt.Println(h.users)
		h.Unregister(u)

	}
	fmt.Println("Hello, playground")
}
```
