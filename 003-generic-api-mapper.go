// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type registry map[reflect.Type]func(reflect.Value) any

var r = make(registry)

func load(v any) func(reflect.Value) any {
	fn, ok := r[reflect.TypeOf(v)]
	if !ok {
		panic(fmt.Errorf("type %T is not registered", v))
	}
	return fn
}

func store[T, V any](fn func(T) V) {
	var tt T
	var t = reflect.TypeOf(tt)
	r[t] = func(v reflect.Value) any {
		in, ok := reflect.TypeAssert[T](v)
		if !ok {
			panic(fmt.Errorf("want %T, got %v", tt, v))
		}
		return fn(in)
	}
}

func main() {
	store(toUserAPI)
	ent := &UserEntity{Name: "John"}
	debugOutput(mapper(ent))
	debugOutput(mapper(*ent))
	debugOutput(mapper([]UserEntity{*ent}))
	debugOutput(mapper([]*UserEntity{ent}))
}

func debugOutput(v any) {
	fmt.Printf("type is: %T\n", v)
	fmt.Printf("value is: %#v\n", v)
	fmt.Println()
}

func mapper(val any) any {
	v := reflect.ValueOf(val)

	fmt.Println(v.Kind())
	// Check if the Kind of the value is a slice
	switch v.Kind() {
	case reflect.Slice:
		var result []any
		for i := 0; i < v.Len(); i++ {
			result = append(result, mapper(v.Index(i).Interface()))
		}
		return result
	case reflect.Struct:
		nv := reflect.New(v.Type())
		nv.Elem().Set(v)
		return mapper(nv.Interface())
	case reflect.Pointer:
		el := v.Elem()
		if el.Kind() != reflect.Struct {
			panic("not a struct")
		}
		return load(val)(v)
	}
	return nil
}

func toUserAPIValue(val reflect.Value) any {
	v, ok := reflect.TypeAssert[*UserEntity](val)
	if !ok {
		panic(fmt.Errorf("want %T, got %T", &UserEntity{}, val))
	}
	return toUserAPI(v)
}

func toUserAPI(e *UserEntity) *UserAPI {
	return &UserAPI{
		Name: e.Name,
	}
}

type UserEntity struct {
	Name  string
	Age   int
	Hobby []*Hobby
}

type UserAPI struct {
	Name  string
	Age   int
	Hobby []*HobbyAPI
}

type Hobby struct {
	Name string
}

type HobbyAPI struct {
	Name string
}
