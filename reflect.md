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


## Nested custom tag
```go
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const EXAMPLE = "example"

type Request struct {
	Name  string   `example:"hello"`
	IDs   []string `example:"1,2,3"`
	Age   int      `example:"1"`
	Count int64    `example:"1000"`
	Pages []int    `example:"1,2,3,4"`
	// Car   *Car
	AnotherCar Car
}

type Car struct {
	Name string `example:"audi"`
	Brand
}

type Brand struct {
	Name string `example:"expensive"`
}

var typeOfSliceInt = reflect.TypeOf([]int{})
var typeOfSliceString = reflect.TypeOf([]string{})

func main() {

	/*
		res := deriveFromExample(req)
		switch v, ok := res.(Request); ok {
		case true:
			fmt.Printf("%#v\n", v)
			fmt.Printf("%#v\n", v.Car)
		}
	*/
	req := Request{}
	result := ParseExampleTag(req)
	switch v := result.(type) {
	case Request:
		fmt.Printf("%+v\n", v)
		fmt.Printf("%+v\n", v.AnotherCar)
	}

}

func ParseExampleTag(in interface{}) interface{} {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	var parse func(tEl reflect.Type, vEl reflect.Value) reflect.Value
	parse = func(tEl reflect.Type, vEl reflect.Value) reflect.Value {
		for i := 0; i < tEl.NumField(); i++ {
			tField := tEl.Field(i)
			vField := vEl.Field(i)
			tagValue := tField.Tag.Get(EXAMPLE)
			switch k := vField.Kind(); k {
			case reflect.Slice:
				switch vField.Type() {
				case typeOfSliceInt:
					value := strings.Split(tagValue, ",")
					result := make([]int, len(value))
					for i, v := range value {
						out, _ := strconv.Atoi(v)
						result[i] = out

					}
					vField.Set(reflect.ValueOf(result))
					fmt.Println("is string int")

				case typeOfSliceString:
					value := strings.Split(tagValue, ",")
					vField.Set(reflect.ValueOf(value))
					fmt.Println("is string slice")
				default:
					fmt.Println("not handled")
				}

			case reflect.Int:
				value, _ := strconv.Atoi(tagValue)
				vField.Set(reflect.ValueOf(value))
			case reflect.Int8:
				value, _ := strconv.ParseInt(tagValue, 10, 8)
				vField.Set(reflect.ValueOf(value))
			case reflect.Int16:
				value, _ := strconv.ParseInt(tagValue, 10, 16)
				vField.Set(reflect.ValueOf(value))
			case reflect.Int32:
				value, _ := strconv.ParseInt(tagValue, 10, 32)
				vField.Set(reflect.ValueOf(value))
			case reflect.Int64:
				value, _ := strconv.ParseInt(tagValue, 10, 64)
				vField.Set(reflect.ValueOf(value))
			case reflect.Ptr:
				elem := vField.Type().Elem()
				instance := reflect.New(elem)
				fmt.Println("got new instance", instance.Interface())
				// result := parse(reflect.TypeOf(tField), ptrInstance)
			case reflect.Struct:
				elem := vField.Type()
				instance := reflect.New(elem)
				out := ParseExampleTag(instance)
				fmt.Println(out)
				// t := reflect.TypeOf(tField)
				// v := reflect.ValueOf(vField)
				// out := parse(t, instance.Elem()).Interface()

				// fmt.Println(out)
				// vField.Set(reflect.ValueOf())
			case reflect.String:
				vField.SetString(tagValue)
			default:
				fmt.Println("unhandled", vField, tField, k)
				vField.Set(reflect.ValueOf(tagValue))
			}
		}
		return vEl
	}

	switch k := t.Kind(); k {
	case reflect.Ptr:
		tEl := t.Elem()
		vEl := v.Elem()
		return parse(tEl, vEl).Interface()
	case reflect.Struct:
		fmt.Println("is struct")
		instance := reflect.New(t)
		instanceElem := instance.Elem()
		return parse(t, instanceElem).Interface()
	default:
		fmt.Println("not handled", k)
	}
	return nil
}

func deriveFromExample(req interface{}) interface{} {

	typeOf := reflect.TypeOf(req)
	valueOf := reflect.ValueOf(req)
	if valueOf.CanSet() {
		return deriveFromExample(valueOf.Interface())
	}
	switch k := typeOf.Kind(); k {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println("is pointer")
		return deriveFromExample(valueOf.Interface())
	case reflect.Struct:
		instance := reflect.New(typeOf)
		instanceElem := instance.Elem()

		for i := 0; i < typeOf.NumField(); i++ {
			vfield := valueOf.Field(i)
			tfield := typeOf.Field(i)

			instanceField := instanceElem.Field(i)

			tagValue := tfield.Tag.Get("example")
			switch k := vfield.Kind(); k {
			case reflect.String:
				instanceField.SetString(tagValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, _ := strconv.ParseInt(tagValue, 10, 64)
				instanceField.SetInt(v)
			case reflect.Slice:
				values := strings.Split(tagValue, ",")
				slice := reflect.MakeSlice(
					tfield.Type,
					len(values),
					cap(values))
				switch t := tfield.Type; t {
				case typeOfSliceInt:
					for j := 0; j < len(values); j++ {
						v := slice.Index(j)
						i, _ := strconv.Atoi(values[j])
						v.Set(reflect.ValueOf(i))
					}
				default:
					for j := 0; j < len(values); j++ {
						v := slice.Index(j)
						v.Set(reflect.ValueOf(values[j]))
					}

				}

				instanceField.Set(slice)
			case reflect.Struct:
				instanceField.Set(reflect.ValueOf(deriveFromExample(vfield.Interface())))
			case reflect.Ptr:
				elem := vfield.Type().Elem()
				ptrInstance := reflect.New(elem)
				ptrInstanceElem := ptrInstance.Elem()
				for j := 0; j < elem.NumField(); j++ {
					elemField := elem.Field(j)
					elemTagValue := elemField.Tag.Get("example")
					ptrInstanceElem.Field(j).Set(reflect.ValueOf(elemTagValue))
				}
				instanceField.Set(ptrInstance)
			default:
				fmt.Printf("type %s is not handled\n", k)
			}
		}
		return instanceElem.Interface()
	default:
		fmt.Println("not handled", k)
		return nil
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
