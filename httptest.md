# Creating 1mb size file to test limit
```go
	var sb strings.Builder
	sb.WriteString("user_id")
	sb.WriteString("\n")

	var i int
	for sb.Len() < 1024*1024 {
		sb.WriteString(fmt.Sprint(i))
		sb.WriteString("\n")
	}

	fileLargerThan1MB := path.Join(t.TempDir(), "1mb.csv")

	if err := os.WriteFile(fileLargerThan1MB, []byte(sb.String()), 0644); err != nil {
		t.Fatalf("failed to write temp csv: %v", err)
	}

```

Handler:

```go
const MaxFileSizeMB = 1 << 20
// ...
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSizeMB)
  // Just this is not enough.
	if err := r.ParseMultipartForm(MaxFileSizeMB); err != nil {
		_ = encodeError(w, err)
		return err
	}

```
