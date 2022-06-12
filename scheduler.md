# Thread-safe scheduler

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

func main() {
	sch, stop := NewScheduler()
	defer stop()

	sch.Schedule("echo", time.Now().Add(1*time.Second), func() {
		fmt.Println("after: hello world")
	})
	sch.Schedule("echo", time.Now().Add(2*time.Second), func() {
		fmt.Println("after: overwrite")
	})

	sch.Schedule("sth", time.Now().Add(1*time.Second), func() {
		fmt.Println("after: something else")
	})
	for _, task := range sch.List(-1) {
		fmt.Println("list:", task.name, task.runAt.Sub(time.Now()))
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Hello, 世界")
}

type Task struct {
	name    string
	runAt   time.Time
	runFunc *time.Timer
}

type Scheduler struct {
	tasks      sync.Map
	ch         chan *Task
	init, quit sync.Once
	done       chan struct{}
	wg         sync.WaitGroup
}

func NewScheduler() (*Scheduler, func()) {
	s := &Scheduler{
		ch:   make(chan *Task),
		done: make(chan struct{}),
	}
	return s, s.stop
}

func (s *Scheduler) Schedule(name string, runAt time.Time, fn func()) bool {
	s.init.Do(func() {
		s.loopAsync()
	})
	select {
	case <-s.done:
		return false
	case s.ch <- &Task{name: name, runAt: runAt, runFunc: time.AfterFunc(runAt.Sub(time.Now()), fn)}:
		return true
	}
}

func (s *Scheduler) loopAsync() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		s.loop()
	}()
}

func (s *Scheduler) loop() {
	for {
		select {
		case <-s.done:
			return
		case task := <-s.ch:
			if t, found := s.tasks.LoadAndDelete(task.name); found {
				t.(*Task).runFunc.Stop()
			}
			_, _ = s.tasks.LoadOrStore(task.name, task)
		}
	}
}

func (s *Scheduler) stop() {
	s.quit.Do(func() {
		close(s.done)
		s.wg.Wait()
	})
}

func (s *Scheduler) List(n int) []Task {
	var tasks []Task
	s.tasks.Range(func(key, value any) bool {
		if n == 0 {
			return false
		}

		task, ok := value.(*Task)
		if !ok {
			return ok
		}
		tasks = append(tasks, Task{
			name:  task.name,
			runAt: task.runAt,
		})
		n--
		return n != 0
	})

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].runAt.Before(tasks[j].runAt)
	})

	return tasks
}
```
