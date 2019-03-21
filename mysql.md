## Dynamic SQL builder

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

func flatten(interfaces ...[]interface{}) (result []interface{}) {
	for _, values := range interfaces {
		for _, value := range values {
			result = append(result, value)
		}
	}
	return
}

func mapString(values []string, fn func(string) string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}

type SQLFields map[string]interface{}

func (s SQLFields) Keys() []string {
	var keys []string
	for key, value := range s {
		if value != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s SQLFields) Values() []interface{} {
	var values []interface{}
	for _, value := range s {
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func (s SQLFields) Update(initial string) (string, error) {
	keys := s.Keys()
	if len(keys) == 0 {
		return "", errors.New("values cannot be empty")
	}
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = fmt.Sprintf("%s = ?", key)
	}
	return fmt.Sprintf(initial, strings.Join(result, ", ")), nil
}

func (s SQLFields) Insert(initial string) (string, error) {
	keys := s.Keys()
	if len(keys) == 0 {
		return "", errors.New("fields cannot be empty")
	}
	placeholders := mapString(keys, func(s string) string {
		return "?"
	})
	return fmt.Sprintf(initial,
		strings.Join(keys, ", "),
		strings.Join(placeholders, ", "),
	), nil
}

func (s SQLFields) InsertOnDuplicateKeys(initial string) (string, error) {
	keys := s.Keys()
	if len(keys) == 0 {
		return "", errors.New("fields cannot be empty")
	}

	placeholders := mapString(keys, func(s string) string {
		return "?"
	})
	updates := mapString(keys, func(s string) string {
		return fmt.Sprintf("%s = ?", s)
	})
	return fmt.Sprintf(initial,
		strings.Join(keys, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	), nil
}

func main() {
	f := SQLFields{
		"name":   "john",
		"age":    100,
		"career": "",
	}
	{
		output, _ := f.Update("UPDATE employee SET %s WHERE name = ?")
		fmt.Println("stmt:", output)
	}
	{
		output, _ := f.Insert("INSERT INTO employee (id, %s) VALUES (HEX(?), %s)")
		fmt.Println("stmt:", output)
	}
	{
		output, _ := f.InsertOnDuplicateKeys("INSERT INTO employee (id, %s) VALUES (HEX(?), %s) ON DUPLICATE KEY UPDATE %s")
		fmt.Println("stmt:", output)
		fmt.Println(flatten([]interface{}{"id"}, f.Values(), f.Values()))
	}
	f = SQLFields{}
	{
		output, err := f.Update("UPDATE employee SET %s WHERE name = ?")
		fmt.Println("stmt:", output, err)
	}

	fmt.Println(flatten([]interface{}{1, 2, 3}, []interface{}{"car", "paper"}))

}
```

## Alternative

```go
package main

import (
	"fmt"
	"strings"
)

func flatten(interfaces ...[]interface{}) (result []interface{}) {
	for _, values := range interfaces {
		for _, value := range values {
			result = append(result, value)
		}
	}
	return
}

type SQLFields map[string]interface{}

func (s SQLFields) keys() []string {
	var keys []string
	for key, value := range s {
		if value != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s SQLFields) MapKeys(fn func(string) string) string {
	keys := s.keys()
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = fn(key)
	}
	return strings.Join(result, ", ")
}

func (s SQLFields) Keys() string {
	return strings.Join(s.keys(), ", ")
}

func (s SQLFields) Values() []interface{} {
	var values []interface{}
	for _, value := range s {
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func (s SQLFields) Update() string {
	return s.MapKeys(func(s string) string {
		return fmt.Sprintf("%s = ?", s)
	})
}

func (s SQLFields) Placeholder() string {
	return s.MapKeys(func(s string) string {
		return "?"
	})
}

func main() {
	f := SQLFields{
		"name":   "john",
		"age":    100,
		"career": "",
	}

	fmt.Println("keys:", f.Keys())
	fmt.Println("values:", f.Values())
	fmt.Println("update:", f.Update())
	fmt.Println("Placeholder:", f.Placeholder())
	fmt.Printf("INSERT INTO employee (%s) VALUES (%s)\n", f.Keys(), f.Placeholder())
	fmt.Printf("INSERT INTO employee (%s) VALUES (%s)\n",
		f.MapKeys(func(s string) string { return s }),
		f.MapKeys(func(s string) string { return "?" }),
	)

	fmt.Println(flatten([]interface{}{1, 2, 3}, []interface{}{"car", "paper"}))
	fmt.Println(flatten(nil))
}
```


## Another alternative

```go
package main

import (
	"fmt"
	"sort"
	"strings"
)

type Field interface {
	Name() string
	Value() interface{}
}

type SqlString struct {
	name  string
	value string
}

func (s SqlString) Name() string {
	return s.name
}

func (s SqlString) Value() interface{} {
	return s.value
}

type SqlInt64 struct {
	name  string
	value int64
}

func (s SqlInt64) Name() string {
	return s.name
}

func (s SqlInt64) Value() interface{} {
	return s.value
}

type SqlFloat64 struct {
	name  string
	value float64
}

func (s SqlFloat64) Name() string {
	return s.name
}

func (s SqlFloat64) Value() interface{} {
	return s.value
}

type SqlBool struct {
	name  string
	value bool
}

func (s SqlBool) Name() string {
	return s.name
}

func (s SqlBool) Value() interface{} {
	return s.value
}

func NamedString(name, value string) SqlString {
	return SqlString{name, value}
}

func NamedInt64(name string, value int64) SqlInt64 {
	return SqlInt64{name, value}
}

func MapString(keys []string, fn func(string) string) []string {
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = fn(key)
	}
	return result
}

func Flatten(interfaces ...[]interface{}) []interface{} {
	var result []interface{}
	for _, i := range interfaces {
		result = append(result, i...)
	}
	return result
}

func Parse(fields ...Field) ([]string, []interface{}) {
	SortFields(fields)
	keys := make([]string, len(fields))
	values := make([]interface{}, len(fields))
	for i, field := range fields {
		keys[i] = field.Name()
		values[i] = field.Value()
	}
	return keys, values
}
func SortFields(fields []Field) {
	sort.Slice(fields, func(i, j int) bool {
		// Left is smaller than right - ascending order.
		return fields[i].Name() < fields[j].Name()
	})
}

func MapFields(fields []Field, fn func(Field) string) []string {
	SortFields(fields)
	result := make([]string, len(fields))
	for i, f := range fields {
		result[i] = fn(f)
	}
	return result
}

func Equals(field string) string {
	return fmt.Sprintf("%s = ?", field)
}

func Where(fields ...string) string {
	result := MapString(fields, Equals)
	sort.Strings(result)
	return strings.Join(result, " AND ")
}
func Select(fields ...string) string {
	sort.Strings(fields)
	return strings.Join(fields, ", ")
}
func main() {
	fields := []Field{NamedString("name", "john"), NamedInt64("age", 10)}
	keys, values := Parse(fields...)
	fmt.Println(keys, values)
	placeholders := MapString(keys, func(string) string {
		return "?"
	})
	fmt.Println(placeholders)

	equals := MapString(keys, Equals)
	fmt.Println(equals)
	mappedKeys := MapFields(fields, func(f Field) string {
		name := f.Name()
		switch f.(type) {
		case SqlString:
			return fmt.Sprintf(`SET %s = COALESCE(NULLIF(NULLIF(?, %s), ""), %s)`, name, name, name)
		case SqlInt64:
			return fmt.Sprintf(`SET %s = COALESCE(NULLIF(NULLIF(?, %s), 0), %s)`, name, name, name)
		default:
			return ""
		}
	})
	fmt.Println(mappedKeys)

	result := Flatten(values, values)
	fmt.Println(result)

	fmt.Println(Where("name", "age", "car"))
}
```
