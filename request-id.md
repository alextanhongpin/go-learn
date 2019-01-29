## Request ID

```go
type key int
const requestIDKey key = 0

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
    reqID := req.Header.Get("X-Request-ID")
    if reqID == "" {
        reqID = generateRandomID()
    }

    return context.WithValue(ctx, requestIDKey, reqID)
}

func requestIDFromContext(ctx context.Context) string {
    return ctx.Value(requestIDKey).(string)
}
```

Middleware:

```
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := newContextWithRequestID(req.Context(), req)
        next.ServeHTTP(rw, req.WithContext(ctx))
    })
}
```

Usage:

```
func handler(rw http.ResponseWriter, req *http.Request) {
    reqID := requestIDFromContext(req.Context())
    fmt.Fprintf(rw, "Hello request ID %v\n", reqID)
}
```
