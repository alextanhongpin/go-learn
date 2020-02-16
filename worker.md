## Background Worker

```
interface 
- start() // Start the worker.
- stop() // Stop the worker.
- send() // Send the message to the worker.
```

Design considerations:
- Send payload to process in the background. How to make it reliable?
- Store the unprocessed ones.
- Process it.
- If the server restarts, load the unprocessed data and trigger them.

```
worker.send(payload)
```

Register handlers during construction to handle polymorphism:
```
worker.callbacks[task] = fn
```

Process workers:
```
worker.process(callback())
```

## Single worker with buffer

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ch     chan interface{}
	wg     sync.WaitGroup
	once   sync.Once
	open bool
	rw     sync.RWMutex
}

func NewWorker(size int) (*Worker, func()) {
	w := &Worker{
		ch: make(chan interface{}, size),
		open: true,
	}
	return w, w.start()
}

func (w *Worker) Send(payload interface{}) {
	w.rw.RLock()
	open := w.open
	w.rw.RUnlock()
	if !open {
		return
	}
	w.ch <- payload
}

func (w *Worker) start() func() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case work, ok := <-w.ch:
				if !ok {
					return
				}
				fmt.Println(work)
			}
		}
	}()
	return func() {
		w.once.Do(func() {
			close(w.ch)

			w.rw.Lock()
			w.open = false
			w.rw.Unlock()

			w.wg.Wait()
		})

	}
}

func main() {
	worker, cancel := NewWorker(2)
	defer cancel()

	for i := 0; i < 20; i++ {
		worker.Send(i)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Hello, playground")
}
```

## With atomic compare to check channels state
```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	ch   chan interface{}
	wg   sync.WaitGroup
	open int32
}

func NewWorker(size int) (*Worker, func()) {
	w := &Worker{
		ch:   make(chan interface{}, size),
		open: 1,
	}
	return w, w.start()
}

func (w *Worker) Send(payload interface{}) {
	if !w.isOpened() {
		return
	}
	w.ch <- payload
}

func (w *Worker) isOpened() bool {
	return atomic.LoadInt32(&w.open) == 1
}

func (w *Worker) start() func() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case work, ok := <-w.ch:
				if !ok {
					return
				}
				fmt.Println(work)
			}
		}
	}()
	return func() {
		if atomic.CompareAndSwapInt32(&w.open, 1, 0) {
			close(w.ch)
			w.wg.Wait()
		}
	}
}

func main() {
	worker, cancel := NewWorker(2)
	defer cancel()

	for i := 0; i < 20; i++ {
		worker.Send(i)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("Hello, playground")
}
```
