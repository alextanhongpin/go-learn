## Token TTL

```go
package main

import (
	"fmt"
	"time"
)

type Token struct {
	token     string
	createdAt time.Time
}

// Allows the duration to be customised. If we use a fixed value, just wrap the whole operation
// in another function.
func (t *Token) CheckExpire(ttl time.Duration) bool {
	return time.Since(t.createdAt) < ttl
}

const defaultTTL = 1 * time.Minute

func checkTokenExpire(t Token) bool {
	return t.CheckExpire(defaultTTL)
}

type TokenTTL struct {
	token     string
	createdAt time.Time
	// This is hardcoded (?). The advantage is we can define multiple TTL for different configuration.
	ttl time.Duration
}

func (t *TokenTTL) CheckExpire() bool {
	return time.Since(t.createdAt) < t.ttl
}

func main() {
	fmt.Println("Hello, playground")
}
```
