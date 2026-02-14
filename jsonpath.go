package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrNotFound = errors.New("not found")

func main() {
	data := map[string]any{
		"name": map[string]any{
			"age":     1,
			"hobbies": []any{"1324", "2222"},
			"meta": []any{
				map[string]any{
					"foo": "bar",
					"bar": "foo",
				},
				map[string]any{
					"foo": "baz",
				},
			},
		},
	}
	fmt.Println(set(data, "name.hobbies[1]", "333333"))
	fmt.Println(set(data, "age", 12))
	fmt.Println(get(data, "name.hobbies[0]"))
	fmt.Println(get(data, "name.hobbies[1]"))
	fmt.Println(get(data, "name.age[1]"))
	fmt.Println(get(data, "name.hobbies.car.foo"))
	fmt.Println(set(data, "name.meta[].baz.far", "haha"))
	fmt.Println(get(data, "age"))
	fmt.Println(data)
}

func get(v any, path string) (any, error) {
	paths, err := parse(path)
	if err != nil {
		return nil, err
	}

	for i, p := range paths {
		v, err = getter(v, p)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, buildKey(paths[:i+1]))
		}
	}

	return v, nil
}

func set(v any, path string, value any) error {
	paths, err := parse(path)
	if err != nil {
		return err
	}

	for i, p := range paths[:len(paths)-1] {
		v, err = getter(v, p)
		if err != nil {
			return fmt.Errorf("%w: %s", err, buildKey(paths[:i+1]))
		}
	}

	last := paths[len(paths)-1]

	return setter(v, last, value)
}

func getter(data any, key any) (any, error) {
	switch v := key.(type) {
	case string:
		m, ok := data.(map[string]any)
		if !ok {
			return nil, errors.New("not a map")
		}
		value, ok := m[v]
		if !ok {
			return nil, ErrNotFound
		}
		return value, nil
	case int:
		s, ok := data.([]any)
		if !ok {
			return nil, errors.New("not a slice")
		}
		if len(s) < v {
			return nil, errors.New("index out of range")
		}
		return s[v], nil
	case []int:
		s, ok := data.([]any)
		if !ok {
			return nil, errors.New("not a slice")
		}
		return s, nil
	default:
		return nil, errors.ErrUnsupported
	}
}

func setter(data any, key any, value any) error {
	switch v := key.(type) {
	case string:
		if d, ok := data.([]any); ok {
			for _, dd := range d {
				if err := setter(dd, key, value); err != nil {
					return err
				}
			}
			return nil
		}

		m, ok := data.(map[string]any)
		if !ok {
			return errors.New("not a map")
		}
		m[v] = value
	case int:
		s, ok := data.([]any)
		if !ok {
			return errors.New("not a slice")
		}
		s[v] = value
	default:
		return fmt.Errorf("unknown key: %T", key)
	}

	return nil
}

func split(key string) ([]any, error) {
	i := strings.Index(key, "[")
	j := strings.Index(key, "]")
	if i > -1 && i < j && j == len(key)-1 {
		a := key[:i]
		b := key[i+1 : j]
		if b == "" {
			return []any{a, []int{}}, nil
		}
		n, err := strconv.Atoi(b)
		if err != nil {
			return nil, errors.New("invalid index")
		}
		if strconv.Itoa(n) != b {
			return nil, errors.New("invalid index")
		}
		if n < 0 {
			return nil, errors.New("negative index")
		}

		return []any{a, n}, nil
	}

	return []any{key}, nil
}

func parse(path string) ([]any, error) {
	var paths []any
	parts := strings.Split(path, ".")
	for i, p := range parts {
		part, err := split(p)
		if err != nil {
			key := strings.Join(parts[:i+1], ".")
			return nil, fmt.Errorf("%w: %s", err, key)
		}
		paths = append(paths, part...)
	}

	return paths, nil
}

func buildKey(parts []any) string {
	var sb strings.Builder
	for _, p := range parts {
		switch v := p.(type) {
		case string:
			sb.WriteString(v)
			sb.WriteString(".")
		case int:
			sb.WriteString(strconv.Itoa(v))
			sb.WriteString(".")
		}
	}
	return strings.TrimSuffix(sb.String(), ".")
}
