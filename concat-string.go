// This program demonstrates on how to concat string performantly in go
package main

import (
	"bytes"
	"log"
)

func main() {

	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString("_")
	b.WriteString("world")
	log.Println(b.String())

}
