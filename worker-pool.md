# Worker pool

Add multiple background workers for long running processes:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	Response interface{}
	Err      error
}

type Task interface {
	Execute() Result
}

type DelayTask struct{}

func (d *DelayTask) Execute() Result {
	time.Sleep(100 * time.Millisecond)
	return Result{
		Response: "done",
	}
}

type WorkerPool struct {
	quit chan interface{}
	wg   sync.WaitGroup
	mu   sync.Mutex
	in   chan Task
	out  chan Result
}

func NewWorkerPool(n int) *WorkerPool {
	wp := &WorkerPool{
		quit: make(chan interface{}),
		in:   make(chan Task, 10),
	}

	// Spawn n background workers.
	for i := 0; i < n; i++ {
		go wp.loop()
	}
	return wp
}

func (w *WorkerPool) loop() {
	for {
		select {
		case <-w.quit:
			return
		case task, ok := <-w.in:
			if !ok {
				return
			}
			res := task.Execute()
			fmt.Println("executed", res)
		default:
		}
	}
}

func (w *WorkerPool) AddTask(tasks ...Task) {
	for _, task := range tasks {
		select {
		case <-w.quit:
			return
		case w.in <- task:
			fmt.Println("send", task)
		}
	}
}

func (w *WorkerPool) Stop() {
	close(w.quit)
}

func main() {
	wp := NewWorkerPool(3)

	tasks := []Task{&DelayTask{}, &DelayTask{}, &DelayTask{}}
	wp.AddTask(tasks...)
	go func() {
		time.Sleep(1 * time.Second)
		tasks := []Task{&DelayTask{}, &DelayTask{}, &DelayTask{}}
		wp.AddTask(tasks...)
	}()

	time.Sleep(3 * time.Second)
	fmt.Println("exiting")
}
```
