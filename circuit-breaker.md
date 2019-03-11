```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task func() (interface{}, error)

type State struct {
	sync.RWMutex
	successCounter   int
	failureCounter   int
	failureThreshold int
	successThreshold int
	breakTime        time.Time
	timeout          time.Duration
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

func (s *State) IsTimeoutExpired() bool {
	return time.Since(s.breakTime) > s.timeout
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
	// do/ 	if operation succeeds
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
	state.breakTime = time.Now()
	return &Opened{state}
}

func (o *Opened) Next() CircuitBreaker {
	// timeout timer expired
	if o.state.IsTimeoutExpired() {
		fmt.Println("is half-opened")
		return NewHalfOpened(o.state)
	}
	return o
}

func (o *Opened) Handle(task Task) (interface{}, error) {
	// do /return failure
	return nil, errors.New("timeout")
}

type HalfOpened struct {
	state  *State
	failed int32
}

func NewHalfOpened(state *State) *HalfOpened {
	// entry/ reset success counter
	state.ResetSuccessCounter()
	return &HalfOpened{state, 0}
}

func (h *HalfOpened) Handle(task Task) (interface{}, error) {
	// do/ 	if operation succeeds
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
	ctx CircuitBreaker
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
	res, err := c.ctx.Handle(task)
	if err != nil {
		c.ctx = c.ctx.Next()
		return nil, err
	}
	c.ctx = c.ctx.Next()
	return res, err
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
}
```
