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
	defer c.WaitGroup.Done()
	
	t := time.NewTicker(duration)
	defer t.Stop()
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

## With ExpireAt

```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrItemDoesNotExist = errors.New("item does not exist")
	ErrItemExpired      = errors.New("item expired")
)

type Item struct {
	ExpireAt time.Time
	Key      string
	Value    interface{}
}

func NewItem(key string, value interface{}, ttl time.Duration) Item {
	return Item{
		ExpireAt: time.Now().Add(ttl),
		Key:      key,
		Value:    value,
	}
}

type Cache struct {
	sync.RWMutex
	sync.Once
	sync.WaitGroup
	quit  chan interface{}
	items map[string]Item
}

func NewCache(cleanupPeriod time.Duration) *Cache {
	cache := &Cache{
		quit:  make(chan interface{}),
		items: make(map[string]Item),
	}
	cache.WaitGroup.Add(1)
	go cache.worker(cleanupPeriod)
	return cache
}

func (c *Cache) worker(cleanupPeriod time.Duration) {
	// Defer sequence matter. The t.Stop will be called before c.WaitGroup.
	defer c.WaitGroup.Done()

	t := time.NewTicker(cleanupPeriod)
	defer t.Stop()
	for {
		select {
		case <-c.quit:
			fmt.Println("quit")
			return
		case <-t.C:
			fmt.Println("cleanup")
			c.Lock()
			for key, item := range c.items {
				if item.ExpireAt.Before(time.Now()) {
					fmt.Println("clear", key)
					delete(c.items, key)
				}
			}
			c.Unlock()

		}
	}
}

func (c *Cache) Stop() {
	c.Once.Do(func() {
		close(c.quit)
		c.Wait()
		fmt.Println("terminated")
	})
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.Lock()
	c.items[key] = NewItem(key, value, duration)
	c.Unlock()
}

func (c *Cache) Get(key string) (*Item, error) {
	c.RLock()
	item, exist := c.items[key]
	c.RUnlock()
	if !exist {
		return nil, ErrItemDoesNotExist
	}

	if item.ExpireAt.Before(time.Now()) {
		c.Delete(key)
		return nil, ErrItemExpired
	}
	return &item, nil
}

func (c *Cache) Delete(key string) {
	c.Lock()
	delete(c.items, key)
	c.Unlock()
}

func main() {
	defer fmt.Println("one")
	defer fmt.Println("two")

	cleanupEvery := 1 * time.Second
	cache := NewCache(cleanupEvery)
	defer cache.Stop()
	cache.Set("hello", "world", 1*time.Second)
	val, err := cache.Get("hello")
	fmt.Println(val, err)

	time.Sleep(5 * time.Second)
	val, err = cache.Get("hello")
	fmt.Println(val)

	fmt.Println("Hello, playground")
}
```

## Generic Cache

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

func main() {
	c := new(Cache)
	ctx := context.Background()
	fmt.Println(CacheDecorator("users", c, ListUsers, ctx, &UserFilter{Name: "john"}))
	fmt.Println(CacheDecorator("users", c, ListUsers, ctx, &UserFilter{Name: "john"}))
	fmt.Println(CacheDecorator("users", c, ShowUser, ctx, "1"))
	fmt.Println(CacheDecorator("users", c, ShowUser, ctx, "1"))
	fmt.Println(CacheDecorator("users", c, ShowUser, ctx, "2"))
}

type Cacheable[T any, U any] func(ctx context.Context, req T) (U, error)

type Cache struct {
	sync.Map
}

func CacheDecorator[T any, U any](prefix string, cache *Cache, fn Cacheable[T, U], ctx context.Context, req T) (u U, err error) {
	key := fmt.Sprintf("%s:%#v", prefix, req)

	v, found := cache.Load(key)
	if found {
		fmt.Println("cache hit", key)
		if res, ok := v.(U); ok {
			return res, nil
		}
		err = fmt.Errorf("cast error: key=%s, want=%T, got=%T", key, u, v)
		return
	}

	val, err := fn(ctx, req)
	if err != nil {
		return u, err
	}
	v, loaded := cache.LoadOrStore(key, val)
	if loaded {
		if res, ok := v.(U); ok {
			return res, nil
		}

		err = fmt.Errorf("cast error: key=%s, want=%T, got=%T", key, u, v)
		return
	}

	return val, nil
}

type UserFilter struct {
	Name string
}

type User struct {
	Age  int
	Name string
}

func ListUsers(ctx context.Context, filter *UserFilter) ([]User, error) {
	return nil, nil
}

func ShowUser(ctx context.Context, id string) (*User, error) {
	switch id {
	case "1":
		return &User{Name: "john", Age: 20}, nil
	case "2":
		return &User{Name: "jane", Age: 30}, nil
	}

	return nil, errors.New("user not found")
}
```
