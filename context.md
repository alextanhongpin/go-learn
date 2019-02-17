```go
package main

import (
	"context"
	"fmt"
)

type contextKey string

const key = contextKey("key")

func main() {
	ctx := context.Background()
	{
		v, ok := ctx.Value(key).(string)
		fmt.Println(v, ok)

	}
	ctx = context.WithValue(ctx, key, "hello")
	{
		v, ok := ctx.Value(key).(string)
		fmt.Println(v, ok)
	}
}

```
