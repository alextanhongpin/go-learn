```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

/** Migrations

create table users (
	id int generated always as identity,
	name text not null,
	primary key (id),
	unique (name)
);

create table books (
	id int generated always as identity,
	title text,
	user_id int not null,
	primary key (id),
	foreign key (user_id) references users(id)
);

insert into users(name) values ('john appleseed');
insert into books(title, user_id) values ('the meaning of life', 1);
*/

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
)

type User struct {
	ID   int
	Name string
}

func (u *User) Fields() []any {
	return []any{&u.ID, &u.Name}
}

type Book struct {
	ID     int
	Title  string
	UserID int
}

func (b *Book) Fields() []any {
	return []any{&b.ID, &b.Title, &b.UserID}
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	id := 1

	u := &User{}
	if err := db.
		QueryRow(`SELECT * FROM users WHERE id = $1`, id).
		Scan(u.Fields()...); err != nil {
		log.Fatalf("failed to query row: %s", err)
	}
	fmt.Printf("%+v\n", u)

	{
		b := &Book{}
		u := &User{}

		// Are the orders of the columns fixed?
		fields := append(u.Fields(), b.Fields()...)
		stmt, args := WhereStmt(`
				SELECT u.*, b.*
				FROM users u
				JOIN books b ON (b.user_id = u.id)
			`,
			map[string]any{
				"u.id":      1,
				"b.user_id": 1,
			},
		)

		fmt.Println("stmt:", stmt)
		if err := db.
			QueryRow(stmt, args...).
			Scan(fields...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}
		fmt.Printf("user: %+v\n", u)
		fmt.Printf("book: %+v\n", b)
	}
}

func Where(kv map[string]any) (string, []any) {
	where := make([]string, 0, len(kv))
	args := make([]any, 0, len(kv))
	for k, v := range kv {
		where = append(where, fmt.Sprintf("%s = $%d", k, len(args)+1))
		args = append(args, v)
	}

	return strings.Join(where, " AND "), args
}

func WhereStmt(stmt string, kv map[string]any) (string, []any) {
	where, args := Where(kv)

	return fmt.Sprintf("%s WHERE %s", stmt, where), args
}
```

## Dynamic?

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

/** Migrations

create table users (
	id int generated always as identity,
	name text not null,
	primary key (id),
	unique (name)
);

create table books (
	id int generated always as identity,
	title text,
	user_id int not null,
	primary key (id),
	foreign key (user_id) references users(id)
);

insert into users(name) values ('john appleseed');
insert into books(title, user_id) values ('the meaning of life', 1);
*/

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
)

type User struct {
	ID   int
	Name string
}

func (u *User) Fields() []any {
	return []any{&u.ID, &u.Name}
}

type Book struct {
	ID     int
	Title  string
	UserID int
}

func (b *Book) Fields() []any {
	return []any{&b.ID, &b.Title, &b.UserID}
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	id := 1

	u := &User{}
	if err := db.
		QueryRow(`SELECT * FROM users WHERE id = $1`, id).
		Scan(u.Fields()...); err != nil {
		log.Fatalf("failed to query row: %s", err)
	}
	fmt.Printf("%+v\n", u)

	{
		b := &Book{}
		u := &User{}

		// Are the orders of the columns fixed?
		fields := append(u.Fields(), b.Fields()...)
		stmt, args := WhereStmt(`
				SELECT u.*, b.*
				FROM users u
				JOIN books b ON (b.user_id = u.id)
			`,
			map[string]any{
				"u.id":      1,
				"b.user_id": 1,
			},
		)

		fmt.Println("stmt:", stmt)
		if err := db.
			QueryRow(stmt, args...).
			Scan(fields...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}
		fmt.Printf("user: %+v\n", u)
		fmt.Printf("book: %+v\n", b)
	}

	{
		qb := NewQueryBuilder(&User{})

		stmt, args, entity, columns := qb.
			SetAlias("u").
			SetTable("users").
			Where(map[string]any{
				"id": id,
			}).
			OrderBy("id", "name").
			SetLimit(10).
			SetOffset(0).
			BuildSelect()

		fmt.Println(stmt, args)
		if err := db.
			QueryRow(stmt, args...).
			Scan(columns...); err != nil {
			log.Fatalf("failed to query row: %s", err)
		}

		fmt.Printf("%+v\n", entity)
	}
}

func Where(kv map[string]any) (string, []any) {
	where := make([]string, 0, len(kv))
	args := make([]any, 0, len(kv))
	for k, v := range kv {
		where = append(where, fmt.Sprintf("%s = $%d", k, len(args)+1))
		args = append(args, v)
	}

	return strings.Join(where, " AND "), args
}

func WhereStmt(stmt string, kv map[string]any) (string, []any) {
	where, args := Where(kv)

	return fmt.Sprintf("%s WHERE %s", stmt, where), args
}

type QueryBuilder[T any] struct {
	entity  T
	alias   string
	table   string
	stmt    string
	columns []string
	where   map[string]any
	order   []string
	limit   *int
	offset  *int
}

func NewQueryBuilder[T any](t T) *QueryBuilder[T] {
	table := strings.ToLower(getType(t))
	return &QueryBuilder[T]{
		entity:  t,
		alias:   table,
		table:   table, // Naive.
		columns: getFields(t, "sql"),
		where:   map[string]any{},
		order:   []string{},
	}
}

func (qb *QueryBuilder[T]) SetAlias(alias string) *QueryBuilder[T] {
	if alias == "" {
		panic("qb: alias cannot be empty")
	}
	qb.alias = alias
	return qb
}

func (qb *QueryBuilder[T]) SetTable(table string) *QueryBuilder[T] {
	if table == "" {
		panic("qb: table cannot be empty")
	}
	qb.table = table
	return qb
}

func (qb *QueryBuilder[T]) Where(m map[string]any) *QueryBuilder[T] {
	for k, v := range m {
		qb.where[k] = v
	}

	return qb
}

func (qb *QueryBuilder[T]) OrderBy(order ...string) *QueryBuilder[T] {
	for _, v := range order {
		qb.order = append(qb.order, v)
	}

	return qb
}

func (qb *QueryBuilder[T]) SetLimit(limit int) *QueryBuilder[T] {
	if limit == 0 {
		panic("qb: limit cannot be 0")
	}
	qb.limit = &limit
	return qb
}

func (qb *QueryBuilder[T]) SetOffset(offset int) *QueryBuilder[T] {
	qb.offset = &offset
	return qb
}

func (qb *QueryBuilder[T]) BuildSelect() (string, []any, T, []any) {
	columns := make([]string, len(qb.columns))
	for i, col := range qb.columns {
		if strings.Contains(col, ".") {
			columns[i] = col
		} else {
			columns[i] = fmt.Sprintf("%s.%s", qb.alias, col)
		}
	}

	where := make(map[string]any)
	for k, v := range qb.where {
		if strings.Contains(k, ".") {
			where[k] = v
		} else {
			where[fmt.Sprintf("%s.%s", qb.alias, k)] = v
		}
	}

	stmt := fmt.Sprintf("SELECT %s FROM %s AS %s", strings.Join(columns, ", "), qb.table, qb.alias)
	stmt, args := WhereStmt(stmt, where)
	if len(qb.order) > 0 {
		order := make([]string, len(qb.order))
		for i, o := range qb.order {
			if strings.Contains(o, ".") {
				order[i] = o
			} else {
				order[i] = fmt.Sprintf("%s.%s", qb.alias, o)
			}
		}
		stmt = fmt.Sprintf("%s ORDER BY %s", stmt, strings.Join(order, ", "))
	}

	return stmt, args, qb.entity, getFieldAddress(qb.entity)
}

func getFieldAddress[T any](t T) []any {
	v := reflect.Indirect(reflect.ValueOf(t))

	fields := make([]any, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fields[i] = f.Addr().Interface()
	}

	return fields
}

func getType(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func getFields(t any, tag string) []string {
	v := reflect.Indirect(reflect.ValueOf(t)).Type()

	fields := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		name := f.Tag.Get(tag)
		if name == "" {
			name = f.Name
		}
		fields[i] = name
	}

	return fields
}
```
