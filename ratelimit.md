# Ratelimit

Basic rate-limiting with go.

```go
package ratelimit

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limiter interface {
	Add(id string) *rate.Limiter
	Get(id string) *rate.Limiter
	Start()
	Stop()
}

type visitor struct {
	limiter   *rate.Limiter
	updatedAt time.Time
}

type IPLimiter struct {
	sync.RWMutex
	ips map[string]*visitor

	callsPerSecond   int
	burst            int
	cleanupDuration  time.Duration // Frequency to cleanup the user's ip.
	inactiveDuration time.Duration // If the visitor's last seen is above the inactive duration, the visitor will be removed.
	quit             chan struct{}
}

func NewIPLimiter(callsPerSecond, burst int, cleanup, inactive time.Duration) *IPLimiter {
	return &IPLimiter{
		callsPerSecond:   callsPerSecond, // e.g. 2 API calls per second.
		burst:            burst,          // e.g. 5 bursts allowed.
		ips:              make(map[string]*visitor),
		cleanupDuration:  cleanup,  // E.g. 5 minute to run the cron every 5 minutes to check for inactive IP.
		inactiveDuration: inactive, // E.g. 5 minute to remove users that are no longer active for 5 minutes.
		quit:             make(chan struct{}),
	}
}

func (i *IPLimiter) Add(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(i.callsPerSecond), i.burst)
	i.Lock()
	i.ips[ip] = &visitor{limiter, time.Now()}
	i.Unlock()
	return limiter
}

func (i *IPLimiter) Get(ip string) *rate.Limiter {
	i.Lock()
	v, exist := i.ips[ip]
	if !exist {
		i.Unlock()
		return i.Add(ip)
	}
	v.updatedAt = time.Now()
	i.Unlock()
	return v.limiter
}

func (i *IPLimiter) cleanup() {
	ticker := time.NewTicker(i.cleanupDuration)
	for {
		select {
		case <-ticker.C:
			i.Lock()
			for ip, v := range i.ips {
				if time.Since(v.updatedAt) > i.inactiveDuration {
					log.Println("ip", ip, "is removed")
					delete(i.ips, ip)
				}
			}
			i.Unlock()
		case <-i.quit:
			return
		}
	}
}

func (i *IPLimiter) Start() {
	log.Println("rateLimiter: start")
	go i.cleanup()
}

func (i *IPLimiter) Stop() {
	log.Println("rateLimiter: stop")
	close(i.quit)
}
```
