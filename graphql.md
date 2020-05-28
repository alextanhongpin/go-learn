## Setup graphiql UI for golang

```go
http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
}))
  
  
// Does not require us to bin-pack the static files when deploying.
var page = []byte(`
<html>
  <head>
    <title>Simple GraphiQL Example</title>
    <link href="https://unpkg.com/graphiql/graphiql.min.css" rel="stylesheet" />
  </head>
  <body style="margin: 0;">
    <div id="graphiql" style="height: 100vh;"></div>

    <script
      crossorigin
      src="https://unpkg.com/react/umd/react.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/graphiql/graphiql.min.js"
    ></script>

    <script>
      const graphQLFetcher = graphQLParams =>
        fetch('http://localhost:8080/query', {
          method: 'post',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(graphQLParams),
        })
          .then(response => response.json())
          .catch(() => response.text());
      ReactDOM.render(
        React.createElement(GraphiQL, { fetcher: graphQLFetcher }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>`)
```
