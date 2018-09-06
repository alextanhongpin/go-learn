Sample send/receive pattern using channels. As long as the struct shares the same channel, they could be decoupled and share the same instance of the data.
```go
package main

import (
	"fmt"
	"log"
	"time"
)

type Receiver struct {
	count uint
	ch    <-chan int
	quit  chan struct{}
}

func (r *Receiver) Start() {
	go func() {
		for {
			select {
			case v, ok := <-r.ch:
				if !ok {
					log.Println("not ok", v)
					return
				}
				r.count++
				log.Println("receive:", v)
			case <-r.quit:
				return
			}
		}
	}()
}

func (r *Receiver) Stop() {
	close(r.quit)
	log.Println("receiver stop")
}

func NewReceiver(ch <-chan int) *Receiver {
	return &Receiver{
		ch:   ch,
		quit: make(chan struct{}),
	}
}

type Sender struct {
	ch    chan<- int
	count uint
}

func (s *Sender) Send(i int) {
	select {
	case s.ch <- i:
		s.count++
		log.Println("send:", i)
	}
}

func NewSender(ch chan<- int) *Sender {
	return &Sender{
		ch: ch,
	}
}

func main() {

	ch := make(chan int, 2)
	rcv := NewReceiver(ch)
	snd := NewSender(ch)

	rcv.Start()
	// defer rcv.Stop()

	snd.Send(1)
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			go snd.Send(i)
		}

		defer rcv.Stop()

		time.Sleep(2 * time.Second)
		for i := 0; i < 10; i++ {
			go snd.Send(i)
		}
	}()

	fmt.Scanln()
	log.Println("receive count:", rcv.count)
	log.Println("send count:", snd.count)
}
```
