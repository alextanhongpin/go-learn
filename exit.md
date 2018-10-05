The defer statement will not be executed, since os.Exit will abruptly exit the program. This is common mistake when you want to perform graceful shutdown.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("Hello, playground")
	os.Exit(0)
}
```
