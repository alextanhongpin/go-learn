Good example at the official godocs: https://godoc.org/golang.org/x/sync/semaphore.

```
package main

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Semaphore struct {
	lock uint32
}

func (s *Semaphore) Lock() bool {
	return atomic.CompareAndSwapUint32(&s.lock, 0, 1)
}

func (s *Semaphore) Unlock() bool {
	return atomic.CompareAndSwapUint32(&s.lock, 1, 0)
}

var s Semaphore

func doWork(s *Semaphore) {
	if !s.Lock() {
		log.Println("pending...")
		return
	}
	defer s.Unlock()

	log.Println("start")
	dur := rand.Intn(350) + 150
	time.Sleep(time.Duration(dur) * time.Millisecond)
	log.Println("stop")

}
func main() {
	var ( 
		wg  sync.WaitGroup
		num = 10
	)
	wg.Add(num)

	for i := 0; i < num; i++ {

		go func() {
			defer wg.Done()
			dur := rand.Intn(750) + 250
			time.Sleep(time.Duration(dur) * time.Millisecond)
			doWork(&s)
		}()
	}

	wg.Wait()
	log.Println("terminating...")

}

```
