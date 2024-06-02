# How to run python in golang?

- use exec.Command, but requires python installed
- build a python binary, then embed thr binary in memory. works on linux only using memfd
- grpc server
- http server
- wasi? but there are no good library yet

run on browser instead, and send the results to the server
