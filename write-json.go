// This program demonstrates how to write a struct to json file

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Point represents the schema of our json output
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func writeJSON(file string, obj interface{}, pretty bool) (err error) {
	var bytes []byte
	if pretty {
		bytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = json.Marshal(obj)
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bytes, 0644)
}

func main() {
	points := []Point{Point{0, 0}, Point{1, 1}}
	err := writeJSON("out.json", points, true)
	if err != nil {
		log.Fatalf("error writing to json: %v\n", err)
	}
}
