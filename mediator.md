```go
package main

import (
	"fmt"
	"time"
)

// Few ways to for Service B to call Service A.
// Embed Service A in Service B.
// Put the services together in a struct. Might not work well if they represent different model.
// Create another struct that hold both Service A and B. Compose the function call together. E.g. Employee + Organization service becomes OrganizationEmployees. 
// This might not work for infra, e.g. Login + Email service.
// This e.g. Service B sends the request to a channel, which is processed and calls service A and the result returned in Service B.
type Mediator chan MediatorRequest

type MediatorEvent string
type MediatorRequest struct {
	Event MediatorEvent
	Ch    chan interface{}
}

const (
	ServiceAWork        = MediatorEvent("ServiceA:Work")
	ServiceAAnotherWork = MediatorEvent("ServiceA:AnotherWork")
)

type ServiceA struct{}

func (a *ServiceA) Work() string {
	fmt.Println("service A work")
	return "done by service A"
}

type ServiceB struct {
	Mediator
}

func (b *ServiceB) Work(i int) {
	// ch := mediate(b.Mediator, ServiceAAnotherWork)
	ch := mediate(b.Mediator, ServiceAWork)
	defer close(ch)
	select {
	case res := <-ch:
		fmt.Println("service B:", res, i)
	case <-time.After(1 * time.Second):
		fmt.Println("awaited for 1 second")
		return
	}
}

func mediate(mediator Mediator, event MediatorEvent) chan interface{} {
	ch := make(chan interface{})
	go func() {
		req := MediatorRequest{
			Event: event,
			Ch:    ch,
		}
		select {
		case mediator <- req:
		case <-time.After(1 * time.Second):
			fmt.Println("no work done")
			return
		}
	}()
	return ch

}

func main() {
	mediator := make(chan MediatorRequest)
	svcA := new(ServiceA)
	svcB := &ServiceB{mediator}
	done := make(chan interface{})
	defer close(done)

	go func() {
		for req := range mediator {
			select {
			case <-done:
				return
			default:
			}
			switch req.Event {
			case ServiceAWork:
				res := svcA.Work()
				select {
				case <-done:
					return
				case req.Ch <- res:
				case <-time.After(1 * time.Second):
					return
				}
			case ServiceAAnotherWork:
				fmt.Println("got another work, but not processing")
				continue
			default:
				fmt.Println(req, "does not exist")
				continue
			}
		}
	}()

	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		svcB.Work(i)
	}

	fmt.Println("Hello, playground")
}

```
