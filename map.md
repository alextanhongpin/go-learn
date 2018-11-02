# Selecting random items from map
```
package main

import (
	"math/rand"
	"testing"
)

var m = map[string]int{"s": 18, "v": 21, "n": 13, "q": 16, "l": 11, "r": 17, "z": 25, "h": 7, "i": 8, "k": 10, "m": 12, "o": 14, "p": 15, "t": 19, "u": 20, "d": 3, "e": 4, "w": 22, "x": 23, "c": 2, "f": 5, "g": 6, "j": 9, "y": 24, "a": 0, "b": 1}

func BenchmarkRandomMapLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomMapLoop()
	}
}
func BenchmarkRandomMapSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomSliceLoop()
	}
}

func randomMapLoop() {
	scores := make(map[string]int)
	for i := 0; i < 100; i++ {
		n := rand.Intn(len(m))
		var key string
		for key = range m {
			if n == 0 {
				break
			}
			n--
		}
		scores[key]++
	}
}

func randomSliceLoop() {
	scores := make(map[string]int)
	// var keys []string
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	for i := 0; i < 100; i++ {
		scores[keys[rand.Intn(len(m))]]++
	}
}
```

Output:

```
:!go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/alextanhongpin/balancer
BenchmarkRandomMapLoop-4           50000             36690 ns/op            1506 B/op          2 allocs/op
BenchmarkRandomMapSlice-4         100000             12307 ns/op            1921 B/op          3 allocs/op
PASS
ok      github.com/alextanhongpin/balancer      3.554s
```
