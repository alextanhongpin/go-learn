package main

import (
	"encoding/json"
	"os"
)

const max = 20

type Bin struct {
	Items []int
	Space int
}

func main() {
	items := []int{20, 1, 3, 5, 12, 17, 7}
	bins := make([]*Bin, 0)

	for _, item := range items {
		var done bool
		for _, bin := range bins {
			if bin.Space >= item {
				bin.Items = append(bin.Items, item)
				bin.Space -= item
				done = true
				break
			}
		}

		if !done {
			bins = append(bins, &Bin{
				Items: []int{item},
				Space: max - item,
			})
		}
	}

	json.NewEncoder(os.Stdout).Encode(bins)
}
