# Load balancer

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	{
		lb := &LoadBalancerNode{
			threshold: 3,
			period:    1 * time.Second,
		}
		for i := 0; i < 5; i++ {
			err := lb.Do(func() error {
				return errors.New("bad error")
			})
			fmt.Println(i, err)
		}
		time.Sleep(1 * time.Second)

		for i := 0; i < 5; i++ {
			err := lb.Do(func() error {
				return errors.New("bad error")
			})
			fmt.Println(i, err)
		}
	}

	lb := LoadBalancer{
		nodes: []*LoadBalancerNode{
			&LoadBalancerNode{
				threshold: 3,
				period:    1 * time.Second,
			},
			&LoadBalancerNode{
				threshold: 10,
				period:    1 * time.Second,
			},
		},
	}
	for i := 0; i < 5; i++ {
		err := lb.Do(func() error {
			if i == 4 {
				return nil
			}
			return errors.New("bad error")
		})
		fmt.Println(i, err)
		fmt.Println()
	}
	fmt.Println("Hello, 世界")

}

type LoadBalancer struct {
	nodes []*LoadBalancerNode
}

func (l *LoadBalancer) Do(fn func() error) error {
	for i := range l.nodes {
		fmt.Println("trying node", i)
		if err := l.nodes[i].Do(fn); err == nil {
			return nil
		}
		fmt.Println("node", i, "unhealthy")
	}
	return errors.New("all unhealthy")
}

type LoadBalancerNode struct {
	count     int
	refreshAt time.Time

	threshold int
	period    time.Duration
}

func (l *LoadBalancerNode) Reset() {
	l.refreshAt = time.Now()
	l.count = 0
}

func (l *LoadBalancerNode) Incr() {
	l.count++
}

func (l *LoadBalancerNode) Allow() bool {
	return l.count < l.threshold
}

func (l *LoadBalancerNode) Do(fn func() error) error {
	if time.Now().After(l.refreshAt) {
		l.Reset()
	}
	if !l.Allow() {
		return errors.New("not allowed")
	}
	if err := fn(); err != nil {
		l.Incr()
		return err
	}
	return nil
}

```
