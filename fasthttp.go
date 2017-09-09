package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "hello world")
}
func main() {
	log.Println("listening to port *:8080. Press ctrl + c to cancel.")
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}

// wrk -d30s -c10 -t5 http://localhost:8080
// Running 30s test @ http://localhost:8080
//   5 threads and 10 connections
//   Thread Stats   Avg      Stdev     Max   +/- Stdev
//     Latency   176.46us    1.03ms  76.04ms   99.16%
//     Req/Sec    16.38k     1.84k   31.33k    86.67%
//   2445516 requests in 30.01s, 340.50MB read
// Requests/sec:  81479.45
// Transfer/sec:     11.34MB
