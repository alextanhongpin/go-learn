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
