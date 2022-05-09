# To debug memory usage when building

```bash
GODEBUG=gctrace=1 go build ./cmd/server
```

You can also run this to check memory used when compiling:

```bash
go build -toolexec '/usr/bin/time -v'

# For macOS
go build -toolexec '/usr/bin/time -alh'
```

## Clear cache
```bash
go clean --cache
```

## Check storage used by cache


```bash
du -hs $(go env GOCACHE)
```
