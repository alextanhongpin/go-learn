## Naive Data Loader implementation


```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Request struct {
	ch  chan Result
	key string
}

type Status int

const (
	None Status = iota
	Pending
	Success
	Failed
)

type Result struct {
	key    string
	val    interface{}
	err    error
	status Status
}

type Loader struct {
	mu            sync.RWMutex
	wg            sync.WaitGroup
	batchFn       BatchFn
	cond          *sync.Cond
	cache         map[string]Result
	batchDuration time.Duration
	done          chan bool
}

type BatchFn func(ctx context.Context, keys []string) ([]Result, error)

func NewLoader(batchFn BatchFn) *Loader {
	l := &Loader{
		batchFn:       batchFn,
		batchDuration: 16 * time.Millisecond,
		cache:         make(map[string]Result),
		cond:          sync.NewCond(&sync.Mutex{}),
		done:          make(chan bool),
	}
	go l.pool()
	return l
}

func (l *Loader) pool() {
	for {
		select {
		case <-time.After(l.batchDuration):
			l.batch()
		case <-l.done:
			return
		}
	}
}

func (l *Loader) batch() {
	// Implement batch find here.
	l.mu.Lock()
	defer l.mu.Unlock()

	var keys []string
	for key := range l.cache {
		result := l.cache[key]
		if result.status == Pending {
			keys = append(keys, key)
		}
	}
	fmt.Println("batching", keys)

	items, err := l.batchFn(context.Background(), keys)
	if err != nil {
		l.cond.L.Lock()
		for _, key := range keys {
			l.cache[key] = Result{
				key:    key,
				val:    nil,
				err:    err,
				status: Failed,
			}
		}
		log.Println("broadcasting result 1")
		l.cond.Broadcast()
		l.cond.L.Unlock()
		return
	}

	itemByID := make(map[string]Result)
	for _, item := range items {
		itemByID[item.key] = item
	}

	l.cond.L.Lock()
	for _, key := range keys {
		item, exists := itemByID[key]
		if exists {
			l.cache[key] = item
		} else {
			l.cache[key] = Result{
				status: Failed,
				err:    errors.New("does not exists"),
				key:    key,
			}
		}
	}
	log.Println("broadcasting result 2")
	l.cond.Broadcast()
	l.cond.L.Unlock()
}

func (l *Loader) Load(key string) (interface{}, error) {
	l.wg.Add(1)
	defer l.wg.Done()

	fmt.Println("Load", key)

	// First, check if the key exists.
	//l.mu.Lock()
	//defer l.mu.Unlock()
	l.mu.RLock()
	result, exists := l.cache[key]
	l.mu.RUnlock()

	condition := func() bool {
		l.mu.RLock()
		result := l.cache[key]
		l.mu.RUnlock()
		return result.status == Pending
	}

	if exists {
		if result.status == Success && result.status == Failed {
			log.Println("cache hit true for", key)
			return result.val, result.err
		}

		// L must be locked when waiting for the status to change.
		l.cond.L.Lock()

		log.Println("alraedy fetching, pending success")
		for condition() {
			l.cond.Wait()
		}
		// Obtain the latest status.
		result = l.cache[key]
		l.cond.L.Unlock()
		return result.val, result.err
	}

	l.mu.Lock()
	if _, exists := l.cache[key]; !exists {
		l.cache[key] = Result{status: Pending}
	}
	l.mu.Unlock()

	l.cond.L.Lock()
	for condition() {
		log.Println("waiting for status pending to change now...", l.cache[key])
		l.cond.Wait()
	}

	l.mu.RLock()
	result = l.cache[key]
	l.mu.RUnlock()

	l.cond.L.Unlock()
	fmt.Printf("Got result: %+v\n", result)

	return result.val, result.err
}

func (l *Loader) Close() {
	l.wg.Wait()
	close(l.done)
}

func main() {
	users, err := findUsers(context.Background(), 5)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("users", users)
}

type Country struct {
	ID   int64
	Name string
}
type User struct {
	CountryID int64
	Country   Country
}

func findUsers(ctx context.Context, n int) ([]User, error) {
	users := make([]User, n)
	l := NewLoader(func(ctx context.Context, keys []string) ([]Result, error) {
		countries, err := findCountryByIDs(ctx, keys)
		if err != nil {
			return nil, err
		}
		result := make([]Result, len(keys))
		for i, country := range countries {
			result[i] = Result{
				val:    country,
				key:    fmt.Sprint(country.ID),
				err:    nil,
				status: Success,
			}
		}
		return result, nil
	})
	defer l.Close()

	g := new(errgroup.Group)
	for i := range users {
		// Important, closure.
		i := i
		g.Go(func() error {
			// This will demonstrate splitting into two batches.
			var duration time.Duration
			if i < 3 {
				duration = 10 * time.Millisecond
			} else {
				duration = 40 * time.Millisecond
			}
			time.Sleep(duration)
			var n int
			if i < 3 {
				n = 1
			} else {
				n = 2
			}
			c, err := l.Load(fmt.Sprint(n))
			if err != nil {
				return err
			}
			users[i].Country, _ = c.(Country)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return users, nil
}

func findCountry(ctx context.Context, id int64) (Country, error) {
	switch id {
	case 1:
		return Country{ID: id, Name: "Malaysia"}, nil
	case 2:
		return Country{ID: id, Name: "Singapore"}, nil
	case 3:
		return Country{ID: id, Name: "Japan"}, nil
	default:
		return Country{ID: id, Name: "None"}, nil
	}
}

func findCountryByIDs(ctx context.Context, ids []string) ([]Country, error) {
	// In SQL, this will be an IN statement.
	countries := make([]Country, len(ids))

	for i, id := range ids {
		countryID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			panic(err)
		}
		countries[i], err = findCountry(ctx, countryID)
		if err != nil {
			return nil, err
		}
	}
	return countries, nil
}
```
