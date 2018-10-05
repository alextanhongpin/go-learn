## Base62 

```go
package main

import (
	"fmt"
	"math"
)

// Base62 character set, [a-zA-Z0-9].
var base62Chars = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

var base62Map map[rune]int

const Base62 int64 = 62

func init() {
	base62Map = make(map[rune]int)
	for k, v := range base62Chars {
		base62Map[v] = k
	}
}

func main() {
	hash := encode(100)
	fmt.Println(hash)
	fmt.Println(decode(hash))
}

func encode(in int64) string {
	var out []rune
	for in > 0 {
		out = append([]rune{base62Chars[in%62]}, out...)
		in /= 62
	}
	return string(out)
}

func decode(in string) int64 {
	var sum int64
	for i, v := range in {
		pow := int64(len(in) - 1 - i)
		p := math.Pow(float64(Base62), float64(pow))
		sum += int64(p) * int64(base62Map[v])
	}
	return sum
}
```

## Getting rune value as int

```go
func main() {
	// Getting rune as int.
	var total int
  // Note that "hello world" and "helol world" will return the same total.
	for _, v := range []rune("helol world") {
		total += int(v - '0')
	}
	fmt.Println("total", total)
}
```

## Base62 with big.Int

```go
package main

import (
	"fmt"
	"math/big"
)

// Base62 character set, [a-zA-Z0-9].
var base62Chars = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

var base62Map map[rune]int

var Base62 = big.NewInt(62)

func init() {
	base62Map = make(map[rune]int)
	for k, v := range base62Chars {
		base62Map[v] = k
	}
}

func main() {
	m := big.NewInt(100)
	hash := encode(m)
	fmt.Println(hash)
	fmt.Println(decode(hash))
}

func encode(in *big.Int) string {
	var out []rune
	for in.Cmp(big.NewInt(0)) > 0 {
		_, mod := in.DivMod(in, Base62, Base62)
		out = append([]rune{base62Chars[int(mod.Int64())]}, out...)
	}
	return string(out)
}

func decode(in string) *big.Int {
	sum := big.NewInt(0)
	for i, v := range in {
		pow := int64(len(in) - 1 - i)
		base62 := big.NewInt(62)
		base62.Exp(base62, big.NewInt(pow), nil)
		val := big.NewInt(int64(base62Map[v]))
		sum.Add(sum, base62.Mul(base62, val))
	}
	return sum
}
```
