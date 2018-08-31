Sample code to measure the latency using LRU cache.

```golang
package models

import (
	"log"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	maxCache = 128
)

type Metric struct {
	Ping *lru.Cache

	// CachedPing is a cache of the latest average ping to avoid recomputing.
	CachedPing time.Duration
}

func NewMetric() *Metric {
	// Setup lru cache here
	ping, err := lru.New(maxCache)
	if err != nil {
		log.Fatal(err)
	}
	return &Metric{
		Ping:       ping,
		CachedPing: time.Duration(1 * time.Millisecond), // Populate fake value
	}
}

func (m *Metric) IncrPing(t time.Duration) {
	m.Ping.Add(t, nil)
	// Update the ping directly
	newPing := (m.CachedPing*time.Duration(m.Ping.Len()-1) + t) / time.Duration(m.Ping.Len())
	m.CachedPing = newPing
}

func (m *Metric) AveragePing() time.Duration {
	var total time.Duration
	for _, k := range m.Ping.Keys() {
		total += k.(time.Duration)
	}
	average := total / time.Duration(m.Ping.Len())

	// Update the local cache before returning
	m.CachedPing = average

	return average
}
```
