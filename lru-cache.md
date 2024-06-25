```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"container/list"
	"fmt"
)

func main() {
	lru := New(3)
	lru.Add("hallo", "world")
	lru.Add("one", "one")
	lru.Add("two", "two")
	lru.Add("one", "one")

	lru.Print()
	fmt.Println()

	lru.Add("three", "three")
	lru.Print()
	fmt.Println()

	lru.Add("one", "one")
	lru.Add("one", "one")
	lru.Add("one", "one")
	lru.Print()
	fmt.Println()

	lru.Add("three", "three")
	lru.Add("three", "three")
	lru.Add("two", "two")
	lru.Print()
	fmt.Println()
}

func New(cap int) *LRUCache {
	return &LRUCache{
		ll:  list.New(),
		kv:  make(map[string]*list.Element),
		cap: cap,
	}
}

type LRUCache struct {
	ll  *list.List
	kv  map[string]*list.Element
	cap int
}

func (c *LRUCache) Add(key string, val any) {
	el, ok := c.kv[key]
	if ok {
		c.ll.Remove(el)
	}
	c.kv[key] = c.ll.PushBack(val)
	if c.ll.Len() > c.cap {
		c.ll.Remove(c.ll.Front())
	}
}

func (c *LRUCache) Get(key string) (any, bool) {
	el, ok := c.kv[key]
	if ok {
		val := c.ll.Remove(el)
		c.ll.PushBack(val)
	}
	return el.Value, ok
}

func (c *LRUCache) Print() {
	// Iterate through list and print its contents.
	for e := c.ll.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
```
