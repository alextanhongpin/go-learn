package main

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// People represents a people model
type People struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

var people map[string]interface{}

func main() {
	p := People{"john", 1}

	// Getting the element of the struct
	s := reflect.ValueOf(&p).Elem()
	log.Printf("elem     = %v\n", s)
	log.Printf("numField = %v\n", s.NumField())
	log.Printf("type     = %v\n\n", s.Type())

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := field.Type()
		fieldKind := field.Kind()
		fieldName := s.Type().Field(i).Name
		fieldTag := s.Type().Field(i).Tag.Get("json")
		log.Printf("index     = %v\n", i)
		log.Printf("field     = %v\n", field)
		log.Printf("fieldType = %v\n", fieldType)
		log.Printf("fieldKind = %v\n", fieldKind)
		log.Printf("fieldName = %v\n", fieldName)
		log.Printf("fieldTag  = %v\n\n", fieldTag)
	}

	// Example of setting struct value through reflect
	p = People{}
	structElem := reflect.ValueOf(&p).Elem()
	structFieldValue := structElem.FieldByName("Name")
	log.Println(structFieldValue)

	if !structFieldValue.IsValid() {
		log.Println("No such field!")
	}
	if !structFieldValue.CanSet() {
		log.Println("cannot set field name")
	}
	structFieldType := structFieldValue.Type()
	log.Println(structFieldType)
	val := reflect.ValueOf("car")
	log.Println(val)
	if structFieldType != val.Type() {
		log.Println("provided value does not match obj field type")
	}
	structFieldValue.Set(val)
	log.Println(p, p.Name)

	// Using Hashicorp's library to convert map to struct
	// Example struct
	people = make(map[string]interface{})
	people["name"] = "car" // Lowercase works
	people["Age"] = 1

	peeps := People{}
	err := mapstructure.Decode(people, &peeps)
	if err != nil {
		log.Println(err)
	}
	log.Println("peeps:", peeps)
}
