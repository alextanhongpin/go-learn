// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	g := Graph[int]{
		0: []int{1, 2},
		2: []int{4},
		3: []int{4},
		4: []int{1},
		//5: []int{4},
	}
	fmt.Println("Hello, 世界", g.TopologicalSort())
	fmt.Println(g.Invert())
}

type Graph[T comparable] map[T][]T

func (g Graph[T]) Invert() Graph[T] {
	o := make(Graph[T])
	for node, edges := range g {
		if _, ok := o[node]; !ok {
			o[node] = []T{}
		}
		for _, edge := range edges {
			o[edge] = append(o[edge], node)
		}
	}
	return o
}

func (g Graph[T]) TopologicalSort() []T {
	inDegrees := make(map[T]int)
	vertices := make(map[T]bool)
	for node, edges := range g {
		vertices[node] = true
		for _, edge := range edges {
			inDegrees[edge]++
			vertices[edge] = true
		}
	}
	numVertices := len(vertices)
	var queue []T
	for node := range vertices {
		if inDegrees[node] == 0 {
			queue = append(queue, node)
		}
	}
	var sorted []T
	for len(queue) > 0 {
		var h T
		h, queue = queue[0], queue[1:]
		numVertices--
		sorted = append(sorted, h)
		for _, node := range g[h] {
			inDegrees[node]--
			if inDegrees[node] == 0 {
				queue = append(queue, node)
			}
		}
	}
	if numVertices == 0 {
		return sorted
	}
	return nil
}
