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

