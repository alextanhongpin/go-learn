package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println(time.Now())
	log.Println(time.Now().String())
	log.Println(time.Now().UTC())
	log.Println(time.Now().UTC().String())
	log.Println(time.Now().Unix())
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, time.Now().UTC().String())
	if err != nil {
		log.Println(err)
	}
	log.Println(t)
	fmt.Println(time.Now().Format("2006.01.02"))
}
