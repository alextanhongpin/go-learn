# Basic hashing, encoding and encryption

```go
package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"log"
	"strconv"
)

func main() {
	a := adler32.New()
	a.Write([]byte("1"))
	fmt.Printf("adler32: %x\n", a.Sum(nil))

	m := make(map[uint32]struct{})
	for i := 0; i < 1; i++ {
		c := crc32.Checksum([]byte(strconv.Itoa(i)), crc32.MakeTable(0xD5828281))
		fmt.Printf("crc32: %x\n", c)
		if _, found := m[c]; found {
			log.Fatal("collision")
		}
		m[c] = struct{}{}
	}

	algos := []struct {
		hash func() hash.Hash
		name string
	}{
		{md5.New, "md5"},
		{sha1.New, "sha1"},
		{sha256.New, "sha256"},
	}
	for _, algo := range algos {
		fmt.Println(algo.name+":", hasher(algo.hash, []byte(""), []byte("")))
	}
}

func hasher(algo func() hash.Hash, key, message []byte) string {
	mac := hmac.New(algo, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(expectedMAC)
}
```
