## Http Server Cancellation

In your `server.go`:

```go
	r.GET("/long", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()
		start := time.Now()

		select {
		case <-ctx.Done():
			stdlog.Println("cancelled after", time.Since(start))
			break
		case <-time.After(time.Second * 5):
			break
		}

		fmt.Fprintf(w, "completed after %v", time.Since(start))
	})
```

In your `client.go`:

```go
package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Result struct {
	Response *http.Response
	Error    error
}

func main() {
	log.Println("initialized")
	start := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()

	c := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/long", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.WithContext(ctx)
	ch := make(chan Result, 1)

	go func() {
		res, err := c.Do(req)
		ch <- Result{
			Response: res,
			Error:    err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("cancelled request after %v\n", time.Since(start))
			return
		case res := <-ch:
			if res.Error != nil {
				log.Println(res.Error)
				return
			}
			defer res.Response.Body.Close()
			body, err := ioutil.ReadAll(res.Response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf(`got body "%s" after %v\n`, string(body), time.Since(start))
			return
		}
	}

	log.Println("completed after", time.Since(start))
}
```

