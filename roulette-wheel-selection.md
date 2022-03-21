Usecases
- distributing rewards - we don't want to give the highest tier equally to just about everyone. 

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// This is the rewards that can be received by the participants, e.g. cash amount in SGD
	rewards := []int64{20, 40, 60, 80, 100}

	// When the distribution is equal, users that spins the wheel will receive an equal chance to win any reward amount.
	equalDistributions := []int64{1, 1, 1, 1, 1}

	// When the distribution is weighted, users that spins the wheel has a higher change to hit the reward with the higher weight.
	// So in this example, users are more likely to earn 40 SGD than 100 SGD.
	weightedDistributions := []int64{15, 20, 15, 2, 1}

	printRewardDistribution(rewards, equalDistributions)
	printRewardDistribution(rewards, weightedDistributions)
}

func printRewardDistribution(rewards, distributions []int64) {
	fmt.Println("Rewards:", rewards)
	fmt.Println("Distributions:", distributions)
	w := NewRouletteWheel(rewards, distributions)
	counter := make(map[int64]int64)
	for i := 0; i < 10000; i++ {
		got := w.Choice()
		counter[got]++
	}
	var totalReward int64
	for _, r := range rewards {
		fmt.Println(r, "SGD:", counter[r])
		totalReward += counter[r] * r
	}
	fmt.Println("Total reward:", totalReward)
	fmt.Println()
}

type RouletteWheel struct {
	value         []int64
	distributions []int64
}

func NewRouletteWheel(value []int64, distributions []int64) *RouletteWheel {
	if len(value) != len(distributions) {
		panic("mismatch args length")
	}

	for i := range distributions[1:] {
		distributions[i+1] += distributions[i]
	}
	return &RouletteWheel{
		value:         value,
		distributions: distributions,
	}
}

func (rw *RouletteWheel) Choice() int64 {
	n := rand.Int63n(rw.distributions[len(rw.distributions)-1])
	for i, d := range rw.distributions {
		if n < d {
			return rw.value[i]
		}
	}
	return rw.value[len(rw.value)-1]
}

```


References
- https://rocreguant.com/roulette-wheel-selection-python/2019/
- https://www.keithschwarz.com/darts-dice-coins/
- https://en.wikipedia.org/wiki/Fitness_proportionate_selection
- https://cybernetist.com/2019/01/24/random-weighted-draws-in-go/
- https://dev.to/trekhleb/weighted-random-algorithm-in-javascript-1pdc
