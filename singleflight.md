## Singleflight

With singleflight, we can prevent thundering herd to our local database/cache by ensuring the ongoing requests are not duplicated. Consider the scenario when the cache is not set, and there are multiple concurrent calls to the service. Instead of waiting for the cache to be filled first, there will be multiple expensive calls to the database/cache that can bring the system down. Singleflight prevents this by "queuing" the concurrent requests and ensuring they are only called once. This will only work if the request is taking a noticable long time.

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	n := 10
	var wg sync.WaitGroup
	wg.Add(n)

	var group singleflight.Group
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			wrapper(&group)
		}()

	}
	wg.Wait()
	fmt.Println("terminating")
}

func wrapper(group *singleflight.Group) {
	// Read from cache here.
	v, err, shared := group.Do("hello", func() (interface{}, error) {
		fmt.Println("singleflight start")
		return expensiveWork()
	})
	fmt.Println(v, err, shared)
	// Set cache here.
}

func expensiveWork() (interface{}, error) {
	fmt.Println("starting work")
	// Expensive operation.
	time.Sleep(1 * time.Second)
	fmt.Println("completing work")
	return 1, nil
}

```

https://rodaine.com/2018/08/x-files-sync-golang/
https://rodaine.com/2017/05/x-files-time-rate-golang/
https://topic.alibabacloud.com/a/golang-singleflight_1_38_30919329.html
