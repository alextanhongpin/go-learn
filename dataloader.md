## Naive Data Loader implementation


```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	Pending Status = iota
	Success
	Failed
)

type Result struct {
	val    interface{}
	err    error
	status Status
}

type Loader struct {
	wg            sync.WaitGroup
	mu            sync.Mutex
	cache         map[string]Result
	keys          []Request
	keyCh         chan Request
	batchDuration time.Duration
	done          chan bool
}

func NewLoader() *Loader {
	l := &Loader{
		cache:         make(map[string]Result),
		keys:          make([]Request, 0),
		keyCh:         make(chan Request),
		batchDuration: 16 * time.Millisecond,
		done:          make(chan bool),
	}
	go l.pool()
	return l
}

func (l *Loader) pool() {
	for {
		select {
		case <-time.After(l.batchDuration):
			l.batchFn()
		case key := <-l.keyCh:
			l.keys = append(l.keys, key)
		case <-l.done:
			return
		}
	}
}

func (l *Loader) batchFn() {
	// Implement batch find here.
	l.mu.Lock()
	keys := l.keys
	fmt.Println("batching", keys)
	l.keys = make([]Request, 0)
	l.mu.Unlock()

	// Get the unique keys (deduplicate)
	// Find results for unique key.
	// Group result by key
	// For each key, return the result

	for i, key := range keys {
		key.ch <- Result{
			val: Country{
				ID:   fmt.Sprint(i),
				Name: key.key,
			},
			status: Success,
		}
	}
}

func (l *Loader) Load(key string) (interface{}, error) {
	l.wg.Add(1)
	defer l.wg.Done()
	fmt.Println("load", key)

	l.mu.Lock()
	if result, exists := l.cache[key]; exists {
		if result.status != Pending {
			log.Println("cache hit", key)
			l.mu.Unlock()
			return result.val, result.err
		} else {
			log.Println("already fetching")
		}
	} else {
		log.Println("not exist", key)
		l.cache[key] = Result{status: Pending}
	}
	l.mu.Unlock()

	// Create a new channel to wait for the result.
	ch := make(chan Result)

	// Send for batching.
	l.keyCh <- Request{ch: ch, key: key}

	// Receive the result of the channel (should include error)
	out := <-ch

	// Cache result.
	l.mu.Lock()
	l.cache[key] = out
	l.mu.Unlock()

	return out.val, out.err
}

func (l *Loader) Close() {
	l.wg.Wait()
	close(l.done)
}

func main() {
	u, _ := findUser(context.Background())
	fmt.Println("user", u)

	users, err := findUsers(context.Background(), 5)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("users", users)
}

type Country struct {
	ID   string
	Name string
}
type User struct {
	CountryID string
	Country   Country
}

func preloadUserCountry(ctx context.Context, users []User) ([]User, error) {
	mapByCountryID := make(map[string]bool)
	for _, user := range users {
		mapByCountryID[user.CountryID] = true
	}
	var ids []string
	for id := range mapByCountryID {
		ids = append(ids, id)
	}

	countries, err := findCountryByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	countryByID := make(map[string]Country)
	for _, country := range countries {
		countryByID[country.ID] = country
	}

	for i := range users {
		users[i].Country = countryByID[users[i].CountryID]
	}

	return users, nil
}

func findUser(ctx context.Context) (User, error) {
	var u User
	users, err := preloadUserCountry(ctx, []User{u})
	if err != nil {
		return User{}, err
	}
	return users[0], nil
}

func findUsers(ctx context.Context, n int) ([]User, error) {
	users := make([]User, n)
	users[0].CountryID = "1"
	l := NewLoader()
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
			c, err := l.Load(users[i].CountryID + fmt.Sprint(i))
			if err != nil {
				return err
			}
			users[i].Country, _ = c.(Country)
			if i == 1 {
				fmt.Println("is it exiting?")
				return errors.New("bad")
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return users, nil
}

func findCountry(ctx context.Context, id string) (Country, error) {
	switch id {
	case "1":
		return Country{ID: id, Name: "Malaysia"}, nil
	default:
		return Country{ID: id, Name: "None"}, nil
	}
}

func findCountryByIDs(ctx context.Context, ids []string) ([]Country, error) {
	// In SQL, this will be an IN statement.
	countries := make([]Country, len(ids))
	var err error
	for i, id := range ids {
		countries[i], err = findCountry(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return countries, nil
}
```
