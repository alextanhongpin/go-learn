Implement event listener in golang.

```go
package main

import (
	"fmt"
	"time"
)

// Dog has a name and a list of people caring for him and watching
// what he does
type Dog struct {
	name    string
	sitters map[string][]chan string
}

// AddSitter adds an event listener to the Dog struct instance
func (d *Dog) AddSitter(e string, ch chan string) {
	if d.sitters == nil {
		d.sitters = make(map[string][]chan string)
	}
	if _, ok := d.sitters[e]; ok {
		d.sitters[e] = append(d.sitters[e], ch)
	} else {
		d.sitters[e] = []chan string{ch}
	}
}

// RemoveSitter removes an event listener from the Dog struct instance
func (d *Dog) RemoveSitter(e string, ch chan string) {
	if _, ok := d.sitters[e]; ok {
		for i := range d.sitters[e] {
			if d.sitters[e][i] == ch {
				d.sitters[e] = append(d.sitters[e][:i], d.sitters[e][i+1:]...)
				break
			}
		}
	}
}

// Emit emits an event on the Dog struct instance
func (d *Dog) Emit(e, response string) {
	if _, ok := d.sitters[e]; ok {
		for _, handler := range d.sitters[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}

func main() {
	doggie := Dog{"Wurf", nil}

	o := make(chan string)

	doggie.AddSitter("bark", o)
	doggie.AddSitter("poop", o)
	doggie.AddSitter("hungry", o)

	go func() {
		for {
			msg := <-o
			fmt.Println("Wurf:", msg)
		}
	}()

	fmt.Println("The dog barked")
	doggie.Emit("bark", "Told not to bark!")

	fmt.Println("The dog pooped!")
	doggie.Emit("poop", "Picked up poop!")

	fmt.Println("The dog is hungry")
	doggie.Emit("hungry", "Feed the dog!")

	time.Sleep(3 * time.Second)

	doggie.RemoveSitter("poop", o)

	// Hired a dog sitter to pick up poop
	dogsitter := make(chan string)
	doggie.AddSitter("poop", dogsitter)
	fmt.Println("Hired a dogsitter to pick up poop")

	go func() {
		for {
			msg := <-dogsitter
			fmt.Println("Dogsitter:", msg)
		}
	}()

	fmt.Println("Dog barked!")
	doggie.Emit("bark", "Told not to bark!")

	fmt.Println("Dog has pooped!")
	doggie.Emit("poop", "Picked up poop!")

	fmt.Println("Dog has pooped again!")
	doggie.Emit("poop", "Picked up poop, again!")

	fmt.Println("The dog is hungry!")
	doggie.Emit("hungry", "Feed the dog!")

	fmt.Scanln()
}

```
