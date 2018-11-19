# Example of repeating job at interval with golang.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func cleanup(i int) {
	if i > 5 {
		fmt.Println("Exiting")
		return
	}
	fmt.Println("performing cleanup")
	time.AfterFunc(time.Duration(1*time.Second), func() {
		cleanup(i + 1)
	})
}

func cleanup2() func() {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(2 * time.Second)
	// Seems like without the quit channel, the goroutines might "leak",
	// since it is not closed correctly.
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("graceful shutdown")
				close(quit)
				return
			case <-ticker.C:
				fmt.Println("Executing")
			}
		}
	}()
	return func() {
		cancel()
		<-quit
		fmt.Println("terminating cleanup2")
	}
}

func main() {
	fmt.Println("Hello, playground")
	go cleanup(1)

	done := cleanup2()
	defer done()
	time.Sleep(10 * time.Second)
	fmt.Println("program end")
}
```
