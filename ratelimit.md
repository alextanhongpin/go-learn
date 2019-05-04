# Ratelimit

Basic rate-limiting with go.

```go
package ratelimit

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limiter interface {
	Add(id string) *rate.Limiter
	Get(id string) *rate.Limiter
	Start()
	Stop()
}

type visitor struct {
	limiter   *rate.Limiter
	updatedAt time.Time
}

type IPLimiter struct {
	sync.RWMutex
	ips map[string]*visitor

	callsPerSecond   int
	burst            int
	cleanupDuration  time.Duration // Frequency to cleanup the user's ip.
	inactiveDuration time.Duration // If the visitor's last seen is above the inactive duration, the visitor will be removed.
	quit             chan struct{}
}

func NewIPLimiter(callsPerSecond, burst int, cleanup, inactive time.Duration) *IPLimiter {
	return &IPLimiter{
		callsPerSecond:   callsPerSecond, // e.g. 2 API calls per second.
		burst:            burst,          // e.g. 5 bursts allowed.
		ips:              make(map[string]*visitor),
		cleanupDuration:  cleanup,  // E.g. 5 minute to run the cron every 5 minutes to check for inactive IP.
		inactiveDuration: inactive, // E.g. 5 minute to remove users that are no longer active for 5 minutes.
		quit:             make(chan struct{}),
	}
}

func (i *IPLimiter) Add(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(i.callsPerSecond), i.burst)
	i.Lock()
	i.ips[ip] = &visitor{limiter, time.Now()}
	i.Unlock()
	return limiter
}

func (i *IPLimiter) Get(ip string) *rate.Limiter {
	i.Lock()
	v, exist := i.ips[ip]
	if !exist {
		i.Unlock()
		return i.Add(ip)
	}
	v.updatedAt = time.Now()
	i.Unlock()
	return v.limiter
}

func (i *IPLimiter) cleanup() {
	ticker := time.NewTicker(i.cleanupDuration)
	for {
		select {
		case <-ticker.C:
			i.Lock()
			for ip, v := range i.ips {
				if time.Since(v.updatedAt) > i.inactiveDuration {
					log.Println("ip", ip, "is removed")
					delete(i.ips, ip)
				}
			}
			i.Unlock()
		case <-i.quit:
			return
		}
	}
}

func (i *IPLimiter) Start() {
	log.Println("rateLimiter: start")
	go i.cleanup()
}

func (i *IPLimiter) Stop() {
	log.Println("rateLimiter: stop")
	close(i.quit)
}
```


## Multi Ratelimiter


This multi ratelimiter accepts multiple rate-limit rules, and if one of them is fulfilled, then an error RateExceededError (429) can be returned. The solutions attempt to store all the timestamp of the event for the given identifier (client ip/path combination), which can grow in size. To keep the memory usage low, we can:

- stop adding event timestamps once it has exceeded the rate limit, and set an expire time in the future to indicate when the service will be available (the smallest timestamp plus the largest period of the rule)
- clear it up at the interval of 1 * largest period. If the largest period is 1 hour, then it makes sense to retain the timestamps for at least 1 hour before clearing them up. In order to get the largest period, we sort the rules in ascending order and take the last rule's period. Then, for each map of client id events, we take the current time minus the largest period `t_exp`, loop through the events and find the first position of the slice array that has value greater than the `t_exp`. From there, we only take the slice index from that position. This is a more efficient way of filtering, since the timestamps are appended, they are already sorted in ascending order.

```go
package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Rule struct {
	Threshold int
	Period    time.Duration
	Name      string
}

func NewRule(threshold int, period time.Duration, name string) *Rule {
	return &Rule{
		Threshold: threshold,
		Period:    period,
		Name:      name,
	}
}

func (r *Rule) IsThresholdExceeded(events []time.Time) bool {
	var count int
	for _, evt := range events {
		if time.Since(evt) < r.Period {
			count++
		}
	}
	return count > r.Threshold
}

type RateLimiter struct {
	sync.RWMutex
	m     map[string][]time.Time
	rules []*Rule
}

func NewRateLimiter(rules []*Rule) *RateLimiter {
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Period < rules[j].Period
	})
	return &RateLimiter{
		m:     make(map[string][]time.Time),
		rules: rules,
	}
}

func (r *RateLimiter) Clean() func(context.Context) {
	var duration time.Duration
	if len(r.rules) == 0 {
		duration = time.Minute
	} else {
		duration = r.rules[len(r.rules)-1].Period
	}
	t := time.NewTicker(duration)
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-t.C:
				r.Lock()
				// Get the current time where all the counts will be invalid.
				exp := time.Now().Add(-duration)
				for id, evts := range r.m {
					var pos int
					for i, evt := range evts {
						if evt.After(exp) {
							pos = i
							break
						}
					}
					fmt.Println("clearing", pos, len(evts))
					r.m[id] = evts[pos:]
					fmt.Println("left", len(r.m[id]))
				}
				r.Unlock()
			}
		}
	}()
	return func(ctx context.Context) {
		sig := make(chan struct{}, 1)
		go func() {
			close(done)
			wg.Done()
			close(sig)
		}()
		select {
		case <-sig:
			fmt.Println("graceful shutdown")
			return
		case <-ctx.Done():
			return
		}
	}
}

func (r *RateLimiter) Add(id string) {
	r.Lock()
	if _, exist := r.m[id]; !exist {
		r.m[id] = make([]time.Time, 0)
	}
	// Before appending, we can remove those that expired first. This is done again by taking the largest period.
	r.m[id] = append(r.m[id], time.Now())
	r.Unlock()
}

func (r *RateLimiter) Allow(id string) bool {
	r.RLock()
	events, exist := r.m[id]
	r.RUnlock()
	if !exist {
		return true
	}
	for _, rule := range r.rules {
		exceeded := rule.IsThresholdExceeded(events)
		if exceeded {
			fmt.Println(rule.Name)
			return false
		}
	}
	return true
}

func main() {
	rules := []*Rule{
		NewRule(1, time.Second, "1 call per second"),
		// The rule for the longer duration is that it should made less calls than the shortest for the same period.
		// The above rule limits up to 3600 calls per hour. But the bottom conditions allows only 30 per hour.
		NewRule(30, time.Hour, "30s call per hour"),
	}
	r := NewRateLimiter(rules)
	shutdown := r.Clean()

	id := "1"
	r.Add(id)
	fmt.Println(r.Allow(id))
	r.Add(id)
	fmt.Println(r.Allow(id))

	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow(id))

	// Waste all quota.
	for i := 0; i < 25; i++ {
		r.Add(id)
	}
	// To bypass the 1 call per second limit.
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow(id))

	for i := 0; i < 5; i++ {
		r.Add(id)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow(id))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdown(ctx)
}
```
## Others

Take a look at other ratelimiting implementation:

- leaky bucket 
- token bucket
- fixed window counter
- sliding window log
- sliding window

## Leaky Bucket

```go
package main

import (
	"fmt"
	"time"
)

type LeakyBucket struct {
	requestsPerSecond int64
	nextAllowedTime   time.Time
	interval          time.Duration
}

func NewRateLimiter(requestsPerSecond int64) *LeakyBucket {
	return &LeakyBucket{
		requestsPerSecond: requestsPerSecond,
		interval:          time.Second / time.Duration(requestsPerSecond),
	}
}

func (l *LeakyBucket) Allow() bool {
	now := time.Now()
	if now.Equal(l.nextAllowedTime) || now.After(l.nextAllowedTime) {
		l.nextAllowedTime = time.Now().Add(l.interval)
		return true
	}
	return false
}

func main() {
	r := NewRateLimiter(5)
	fmt.Println(r.Allow())
	time.Sleep(300 * time.Millisecond)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(300 * time.Millisecond)
	fmt.Println(r.Allow())
}
```

## Token Bucket

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	requestsPerSecond int64
	sync.RWMutex
	counter int64
}

func NewRateLimiter(requestsPerSecond int64) *TokenBucket {
	return &TokenBucket{
		requestsPerSecond: requestsPerSecond,
		counter:           requestsPerSecond,
	}
}

func (t *TokenBucket) Allow() bool {
	t.RLock()
	counter := t.counter
	t.RUnlock()
	if counter <= 0 {
		return false
	}
	t.Lock()
	t.counter--
	t.Unlock()
	return t.counter < t.requestsPerSecond
}

func (t *TokenBucket) Start() func(context.Context) {
	ticker := time.NewTicker(time.Second / time.Duration(t.requestsPerSecond))
	var wg sync.WaitGroup

	wg.Add(1)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				t.Lock()
				t.counter++
				t.Unlock()
				return
			}
		}
	}()
	return func(ctx context.Context) {
		sig := make(chan struct{})
		go func() {
			close(done)
			wg.Wait()
			close(sig)
		}()
		select {
		case <-ctx.Done():
			return
		case <-sig:
			return
		}
	}
}

func main() {
	r := NewRateLimiter(5)
	shutdown := r.Start()

	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdown(ctx)
}
```

## Multi-rate limiter 

```go
package main

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiters struct {
	limiters []RateLimiter
}

func (l *multiLimiters) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

type ApiConnection struct {
	rateLimiter RateLimiter
}

func Open() *ApiConnection {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)   // Limit per second with no burstiness
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // Limit per second with burstiness of 10
	&ApiConnection{
		// 1 per second.
		// rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

func (a *ApiConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func (a *ApiConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Hello, playground")
}
```


# Different Algorithms

NOTE: The implementations below are not production ready - they are not concurrent safe (no mutex/channel applied), not distributed (only works on local machine, should create an interface that calls redis/datastore instead) and could have been written in a better way. 

## Fixed Window Counter

For every time window, keep a counter and increment them everytime they are called.

Given a period of 5 seconds,
when the time is `<time>`
then the time window should be `<time window>`

- time: 00:03, time window: 00:00, counter: 1
- time: 00:05, time window: 00:05, counter: 1
- time: 00:06, time window: 00:05, counter: 2
- time: 00:07, time window: 00:05, counter: 3
- time: 00:11, time window: 00:10, counter: 1

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func moduloTime(ts, window int64) int64 {
	return ts - (ts % window)
}

type FixedWindowCounter struct {
	requestsPerSecond int64
	sync.RWMutex
	counter map[int64]int64
}

func NewRateLimiter(n int64) *FixedWindowCounter {
	return &FixedWindowCounter{
		requestsPerSecond: n,
		counter:           make(map[int64]int64),
	}
}

func (f *FixedWindowCounter) getTimeWindow() int64 {
	return moduloTime(time.Now().Unix(), f.requestsPerSecond)
}

func (f *FixedWindowCounter) Allow() bool {
	f.RLock()
	count := f.counter[f.getTimeWindow()]
	f.counter[f.getTimeWindow()] += 1
	f.RUnlock()
	return count < f.requestsPerSecond
}

func main() {
	r := NewRateLimiter(5)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
}
```


## Sliding Window Log
```go
package main

import (
	"fmt"
	"log"
	"time"
)

// round the unix time (seconds) to the nearest window.
func modulo(now int64, window int) int64 {
	return now - (now % int64(window))
}

type SlidingWindowLog struct {
	maxRequestPerSec int
	logs             []time.Time
}

func NewSlidingWindowLog(n int) *SlidingWindowLog {
	return &SlidingWindowLog{
		maxRequestPerSec: n,
		logs:             make([]time.Time, 0),
	}
}

func (s *SlidingWindowLog) Allow() bool {
	curr := time.Now().Add(-1 * time.Second)
	var pos int
	for i, t := range s.logs {
		pos = i
		if t.After(curr) {
			break
		}
	}
	s.logs = s.logs[pos:]
	s.logs = append(s.logs, time.Now())
	log.Println(s.logs, pos)
	return len(s.logs) <= s.maxRequestPerSec
}

func main() {
	r := NewSlidingWindowLog(5)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(5 * time.Second)
	fmt.Println(r.Allow())
}
```

## Sliding Window

NOTE: This is an unoptimized version. We can skip storing the timestamp of the events and instead store just the count.

```go
package main

import (
	"fmt"
	"time"
)

// round the unix time (seconds) to the nearest window.
func modulo(now int64, window int) int64 {
	return now - (now % int64(window))
}

type SlidingWindowLog struct {
	maxRequestPerSec int
	logs             map[int64][]time.Time
}

func NewSlidingWindowLog(n int) *SlidingWindowLog {
	return &SlidingWindowLog{
		maxRequestPerSec: n,
		logs:             make(map[int64][]time.Time, 0),
	}
}

func (s *SlidingWindowLog) Allow() bool {
	now := time.Now()
	curr := modulo(now.Unix(), 1)
	prev := modulo(now.Add(-1*time.Second).Unix(), 1)
	prevprev := modulo(now.Add(-2*time.Second).Unix(), 1)
	for k := range s.logs {
		if k <= prevprev {
			delete(s.logs, k)
		}
	}

	logsNow := s.logs[curr]
	logsPrev := s.logs[prev]

	total := int64(len(logsPrev))*(int64(1)-(now.Unix()-curr)) + int64(len(logsNow))
	s.logs[curr] = append(s.logs[curr], time.Now())
	fmt.Println(s.logs)
	return total <= int64(s.maxRequestPerSec)
}

func main() {
	r := NewSlidingWindowLog(5)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(5 * time.Second)
	fmt.Println(r.Allow())
	time.Sleep(5 * time.Second)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
}
```

## Improved sliding window
```go
package main

import (
	"fmt"
	"time"
)

// round the unix time (seconds) to the nearest window.
func modulo(now int64, window int) int64 {
	return now - (now % int64(window))
}

type SlidingWindowLog struct {
	maxRequestPerSec int
	logs             map[int64]int
}

func NewSlidingWindowLog(n int) *SlidingWindowLog {
	return &SlidingWindowLog{
		maxRequestPerSec: n,
		logs:             make(map[int64]int, 0),
	}
}

func (s *SlidingWindowLog) Allow() bool {
	now := time.Now()
	curr := modulo(now.Unix(), 1)
	prev := modulo(now.Add(-1*time.Second).Unix(), 1)
	prevprev := modulo(now.Add(-2*time.Second).Unix(), 1)
	for k := range s.logs {
		if k <= prevprev {
			delete(s.logs, k)
		}
	}

	logsNow := s.logs[curr]
	logsPrev := s.logs[prev]

	total := int64(logsPrev)*(int64(1)-(now.Unix()-curr)) + int64(logsNow)
	s.logs[curr]++
	fmt.Println(s.logs)
	return total <= int64(s.maxRequestPerSec)
}

func main() {
	r := NewSlidingWindowLog(5)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(5 * time.Second)
	fmt.Println(r.Allow())
	time.Sleep(5 * time.Second)
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
}
```

## References

- https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
