# SQL-Like filters with go

- [how ent does it](https://entgo.io/blog/2021/07/01/automatic-graphql-filter-generation/)
- my own lib [github.com/alextanhongpin/goql](https://github.com/alextanhongpin/goql)
- alternative is to do code generation base on a type

```go
type User struct {
	Name string
	Age int
}

// Produces

type UserFilter struct {
	NameEq string `json:"name.eq"
	NameNeq string `json:"name.neq"
	...
}
```

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/gorilla/schema"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Name     *Where[*null.String] `json:"name" schema:"name"`
	Age      *Where[*null.Int]    `json:"age" schema:"age"`
	Married  *Where[*null.Bool]   `json:"married" schema:"married"`
	Birthday *Where[*null.Time]   `json:"birthday" schema:"birthday"`
	Height   *Where[*null.Float]  `json:"height" schema:"height"`
}

type UserFilter struct {
	And []User
	Or  []User
	Not []User
}

func main() {
	b := []byte(fmt.Sprintf(`{
		"name": "is:null", 
		"age": "gt:10",
		"married": "eq:true",
		"birthday": "gt:%s"
	}`, time.Now().Format(time.RFC3339)))
	var u User
	if err := json.Unmarshal(b, &u); err != nil {
		panic(err)
	}
	debug(u)

	v := make(url.Values)
	v.Set("name", "eq:john appleseed")
	v.Set("age", "gt:17")
	v.Set("married", "eq:true")
	v.Set("birthday", fmt.Sprintf("gt:%s", time.Now().Format(time.RFC3339)))

	var u2 User

	var dec = schema.NewDecoder()
	if err := dec.Decode(&u2, v); err != nil {
		panic(err)
	}
	debug(u2)

}
func debug(u User) {
	fmt.Printf("%+v: %+v\n", u.Name, *u.Name.T)
	fmt.Printf("%+v: %+v\n", u.Age, *u.Age.T)
	fmt.Printf("%+v: %+v\n", u.Married, *u.Married.T)
	fmt.Printf("%+v: %+v\n", u.Birthday, *u.Birthday.T)
	fmt.Printf("%+v\n", u.Height)
}

const (
	IsTrue Is = iota
	IsFalse
	IsNull
	IsNotNull
	IsUnknown
)

var ErrUnknownIsValue = errors.New("unknown 'is' value")

type Is int

func (is Is) String() string {
	// is:not variant of `true` and `false` is redundant.
	text, ok := isText[is]
	if !ok {
		return ""
	}
	return text
}

func (is *Is) MarshalText() ([]byte, error) {
	return []byte(is.String()), nil
}

func (is *Is) UnmarshalText(in []byte) error {
	vl, ok := isFromText[string(in)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownIsValue, string(in))
	}
	*is = vl
	return nil
}

var isText = map[Is]string{
	IsTrue:    "true",
	IsFalse:   "false",
	IsNull:    "null",
	IsNotNull: "notnull",
	IsUnknown: "unknown",
}

var isFromText map[string]Is

func init() {
	isFromText = make(map[string]Is)
	for k, v := range isText {
		isFromText[v] = k
	}
}

type Where[T encoding.TextUnmarshaler] struct {
	Eq, Neq, Lt, Lte, Gt, Gte bool
	Is                        *Is
	T                         T
}

// Where should be strictly unmarshal only operation.
func (w *Where[T]) UnmarshalText(in []byte) error {
	var n T
	w.T = reflect.New(reflect.TypeOf(n).Elem()).Interface().(T)

	vals := bytes.SplitN(in, []byte(":"), 2)

	if len(vals) == 1 {
		w.Eq = true
		return w.T.UnmarshalText(vals[0])
	}

	switch {
	case bytes.Equal(vals[0], []byte("eq")):
		w.Eq = true
	case bytes.Equal(vals[0], []byte("neq")):
		w.Neq = true
	case bytes.Equal(vals[0], []byte("lt")):
		w.Lt = true
	case bytes.Equal(vals[0], []byte("lte")):
		w.Lte = true
	case bytes.Equal(vals[0], []byte("gt")):
		w.Gt = true
	case bytes.Equal(vals[0], []byte("gte")):
		w.Gte = true
	case bytes.Equal(vals[0], []byte("is")):
		is := new(Is)
		if err := is.UnmarshalText(vals[1]); err != nil {
			return err
		}
		w.Is = is
		return nil
	}

	return w.T.UnmarshalText(vals[1])
}
```
