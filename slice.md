```go
var b [64]byte
copy(b[:], "hello world")

// b is now [64]byte

// This converts array b to slice b
fmt.Println(b[:])
```
