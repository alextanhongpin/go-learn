```go
package main

import "fmt"

type point struct {
	x    int
	y    int
	move int
}

func main() {
	start := point{0, 0, 0}
	target := point{3, 2, 0}
	pawns := []point{
		{2, 0, 0},
		{2, -1, 0},
		{2, 1, 0},
		{1, 1, 0},
		{1, 2, 0},
		{0, 2, 0},
		{-2, -2, 0},
	}
	move := knightMove(start, target, pawns)
	fmt.Println("took", move, "moves")
}

// - - - - - - - -
// - - - - 2 - - -
// - - 1 p p - X -
// - - - - p p - -
// - - - N - p - -
// - - - - - p - -
// - p - - - - - -
// - - - - - - - -

func knightMove(start point, tgt point, pawns []point) int {
	combinations := []point{
		{2, 1, 0},
		{2, -1, 0},
		{-2, -1, 0},
		{-2, 1, 0},
		{-1, 2, 0},
		{1, 2, 0},
		{-1, -2, 0},
		{1, -2, 0},
	}

	cache := make(map[int]map[int]bool)
	queue := []point{}
	queue = append(queue, start)

	for _, pwn := range pawns {
		if _, found := cache[pwn.x]; !found {
			cache[pwn.x] = make(map[int]bool)
			cache[pwn.x][pwn.y] = true
		} else {
			cache[pwn.x][pwn.y] = true
		}
	}

	var pos point
	for {
		pos, queue = queue[0], queue[1:]
		for _, com := range combinations {
			pos := pos
			pos.x += com.x
			pos.y += com.y

			if _, found := cache[pos.x]; !found {
				cache[pos.x] = make(map[int]bool)
				cache[pos.x][pos.y] = true
			} else {
				if cache[pos.x][pos.y] {
					continue
				}
			}
			pos.move++
			if pos.x == tgt.x && pos.y == tgt.y {
				return pos.move
			}
			queue = append(queue, pos)
		}
	}
}
```
