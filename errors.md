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


## Error identity

```go
package main

import (
	"errors"
	"fmt"
)

var ErrOriginal = errors.New("original")

type ErrNotFound struct {
	name  string
	error error
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s: not found", e.name)
}

func (e *ErrNotFound) Unwrap() error {
	return e.error
}

func NewErrNotFound(err error, name string) *ErrNotFound {
	return &ErrNotFound{
		name:  name,
		error: err,
	}
}

func main() {
	err := NewErrNotFound(ErrOriginal, "user")

	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrOriginal))
	
	var nferr *ErrNotFound
	ok := errors.As(err, &nferr)
	fmt.Println(ok, nferr)
}
```

## MultiError

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

type MultiError struct {
	errors []error
}

func NewMultiError(errs ...error) *MultiError {
	if errs == nil {
		errs = make([]error, 0)
	}
	return &MultiError{
		errors: errs,
	}
}

func (m *MultiError) Error() string {
	msg := make([]string, len(m.errors))
	for i, err := range m.errors {
		msg[i] = err.Error()
	}
	return strings.Join(msg, "\n")
}

func (m *MultiError) Add(err error) bool {
	if err != nil {
		m.errors = append(m.errors, err)
		return true
	}
	return false
}

func (m *MultiError) AddString(s string) bool {
	if s != "" {
		m.errors = append(m.errors, errors.New(s))
		return true
	}
	return false
}

func main() {
	merr := NewMultiError()
	if merr.Add(errors.New("hello")) {
		fmt.Println("errors added")
	}
	merr.AddString("world")
	fmt.Println(merr)
}
```

## Error handling concurrency

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	ctx := context.Background()
	results, err := Google(ctx, "golang")
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	fmt.Println(elapsed)
	for _, result := range results {
		fmt.Println(result)
	}
}

type Result string

func Google(ctx context.Context, query string) (results []Result, err error) {
	g, ctx := errgroup.WithContext(ctx)

	searches := []Search{Web, Image, Video}
	results = make([]Result, len(searches))
	for i, search := range searches {
		i, search := i, search
		g.Go(func() error {
			result, err := search(ctx, query)
			fmt.Println(result, err)
			if err == nil {
				results[i] = result
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(ctx context.Context, query string) (Result, error) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		if rand.Intn(2) < 1 {
			return Result(""), errors.New("bad request")
		}
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}
```

## Building Error

When creating custom errors, there are useful fields to define
- code: a unique error code, e.g. `user.invalidName` that can be used for localization etc. It is actually the error `id`, but somehow `code` is more often associated with error than `id`
- metadata: the additional information to be passed down for constructing a more meaningful error message. The data is not always known during compile time such as min/max value, and may only be known during run-time. They can be made optional or required. If optional, the client must handle the scenario where the data is not provided.
- kind: a grouping for errors, e.g. `not found`, `created`, `conflict`, etc. This could be for example be mapped to HTTP status codes at the API layer
- message: a readable human error message, usually for logging purposes, and different from application errors that requires translation

```go
package main

import (
	"errors"
	"fmt"
)

// Overriding the interface ensures that the `Build` method must be called.
var ErrNameTooLong ErrorBuilder = NewError("user.invalidName", "Name is too long")

// This makes the Build() method optional.
var ErrNameIsRequired = NewError("user.nameIsRequired", "Name is required")

type ErrorBuilder interface {
	Build(metadata map[string]interface{}) error
}

func NewError(id, msg string) *Error {
	return &Error{id: id, message: msg}
}

type Error struct {
	id       string
	message  string
	metadata map[string]interface{}
}

// Build uses a value receiver to avoid mutating the original error.
// It returns a pointer receiver, so that the errors.As can be fulfilled.
func (e Error) Build(metadata map[string]interface{}) error {
	e.metadata = metadata
	return &e
}

func (e *Error) Is(other error) bool {
	err, ok := other.(*Error)
	if !ok {
		return false
	}
	return e.id == err.id
}

func (e Error) Error() string {
	return e.message
}

func main() {
	errorMatch()
	errorDoesNotMatch()
	errorBuild()
	errorIsSentinel()
}

func errorMatch() {
	err := &Error{message: "bad request"}

	var e *Error
	if errors.As(err, &e) {
		fmt.Println("match", e)
	} else {
		fmt.Println("not match")
	}
}

func errorDoesNotMatch() {
	err2 := errors.New("hello")
	var e *Error
	if errors.As(err2, &e) {
		fmt.Println("match", e)
	} else {
		fmt.Println("not match")
	}
}

func errorBuild() {
	err := ErrNameTooLong.Build(map[string]interface{}{
		"name": "john",
	})
	var e *Error
	if errors.As(err, &e) {
		fmt.Println("match", e, e.metadata)
	} else {
		fmt.Println("not match")
	}

	// The original error remains immutable.
	fmt.Printf("%#v\n", ErrNameTooLong)
}

func errorIsSentinel() {
	err := ErrNameTooLong.Build(map[string]interface{}{
		"name": "john",
	})
	if errors.Is(err, ErrNameTooLong.Build(nil)) {
		fmt.Println("match", err)
	} else {
		fmt.Println("not match")
	}

	err = ErrNameIsRequired
	if errors.Is(err, ErrNameIsRequired) {
		fmt.Println("match", err)
	} else {
		fmt.Println("not match")
	}
}
```
