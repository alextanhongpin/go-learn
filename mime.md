# How to detect type of image
```go
package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
)

func main() {
	f, err := os.Open("sample.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fileType := mime.TypeByExtension(path.Ext(f.Name()))
	if fileType == "" {
		// SVG is detected as `text/xml`, when it's supposed to be `image/svg+xml`
		// https://github.com/golang/go/issues/15888
		fileType = http.DetectContentType(b)
	}
	fmt.Println("filetype:", fileType)
}
```
