## Sample in-memory expiring cache for golang

```golang
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cache := NewCache(1 * time.Second)
	defer cache.Stop()
	cache.Set("hello", "world", 200*time.Millisecond)
	cache.Set("foo", "bar", 1*time.Second)
	if res := cache.Get("foo"); res != nil {
		fmt.Println(res.(string))
	}
	fmt.Println("program terminating")
}

type CacheItem struct {
	createdAt time.Time
	value     interface{}
	ttl       time.Duration
}
type Cache struct {
	sync.RWMutex
	sync.Once
	sync.WaitGroup
	data   map[string]*CacheItem
	quitCh chan interface{}
}

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		quitCh: make(chan interface{}),
		data:   make(map[string]*CacheItem),
	}
	cache.WaitGroup.Add(1)
	go cache.loop(duration)
	return cache
}

func (c *Cache) Stop() {
	c.Once.Do(func() {
		close(c.quitCh)
		c.WaitGroup.Wait()
	})
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	c.data[key] = &CacheItem{value: value, createdAt: time.Now(), ttl: ttl}
	c.Unlock()
}

func (c *Cache) Get(key string) interface{} {
	c.RLock()
	data, exist := c.data[key]
	c.RUnlock()
	if !exist {
		return nil
	}
	if time.Since(data.createdAt) > data.ttl {
		c.Lock()
		delete(c.data, key)
		c.Unlock()
		return nil
	}
	return data.value
}

func (c *Cache) loop(duration time.Duration) {
	t := time.NewTicker(duration)
	defer t.Stop()
	defer c.WaitGroup.Done()
	for {
		select {
		case <-c.quitCh:
			fmt.Println("quitting")
			return
		case <-t.C:
			c.Lock()
			for key, data := range c.data {
				if time.Since(data.createdAt) > data.ttl {
					fmt.Println("clearing", key)
					delete(c.data, key)
				}
			}
			c.Unlock()
		}
	}
}
```
