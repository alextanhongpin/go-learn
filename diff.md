# Diff text tool

Testing out different diff tools for snapshot testing.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = `bca hello world`
	text2 = `abc hi world`
)

func main() {
	dmp := diffmatchpatch.New()

	fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(text1, text2)
	diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)
	fmt.Println()
	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println(diff_text(diffs))

	edits := myers.ComputeEdits(span.URIFromPath("a.txt"), text1, text2)
	diff := fmt.Sprint(gotextdiff.ToUnified("a.txt", "b.txt", text1, edits))
	fmt.Println()
	fmt.Println(diff)

	fmt.Println()
	fmt.Println(cmp.Diff(text1, text2))
}

func diff_text(diffs []diffmatchpatch.Diff) string {
	var b strings.Builder

	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			b.WriteByte('+')
			b.WriteString(text)
			b.WriteByte('+')
		case diffmatchpatch.DiffDelete:
			b.WriteByte('-')
			b.WriteString(text)
			b.WriteByte('-')
		case diffmatchpatch.DiffEqual:
			b.WriteString(text)
		}
	}
	return b.String()
}
```
