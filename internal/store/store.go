package store

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"
)

const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00" // .000000Z = keep trailing zeros

type Datetime time.Time

func (t *Datetime) Scan(v any) error {
	vt, err := time.Parse(RFC3339Micro, v.(string))
	if err != nil {
		return err
	}
	*t = Datetime(vt)
	return nil
}

func (t *Datetime) Value() (driver.Value, error) {
	return time.Time(*t).UTC().Format(RFC3339Micro), nil
}

// Repository store operations.
type Repository interface {
	UsersRepository
}

// user repository

var (
	ErrUserNotFound = errors.New("user not found")
)

// UsersRepository defines the user store operations.
type UsersRepository interface {
	InsertUser(ctx context.Context, params AddUser) (User, error)
	GetUser(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type AddUser struct {
	UserID       string
	Email        string
	PasswordHash string
}

type User struct {
	UserID       string
	Email        string
	PasswordHash string
	CreatedAt    Datetime
}

type InsertUserParams struct {
	UserID       string
	Email        string
	PasswordHash string
}
