```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Test struct {
	Name          string `json:"name"`
	Age           int    `json:"age,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

func main() {
	newTest := Test{
		Name:          "john",
		Age:           100,
		EmailVerified: true,
	}
	out := encode(newTest)
	fmt.Println("encode success:", out.Encode())
	var res Test
	if err := decode(&res, out); err != nil {
		log.Fatal(err)
	}
	fmt.Println("decode success:", res)

}

func decode(in interface{}, u url.Values) error {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	// TODO: Return error
	if !v.CanSet() {
		return errors.New("must be pointer")
	}
	switch t.Kind() {
	case reflect.Struct:
		// Todo, check if it's pointer and struct
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			tags := strings.Split(tag, ",")
			if len(tags) == 0 {
				continue
			}
			fieldName := f.Name
			tagValue := u.Get(tags[0])
			if len(strings.TrimSpace(tagValue)) == 0 {
				continue
			}
			fieldByName := v.FieldByName(fieldName)
			if !fieldByName.IsValid() || !fieldByName.CanSet() {
				continue
			}
			switch f.Type.Kind() {
			case reflect.String:
				fieldByName.SetString(tagValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				tmp, _ := strconv.Atoi(tagValue)
				fieldByName.SetInt(int64(tmp))
			case reflect.Bool:
				fieldByName.SetBool(tagValue == "true")
			}
		}
		return nil
	default:
		return nil
	}
}

func encode(in interface{}) url.Values {
	u := url.Values{}
	t := reflect.TypeOf(in)
	switch t.Kind() {
	// Only handles if it is a struct
	case reflect.Struct:
		vin := reflect.ValueOf(in)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			v := vin.Field(i)
			tag := f.Tag.Get("json")
			tags := strings.Split(tag, ",")

			// Skip if no tags
			if len(tags) == 0 {
				continue
			}
			name := tags[0]
			if len(strings.TrimSpace(name)) == 0 {
				continue
			}
			omitempty := strings.HasSuffix(tag, ",omitempty")

			z := reflect.Zero(v.Type())
			isZero := z.Interface() == v.Interface()
			if isZero && omitempty {
				continue
			}

			var tmp interface{}
			switch f.Type.Kind() {
			case reflect.String:
				tmp = v.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				tmp = v.Int()
			case reflect.Bool:
				tmp = v.Bool()
			}
			u.Add(name, fmt.Sprint(tmp))
		}
		return u
	default:
		return u
	}

}
```



## Setting nested struct
```go
package main

import (
	"fmt"
	"reflect"
)

type B struct {
	C string
}
type A struct {
	ID string
	B
}

func main() {
	// Setting nested struct.
	var a A
	v := reflect.ValueOf(&a)
	v.Elem().FieldByIndex([]int{0}).Set(reflect.ValueOf("hello"))
	v.Elem().FieldByIndex([]int{1}).Set(reflect.ValueOf(B{"car"}))
	fmt.Println("Hello, playground", a)
}
```

## Parsing nested struct

```go
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type B struct {
	C string `example:"this is c"`
}

type D struct {
	Name string `example:"john"`
	*B
	Strings []int `example:"car,paper,this"`
}

type A struct {
	ID string `example:"this is id"`
	*B
	Numbers []int `example:"1,2,3"`
	Strings []int `example:"car,paper,this"`
	D
	IsValid bool `example:"true"`
	IsAge   bool `example:"false"`
	IsNone  bool
	Name    string   `json:"name" example:"hello"`
	IDs     []string `json:"ids" example:"1,2,3"`
	Age     int      `json:"age" example:"1"`
	Count   int64    `json:"count" example:"1000"`
	Pages   []int    `json:"pages" example:"1,2,3,4"`
}

var tagName = "example"

func main() {
	// Setting nested struct.
	{

		var a A
		out := ParseTag(a)
		fmt.Println(out)
		switch o := out.(type) {
		case A:
			fmt.Printf("%#v\n", o)
			fmt.Printf("%#v\n", o.B)
			fmt.Printf("%#v\n", o.D)
			fmt.Printf("%#v\n", o.D.B)
		}
	}
	{
		var a A
		out := ParseTag(&a)
		fmt.Println(out)
		switch o := out.(type) {
		case A:
			fmt.Printf("%#v\n", o)
			fmt.Printf("%#v\n", o.B)
			fmt.Printf("%#v\n", o.D)
			fmt.Printf("%#v\n", o.D.B)
		}
		b, _ := json.Marshal(out)
		fmt.Println(string(b))
	}
}

var reflectSliceIntType = reflect.TypeOf([]int{})

func ParseTag(in interface{}) interface{} {
	var parse func(in interface{}) reflect.Value
	parse = func(in interface{}) reflect.Value {
		v := reflect.ValueOf(in)
		switch k := v.Kind(); k {
		case reflect.Struct:

			t := reflect.TypeOf(in)
			newEl := reflect.New(t)
			el := newEl.Elem()

			for i := 0; i < el.NumField(); i++ {
				r := el.FieldByIndex([]int{i})
				tagValue := t.Field(i).Tag.Get(tagName)

				switch r.Kind() {
				case reflect.Struct:
					newStruct := reflect.New(reflect.TypeOf(r.Interface()))
					out := parse(newStruct.Elem().Interface())
					r.Set(reflect.ValueOf(out.Interface()))
				case reflect.Ptr:
					newStruct := reflect.New(reflect.TypeOf(r.Interface()).Elem())
					out := parse(newStruct.Elem().Interface())
					r.Set(out.Addr())
				case reflect.Bool:
					r.SetBool(tagValue == "true")
				case reflect.Slice:

					val := strings.Split(tagValue, ",")

					switch r.Type() {
					case reflectSliceIntType:
						out := make([]int, len(val))
						for i, v := range val {
							out[i], _ = strconv.Atoi(v)
						}
						r.Set(reflect.ValueOf(out))

					default:
						r.Set(reflect.ValueOf(val))
					}
				case reflect.Int:
					val, _ := strconv.Atoi(tagValue)
					r.Set(reflect.ValueOf(val))
				case reflect.Int8:
					val, _ := strconv.ParseInt(tagValue, 10, 8)
					r.Set(reflect.ValueOf(val))
				case reflect.Int16:
					val, _ := strconv.ParseInt(tagValue, 10, 16)
					r.Set(reflect.ValueOf(val))
				case reflect.Int32:
					val, _ := strconv.ParseInt(tagValue, 10, 32)
					r.Set(reflect.ValueOf(val))
				case reflect.Int64:
					val, _ := strconv.ParseInt(tagValue, 10, 64)
					r.Set(reflect.ValueOf(val))
				default:
					r.Set(reflect.ValueOf(tagValue))
				}
			}
			return el
		case reflect.Ptr:
			r := reflect.New(reflect.TypeOf(in).Elem())
			out := parse(r.Elem().Interface())
			return out
		default:
			return v
		}
	}
	out := parse(in)
	return out.Interface()
}
```

## Creating new instance reflect

Pointer:
```go
real := new(A)
reflected := reflect.New(reflect.TypeOf(real).Elem()).Elem().Interface()
fmt.Println(real)
fmt.Println(reflected)
```

Struct:

```go
real := A{}
reflected := reflect.New(reflect.TypeOf(real)).Elem().Interface()
fmt.Println(real)
fmt.Println(reflected)
```

## Check empty values with reflect

Alternative is to just check nil, 0 or empty string:

```go
package gen

import "reflect"

func IsZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
```
