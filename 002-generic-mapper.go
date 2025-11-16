// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	Register(toInt)
	Register(toAge)
	Register(toUserAPI)
	fmt.Println(Map[int]("1234"))
	fmt.Println(Map[Age](100))
	fmt.Println(Map[*UserAPI](&UserDB{
		Name: "John",
	}))
	fmt.Println("Hello, 世界")
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

type Age int

func toAge(n int) Age {
	return Age(n)
}

type UserDB struct {
	Name string
}
type UserAPI struct {
	Name string `json:"name"`
}

func toUserAPI(dto *UserDB) *UserAPI {
	return &UserAPI{
		Name: dto.Name,
	}
}

type registry struct {
	m map[reflect.Type]map[reflect.Type]func(reflect.Value) reflect.Value
}

var r registry

func init() {
	r.m = make(map[reflect.Type]map[reflect.Type]func(reflect.Value) reflect.Value)
}

func Map[T any](val any) T {
	var tt = val
	var vv T
	t, v := reflect.TypeOf(tt), reflect.TypeOf(vv)
	m, ok := r.m[t]
	if !ok {
		panic(fmt.Errorf("input type %T not registered", tt))
	}
	fn, ok := m[v]
	if !ok {
		panic(fmt.Errorf("output type %T not registered", vv))
	}
	out := fn(reflect.ValueOf(val))
	res, ok := reflect.TypeAssert[T](out)
	if !ok {
		panic(fmt.Errorf("want %T, got %T", vv, res))
	}
	return res
}

func Register[T, V any](fn func(T) V) {
	var tt T
	var vv V
	t, v := reflect.TypeOf(tt), reflect.TypeOf(vv)
	if _, ok := r.m[t]; !ok {
		r.m[t] = make(map[reflect.Type]func(reflect.Value) reflect.Value)
	}
	r.m[t][v] = func(val reflect.Value) reflect.Value {
		t, ok := reflect.TypeAssert[T](val)
		if !ok {
			panic(fmt.Errorf("want %T, got %T", tt, val))
		}
		return reflect.ValueOf(fn(t))
	}
}
