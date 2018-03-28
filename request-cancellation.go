package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Alternative:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	// ctx, cancel := context.WithCancel(context.Background())
	// time.AfterFunc(2*time.Millisecond, cancel)
	defer cancel()

	req, err := http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))
}
