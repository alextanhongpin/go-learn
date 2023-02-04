```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrTooManyRequests = errors.New("too many requests")

type Task func() (interface{}, error)

type State struct {
	sync.RWMutex
	successCounter   int
	failureCounter   int
	failureThreshold int
	successThreshold int
	timer            time.Time
	timeout          time.Duration
}

func (s *State) StartTimeoutTimer() {
	s.Lock()
	s.timer = time.Now()
	s.Unlock()
}

func (s *State) IncrementSuccessCounter() {
	s.Lock()
	s.successCounter++
	s.Unlock()
}

func (s *State) ResetSuccessCounter() {
	s.Lock()
	s.successCounter = 0
	s.Unlock()
}

func (s *State) IncrementFailureCounter() {
	s.Lock()
	s.failureCounter++
	s.Unlock()
}

func (s *State) ResetFailureCounter() {
	s.Lock()
	s.failureCounter = 0
	s.Unlock()
}

func (s *State) IsTimeoutTimerExpired() bool {
	s.RLock()
	timer, timeout := s.timer, s.timeout
	s.RUnlock()
	return time.Since(timer) > timeout
}

func (s *State) IsFailureThresholdExceeded() bool {
	s.RLock()
	failureCounter, failureThreshold := s.failureCounter, s.failureThreshold
	s.RUnlock()
	return failureCounter > failureThreshold
}

func (s *State) IsSuccessThresholdExceeded() bool {
	s.RLock()
	successCounter, successThreshold := s.successCounter, s.successThreshold
	s.RUnlock()
	return successCounter > successThreshold
}

type CircuitBreaker interface {
	Next() CircuitBreaker
	Handle(Task) (interface{}, error)
}

type Closed struct {
	state *State
}

func NewClosed(state *State) *Closed {
	// entry/reset failure counter
	state.ResetFailureCounter()
	return &Closed{state}
}

func (c *Closed) Next() CircuitBreaker {
	// failure threshold reached
	if c.state.IsFailureThresholdExceeded() {
		fmt.Println("is opened")
		return NewOpened(c.state)
	}
	return c
}

func (c *Closed) Handle(task Task) (interface{}, error) {
	// do/	if operation succeeds
	// 		return result
	// 	else
	// 		increment failure counter
	//		return failure
	res, err := task()
	if err != nil {
		c.state.IncrementFailureCounter()
		return nil, err
	}
	return res, nil
}

type Opened struct {
	state *State
}

func NewOpened(state *State) *Opened {
	// entry/ start timeout timer
	state.StartTimeoutTimer()
	return &Opened{state}
}

func (o *Opened) Next() CircuitBreaker {
	// timeout timer expired
	if o.state.IsTimeoutTimerExpired() {
		fmt.Println("is half-opened")
		return NewHalfOpened(o.state)
	}
	return o
}

func (o *Opened) Handle(task Task) (interface{}, error) {
	// do /return failure
	return nil, ErrTooManyRequests
}

type HalfOpened struct {
	state  *State
	failed int32
}

func NewHalfOpened(state *State) *HalfOpened {
	// entry/ reset success counter
	state.ResetSuccessCounter()
	return &HalfOpened{state: state}
}

func (h *HalfOpened) Handle(task Task) (interface{}, error) {
	// do/	if operation succeeds
	// 		increment success counter
	// 		return result
	// 	else
	// 		return failure
	res, err := task()
	if err != nil {
		atomic.CompareAndSwapInt32(&h.failed, 0, 1)
		return nil, err
	}
	atomic.CompareAndSwapInt32(&h.failed, 1, 0)
	h.state.IncrementSuccessCounter()
	return res, err
}

func (h *HalfOpened) Next() CircuitBreaker {
	// success count threshold reached
	if h.state.IsSuccessThresholdExceeded() {
		fmt.Println("success count threshold reached")
		fmt.Println("is closed")
		return NewClosed(h.state)
	}
	if atomic.LoadInt32(&h.failed) == 1 {
		// operation failed
		fmt.Println("operation failed")
		return NewOpened(h.state)
	}

	return h
}

type CircuitBreakerImpl struct {
	CircuitBreaker
}

func NewDefaultState() *State {
	return &State{
		successCounter:   5,
		successThreshold: 5,
		failureCounter:   5,
		failureThreshold: 5,
		timeout:          5 * time.Second,
	}
}
func NewCircuitBreaker(state *State) *CircuitBreakerImpl {
	if state == nil {
		state = NewDefaultState()
	}
	cb := NewClosed(state)
	return &CircuitBreakerImpl{cb}
}

func (c *CircuitBreakerImpl) Handle(task Task) (interface{}, error) {
	c.CircuitBreaker = c.Next()
	return c.CircuitBreaker.Handle(task)
}

func main() {
	state := NewDefaultState()
	state.timeout = 1 * time.Second
	cb := NewCircuitBreaker(state)
	for i := 0; i < 10; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("some error")
		})
		fmt.Println(res, err)
	}
	fmt.Println("sleep 1,1 seconds")
	time.Sleep(1100 * time.Millisecond)

	for i := 0; i < 3; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("another error")
		})
		fmt.Println(res, err)
	}

	fmt.Println("sleep 1.1 seconds")
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 15; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return true, nil
		})
		fmt.Println(res, err)
	}
	for i := 0; i < 20; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("some error")
		})
		fmt.Println(res, err)
	}
}
```

## Circuit Breaker Rewrite

```go
package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var ErrTooManyRequests = errors.New("too many requests")

type State struct {
	failureCount          int64
	failureCountThreshold int64
	successCount          int64
	successCountThreshold int64
	timeoutTimer          time.Time
	timeoutDuration       time.Duration
	sync.RWMutex
}

func (s *State) ResetFailureCount() {
	s.Lock()
	s.failureCount = 0
	s.Unlock()
}

func (s *State) IncrementFailureCount() {
	s.Lock()
	s.failureCount += 1
	s.Unlock()
}

func (s *State) StartTimeoutTimer() {
	s.Lock()
	s.timeoutTimer = time.Now()
	s.Unlock()
}

func (s *State) ResetSuccessCount() {
	s.Lock()
	s.successCount = 0
	s.Unlock()
}

func (s *State) IncrementSuccessCount() {
	s.Lock()
	s.successCount += 1
	s.Unlock()
}

func (s *State) IsSuccessCountThresholdReached() bool {
	s.RLock()
	count, threshold := s.successCount, s.successCountThreshold
	s.RUnlock()
	return count >= threshold
}

func (s *State) IsFailureCountThresholdReached() bool {
	s.RLock()
	count, threshold := s.failureCount, s.failureCountThreshold
	s.RUnlock()
	return count >= threshold
}

func (s *State) IsTimeoutTimerExpired() bool {
	return time.Since(s.timeoutTimer) > s.timeoutDuration
}

func NewState(
	successCountThreshold,
	failureCountThreshold int64,
	timeoutDuration time.Duration,
) *State {
	return &State{
		successCountThreshold: successCountThreshold,
		failureCountThreshold: failureCountThreshold,
		timeoutDuration:       timeoutDuration,
	}
}

type CircuitBreaker interface {
	Next() CircuitBreaker
	Do(func() error) error
}

type Closed struct {
	state *State
	name  string // To identify the state.
}

func NewClosed(state *State) *Closed {
	state.ResetFailureCount()
	return &Closed{state: state}
}

func (c *Closed) Do(fn func() error) error {
	err := fn()
	if err != nil {
		c.state.IncrementFailureCount()
		return err
	}
	return nil
}

func (c *Closed) Next() CircuitBreaker {
	if c.state.IsFailureCountThresholdReached() {
		return NewOpen(c.state)
	}
	return c
}

type Open struct {
	state *State
}

func NewOpen(state *State) *Open {
	state.StartTimeoutTimer()
	return &Open{state: state}
}

func (o *Open) Do(fn func() error) error {
	return ErrTooManyRequests
}

func (o *Open) Next() CircuitBreaker {
	if o.state.IsTimeoutTimerExpired() {
		return NewHalfOpen(o.state)
	}
	return o
}

type HalfOpen struct {
	// error error
	error int64
	state *State
}

func NewHalfOpen(state *State) *HalfOpen {
	state.ResetSuccessCount()
	return &HalfOpen{state: state}
}

func (h *HalfOpen) Do(fn func() error) error {
	err := fn()
	if err != nil {
		atomic.CompareAndSwapInt64(&h.error, 0, 1)
		// h.error = err
		return err
	}
	h.state.IncrementSuccessCount()
	// h.error = nil
	atomic.CompareAndSwapInt64(&h.error, 1, 0)
	return nil
}

func (h *HalfOpen) Next() CircuitBreaker {
	if atomic.LoadInt64(&h.error) == 1 {
		// if h.error != nil {
		return NewOpen(h.state)
	}
	if h.state.IsSuccessCountThresholdReached() {
		return NewClosed(h.state)
	}
	return h
}

type CircuitBreakerImpl struct {
	state CircuitBreaker
}

func NewCircuitBreaker(successCountThreshold, failureCountThreshold int64, timeoutDuration time.Duration) *CircuitBreakerImpl {
	state := NewState(successCountThreshold, failureCountThreshold, timeoutDuration)
	return &CircuitBreakerImpl{
		state: NewClosed(state),
	}
}

func (c *CircuitBreakerImpl) Do(fn func() error) error {
	c.state = c.state.Next()
	// Print current state.
	fmt.Println("[CircuitBreakerState]:", reflect.TypeOf(c.state).String())
	return c.state.Do(fn)
}

func main() {
	var (
		successThreshold int64 = 5
		failureThreshold int64 = 5
		timeoutDuration        = 1 * time.Second
	)
	cb := NewCircuitBreaker(successThreshold, failureThreshold, timeoutDuration)

	// Trigger error 6 times.
	for i := 0; i < 6; i += 1 {
		err := cb.Do(func() error { return errors.New("bad request") })
		if err != nil {
			fmt.Println(err)
		}
	}

	// Sleep for 1 second to recover.
	time.Sleep(2 * time.Second)
	// Trigger success 6 times.
	for i := 0; i < 6; i += 1 {
		err := cb.Do(func() error { return nil })
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Hello, playground")
}
```


## Rewrite again

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type mockHandler struct {
	err error
}

func (m *mockHandler) Handle(ctx context.Context) error {
	return m.err
}

func main() {
	m := &mockHandler{}
	cb := NewCircuitBreaker(m, SuccessThreshold(1), FailureThreshold(1), TimeoutFn(func(clock Clock) time.Time {
		return time.Now().Add(1 * time.Second)
	}))
	ctx := context.Background()
	fmt.Println(cb.Handle(ctx))

	m.err = errors.New("bad job")
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))

	fmt.Println(cb.Handle(ctx))
	time.Sleep(2 * time.Second)

	m.err = nil
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	m.err = errors.New("something else")
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println("why")
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	time.Sleep(2 * time.Second)
	m.err = nil
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
	fmt.Println(cb.Handle(ctx))
}

var ErrUnavailable = errors.New("circuit: unavailable")

type Handler interface {
	Handle(ctx context.Context) error
}

type Clock interface {
	Now() time.Time
}

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type circuitBreakerOption struct {
	Clock            Clock
	SuccessThreshold int
	FailureThreshold int
	TimeoutFn        func(clock Clock) time.Time
}

func newDefaultOption() *circuitBreakerOption {
	return &circuitBreakerOption{
		Clock:            &clock{},
		SuccessThreshold: 3,
		FailureThreshold: 3,
		TimeoutFn: func(c Clock) time.Time {
			return c.Now().Add(1 * time.Second)
		},
	}
}

type CircuitBreakerOption func(*circuitBreakerOption)

func CustomClock(clock Clock) CircuitBreakerOption {
	return func(o *circuitBreakerOption) {
		o.Clock = clock
	}
}

func SuccessThreshold(n int) CircuitBreakerOption {
	return func(o *circuitBreakerOption) {
		o.SuccessThreshold = n
	}
}

func FailureThreshold(n int) CircuitBreakerOption {
	return func(o *circuitBreakerOption) {
		o.FailureThreshold = n
	}
}

func TimeoutFn(fn func(Clock) time.Time) CircuitBreakerOption {
	return func(o *circuitBreakerOption) {
		o.TimeoutFn = fn
	}
}

type CircuitBreaker struct {
	open     *Open
	closed   *Close
	halfOpen *HalfOpen
	h        Handler
	option   *circuitBreakerOption
}
type clock struct{}

func (c *clock) Now() time.Time { return time.Now() }

func NewCircuitBreaker(h Handler, options ...CircuitBreakerOption) *CircuitBreaker {
	opt := newDefaultOption()
	for _, o := range options {
		o(opt)
	}

	return &CircuitBreaker{
		h:      h,
		closed: &Close{h: h, failureThreshold: opt.FailureThreshold},
		option: opt,
	}
}

func (cb *CircuitBreaker) Valid() bool {
	var n int
	switch {
	case cb.open != nil:
		n++
	case cb.closed != nil:
		n++
	case cb.halfOpen != nil:
		n++
	}
	return n == 1
}

func (cb *CircuitBreaker) State() State {
	if !cb.Valid() {
		panic("circuit breaker: invalid state")
	}

	switch {
	case cb.open != nil:
		return StateOpen
	case cb.closed != nil:
		return StateClosed
	case cb.halfOpen != nil:
		return StateHalfOpen
	default:
		panic("circuit breaker: invalid state")
	}
}

func (cb *CircuitBreaker) Handle(ctx context.Context) error {
	switch cb.State() {
	case StateOpen:
		if cb.open.Next() {
			cb.open = nil
			cb.halfOpen = &HalfOpen{h: cb.h, successThreshold: cb.option.SuccessThreshold}
			fmt.Println("open -> half-open")
			return cb.halfOpen.Handle(ctx)
		}

		return cb.open.Handle(ctx)
	case StateHalfOpen:
		if err := cb.halfOpen.Handle(ctx); err != nil {
			cb.halfOpen = nil
			cb.open = &Open{deadline: cb.option.TimeoutFn(cb.option.Clock)}
			fmt.Println("half-open -> open")
			return err
		}

		if cb.halfOpen.Next() {
			cb.halfOpen = nil
			cb.closed = &Close{h: cb.h, failureThreshold: cb.option.FailureThreshold}
			fmt.Println("half-open -> closed")
		}

		return nil
	case StateClosed:
		if err := cb.closed.Handle(ctx); err != nil {
			if cb.closed.Next() {
				cb.closed = nil
				cb.open = &Open{
					clock:    cb.option.Clock,
					deadline: cb.option.TimeoutFn(cb.option.Clock),
				}
				fmt.Println("closed -> open")
			}
			return err
		}
		return nil
	default:
		panic("circuit breaker: invalid state")
	}

}

type Open struct {
	clock    Clock
	deadline time.Time
}

func (o *Open) Next() bool {
	return o.clock.Now().After(o.deadline)
}

func (o *Open) Handle(ctx context.Context) error {
	return ErrUnavailable
}

type HalfOpen struct {
	successCounter   int
	successThreshold int
	h                Handler
}

func (o *HalfOpen) Next() bool {
	return o.successCounter >= o.successThreshold
}

func (o *HalfOpen) Handle(ctx context.Context) error {
	if err := o.h.Handle(ctx); err != nil {
		return err
	}
	o.successCounter++
	return nil
}

type Close struct {
	failureCounter   int
	failureThreshold int
	h                Handler
}

func (o *Close) Next() bool {
	return o.failureCounter >= o.failureThreshold
}

func (o *Close) Handle(ctx context.Context) error {
	if err := o.h.Handle(ctx); err != nil {
		o.failureCounter++
		return err
	}
	return nil
}
```


## How to implement a thread safe and distributed cb?

### Idea 1
Simple, at the start of the request, load the cb state from the distributed store.

Each state should have a version, e.g. `20230101-<version>`. 

Everytime the state changes, increment the version and save the state in the store, only if the version is greater than the existing store version.

However, if there are concurrent failures, the version can only be updated incrementally.


### Idea 2

Event sourcing. Everytime there is a state change, append to a list of states. Take the majority of the last n states (or last n states based on the circuit breaker timeout) to decide the current state. Prefer to be pessimistic over optimistic, so take the worst possible states first.

If one of the state is opened, then all of them are opened, and can only make new request after the given time. Once they are half-opened, we can perhaps average the number of times the state appeared as the success threshold...
