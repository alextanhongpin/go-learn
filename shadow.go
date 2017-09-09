// This program demonstrates how to exclude fields in the json output
package main

import (
	"encoding/json"
	"log"
)

type UserPrivate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPublic struct {
	*UserPrivate
	Password bool `json:"password,omitempty"`
}

func main() {

	usrPriv := UserPrivate{"john.doe@mail.com", "123456"}
	usrPub := UserPublic{
		UserPrivate: &usrPriv,
	}
	// Convert it to bytes
	out, err := json.Marshal(usrPub)
	if err != nil {
		log.Println(err)
	}
	log.Printf("with shadowing: %s\n", string(out))
}
