# go get not getting latest

Sometimes after pushing a new commit to Github, running `go get -u ...` doesn't retrieve the latest due to caching. Run the command below to ensure the latest can be fetched immediately.

```bash
GOPROXY=direct go get -u github.com/yourpkg
```


# Enable goprivate for all repo from an organization

https://stackoverflow.com/questions/58305567/how-to-set-goprivate-environment-variable

go env -w GOPRIVATE=github.com/<OrgNameHere>/*


## Rename all forked packages


```
gofmt -w -r '"github.com/justwatch/facebook-marketing-api-golang-sdk/fb" -> "github.com/yourorg/facebook-marketing-api-golang-sdk/fb"' marketing/v22/*.go
```
