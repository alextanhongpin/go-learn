There is an issue with this implementation where it will fail if the tokens length is greater than the `chunk_size`. It will get caught in an infinite loop.

We will address it by limiting the number of characters per token.

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := `this is a super long text that you cannot imagine also this is another thing that you cannot comprehend haha paper scissors rock`
	tokens := strings.Fields(text)
	chunk_size := 26
	chunk_overlap := 3
	batch := 0

	i := 0
	curr := 0
	m := make(map[int][]string)
	for i < len(tokens) {
		fmt.Println("i", i)
		if curr < chunk_size {
			curr += len(tokens[i])

			m[batch] = append(m[batch], tokens[i])

			i++
			//		fmt.Println("appending token", tokens[i])
		} else {
			batch += 1
			curr = 0
			for curr < chunk_overlap {
				i--
				curr += len(tokens[i])
			}
			fmt.Println("revert back to", i)
			curr = 0
		}
	}
	fmt.Println("Hello, 世界", m)
}

```
