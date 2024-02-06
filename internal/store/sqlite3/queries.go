package sqlite3

import (
	"context"
	"database/sql"
)

// DBTx common database operations.
type DBTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Queries allows single and transactional queries.
type Queries struct {
	readwrite DBTx
	readonly  DBTx
}

// NewQueries create a new comments query.
func NewQueries(ro, rw DBTx) *Queries {
	return &Queries{
		readonly:  ro,
		readwrite: rw,
	}
}
