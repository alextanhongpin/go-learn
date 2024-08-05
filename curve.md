# Curve

Normally for retry mechanism, we use exponential growth for retrying the failed requests. However, there are some scenarios that requires the initial exponential growth, then decay.

An example is a process that waits for a value to be populated, e.g. cache. When one process is working on populating the cache, other processes will wait and retry. However, the retry period should decrease over time because the chances of finding the item in the cache becomes higher.

To visualize how this chart works, paste this `\ln\left(x+1\right)\exp\left(-0.01x^{2}\right)` in https://www.desmos.com/calculator.

<img width="1455" alt="image" src="https://github.com/user-attachments/assets/4ae0b644-db5f-4eb0-a9dc-656ca1ade371">


The end value should be clipped and set to 0 and the retry terminated when it becomes too small.

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
)

func main() {
	for i := range 100 {
		fmt.Println(exponentialGrowthAndDecay(float64(i)))
	}
	fmt.Println("Hello, 世界")
}

func exponentialGrowthAndDecay(x float64) float64 {
	// ln(x+1) * exp(-0.01x^2)
	// ln = 2.303 * log
	return 2.303 * math.Log(x+1) * math.Exp(-0.01*x*x)
}

```
