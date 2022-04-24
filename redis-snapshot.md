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

Using atomic

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Every struct {
	Seconds   int64
	Threshold int64
}

type Snapshotter struct {
	unix      int64
	counter   int64
	every     []Every
	deltaTick int64
	done      chan bool
	callback  func()
}

func NewSnapshotter(every []Every, callback func()) (*Snapshotter, func()) {
	if len(every) == 0 {
		panic("initialized without periods")
	}

	sort.Slice(every, func(i, j int) bool {
		return every[i].Seconds < every[j].Seconds
	})

	thresholds := make([]int64, len(every))
	for i, e := range every {
		thresholds[i] = e.Threshold
	}

	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i] > thresholds[j]
	})
	for i, e := range every {
		if thresholds[i] != e.Threshold {
			panic("threshold must be in descending values")
		}
	}

	deltaTick := big.NewInt(every[0].Seconds)
	for _, e := range every[1:] {
		deltaTick.GCD(nil, nil, deltaTick, big.NewInt(e.Seconds))
	}

	s := &Snapshotter{
		unix:      time.Now().Unix(),
		counter:   0,
		every:     every,
		deltaTick: deltaTick.Int64(),
		done:      make(chan bool),
		callback:  callback,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.loop()
	}()
	var once sync.Once
	return s, func() {
		once.Do(func() {
			close(s.done)
			wg.Wait()
		})
	}
}

func (s *Snapshotter) isTriggered() bool {
	lastUnix := atomic.LoadInt64(&s.unix)
	elapsedSeconds := int64(time.Since(time.Unix(lastUnix, 0)).Seconds())

	counter := s.Count()
	for _, every := range s.every {
		if elapsedSeconds < every.Seconds {
			return false
		}
		// Execute at every given seconds, if the amount is at least the given threshold.
		// Changing this to <= will hit the condition at least the given seconds with the given threshold.
		if elapsedSeconds == every.Seconds {
			if counter < every.Threshold {
				return false
			}
			return true
		}
	}
	return false
}

func (s *Snapshotter) loop() {
	t := time.NewTicker(time.Duration(s.deltaTick) * time.Second)
	for {
		select {
		case <-t.C:
			if s.isTriggered() {
				s.callback()
				atomic.StoreInt64(&s.counter, 0)
				atomic.StoreInt64(&s.unix, time.Now().Unix())
			}
		case <-s.done:
			return
		}
	}
}

func (s *Snapshotter) Inc() int64 {
	return atomic.AddInt64(&s.counter, 1)
}

func (s *Snapshotter) Count() int64 {
	return atomic.LoadInt64(&s.counter)
}

func main() {
	s, close := NewSnapshotter([]Every{
		{1, 10}, // Execute every 1 second if the count is at least 10
		{2, 5},  // Execute every 2 second if the count is at least 5
		// Second 3 and 4 won't be checked.
		{5, 1}, // Execute every 5 second if the count is at least 1
	}, func() {
		fmt.Println("triggered")
	})
	defer close()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("incrementing")
		for i := 0; i < 5; i++ {
			s.Inc()
		}
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("incrementing")
		for i := 0; i < 5; i++ {
			s.Inc()
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("Hello, 世界")
}
```
