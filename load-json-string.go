package main

import (
	"io/ioutil"
	"log"
)

func main() {
	out, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
}
