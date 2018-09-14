For converting from a string to a byte slice, string -> []byte:
```go
[]byte(str)
```

For converting an array to a slice, [20]byte -> []byte:
```go
arr[:]
```

For copying a string to an array, string -> [20]byte:
```go
copy(arr[:], str)
```

Same as above, but explicitly converting the string to a slice first:
```go
copy(arr[:], []byte(str))
```

## String to Fixed byte

```go
func main() {
    s := "abc"
    var a [20]byte
    
    // Copy the string to a fixed byte
    copy(a[:], s)
    fmt.Println("s:", []byte(s), "a:", a)
}
```


## Quick Sorting

```go
func main() {
	dataset := make(map[string]interface{})
	dataset["zeta"] = "1"
	dataset["zzta"] = "2"
	dataset["zata"] = "3"
	dataset["cyta"] = "4"
	dataset["eeta"] = "5"
	dataset["1eta"] = "5"
	dataset["aleta"] = "5"
	dataset["ata"] = "5"

	sorted := make([][]byte, 0, len(dataset))
	for key := range dataset {
		sorted = append(sorted, []byte(key))
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if bytes.Compare(sorted[i][:], sorted[j][:]) > 0 {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	log.Println(sorted)
	for k, v := range sorted {
		log.Println(k, string(v))
	}
}
```

Output:
```bash
2018/07/23 11:00:11 0 1eta
2018/07/23 11:00:11 1 aleta
2018/07/23 11:00:11 2 ata
2018/07/23 11:00:11 3 cyta
2018/07/23 11:00:11 4 eeta
2018/07/23 11:00:11 5 zata
2018/07/23 11:00:11 6 zeta
2018/07/23 11:00:11 7 zzta
```

## Concatenating bytes

```
package main

import (
	"fmt"
)

func main() {
	b := append(append([]byte{}, []byte("hello ")...), []byte("world")...)
	fmt.Println(string(b))
	
}
```

Output:

```
hello world
```
