# Background Channel Pattern

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	r := NewReceiver(10)
	r.Start()

	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < 100; i++ {
		fmt.Println("sending", i, "is", r.Send(i))
	}

	go func() {
		time.Sleep(1 * time.Second)
		r.Stop()
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("terminating")
}

type Receiver struct {
	ch   chan int
	quit chan struct{}
}

func NewReceiver(buffer int) *Receiver {
	return &Receiver{
		ch:   make(chan int, buffer),
		quit: make(chan struct{}),
	}
}

func (r *Receiver) Start() {
	go func() {
		for {
			select {
			case <-r.quit:
				break
			case v, ok := <-r.ch:
				if !ok {
					break
				}
				// Fake processing time.
				// time.Sleep(100 * time.Millisecond)
				fmt.Println("got", v)
			}
		}
	}()
}

func (r *Receiver) Stop() {
	// Compare to r.quit<-struct{}, which is better?
	close(r.quit)
}

func (r *Receiver) Send(i int) bool {
	j := jitter(50, 1000)
	fmt.Println("jitter is", j)
	select {
	case <-r.quit:
		fmt.Println("quit")
		return false
	case r.ch <- i:
		fmt.Println("send", i)
		return true
		// This pattern will drop the message if it fails to deliver it within the set time duration.
	case <-time.After(j):
		return false
		// This pattern will just drop the message if it's unable to send to the channel.
		//default:
		//	return false
	}
}

func jitter(min, max int) time.Duration {
	duration := min + rand.Intn(max)
	return time.Duration(duration) * time.Millisecond
}
```
