# Return value is faster than return pointer
```go
package main_test

import (
	"log"
	"testing"
)

type Value struct {
	Name      string
	Age       int8
	IsMarried bool
}

func ReturnValue(name string, age int8) Value {
	return Value{name, age, false}
}

func ReturnPointer(name string, age int8) *Value {
	return &Value{name, age, false}
}

func BenchmarkReturnValue(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v := ReturnValue("john doe", 20)
		if v.Name != "john doe" {
			log.Fatal("not equal")
		}
		if v.Age != 20 {
			log.Fatal("not equal")
		}
		v.Name = "jane doe"
		if v.Name != "jane doe" {
			log.Fatal("not equal")
		}
	}
}

func BenchmarkReturnPointer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v := ReturnPointer("john doe", 20)
		if v.Name != "john doe" {
			log.Fatal("not equal")
		}
		if v.Age != 20 {
			log.Fatal("not equal")
		}
		v.Name = "jane doe"
		if v.Name != "jane doe" {
			log.Fatal("not equal")
		}
	}
}
```

Output:

```
:!go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/alextanhongpin/go-bench/map-int-vs-string
BenchmarkReturnValue-4          2000000000               1.41 ns/op
BenchmarkReturnPointer-4        500000000                3.21 ns/op
PASS
ok      github.com/alextanhongpin/go-bench/map-int-vs-string    4.892s
```

## Using integer as key for map is faster than string/rune.

