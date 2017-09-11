// This program demonstrates how to concat two arrays in go

package main

import "log"

func main() {

	a := []string{"a", "b"}
	b := []string{"1", "2"}

	log.Println(append(a, b...))
}
