```go
// Package null is a dummy function used to coordinate multiple unrelated work
package null

type Fn func() error
```

If you have a list of unrelated work, which you want to run in separate goroutines, you might have to do this:


```go
var wg sync.WaitGroup
// Bad. Hardcoding this will lead to trouble when you need to add more work.
wg.Add(4)

go func() {
  defer wg.Done()
  workA()
}()

go func() {
  defer wg.Done()
  workB()
}()

go func() {
  defer wg.Done()
  workC()
}()

go func() {
  defer wg.Done()
  workD()
}()

// Wait for all of them to finish
wg.Wait()
```

Another way to do it:

```go
fns := []null.Fn{
  null.Fn(func() error { return workA() }),
  null.Fn(func() error { return workB() }),
  null.Fn(func() error { return workC() }),
  null.Fn(func() error { return workD() }),
}

var wg sync.WaitGroup
wg.Add(len(fns))

for _, fn := range fns {
  go func(f null.Fn) {
    defer wg.Done()
    f()
  }(fn)
}

wg.Wait()
```
