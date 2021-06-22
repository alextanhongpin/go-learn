## How to output UUID as a short string

Context: When doing cursor pagination in GraphQL, I want to show the uuid as base64 encoded string using uuid. However, the UUID is generally too long.

Some other options:
- using hashid, but the output is still too long

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New()
	fmt.Println(id) // 4e77153e-8194-4e3d-99df-a721ac539ea6
	b, err := id.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b)

	// There's a difference between using []byte(uuid) vs binary uuid.
	out := base64.RawURLEncoding.EncodeToString([]byte(id.String()))
	fmt.Println(out) // NGU3NzE1M2UtODE5NC00ZTNkLTk5ZGYtYTcyMWFjNTM5ZWE2

	out = base64.RawURLEncoding.EncodeToString(b)
	fmt.Println(out) // TncVPoGUTj2Z36chrFOepg
}

```
This is the same:

```go
package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	//	id := uuid.New()
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")
	fmt.Println(id)
	b, _ := id.MarshalBinary()
	fmt.Println(base64.StdEncoding.EncodeToString(b))
	
	// Same as...
	b, _ = hex.DecodeString(strings.ReplaceAll(id.String(), "-", ""))
	fmt.Println(fmt.Sprintf("%x", b))
	fmt.Println(base64.StdEncoding.EncodeToString(b))
}
```

## Shorter cursor

```go
package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	t := time.Now().UnixNano() // 1257894000000000000
	s := fmt.Sprint(t)
	fmt.Println(s)
	// Good. Shorter hex.
	b, _ := hex.DecodeString(s)
	fmt.Println(base64.StdEncoding.EncodeToString(b)) // EleJQAAAAAAA
	
	// Bad. Longer string.
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))) // MTI1Nzg5NDAwMDAwMDAwMDAwMA==

}
```

References:
- https://stackoverflow.com/questions/37934162/output-uuid-in-go-as-a-short-string
