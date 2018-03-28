package main

import (
	"net"
	"net/http"
	"time"
)

func main() {
	// c := &http.Client{
	// 	Timeout: 15 * time.Second,
	// }
	// resp, err := c.Get("http://www.google.com")

	c := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second, // limits the time spent establishing a TCP connection (if a new one is needed).
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second, // limits the time spent performing the TLS handshake.
			ResponseHeaderTimeout: 10 * time.Second, // limits the time spent reading the headers of the response.
			ExpectContinueTimeout: 1 * time.Second,  // limits the time the client will wait between sending the request headers when including an Expect: 100-continue and receiving the go-ahead to send the body
		},
	}
}
