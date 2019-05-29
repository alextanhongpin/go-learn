Checks if the data can be synced based on the changes made per time range. Can be used for syncing algorithm, or online-offline transitioning.

```go
package main

import (
	"fmt"
	"time"
)

type Tracker struct {
	counter   int64
	start     time.Time
	threshold int64
	duration  time.Duration
}

func (t *Tracker) ShouldUpdate() bool {
	return t.counter > t.threshold || time.Since(t.start) > t.duration
}

func (t *Tracker) Reset() {
	t.start = time.Now()
	t.counter = 0
}

func (t *Tracker) Increment() {
	t.counter++
}

func NewTracker(threshold int64, duration time.Duration) *Tracker {
	tracker := &Tracker{
		threshold: threshold,
		duration:  duration,
	}
	tracker.Reset()
	return tracker
}

func main() {
	// Update if there are 10 calls made or the every 1 second.
	tracker := NewTracker(10, 1*time.Second)
	fmt.Println(tracker.ShouldUpdate())
	for i := 0; i < 20; i++ {
		tracker.Increment()
	}
	fmt.Println(tracker.ShouldUpdate())
	tracker.Reset()
	fmt.Println(tracker.ShouldUpdate())
	time.Sleep(2 * time.Second)
	fmt.Println(tracker.ShouldUpdate())

	// TODO: Chain it with multiple conditions
	NewTracker(10000, 10*time.Second) // Update if there are at least 10,000 calls made in the last second.
	NewTracker(100, 10*time.Minute)   // Update if there are at least 100 calls made in the last minute.
	NewTracker(10, 1*time.Hour)       // Update if there are at least 10 call made in the last hour.
	NewTracker(1, 1*time.Day)         // Update if there are at least 1 call in a dall.
}
```
