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
	if counter == 0 {
		return false
	}
	t.Lock()
	t.counter--
	t.Unlock()
	return t.counter < t.requestsPerSecond
}

func (t *TokenBucket) Start() func(context.Context) {
	// Note: We can just refill with with the requestsPerSecond every second.
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
				t.RLock()
				counter := t.counter
				if counter == t.requestsPerSecond {
					return
				}
				t.RUnlock()
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

Optimized token bucket with lazy refill:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	requestsPerSecond int64
	sync.RWMutex
	counter        int64
	lastRefillTime time.Time
}

func NewRateLimiter(requestsPerSecond int64) *TokenBucket {
	return &TokenBucket{
		requestsPerSecond: requestsPerSecond,
		counter:           requestsPerSecond,
	}
}

func (t *TokenBucket) Allow() bool {
	t.lazyRefill()
	t.RLock()
	counter := t.counter
	t.RUnlock()
	if counter == 0 {
		return false
	}
	t.Lock()
	t.counter--
	t.Unlock()
	return t.counter < t.requestsPerSecond
}

func (t *TokenBucket) lazyRefill() {
	elapsedSeconds := time.Now().Sub(t.lastRefillTime).Seconds()
	c := int64(elapsedSeconds) * t.requestsPerSecond
	if c > 0 {
		t.Lock()
		t.counter = min(t.counter+c, t.requestsPerSecond)
		t.Unlock()
		t.lastRefillTime = time.Now()
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	r := NewRateLimiter(5)

	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow())
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
	"sync"
	"time"
)

type SlidingWindowLog struct {
	requestsPerSecond int64
	sync.RWMutex
	events []time.Time
}

func NewRateLimiter(requestsPerSecond int64) *SlidingWindowLog {
	return &SlidingWindowLog{
		requestsPerSecond: requestsPerSecond,
		events:            make([]time.Time, 0),
	}
}

func (s *SlidingWindowLog) Allow() bool {
	s.Lock()
	now := time.Now().Add(-1 * time.Second)
	var pos int
	for i, evt := range s.events {
		pos = i
		if evt.After(now) {
			break
		}
	}
	s.events = s.events[pos:]
	s.events = append(s.events, time.Now())
	counter := int64(len(s.events))
	s.Unlock()
	return counter < s.requestsPerSecond
}

func main() {
	r := NewRateLimiter(5)

	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	fmt.Println(r.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(r.Allow())
}

```

## Sliding Window

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindow struct {
	requestsPerSecond int64
	sync.RWMutex
	events map[int64]int64
}

func NewRateLimiter(requestsPerSecond int64) *SlidingWindow {
	return &SlidingWindow{
		requestsPerSecond: requestsPerSecond,
		events:            make(map[int64]int64),
	}
}

func (s *SlidingWindow) Allow() bool {
	now := time.Now()
	curr := now.Unix()
	prev := curr - 1

	s.Lock()
	for ts := range s.events {
		if ts < prev {
			delete(s.events, ts)
		}
	}
	prevCount := s.events[prev]
	s.events[curr]++
	currCount := s.events[curr]
	fmt.Println(s.events)
	s.Unlock()

	if prevCount == 0 {
		return currCount < s.requestsPerSecond
	}
	delta := 1 - (float64(now.UnixNano()-(now.Unix()*1e9)) / 1e9)
	counter := int64(float64(prevCount)*delta + float64(currCount))
	return counter < s.requestsPerSecond
}

func main() {
	rl := NewRateLimiter(5)
	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	fmt.Println(rl.Allow())
	time.Sleep(1 * time.Second)
	fmt.Println(rl.Allow())
	time.Sleep(3 * time.Second)
	fmt.Println(rl.Allow())
}
```

## Rate Limit Manager for ClientIP

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SlidingWindow struct {
	requestsPerSecond int64
	sync.RWMutex
	events map[int64]int64
}

type RateLimiter interface {
	Allow() bool
}

func NewRateLimiter(requestsPerSecond int64) *SlidingWindow {
	return &SlidingWindow{
		requestsPerSecond: requestsPerSecond,
		events:            make(map[int64]int64),
	}
}

func (s *SlidingWindow) Allow() bool {
	now := time.Now()
	curr := now.Unix()
	prev := curr - 1

	s.Lock()
	for ts := range s.events {
		if ts < prev {
			delete(s.events, ts)
		}
	}
	prevCount := s.events[prev]
	s.events[curr]++
	currCount := s.events[curr]
	fmt.Println(s.events)
	s.Unlock()

	if prevCount == 0 {
		return currCount < s.requestsPerSecond
	}
	delta := 1 - (float64(now.UnixNano()-(now.Unix()*1e9)) / 1e9)
	counter := int64(float64(prevCount)*delta + float64(currCount))
	return counter < s.requestsPerSecond
}

type ClientRateLimiter struct {
	limiter        RateLimiter
	lastActiveTime time.Time
}

func NewClientRateLimiter(requestsPerSecond int64) *ClientRateLimiter {
	return &ClientRateLimiter{
		limiter:        NewRateLimiter(requestsPerSecond),
		lastActiveTime: time.Now(),
	}
}
func (c *ClientRateLimiter) Elapsed(duration time.Duration) bool {
	fmt.Println(time.Since(c.lastActiveTime))
	return time.Since(c.lastActiveTime) > duration
}
func (c *ClientRateLimiter) Allow() bool {
	return c.limiter.Allow()
}

type RateLimitManager struct {
	requestsPerSecond int64
	sync.RWMutex
	m map[string]*ClientRateLimiter
}

func NewRateLimitManager(requestsPerSecond int64) *RateLimitManager {
	return &RateLimitManager{
		requestsPerSecond: requestsPerSecond,
		m:                 make(map[string]*ClientRateLimiter),
	}
}

func (r *RateLimitManager) Allow(clientID string) bool {
	r.Lock()
	limiter, exist := r.m[clientID]
	if !exist {
		r.m[clientID] = NewClientRateLimiter(r.requestsPerSecond)
		limiter = r.m[clientID]
	}
	limiter.lastActiveTime = time.Now()
	r.Unlock()
	return limiter.Allow()
}
func (r *RateLimitManager) Clear() func(context.Context) {
	done := make(chan struct{})
	duration := time.Duration(r.requestsPerSecond*2) * time.Second
	t := time.NewTicker(duration)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("graceful shutdown")
				return
			case <-t.C:
				r.Lock()
				for id, client := range r.m {
					fmt.Println("clearing")
					if client.Elapsed(duration) {
						fmt.Println("clearing", id)
						delete(r.m, id)
					}
				}
				r.Unlock()
			}
		}
	}()
	return func(ctx context.Context) {
		signal := make(chan struct{})
		go func() {
			close(done)
			fmt.Println("closing done")
			wg.Wait()
			close(signal)
		}()
		select {
		case <-ctx.Done():
			return
		case <-signal:
			fmt.Println("signal received")
			return
		}
	}
}

func main() {
	rl := NewRateLimitManager(5)
	shutdown := rl.Clear()
	ip := "0.0.0.0"
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	time.Sleep(1 * time.Second)
	fmt.Println(rl.Allow(ip))
	fmt.Println("sleep 15s")
	time.Sleep(20 * time.Second)
	fmt.Println(rl.Allow(ip))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdown(ctx)
	fmt.Println("shutting down")
}
```

## Rate Limit Manager, Standalone shutdown

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SlidingWindow struct {
	requestsPerSecond int64
	sync.RWMutex
	events map[int64]int64
}

type RateLimiter interface {
	Allow() bool
}

func NewRateLimiter(requestsPerSecond int64) *SlidingWindow {
	return &SlidingWindow{
		requestsPerSecond: requestsPerSecond,
		events:            make(map[int64]int64),
	}
}

func (s *SlidingWindow) Allow() bool {
	now := time.Now()
	curr := now.Unix()
	prev := curr - 1

	s.Lock()
	for ts := range s.events {
		if ts < prev {
			delete(s.events, ts)
		}
	}
	prevCount := s.events[prev]
	s.events[curr]++
	currCount := s.events[curr]
	fmt.Println(s.events)
	s.Unlock()

	if prevCount == 0 {
		return currCount < s.requestsPerSecond
	}
	delta := 1 - (float64(now.UnixNano()-(now.Unix()*1e9)) / 1e9)
	counter := int64(float64(prevCount)*delta + float64(currCount))
	return counter < s.requestsPerSecond
}

type ClientRateLimiter struct {
	limiter        RateLimiter
	lastActiveTime time.Time
}

func NewClientRateLimiter(requestsPerSecond int64) *ClientRateLimiter {
	return &ClientRateLimiter{
		limiter:        NewRateLimiter(requestsPerSecond),
		lastActiveTime: time.Now(),
	}
}
func (c *ClientRateLimiter) Elapsed(duration time.Duration) bool {
	fmt.Println(time.Since(c.lastActiveTime))
	return time.Since(c.lastActiveTime) > duration
}
func (c *ClientRateLimiter) Allow() bool {
	return c.limiter.Allow()
}

type RateLimitManager struct {
	requestsPerSecond int64
	wg                sync.WaitGroup
	sync.RWMutex
	m    map[string]*ClientRateLimiter
	quit chan struct{}
}

func NewRateLimitManager(requestsPerSecond int64) *RateLimitManager {
	mgr := &RateLimitManager{
		requestsPerSecond: requestsPerSecond,
		m:                 make(map[string]*ClientRateLimiter),
		quit:              make(chan struct{}),
	}
	mgr.start()
	return mgr
}

func (r *RateLimitManager) Allow(clientID string) bool {
	r.Lock()
	limiter, exist := r.m[clientID]
	if !exist {
		r.m[clientID] = NewClientRateLimiter(r.requestsPerSecond)
		limiter = r.m[clientID]
	}
	limiter.lastActiveTime = time.Now()
	r.Unlock()
	return limiter.Allow()
}

func (r *RateLimitManager) start() {
	duration := time.Duration(r.requestsPerSecond*2) * time.Second
	t := time.NewTicker(duration)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case <-r.quit:
				fmt.Println("graceful shutdown")
				return
			case <-t.C:
				r.Lock()
				for id, client := range r.m {
					fmt.Println("clearing")
					if client.Elapsed(duration) {
						fmt.Println("clearing", id)
						delete(r.m, id)
					}
				}
				r.Unlock()
			}
		}
	}()
}

func (r *RateLimitManager) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	signal := make(chan struct{})
	go func() {
		close(r.quit)
		fmt.Println("closing done")
		r.wg.Wait()
		close(signal)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-signal:
		fmt.Println("signal received")
		return nil
	}
}

func main() {
	rl := NewRateLimitManager(5)
	defer rl.Clear()
	ip := "0.0.0.0"
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	fmt.Println(rl.Allow(ip))
	time.Sleep(1 * time.Second)
	fmt.Println(rl.Allow(ip))
	fmt.Println("sleep 15s")
	time.Sleep(20 * time.Second)
	fmt.Println(rl.Allow(ip))

	fmt.Println("shutting down")
}
```

## Leaky-bucket like rate limiter

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rl := &RateLimiter{
		n:      8,
		period: 5 * time.Second,
		cache:  make(map[string]RateLimitStat),
	}
	now := time.Now()
	n := 0
	for i := 0; i < 33; i++ {
		allow := rl.Allow("1")
		if allow {
			n++
		}
		fmt.Println(allow, time.Since(now))
		delay := time.Duration(100 + rand.Intn(100))
		time.Sleep(delay * time.Millisecond)
	}
	fmt.Println("total", n, "requests in", time.Since(now))
}

type RateLimitStat struct {
	count int64
	until time.Time
}

type RateLimiter struct {
	n      int64
	period time.Duration
	cache  map[string]RateLimitStat
	mu     sync.Mutex
}

func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	state := l.cache[key]
	if time.Now().Sub(state.until) > 0 {
		l.cache[key] = RateLimitStat{count: 1, until: time.Now().Add(l.period)}
		return true
	}

	rate := l.period.Microseconds() / l.n
	/*
		To visualize this, consider having a rate limiter with 5 requests per second.
		Each dash below represents 200 ms, which is the rate.
		At time >200ms and <400ms, we want to ensure that there can be max 2 counts of requests made.
		In short, the rate of requests is smoothen, similar to leaky bucket.
		
		     |
		     V
		0s - - - - - > 1s
		   1 2 3 4 5
	*/
	left := l.n - state.until.Sub(time.Now()).Microseconds()/rate
	if state.count < left {
		state.count++
		l.cache[key] = state
		return true
	}

	return false
}
```
## References

- https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
- https://medium.com/@saisandeepmopuri/system-design-rate-limiter-and-data-modelling-9304b0d18250
- https://hechao.li/2018/06/25/Rate-Limiter-Part1/

