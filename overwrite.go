// This program demonstrates how to overwrite the json tag in the struct
package main

import (
	"encoding/json"
	"log"
)

type BadField struct {
	Name string `json:"NameString"`
}

type GoodField struct {
	*BadField
	BadName string `json:"NameString,omitempty"`
	Name    string `json:"name"`
}

func main() {
	b := BadField{"john.doe"}
	g := GoodField{
		BadField: &b,
		Name:     b.Name,
	}

	out, err := json.Marshal(g)
	if err != nil {
		log.Printf("error unmarshalling: %v\n", err)
	}
	log.Println(string(out))
}
