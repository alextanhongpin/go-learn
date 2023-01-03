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
