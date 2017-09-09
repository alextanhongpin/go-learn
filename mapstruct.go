package main

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var people map[string]interface{}

func main() {
	// Using Hashicorp's library to convert map to struct
	// Example struct
	people = make(map[string]interface{})
	people["name"] = "car" // Lowercase works
	people["Age"] = 1

	peeps := People{}
	err := mapstructure.Decode(people, &peeps)
	if err != nil {
		log.Println(err)
	}
	log.Printf("peeps: %#v\n", peeps)
}
