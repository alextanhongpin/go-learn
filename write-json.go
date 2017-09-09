// This program demonstrates how to write a struct to json file

package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"
)

// Point represents the schema of our json output
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// func writeJSONAlt(file string, obj interface{}, pretty bool) (err error) {
// 	var bytes []byte
// 	if pretty {
// 		bytes, err = json.MarshalIndent(obj, "", "  ")
// 	} else {
// 		bytes, err = json.Marshal(obj)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(file, bytes, 0644)
// }

func writeJSON(w io.WriteCloser, obj interface{}, pretty bool) (err error) {
	var bytes []byte
	if pretty {
		bytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = json.Marshal(obj)
	}
	if err != nil {
		return err
	}
	w.Write(bytes)
	return
}

func writeXML(w io.WriteCloser, obj interface{}, pretty bool) (err error) {
	var bytes []byte
	if pretty {
		bytes, err = xml.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = xml.Marshal(obj)
	}
	if err != nil {
		return err
	}
	w.Write(bytes)
	return
}

func main() {
	points := []Point{Point{0, 0}, Point{1, 1}}
	var w io.WriteCloser
	var err error

	runningInPlayGround := os.Getenv("user") == ""
	if !runningInPlayGround {
		w, err = os.Create("myfile.txt")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		w = os.Stdout
	}
	writeJSON(w, points, true)
	writeXML(w, points, true)
	w.Close()
}
