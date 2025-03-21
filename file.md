```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	f := NewFile("./testdata/hello.txt")
	f.overwrite = true
	defer f.Close()

	d := &dump{
		r:   f,
		out: f,
		raw: os.Stdout,
	}
	type User struct {
		Name string `json:"name"`
	}
	if err := d.Dump(&User{Name: "jane"}); err != nil {
		panic(err)
	}
}

type dump struct {
	r   io.ReadCloser
	raw io.WriteCloser
	out io.WriteCloser
}

func (d *dump) Dump(v any) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	fmt.Println("READ", string(b))
	b, err = json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	_, err = d.raw.Write(b)
	if err != nil {
		return err
	}

	_, err = d.out.Write(b)
	if err != nil {
		return err
	}

	return nil
}

type rwc interface {
	io.Reader
	io.Writer
	io.Closer
}

var _ rwc = (*File)(nil)

type File struct {
	f         *os.File
	name      string
	overwrite bool
	exists    bool
}

func NewFile(name string) *File {
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		panic(err)
	}

	var exists bool
	f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
	if errors.Is(err, os.ErrExist) {
		exists = true
		f, err = os.OpenFile(name, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	return &File{
		f:      f,
		name:   name,
		exists: exists,
	}
}

func (f *File) Write(b []byte) (int, error) {
	if f.exists {
		if f.overwrite {
			// We need to truncate the file content.
			f, err := os.OpenFile(f.name, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			return f.Write(b)
		}

		return 0, nil
	}

	return f.f.Write(b)
}

func (f *File) Read(b []byte) (int, error) {
	return f.f.Read(b)
}

func (f *File) Close() error {
	return f.f.Close()
}
```

## Touch

```go
package main

import (
	"fmt"
	"os"
)

func hasTouch(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	panic(err)
}

func main() {
	err := TouchFile("test.txt")
	fmt.Println(os.IsExist(err))
	fmt.Println(hasTouch("test.txt"))
	fmt.Println(hasTouch("what.txt"))

	_, err = os.Stat("non.txt")
	fmt.Println(os.IsNotExist(err))

	err = touchWithContent("hello.txt", func(f *os.File) error {
		_, err := f.WriteString("Hello, World!")
		return err
	})
	fmt.Println(os.IsExist(err))
	b, err := os.ReadFile("hello.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_EXCL|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func touchWithContent(name string, fn func(f *os.File) error) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_EXCL|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("writing file...")
	if err := fn(file); err != nil {
		return err
	}

	return nil
}

```
