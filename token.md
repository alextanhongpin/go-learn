## Generate random token string

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	token := base64.URLEncoding.EncodeToString(buf)
	fmt.Println(token, len(token))
}
```
