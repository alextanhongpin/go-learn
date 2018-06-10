
## CPU Profiling

Or to look at a 30-second CPU profile:

```bash
$ go tool pprof  http://localhost:6060:/debug/pprof/profile

$ go tool pprof /Users/alextanhongpin/pprof/pprof.samples.cpu.003.pb.gz

$ (pprof) top10
$ (pprof) top5 -cum
```


## Memory Profiling

Then use the pprof tool to look at the heap profile:

```bash
$ go tool pprof http://localhost:6060/debug/pprof/heap

(pprof) top5
(pprof) list FnName
```


One option is ‘–alloc_space’ which tells you how many megabytes have been allocated.

```bash
$ go tool pprof --alloc_space http://localhost:6060/debug/pprof/heap
```

The other – ‘–inuse_space’ tells you know how many are still in use.

```
<!-- $ go tool pprof --inuse_objects http://localhost:6060/debug/pprof/heap -->
$ go tool pprof --inuse_space http://localhost:6060/debug/pprof/heap
```
