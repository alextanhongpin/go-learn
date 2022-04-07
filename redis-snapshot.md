## Redis-like Snapshotting logic

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

type Every struct {
	Frequency int
	Duration  time.Duration
}

type Snapshotter struct {
	sync.RWMutex
	every     []Every
	count     int
	syncFn    func() error
	lastRunAt time.Time
}

func NewSnapshotter(every []Every, syncFn func() error) (*Snapshotter, func()) {
	sort.Slice(every, func(i, j int) bool {
		return every[i].Duration < every[j].Duration
	})
	s := &Snapshotter{
		every:     every,
		syncFn:    syncFn,
		lastRunAt: time.Now(),
	}
	stop := s.start()
	return s, stop
}

func (s *Snapshotter) isTriggered() bool {
	s.RLock()
	defer s.RUnlock()

	for _, e := range s.every {
		if time.Since(s.lastRunAt) < e.Duration {
			return false
		}
		if s.count >= e.Frequency {
			return true
		}
	}
	return false
}
func (s *Snapshotter) start() func() {
	done := make(chan bool)
	go func() {
		t := time.NewTicker(s.every[0].Duration)
		defer t.Stop()

		for {
			select {
			case last := <-t.C:
				fmt.Println(time.Since(s.lastRunAt))
				if s.isTriggered() {
					_ = s.syncFn()
					s.Lock()
					s.count = 0
					s.lastRunAt = last
					s.Unlock()
				}
			case <-done:
				return
			}
		}
	}()

	return func() {
		close(done)
	}
}

func (s *Snapshotter) SetCount(n int) {
	s.Lock()
	s.count = n
	s.Unlock()
}
func (s *Snapshotter) Count() (count int) {
	s.RLock()
	count = s.count
	s.RUnlock()
	return
}

func main() {
	every := []Every{
		Every{9, 1 * time.Second},
		Every{5, 2 * time.Second},
		Every{1, 5 * time.Second},
	}
	s, stop := NewSnapshotter(every, func() error {
		fmt.Println("snapshotting")
		return nil
	})
	defer stop()
	s.SetCount(1)

	time.Sleep(10 * time.Second)
	fmt.Println("Hello, 世界")
}
```
