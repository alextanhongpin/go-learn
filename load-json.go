package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func loadJSON(file string, obj interface{}) error {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &obj)
}

// Point represents the schema of the json we want to load
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	var points []Point
	if err := loadJSON("out.json", &points); err != nil {
		log.Printf("error loading json: %v", err)
	}
	log.Printf("load json: %#v\n", points)
}
