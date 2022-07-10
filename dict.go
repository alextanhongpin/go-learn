# Dict

## Why

There are some dict operations that could have been added to make map operations safer.

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type Status int

const (
	Unknown Status = iota
	Pending
	Failed
	Success
)

var (
	textByStatus = NewDict(map[Status]string{
		Unknown: "unknown",
		Pending: "pending",
		Failed:  "failed",
		Success: "success",
	})
	statusByText = textByStatus.Invert()
)

func (s Status) String() string {
	str, _ := textByStatus.Get(s)
	return str
}

func (s Status) Valid() bool {
	_, ok := textByStatus.Get(s)
	return ok
}

func NewStatus[T string | int](unk T) (Status, bool) {
	switch v := any(unk).(type) {
	case string:
		status, ok := statusByText.Get(v)

		return status, ok
	case int:
		status := Status(v)
		if status.Valid() {
			return status, true
		}

		return 0, false
	}

	return 0, false
}

func main() {
	fmt.Println(NewStatus(1))
	fmt.Println(NewStatus("failed"))
	s, ok := NewStatus("hello")
	fmt.Println("got", s, ok)
}

type Dict[U, V comparable] struct {
	value map[U]V
}

func NewDict[U, V comparable](m map[U]V) *Dict[U, V] {
	return &Dict[U, V]{value: m}
}

func (d *Dict[U, V]) Get(key U) (v V, found bool) {
	v, found = d.value[key]

	return
}

func (d *Dict[U, V]) MustGet(key U) V {
	v, found := d.Get(key)
	if !found {
		panic(fmt.Errorf("invalid dict key"))
	}

	return v
}

func (d *Dict[U, V]) Set(key U, value V) {
	if d.value == nil {
		d.value = make(map[U]V)
	}

	d.value[key] = value

	return
}

func (d *Dict[U, V]) Invert() *Dict[V, U] {
	res := make(map[V]U)

	for k, v := range d.value {
		res[v] = k
	}

	return NewDict[V, U](res)
}
```
