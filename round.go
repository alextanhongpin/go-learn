// This program demonstrates how to round numbers in golang

package main

import (
	"log"
	"math"
)

func round(val float64) float64 {
	_, v := math.Modf(val)
	if math.Abs(v) >= .5 {
		if val >= 0 {
			return math.Ceil(val)
		}
		return math.Floor(val)
	}
	if val >= 0 {
		return math.Floor(val)
	}
	return math.Ceil(val)
}

func main() {
	log.Println(round(123.54))
	log.Println(round(-100.4))
	log.Println(round(-0.4))
	log.Println(round(-0.5))
	log.Println(round(0.5))
	log.Println(round(0.4))
}
