```go
package main

import (
	"context"
	"fmt"
)

type key string

func main() {
	ctx := context.Background()
	var auth key = "auth"

	// Storing the value.
	ctx = context.WithValue(ctx, auth, "123456")

	// Getting the value.
	{

		val, ok := ctx.Value(auth).(string)
		fmt.Println(val, ok)
	}
	
	// Getting with the incorrect type assertion.
	{
		val, ok := ctx.Value(auth).(bool)
		fmt.Println(val, ok)
	}

	// Getting an unknown key.
	{
		var unknown key = "unknown"
		val, ok := ctx.Value(unknown).(string)
		fmt.Println(val, ok)
	}
}
```
