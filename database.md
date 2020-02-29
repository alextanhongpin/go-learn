## Fast json query

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Connect to database.
	db, err := database.New(...)
	if err != nil {
		log.Fatal(err)
	}

	var b []byte
	stmt := `
		SELECT JSON_OBJECT(
			"id", HEX(id),
			"email", email,
			"email_verified", IF(email_verified = 1, true, false) IS true
		) FROM employee
	`
	err = db.QueryRow(stmt).Scan(&b)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal only if you need to work with the data.
	// var m model.Employee
	// if err := json.Unmarshal(b, &m); err != nil {
	//         log.Fatal(err)
	// }
	// log.Println(m)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(b))
	})

	fmt.Println("listening to port *:4000")
	http.ListenAndServe(":4000", nil)
}
```

## With gin

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.New(...)
	if err != nil {
		log.Fatal(err)
	}

	var b []byte
	stmt := `
		SELECT JSON_OBJECT(
			"id", HEX(id),
			"email", email,
			"email_verified", IF(email_verified = 1, true, false) IS true
		) FROM employee
	`
	err = db.QueryRow(stmt).Scan(&b)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		// Send as string...
		// c.String(http.StatusOK, string(b))
		c.JSON(http.StatusOK, m)
	})
	r.Run()
}
```

## Transaction 


```go
package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Tx interface {
	Exec(query string, args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Prepare(query string) (*Stmt, error)
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	Query(query string, args ...interface{}) (*Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(query string, args ...interface{}) *Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
}

type TxFn func(Tx) error

func withTransaction(db *sql.DB, fn TxFn) error {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

func main() {

	fmt.Println("Hello, playground")
}
```

## Sample Unit of Work (?)


```go
package main

import (
	"context"
	"fmt"
)

type Tx interface {
	Commit()
}

type User struct {
	Name string
}
type Photo struct {
	Name string
}

type userRepository interface {
	Save(context.Context, User) error
}

type photoRepository interface {
	Save(context.Context) error
}

type UserRepository struct {
	tx Tx
}

func NewUserRepository(tx Tx) *UserRepository {
	return &UserRepository{tx: tx}
}

func (u *UserRepository) Save(ctx context.Context) error {
	return nil
}

type PhotoRepository struct {
	tx Tx
}

func NewPhotoRepository(tx Tx) *PhotoRepository {
	return &PhotoRepository{tx: tx}
}

func (p *PhotoRepository) Save(ctx context.Context) error {
	return nil
}

func main() {
	db := ... // Connect to database.
	err := withTransaction(db, func(tx Tx) error {
		ctx := context.Background()
		ur := NewUserRepository(tx)
		pr := NewPhotoRepository(tx)
		if err := ur.Save(ctx); err != nil {
			return err
		}
		if err := pr.Save(ctx); err != nil {
			return err
		}
		return nil
	})

	fmt.Println("Hello, playground")
}
```

## Alternative Transaction Pattern with Golang

```go
package main

import (
	"context"
	"fmt"
)

type Tx interface {
	Exec(query string, args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Prepare(query string) (*Stmt, error)
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	Query(query string, args ...interface{}) (*Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(query string, args ...interface{}) *Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
}

type UserRepository struct {
	tx Tx
}

func NewUserRepository(tx Tx) *UserRepository {
	return &UserRepository{tx: tx}
}

func (u *UserRepository) WithTx(tx Tx) *UserRepository {
	return NewUserRepository(tx)
}

type PhotoRepository struct {
	tx Tx
}

func NewPhotoRepository(tx Tx) *PhotoRepository {
	return &PhotoRepository{tx: tx}
}

func (p *PhotoRepository) WithTx(tx Tx) *PhotoRepository {
	return NewPhotoRepository(tx)
}

func main() {
	db := // Initialize db...
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	photoRepo := NewPhotoRepository(db)
	
	withTransaction(db, func(tx Tx) error {
		userTx := userRepo.WithTx(tx)
		photoTx := photoRepo.WithTx(tx)
		// Do something...
		return nil
	})
	fmt.Println("Hello, playground")
}
```
# Key-value store

- https://github.com/gomods/athens
- https://github.com/syndtr/goleveldb
- https://github.com/etcd-io/bbolt
- https://github.com/dgraph-io/badger

# Search

- https://github.com/blevesearch/bleve

# Cache

- https://github.com/coocood/freecache
- https://github.com/allegro/bigcache
- https://github.com/patrickmn/go-cache
