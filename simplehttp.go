package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	log.Println("listening to port *:8080. Press ctrl + c to cancel.")
	http.ListenAndServe(":8080", nil)
}

// Running 30s test @ http://localhost:8080
//   5 threads and 10 connections
//   Thread Stats   Avg      Stdev     Max   +/- Stdev
//     Latency   487.78us    4.89ms 150.47ms   99.25%
//     Req/Sec    11.97k     1.50k   17.80k    84.20%
//   1787208 requests in 30.08s, 218.17MB read
// Requests/sec:  59419.98
// Transfer/sec:      7.25MB
