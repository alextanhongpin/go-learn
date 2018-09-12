Simple logic to calculate moving average using a LRU cache:

```
package models

import (
	"log"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	maxCache = 1000
)

// Metric represents a metric object to store latency.
type Metric struct {
	// Value is a cache of the latest average ping duration to avoid recomputing.
	Value time.Duration

	// Cache stores a list of ping latency.
	Cache *lru.Cache
	// Transactions stores a list of transaction latency.
	// CachedTransaction is a cache of the latest average transaction duration.
}

func NewMetric() *Metric {
	// Setup lru cache here
	ping, err := lru.New(maxCache)
	if err != nil {
		log.Fatal(err)
	}
	return &Metric{
		Cache: ping,
		Value: time.Duration(1 * time.Millisecond), // Populate fake value
	}
}

// Incr the ping latency into the list and updates the cached ping.
func (m *Metric) Incr(t time.Duration) {
	m.Cache.Add(t, nil)
	// Update the ping directly
	newCache := (m.Value*time.Duration(m.Cache.Len()-1) + t) / time.Duration(m.Cache.Len())
	m.Value = newCache
}

// Count recomputes the average ping duration and updates the cache before returning them.
func (m *Metric) Count() time.Duration {
	var total time.Duration
	for _, k := range m.Cache.Keys() {
		total += k.(time.Duration)
	}
	average := total / time.Duration(m.Cache.Len())

	// Update the local cache before returning
	m.Value = average

	return average
}

```
