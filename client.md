## Client API adapter
```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Message string `json:"message"`
}

func main() {
	b, err := json.Marshal(&Request{
		Message: "hello world",
	})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequest("POST", "http://url", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	// Disable keep-alive, to ensure that sockets/file descriptors do not run out.
	req.Close = true
	resp, err := client.Do(req)
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	// When there is a redirection failure, both variables will be non-nil.
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	fmt.Println("Hello, playground")
}

type BooksClientAPI interface {
	FetchBooks()
}

type ClientAdapter struct {
	client *http.Client
}

func (c *ClientAdapter) FetchBooks() {
}

// TODO: Add middleware retries, error
// https://medium.com/@nitishkr88/http-retries-in-go-e622e51d249f
```
