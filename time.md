
## Zero time

This program demonstrates how to handle zero time
```go
package main

import (
	"log"
	"time"
)

type Book struct {
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func main() {

	var b Book
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}

	log.Println(b)
}
```

## Local time
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	singapore := time.FixedZone("Singapore Time", secondsEastOfUTC)
	fmt.Println(time.Now(), time.Now().In(singapore).Format(time.RFC3339))
}
```
