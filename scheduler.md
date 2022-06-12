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

	sch.Schedule("greet", time.Now().Add(1*time.Second), func() {
		fmt.Println("after: hi")
	})
	sch.Schedule("bark", time.Now().Add(1*time.Second), func() {
		fmt.Println("after: woof")
	})
	sch.Unschedule("bark")

	for _, task := range sch.List(-1) {
		fmt.Println("list:", task.name, task.runAt.Sub(time.Now()))
	}
	sch.Schedule("smile", time.Now().Add(4*time.Second), func() {
		fmt.Println("after: :)")
	})
	time.Sleep(3 * time.Second)

	fmt.Println("any tasks left?")
	for _, task := range sch.List(-1) {
		fmt.Println("list:", task.name, task.runAt.Sub(time.Now()))
	}
	time.Sleep(2 * time.Second)
	fmt.Println("program exiting")
}

type Task struct {
	name    string
	runAt   time.Time
	runFunc *time.Timer
}

type Scheduler struct {
	tasks        sync.Map
	scheduleCh   chan *Task
	unscheduleCh chan string
	init, quit   sync.Once
	done         chan struct{}
	wg           sync.WaitGroup
}

func NewScheduler() (*Scheduler, func()) {
	s := &Scheduler{
		scheduleCh:   make(chan *Task),
		unscheduleCh: make(chan string),
		done:         make(chan struct{}),
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
	case s.scheduleCh <- &Task{name: name, runAt: runAt, runFunc: time.AfterFunc(runAt.Sub(time.Now()), func() {
		fn()
		s.tasks.Delete(name)
	})}:
		return true
	}
}
func (s *Scheduler) Unschedule(name string) bool {
	select {
	case <-s.done:
		return false
	case s.unscheduleCh <- name:
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
		case task := <-s.scheduleCh:
			if t, found := s.tasks.LoadAndDelete(task.name); found {
				t.(*Task).runFunc.Stop()
			}
			_, _ = s.tasks.LoadOrStore(task.name, task)
		case key := <-s.unscheduleCh:
			if t, found := s.tasks.LoadAndDelete(key); found {
				t.(*Task).runFunc.Stop()
			}
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
