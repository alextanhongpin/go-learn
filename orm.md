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
