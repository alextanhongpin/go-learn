# Promises with Golang

See [here](https://github.com/chebyrash/promise/blob/master/promise.go)

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"

	"play.ground/promise"
)

func main() {
	p := promise.New(func(resolve func(int), reject func(error)) {
		time.Sleep(1 * time.Second)
		resolve(1)
	}).Then(func(n int) *promise.Promise[int] {
		time.Sleep(1 * time.Second)
		return promise.Resolve(n + 1)
	})

	n, err := p.Await()
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
-- go.mod --
module play.ground
-- promise/promise.go --
package promise

import (
	"errors"
	"sync"
)

var ErrNoResult = errors.New("no result")

type Promise[T any] struct {
	result  T
	err     error
	pending bool
	wg      sync.WaitGroup
	mu      sync.RWMutex
}

func New[T any](resolver Resolver[T]) *Promise[T] {
	p := &Promise[T]{pending: true}
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		resolver(p.resolve, p.reject)
	}()

	return p
}

func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()

	p.mu.RLock()
	result, err := p.result, p.err
	p.mu.RUnlock()

	return result, err
}

func (p *Promise[T]) Then(resolve func(T) *Promise[T]) *Promise[T] {
	t, err := p.Await()
	if err != nil {
		return p
	}
	return resolve(t)
}

func (p *Promise[T]) resolve(t T) {
	p.mu.Lock()
	p.result = t
	p.pending = false
	p.mu.Unlock()
}

func (p *Promise[T]) reject(err error) {
	p.mu.Lock()
	p.err = err
	p.pending = false
	p.mu.Unlock()
}

type Resolver[T any] func(resolve func(T), reject func(error))

func Resolve[T any](t T) *Promise[T] {
	return &Promise[T]{result: t}
}

func Reject[T any](err error) *Promise[T] {
	return &Promise[T]{err: err}
}
```
