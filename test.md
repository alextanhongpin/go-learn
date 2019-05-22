Setting up test in go playground.

```go
package main

import (
	"fmt"
	"testing"
)

func TestCaseA(t *testing.T) {
	fmt.Println("Hello tester")
	t.Fail()
}

func matchString(a, b string) (bool, error) {
	return a == b, nil
}

func main() {
	testSuite := []testing.InternalTest{
		{
			Name: "TestCaseA",
			F:    TestCaseA,
		},
	}
	testing.Main(matchString, testSuite, nil, nil)
}
```
