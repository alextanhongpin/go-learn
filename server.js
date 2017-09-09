const http = require('http')
const port = 3000

const requestHandler = (request, response) => {
  response.end('hello world')
}

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err)
  }

  console.log(`server is listening on ${port}`)
})

// wrk -d30s -c10 -t5 http://localhost:3000
// Running 30s test @ http://localhost:3000
//   5 threads and 10 connections
//   Thread Stats   Avg      Stdev     Max   +/- Stdev
//     Latency   423.82us  556.95us  26.58ms   98.77%
//     Req/Sec     5.04k     0.95k    6.29k    66.00%
//   754350 requests in 30.10s, 79.85MB read
// Requests/sec:  25062.11
// Transfer/sec:      2.65MB
