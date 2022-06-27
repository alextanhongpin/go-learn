# SQL-Like filters with go

- [how ent does it](https://entgo.io/blog/2021/07/01/automatic-graphql-filter-generation/)
- my own lib [github.com/alextanhongpin/goql](https://github.com/alextanhongpin/goql)
- if only one operation per field is supported, we can store the ops separately
- otherwise, just store the ops separately (op: gte, value:17), (gte: 17)
- alternative is to do code generation base on a type, e.g

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
	v.Set("name", "eq:john appleseed") // limited to single op, must be slice to support multiple values
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

## Alternative

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var u UserFilter
	if err := json.Unmarshal([]byte(`
	{
		"and": [{"name": {"eq": "john"}}],
		"or": [{"name": {"neq": "jane"}}]
	}
	`), &u); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", u)
}

type User struct {
	Name Where[string] `json:"name"`
	// Age  []WhereOp[int] `json:"age"`
}

type UserFilter struct {
	And []User
	Or  []User
}

type Where[T any] struct {
	Eq  T   `json:"eq,omitempty"`
	Neq T   `json:"neq,omitempty"`
	In  []T `json:"in,omitempty"`

	And []Where[T] `json:"and,omitempty"`
	Or  []Where[T] `json:"or,omitempty"`
}

// This is less strongly typed ...
type WhereOp[T any] struct {
	Op string `json:"op"`
	T  []T    `json:"value"`
}
```


## Another version

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/schema"
)

func main() {
	now := time.Now()
	now.Add(-24 * time.Hour)
	b := []byte(fmt.Sprintf(`
	{
 		"age": {
			"between": {"left": 13, "right": 30}
  		},
		"name": {
			"like": ["a%%", "b%%", "c%%"]
		},
		"createdAt": {
			"gt": %q
		},
		"or": [{"age": {"gt": 10}}]
	}`, now.Format(time.RFC3339)))

	var u UserWhere
	if err := json.Unmarshal(b, &u); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *u.Age.Between)
	fmt.Printf("%+v\n", u.Name)
	fmt.Printf("%+v\n", u.CreatedAt)
	fmt.Printf("%+v\n\n", u)

	// For most cases, this is more than enough, since most SQL lib doesn't check the types.
	u.Print()

	fmt.Println("\ndecoding url.Values:")
	dec := schema.NewDecoder()

	v := make(url.Values)
	v.Set("age.between.left", "17")
	v.Set("age.between.right", "100")
	v.Add("name.like", "a%")
	v.Add("name.like", "b%")
	v.Add("name.like", "c%")
	v.Add("createdAt.gt", now.Format(time.RFC3339))
	v.Add("and.0.name.eq", "jessie")
	var u2 UserWhere
	if err := dec.Decode(&u2, v); err != nil {
		panic(err)
	}
	u2.Print()

	fmt.Printf("%+v", *u2.And[0].Name)
}

type UserWhere struct {
	UserFilter // Perhaps it's better not to embed it here
	Where[UserFilter]
	Limit  *int
	Offset *int
	Sort   []string
}

type UserFilter struct {
	Age       *WhereOp[int]        `json:"age,omitempty"`
	Name      *WhereOp[string]     `json:"name,omitempty"`
	CreatedAt *WhereOp[*time.Time] `json:"createdAt,omitempty"`
}

func (u UserFilter) Print() {
	if u.Age != nil {
		for k, v := range u.Age.Values() {
			fmt.Printf("age %s %v\n", k, v)
		}
	}
	if u.Name != nil {
		for k, v := range u.Name.Values() {
			fmt.Printf("name %s %v\n", k, v)
		}
	}
	if u.CreatedAt != nil {
		for k, v := range u.CreatedAt.Values() {
			fmt.Printf("createdAt %s %v\n", k, v)
		}
	}
}

type Where[T any] struct {
	And []T `json:"and,omitempty"`
	Or  []T `json:"or,omitempty"`
	Not []T `json:"not,omitempty"`
}

type Between[T any] struct {
	Left  T `json:"left,omitempty"`
	Right T `json:"right,omitempty"`
}

type WhereOp[T any] struct {
	Eq      T           `json:"eq,omitempty"`
	Neq     T           `json:"neq,omitempty"`
	Lt      T           `json:"lt,omitempty"`
	Lte     T           `json:"lte,omitempty"`
	Gt      T           `json:"gt,omitempty"`
	Gte     T           `json:"gte,omitempty"`
	Is      *bool       `json:"is,omitempty"`
	In      []T         `json:"in,omitempty"`
	Between *Between[T] `json:"between,omitempty"`
	Like    []T         `json:"like,omitempty"`
	ILike   []T         `json:"ilike,omitempty"`
	Not     bool        `json:"not,omitempty"`
}

func (w *WhereOp[T]) Values() map[string]any {
	b, err := json.Marshal(w)
	if err != nil {
		panic(err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}

	res := make(map[string]any)
	for k, v := range m {
		if v == nil {
			continue
		}
		res[k] = v
	}

	return res
}
```
