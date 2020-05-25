```go
package main

import (
	"errors"
	"fmt"
)

var ErrOne = errors.New("one")

func main() {
	e1 := ErrOne
	e2 := fmt.Errorf("two: %w", e1)
	e3 := fmt.Errorf("three: %w", e2)

	fmt.Println(e1)
	fmt.Println(e2)
	fmt.Println(e3)
	fmt.Println(errors.Unwrap(e2))
	fmt.Println(errors.Unwrap(e3))
	fmt.Println(errors.Unwrap(errors.Unwrap(e3)))
	fmt.Println(errors.Is(e1, ErrOne))
	fmt.Println(errors.Is(e2, ErrOne))	
	fmt.Println(errors.Is(e3, ErrOne))	
}
```

## Custom Error

```go
package main

import (
	"errors"
	"fmt"
)

var ErrEmpty = errors.New("file is empty")
var ErrExist = errors.New("file exists")

type File struct {
	Name string
}

func NewFile(name string) *File {
	return &File{Name: name}
}

type FileError struct {
	file *File
	err  error
}

func NewFileError(err error, file *File) *FileError {
	return &FileError{
		err:  err,
		file: file,
	}
}

func (f *FileError) Error() string {
	if f.err != nil {
		return f.err.Error()
	}
	return ""
}

func (f *FileError) Unwrap() error {
	return f.err
}

func main() {
	f := NewFile("path.txt")
	err := NewFileError(ErrExist, f)
	err2 := fmt.Errorf("bad request: %w", err)
	
	fmt.Println("err", err)
	fmt.Println(errors.Is(err, ErrExist))
	fmt.Println(errors.Is(err, ErrEmpty))
	fmt.Println(errors.Is(err, err))
	
	fmt.Println("err2", err2)
	fmt.Println(errors.Is(err2, ErrExist))
	fmt.Println(errors.Is(err2, ErrEmpty))
	fmt.Println(errors.Is(err2, err))
	
	var fe *FileError
	if errors.As(err2, &fe) {
		fmt.Println("yes", fe)
	}
}	

```
