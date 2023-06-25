# go get not getting latest

Sometimes after pushing a new commit to Github, running `go get -u ...` doesn't retrieve the latest due to caching. Run the command below to ensure the latest can be fetched immediately.

```bash
GOPROXY=direct go get -u github.com/yourpkg
```
