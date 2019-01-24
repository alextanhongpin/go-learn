Worker Pool implementation

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
type WorkerPool struct {
	wg     *sync.WaitGroup
	cond   *sync.Cond
	quit   chan interface{}
	taskCh chan Task
}

func NewWorkerPool(taskLimit int) *WorkerPool {
	return &WorkerPool{
		quit:   make(chan interface{}, 1),
		taskCh: make(chan Task, taskLimit),
		cond:   sync.NewCond(new(sync.Mutex)),
		wg:     new(sync.WaitGroup),
	}
}

func (w *WorkerPool) Start(n int) *sync.WaitGroup {
	w.wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			w.loop(i)
		}(i)
	}
	fmt.Printf("started %d workers\n", n)
	return w.wg
}

func (w *WorkerPool) AddTask(tasks ...Task) {
	for _, task := range tasks {
		select {
		case <-w.quit:
			return
		case w.taskCh <- task:
			w.cond.Broadcast()
			fmt.Println("received task", task)
		}
	}
}

func (w *WorkerPool) loop(id int) {
	fmt.Println("starting worker", id)
	defer func() {
		fmt.Println("defer done")
		w.wg.Done()
	}()
	for {
		w.cond.L.Lock()
		for len(w.taskCh) == 0 {
			select {
			case <-w.quit:
				w.cond.L.Unlock()
				return
			default:
				fmt.Println("waiting for tasks")
				w.cond.Wait()
				fmt.Println("wait ends", len(w.taskCh))
			}
		}
		w.cond.L.Unlock()
		select {
		case <-w.quit:
			fmt.Println("quitting")
			return
		case task, ok := <-w.taskCh:
			if !ok {
				return
			}
			res := task.Execute()
			fmt.Println(res)
		default:
			fmt.Println("defaults")
		}
	}

}
func (w *WorkerPool) Stop() {
	fmt.Println("Stop")
	w.cond.Broadcast()
	close(w.taskCh)
	close(w.quit)
}

func main() {
	wp := NewWorkerPool(10)

	job := wp.Start(3)
	go func() {
		tasks := []Task{&DelayTask{}}
		wp.AddTask(tasks...)
	}()
	go func() {
		time.Sleep(2 * time.Second)
		tasks := []Task{&DelayTask{}, &DelayTask{}}
		wp.AddTask(tasks...)
	}()
	go func() {
		time.Sleep(5 * time.Second)
		wp.Stop()
	}()

	job.Wait()
	fmt.Println("exiting")
}

type DelayTask struct{}

func (d *DelayTask) Execute() Result {
	time.Sleep(500 * time.Millisecond)
	return Result{
		Response: "done",
	}
}
```
