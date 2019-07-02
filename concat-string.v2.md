```go
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// Before go 1.10.
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	fmt.Println(b.String())

	// This is faster and uses less memory.
	var sb strings.Builder
	sb.WriteString("hello")
	sb.WriteString(" ")
	sb.WriteString("world")
	fmt.Println(sb.String())

}

```
