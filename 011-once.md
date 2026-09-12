execute once by namespace

```
// Lazy
sync.Map + synce.OnceValue/s // LoadOrCreate pattern

// Delayed resolve/reject
sync.Map + context.WithCancelCause

// direct
sync.Map LoarOrStore
```
