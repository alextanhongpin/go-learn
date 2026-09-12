execute once by namespace

```
// Lazy
sync.Map + synce.OnceValue/s // LoadOrCreate pattern

// Delayed resolve/reject
sync.Map + context.WithCancelCause
sync.Map + channel + data/error

// direct
sync.Map LoarOrStore
```

We can also add weak.Pointers to it for stuff like cache etc. Basically things with TTL
