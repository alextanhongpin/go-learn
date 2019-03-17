# Basic query

```go
package main

import "database/sql"

func main() {
	res, err := db.Exec(stmt, param)

	id, err := res.LastInsertId() // (int64, error)

	count, err := res.RowsAffected() // (int64, error)
}

type User struct {
	Name string
	ID   string
}

func GetUser(db *sql.DB) (*User, error) {
	var u User
	err := db.QueryRow(`SELECT name FROM user WHERE id = ? LIMIT 1`, 1).Scan(&u.Name)
	return &u, err
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.QueryRow(`SELECT name FROM user`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Name); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

```

## Sample sql builder in golang
```go
package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const FuncSeparator = ":"
const Placeholder = "?"

func hasFunction(key string) (field, fn string, exist bool) {
	if paths := strings.Split(key, FuncSeparator); len(paths) == 2 {
		field = paths[0]
		fn = paths[1]
		exist = true
		return
	}
	return
}

type SQLBuilder struct {
	sb     strings.Builder
	values []interface{}
}

func stripFunction(key string) string {
	if paths := strings.Split(key, FuncSeparator); len(paths) == 2 {
		return paths[0]
	}
	return key
}

func and(keys ...string) string {
	return strings.Join(keys, " AND ")
}

func updateFields(keys ...string) string {
	return strings.Join(mapString(keys, updateWithFunction), ", ")
}

func insertFields(keys ...string) string {
	return strings.Join(mapString(keys, stripFunction), ", ")
}

func insertPlaceholders(keys ...string) string {
	return strings.Join(mapString(keys, placeholderWithFunction), ", ")
}

func updateWithFunction(key string) string {

	if name, fn, exist := hasFunction(key); exist {
		return fmt.Sprintf("%s = %s(%s)", name, fn, Placeholder)
	}
	return fmt.Sprintf("%s = %s", key, Placeholder)
}

func placeholderWithFunction(key string) string {
	if _, fn, exist := hasFunction(key); exist {
		return fmt.Sprintf("%s(?)", fn)
	}
	return Placeholder
}

func (s *SQLBuilder) Select(fields ...string) *SQLBuilder {
	s.sb.Reset()
	sort.Strings(fields)
	keys := strings.Join(fields, ", ")
	s.sb.WriteString(fmt.Sprintf("SELECT (%s)", keys))
	return s
}

func (s *SQLBuilder) From(table string) *SQLBuilder {
	prev := s.sb.String()
	s.sb.Reset()
	s.sb.WriteString(fmt.Sprintf("%s FROM %s", prev, table))
	return s
}

func (s *SQLBuilder) InsertInto(table string, fields SQLFields) *SQLBuilder {
	// Since this is the first statement, clear the existing values.
	s.values = s.values[:0]
	s.sb.Reset()
	k, v := fields.KeyValues()
	s.values = append(s.values, v...)
	out := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		insertFields(k...),
		insertPlaceholders(k...))
	s.sb.WriteString(out)
	return s
}

func (s *SQLBuilder) OnDuplicateKeyUpdate(fields SQLFields) *SQLBuilder {
	s.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	// NOTE: Map orders will change with every iteration..., the order for keys and values will be different if called separately.
	// keys := fields.MapKeys(updateWithFunction)
	// values := fields.Values()
	k, v := fields.KeyValues()
	s.values = append(s.values, v...)
	s.sb.WriteString(updateFields(k...))
	return s
}

func (s *SQLBuilder) Where(conditions SQLFields) *SQLBuilder {
	k, v := conditions.KeyValues()
	keys := mapString(k, updateWithFunction)
	s.values = append(s.values, v...)
	prev := s.sb.String()
	s.sb.Reset()
	s.sb.WriteString(fmt.Sprintf("%s WHERE %s", prev, and(keys...)))
	return s
}

type SQLFields map[string]interface{}

func (s SQLFields) KeyValues() ([]string, []interface{}) {
	var keys []string
	for k, v := range s {
		if v == "" || v == nil {
			continue
		}
		keys = append(keys, k)

	}
	// Sort to enable prepare statement cache
	sort.Strings(keys)

	var values []interface{}
	for _, key := range keys {
		values = append(values, s[key])
	}
	return keys, values
}

func mapString(in []string, fn func(string) string) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

func main() {
	{
		sb := SQLBuilder{}
		s := sb.InsertInto("users", SQLFields{
			"name":   "hello",
			"age":    10,
			"car":    "",
			"id:HEX": "34",
		}).OnDuplicateKeyUpdate(SQLFields{
			"name":       "hello",
			"updated_at": time.Now(),
		})
		fmt.Println("====TEST INSERT WITH DUPLICATE KEYS====")
		fmt.Println(s.sb.String(), s.values)
	}
	{
		fmt.Println("====TEST SELECT====")
		sb := SQLBuilder{}
		s := sb.Select("name", "age").From("users").Where(SQLFields{"age": 100, "id:HEX": "100"})
		fmt.Println(s.sb.String(), s.values)
	}
}

//TODO: Where, limit, sort
```
