## Naive Send Receive
Sample send/receive pattern using channels. As long as the struct shares the same channel, they could be decoupled and share the same instance of the data.

The issue with this code is it doesn't prevent the user from sending to a closed channel. See the solution below.
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

## Safer Queue

```
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Queue struct {
	sync.RWMutex

	ch     chan int
	closed bool
}

func NewQueue(cap int) *Queue {
	return &Queue{
		ch: make(chan int, cap),
	}
}

func (q *Queue) Close() {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return
	}
	close(q.ch)
	q.closed = true
	log.Println("queue closed")
}

func (q *Queue) Enqueue(i int) error {
	q.RLock()
	defer q.RUnlock()
	if q.closed {
		return fmt.Errorf("queue closed, unable to send")
	}
	q.ch <- i
	log.Println("queue send", i)
	return nil
}

func (q *Queue) Receiver() <-chan int {
	return q.ch
}

type Receiver struct {
	count uint
	q     *Queue
}

func (r *Receiver) Start() {
	go func() {
		for v := range r.q.Receiver() {
			log.Println("receive", v)
			r.count++
		}
	}()
}

func NewReceiver(q *Queue) *Receiver {
	return &Receiver{
		q: q,
	}
}

type Sender struct {
	count uint
	q     *Queue
}

func (s *Sender) Send(i int) {
	if err := s.q.Enqueue(i); err != nil {
		log.Println(err)
		return
	}
	s.count++
}

func NewSender(q *Queue) *Sender {
	return &Sender{
		q: q,
	}
}

func main() {
	var wg sync.WaitGroup
	q := NewQueue(100)

	rcv := NewReceiver(q)
	snd := NewSender(q)

	rcv.Start()

	snd.Send(1)
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			snd.Send(i)
		}

		q.Close()

		time.Sleep(2 * time.Second)
		for i := 0; i < 10; i++ {
			go snd.Send(i)
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("receive count:", rcv.count)
	log.Println("send count:", snd.count)
}
```
