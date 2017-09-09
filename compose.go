// This program demonstrates how to compose multiple struct to be returned in the json
package main

import (
	"encoding/json"
	"log"
)

type User struct {
	Email    string `json:"email"`
}

type Skill struct {
	Name string `json:"name"`
	Level int `json:"level"`
}

type Skills []Skill

func main() {

	usr := User{"john.doe@mail.com"}
	skills := Skills{Skill{"javascript", 1}, {"go", 2}}

	// Convert our composed anyonymous struct to bytes
	out, err := json.Marshal(struct {
		*User
		*Skills `json:"skills"`
	} {
		User: &usr,
		Skills: &skills,
	})
	
	if err != nil {
		log.Println(err)
	}
	log.Printf("with shadowing: %s\n", string(out))
}
