# How to get private module

```bash
GOPRIVATE=github.com/alextanhongpin/testdump        go get github.com/alextanhongpin/testdump/httpdump
```

# How to force fetch the latest module version
```bash
GOPROXY=direct go get -u github.com/alextanhongpin/dump/pkg/reviver
```
