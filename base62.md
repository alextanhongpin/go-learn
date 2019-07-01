## Base62 

```go
package main

import (
	"fmt"
	"math/big"
)

var base62Chars = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

const Base62 int64 = 62

var base62Map map[rune]int

func init() {
	base62Map = make(map[rune]int)
	// Reverse lookup.
	for i, char := range base62Chars {
		base62Map[char] = i
	}
}

func main() {
	var i int64
	for i = 0; i < Base62; i += 1 {
		if decode(encode(i)) != i {
			fmt.Println(i)
		}
	}
	for i = 0; i < Base62; i += 1 {
		j := big.NewInt(i)
		if decodeBigInt(encodeBigInt(j)).Cmp(big.NewInt(i)) != 0 {
			fmt.Println(i)
		}
	}
	fmt.Println(encode(0))
	fmt.Println(encodeBigInt(big.NewInt(0)))
	fmt.Println("terminating")
}

func encode(i int64) string {
	var out []rune
	// Special handling.
	if i == 0 {
		return string(base62Chars[0])
	}
	for i > 0 {
		c := base62Chars[i%Base62]
		// Append character in reverse order.
		out = append([]rune{c}, out...)
		i /= Base62
	}
	return string(out)
}

func decode(str string) int64 {
	var result int64
	for _, s := range str {
		result += result*Base62 + int64(base62Map[s])
	}
	return result
}

func encodeBigInt(i *big.Int) string {
	var out []rune
	var zero = big.NewInt(0)
	if i.Cmp(zero) == 0 {
		return string(base62Chars[0])
	}

	var Base62BigInt = big.NewInt(62)
	for i.Cmp(zero) > 0 {
		_, mod := i.DivMod(i, Base62BigInt, Base62BigInt)
		c := base62Chars[mod.Int64()]
		// Append character in reverse order.
		out = append([]rune{c}, out...)
	}
	return string(out)

}

func decodeBigInt(str string) *big.Int {
	var Base62BigInt = big.NewInt(62)
	sum := big.NewInt(0)
	for _, s := range str {
		sum = sum.Add(sum.Mul(sum, Base62BigInt), big.NewInt(int64(base62Map[s])))
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
