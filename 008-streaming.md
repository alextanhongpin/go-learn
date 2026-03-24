```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"sync"
	"time"
)

func main() {
	r, w := io.Pipe()
	r2, w2 := io.Pipe()
	a1 := New("foo", r, w2)
	a2 := New("bar", r2, os.Stdout)

	var wg sync.WaitGroup
	wg.Go(a1.Run)
	wg.Go(a2.Run)

	for range 3 {
		fmt.Fprint(w, "hello\n")
		time.Sleep(100 * time.Millisecond)
	}
	w.Close()
	wg.Wait()
	fmt.Println("done")
}

func Stdin(r io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				break
			}
		}
	}
}

type Agent struct {
	name   string
	stdin  io.Reader
	stdout io.WriteCloser
}

func New(name string, stdin io.Reader, stdout io.WriteCloser) *Agent {
	return &Agent{name, stdin, stdout}
}

func (a *Agent) Run() {
	for line := range Stdin(a.stdin) {
		// Debug
		if a.stdout != os.Stdout {
			fmt.Fprintf(os.Stdout, "%s: %s\n", a.name, line)
		}
		fmt.Fprintf(a.stdout, "%s: %s\n", a.name, line)
	}
	if a.stdout != os.Stdout {
		a.stdout.Close()
	}
}
```
