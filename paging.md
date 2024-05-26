## Cursor Pagination

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// TODO: Figure this out...
// if it is the first pagination, fetch the first n + 1, n defaults to 20 items
// set the startCursor to the first item in the list; startCursor = 436, activity_id = bank_transfer_payment_approved
// set the endCursor to the nth item in the list: endCursor: 316, activity_id=document_uploaded
// if the n+1 item exist, set hasNextPage to true, hasNextPage: true 
// set the hasPrevPage to false: hasPrevPage: false

type CursorPaginationRequest struct {
	First     int
	After     string
	Before    string
	OrderBy   string
	Direction string
}

type CursorPaginationResponse struct {
	HasNextPage bool
	HasPrevPage bool
	StartCursor string
	EndCursor   string
}

func ToBase64(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

type Cursor map[string]interface{}

func (c Cursor) String() string {
	u := url.Values{}
	for k, v := range c {
		u.Set(k, fmt.Sprint(v))
	}
	return ToBase64(u.Encode())
}

func (c Cursor) Stmt() string {
	result := make([]string, len(c))
	var i int
	for k, v := range c {
		result[i] = fmt.Sprintf("%s >= %v", k, v)
		i++
	}
	return strings.Join(result, " AND ")
}

func NewCursor(b64 string) (Cursor, error) {
	c := make(Cursor, 0)
	b, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return c, err
	}
	m, err := url.ParseQuery(string(b))
	if err != nil {
		fmt.Println("error unmarshaling", err)
		return c, err
	}
	for k, v := range m {
		c[k] = v[0]
	}
	return c, nil
}

func main() {
	startCursor := Cursor{
		"activity_id": "bank_transfer_payment_approve",
		"id":          436,
	}
	endCursor := Cursor{
		"activity_id": "document_uploaded",
		"id":          316,
	}
	res := CursorPaginationResponse{
		HasNextPage: true,
		HasPrevPage: false,
		StartCursor: startCursor.String(),
		EndCursor:   endCursor.String(),
	}
	fmt.Println(res)

	req := CursorPaginationRequest{
		First:     10, // + 1
		After:     endCursor.String(),
		OrderBy:   "created_at",
		Direction: "desc",
	}

	c, err := NewCursor(req.After)
	if err != nil {
		log.Fatal(err)
	}
	if _, exist := c[req.OrderBy]; !exist {
		fmt.Println("perform a fresh query")
	}
	fmt.Println(c.Stmt())
}
```

## Generic Pagination

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"net/url"
)

func main() {
	p := &paginator[any]{
		URL:   "/hello",
		Limit: 20,
	}
	as := make([]any, 21)
	fmt.Println(len(p.Paginate(as)))
	fmt.Println(p.HasNext(as))

	v := url.Values{}
	v.Set("name", "john")
	v.Set("age", "10")
	fmt.Println(p.BuildURL(v))
	fmt.Println("Hello, 世界")

}

var _ Paginator[any] = new(paginator[any])

type Paginator[T any] interface {
	Paginate(v []T) []T
	HasNext(v []T) bool
	BuildURL(v url.Values) string
}

type paginator[T any] struct {
	URL   string
	Limit int
}

func (p *paginator[T]) Paginate(v []T) []T {
	if len(v) > p.Limit {
		return v[:p.Limit]
	}

	return v
}

func (p *paginator[T]) HasNext(v []T) bool {
	return len(v) > p.Limit
}

func (p *paginator[T]) BuildURL(v url.Values) string {
	u, err := url.Parse(p.URL)
	if err != nil {
		panic(err)
	}
	u.RawQuery = v.Encode()
	return u.String()
}

```


## HashID

If all params for cursor is int, use hash id, or now known as [sqids](https://github.com/sqids/sqids-go)
