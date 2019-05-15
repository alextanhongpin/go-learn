## Graceful Shutdown Pattern

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Shutdown func(context.Context)

type Container struct {
	supervisor []Shutdown
}

func NewContainer() *Container {
	return &Container{
		supervisor: make([]Shutdown, 0),
	}
}

func (c *Container) AddShutdown(fn Shutdown) {
	c.supervisor = append(c.supervisor, fn)
}

func (c *Container) Shutdown(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(c.supervisor))
	fmt.Println("shutting down")
	for _, shutdown := range c.supervisor {
		go func(fn Shutdown) {
			defer wg.Done()
			fn(ctx)
			fmt.Println("completed")
		}(shutdown)
	}
	wg.Wait()
}

func main() {
	con := NewContainer()
	con.AddShutdown(func(ctx context.Context) {
		signal := make(chan interface{})
		go func() {
			// Shutting down the operation in the goroutine and signalling them when it's done.
			time.Sleep(5 * time.Second)
			close(signal)
		}()

		select {
		case <-ctx.Done():
			fmt.Println("cancel")
			return
		case <-signal:
			return
		}

	})
	con.AddShutdown(func(context.Context) {
		time.Sleep(1 * time.Second)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con.Shutdown(ctx)
	fmt.Println("Hello, playground")
}
```
