# Formating code

Should be as simple as:
```go
content, err := format.Source(content)
// check error
file.Write(content)
```

References:
https://golang.org/pkg/go/format/

## Format with import

```go
package main

import "golang.org/x/tools/imports"

// FormatSource is gofmt with addition of removing any unused imports.
func FormatSource(source []byte) ([]byte, error) {
	return imports.Process("", source, &imports.Options{
		AllErrors: true, Comments: true, TabIndent: true, TabWidth: 8,
	})
}
```
