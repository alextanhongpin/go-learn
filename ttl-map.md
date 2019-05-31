## Redis-like TTL capability implemented using map

TODO: Select only 20% of the map if the size exceeds >100k. Then do checking at get. Also investigate the penalty for time.After vs unix timestamp.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type item struct {
	value      string
	lastAccess int64
}

type TTLMap struct {
	sync.RWMutex
	data map[string]*item
	ttl  int64
	once sync.Once
	wg   sync.WaitGroup
	quit chan interface{}
}

func NewTTLMap(ttl int64) *TTLMap {
	return &TTLMap{
		data: make(map[string]*item),
		quit: make(chan interface{}),
		ttl:  ttl,
	}
}

func (t *TTLMap) Cleanup(every time.Duration) func(context.Context) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.cleanup(every)
	}()
	return func(ctx context.Context) {
		done := make(chan interface{})
		t.once.Do(func() {
			close(t.quit)
			t.wg.Wait()
			close(done)
		})

		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}
}

func (t *TTLMap) Len() int {
	t.RLock()
	l := len(t.data)
	t.RUnlock()
	return l
}

func (t *TTLMap) Put(key, value string) {
	t.Lock()
	it, ok := t.data[key]
	if !ok {
		it = &item{value: value}
		t.data[key] = it
	}
	it.lastAccess = time.Now().Unix()
	t.Unlock()
}

func (t *TTLMap) Get(key string) (value string) {
	t.Lock()
	if it, ok := t.data[key]; ok {
		value = it.value
		it.lastAccess = time.Now().Unix()
	}
	t.Unlock()
	return
}

func (t *TTLMap) cleanup(every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-t.quit:
			return
		case <-ticker.C:
			t.Lock()
			for k, v := range t.data {
				if time.Now().Unix()-v.lastAccess > t.ttl {
					// log.Println("deleted", k)
					delete(t.data, k)
				}
			}
			t.Unlock()
		}
	}
}
func main() {
	ttlMap := NewTTLMap(int64(time.Second.Seconds()))
	shutdown := ttlMap.Cleanup(time.Second)

	add := func(n int) {
		for i := 0; i < n; i++ {
			val := rand.Intn(100000)
			s := strconv.FormatInt(int64(val), 10)
			ttlMap.Put(s, s)
		}
		log.Println(ttlMap.Len())
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(i*950) * time.Millisecond)
			n := 200 + rand.Intn(800)
			add(n)
		}(i)
	}

	done := make(chan interface{})
	go func() {
		time.Sleep(10 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown(ctx)
		close(done)
	}()
	<-done

	log.Println(ttlMap.Len())
	fmt.Println("done")
}
```

## Modulo time

Both results will give the correct modulo time. This is useful for determining the time window and can be used for time series bucket assignment.
```
package main

import (
	"fmt"
	"log"
	"testing/quick"
)

func main() {
	f := func(a, b int64) bool {
		return one(a, b) == two(a, b)
	}
	if err := quick.Check(f, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("program terminating")
}

func one(ts, duration int64) int64 {
	return ts / duration * duration
}

func two(ts, duration int64) int64 {
	return ts - (ts % duration)
}
```
