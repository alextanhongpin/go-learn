package main

import (
	"log"
	"net/http"
	"time"
)

// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
func main() {
	// The default listenAndServe does not enforce timeout, and should
	// be avoided in production
	// http.ListenAndServe(":8080")
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("listening to port *:8080. press ctrl + c to cancel.")
	log.Fatal(srv.ListenAndServe(":8080"))

	// Enforcing timeouts on client connections helps prevent leaking file descriptors:
	// http: Accept error: accept tcp [::]:80: accept4: too many open files; retrying in 5ms
}
