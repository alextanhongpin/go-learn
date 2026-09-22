You need to strip the prefix to match the path value:

```go
package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	{
		api := http.NewServeMux()
		api.HandleFunc("GET /{name}", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "got user:"+r.PathValue("name"))
		})
		api.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "got users")
		})

		mux.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", api))
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})

	slog.Info("listening to port *:8080")
	http.ListenAndServe(":8080", mux)
}

```
