## Convert MD5 to big.Int

Useful for consistent hashing etc.

```go
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	hash := md5.New()
	hash.Write([]byte("hello world"))
	hashed := hash.Sum(nil)
	hexStr := hex.EncodeToString(hashed)

	i := big.NewInt(0)
	i.SetString(hexStr, 16)
	// This will overflow.
	// i, err := strconv.ParseUint(hexStr, 16, 32)
	fmt.Println(i.String())
}
```

## FNV Hash

Wiki [here](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function). This is used by go's kafka library [here](https://github.com/segmentio/kafka-go/blob/eba9cae7fd57401a8078f6f25bc00cb753cd4f42/balancer.go#L127).

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
)

func main() {
	cache := make(map[uint32]bool)
	for i := 0; i < 26; i++ {
		r := string('a' + rune(i))
		h := fnv.New32()
		h.Write([]byte(r))
		i32 := h.Sum32()
		hash := hex.EncodeToString(h.Sum(nil))
		if cache[i32] {
			fmt.Println("duplicate")
		}
		cache[i32] = true
		fmt.Println(i, r, i32, hash)
	}
}

```
