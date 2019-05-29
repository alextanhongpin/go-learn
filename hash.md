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
