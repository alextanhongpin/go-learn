```go
func cmd(arg string, args ...string) {
	cmd := exec.Command(arg, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	log.Printf("sdtout: %s", stdout.String())
	log.Printf("sderr: %s", stderr.String())
}
```
