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
	Name          string `json:"name,omitempty"`
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
			if tagValue == "" {
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
				b := false
				if tagValue == "true" {
					b = true
				}
				fieldByName.SetBool(b)
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
