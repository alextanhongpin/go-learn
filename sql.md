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
