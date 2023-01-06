# Unit of Work with sqlc

Allows repositories to be cloned with the given transaction context, while fulfilling the repository interface.

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"database/sql"
	"fmt"
)

func main() {
	repo := &Queries{}
	repoCloner := &RepoCloner{q: repo}
	uow := &unitOfWork{}
	uc := &UseCase{
		repo:       repo,
		repoCloner: repoCloner,
		uow:        uow,
	}
	n, err := uc.Exec(context.Background())
	fmt.Println(n, err)
}

type unitOfWork struct{}

func (u *unitOfWork) BeginTx(fn func(tx *sql.Tx) error) error {
	tx := &sql.Tx{}
	return fn(tx)
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) Count(ctx context.Context) (int, error) {
	return 100, nil
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

type UOW interface {
	BeginTx(func(tx *sql.Tx) error) error
}

type Repo interface {
	Count(context.Context) (int, error)
}

type Cloner[T any] interface {
	Clone(tx *sql.Tx) T
}

type RepoCloner struct {
	q *Queries
}

func (r *RepoCloner) Clone(tx *sql.Tx) Repo {
	return r.q.WithTx(tx)
}

type UseCase struct {
	repo       Repo
	repoCloner Cloner[Repo]
	uow        UOW
}

func (u *UseCase) Exec(ctx context.Context) (int, error) {
	var count int
	if err := u.uow.BeginTx(func(tx *sql.Tx) error {
		repo := u.repoCloner.Clone(tx)
		n, err := repo.Count(ctx)
		if err != nil {
			return err
		}

		count = n
		return nil
	}); err != nil {
		return 0, err
	}

	return count, nil
}
```


## Another alternative

```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	uow := &UoW{db: db}
	repo := &postgresRepository{uow: uow}
	uc := &UseCase{repo: repo, uow: uow}
	n, err := uc.Num(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("1. n is:", n)
	n, err = uc.NumTx(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("2. n is:", n)
}

type Repository interface {
	Num(ctx context.Context) (int, error)
}

type postgresRepository struct {
	uow *UoW
}

func (r *postgresRepository) Num(ctx context.Context) (int, error) {
	var n int
	if err := r.uow.DBCtx(ctx).QueryRow(`select 1 + 1`).Scan(&n); err != nil {
		return 0, err
	}

	return n, nil
}

type UseCase struct {
	uow  *UoW
	repo Repository
}

func (uc *UseCase) Num(ctx context.Context) (int, error) {
	n, err := uc.repo.Num(ctx)
	return n, err
}

func (uc *UseCase) NumTx(ctx context.Context) (int, error) {
	var count int
	err := uc.uow.RunInTx(ctx, func(ctx context.Context) error {
		n, err := uc.repo.Num(ctx)
		if err != nil {
			return err
		}

		count = n
		return nil
	})
	return count, err
}

type UoW struct {
	db *sql.DB
	tx *sql.Tx
}

type uowKey string

var key = uowKey("uow")

func (key uowKey) With(ctx context.Context, uow *UoW) context.Context {
	return context.WithValue(ctx, key, uow)
}

func (key uowKey) Value(ctx context.Context) (*UoW, bool) {
	uow, ok := ctx.Value(key).(*UoW)
	return uow, ok
}

func (key uowKey) MustValue(ctx context.Context) *UoW {
	uow, ok := key.Value(ctx)
	if !ok {
		panic("uow: UnitOfWork context not found")
	}

	return uow
}

func (u *UoW) DBCtx(ctx context.Context) DBTX {
	if uow, ok := key.Value(ctx); ok {
		fmt.Println("uow.isTx", uow.IsTx())
		return uow.DB()
	}
	fmt.Println("uow.isTx", u.IsTx())

	return u.DB()
}

func (u *UoW) DB() DBTX {
	if u.IsTx() {
		return u.tx
	}

	return u.db
}

func (u *UoW) IsTx() bool {
	return u.tx != nil
}

func (u *UoW) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if u.IsTx() {
		return errors.New("uow: cannot nest transaction")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx = key.With(ctx, &UoW{tx: tx})
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
```
