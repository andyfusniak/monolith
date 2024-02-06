package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/andyfusniak/monolith/internal/store"
	"github.com/pkg/errors"
)

// Store memberships store.
type Store struct {
	*Queries
	rw *sql.DB
}

// NewStore returns a new store.
func NewStore(ro, rw *sql.DB) store.Repository {
	return &Store{
		rw:      rw,
		Queries: NewQueries(ro, rw),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.rw.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries(tx, tx) // read only queries
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("[sqlite3] tx rollback failed: %v: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// users

// InsertUser adds a new user row to the users table.
func (q *Queries) InsertUser(ctx context.Context, params store.AddUser) (store.User, error) {
	const query = `
insert into users
  (user_id, email, password_hash, created_at)
values
  (:user_id, :email, :password_hash, :created_at)
returning
  user_id, email, password_hash, created_at
`
	r := store.User{}
	now := store.Datetime(time.Now().UTC())
	if err := q.readwrite.QueryRowContext(ctx, query,
		sql.Named("user_id", params.UserID),             // :user_id
		sql.Named("email", params.Email),                // :email
		sql.Named("password_hash", params.PasswordHash), // :password_hash
		sql.Named("created_at", &now),                   // :created_at
	).Scan(
		&r.UserID,       // 0 user_id
		&r.Email,        // 1 email
		&r.PasswordHash, // 2 password_hash
		&r.CreatedAt,    // 3 created_at
	); err != nil {
		return store.User{}, errors.Wrapf(err,
			"[sqlite3:users] query row scan failed query=%q", query)
	}

	return r, nil
}

// GetUser gets a user row by primary key.
func (q *Queries) GetUser(ctx context.Context, userID string) (store.User, error) {
	const query = `
select
  user_id, email, password_hash, created_at
from users
where user_id = :user_id
`
	r := store.User{}
	if err := q.readonly.QueryRowContext(ctx, query,
		sql.Named("user_id", userID), // :user_id
	).Scan(
		&r.UserID,       // 0 user_id
		&r.Email,        // 1 email
		&r.PasswordHash, // 2 password_hash
		&r.CreatedAt,    // 3 created_at
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.User{}, store.ErrUserNotFound
		}

		return store.User{}, errors.Wrapf(err,
			"[sqlite3:users] query row scan failed query=%q", query)
	}

	return r, nil
}

// GetUserByEmail gets a user row by email.
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (store.User, error) {
	const query = `
select
  user_id, email, password_hash, created_at
from users
where email = :email
`
	r := store.User{}
	if err := q.readonly.QueryRowContext(ctx, query,
		sql.Named("email", email), // :email
	).Scan(
		&r.UserID,       // 0 user_id
		&r.Email,        // 1 email
		&r.PasswordHash, // 2 password_hash
		&r.CreatedAt,    // 3 created_at
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.User{}, store.ErrUserNotFound
		}

		return store.User{}, errors.Wrapf(err, "[sqlite3:users] query row scan failed query=%q", query)
	}

	return r, nil
}
