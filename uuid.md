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

References:
- https://stackoverflow.com/questions/37934162/output-uuid-in-go-as-a-short-string
