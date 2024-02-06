package service

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/andyfusniak/base58"
	"github.com/andyfusniak/monolith/internal/store"
	"github.com/pkg/errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserWrongPassword    = errors.New("wrong password")
	ErrUserPasswordTooShort = errors.New("password too short")
)

type User struct {
	ID        string  `json:"user_id"`
	Email     string  `json:"email"`
	CreatedAt ISOTime `json:"created_at"`
}

// CreateUser params.URole should be set to RoleUser or RoleAdmin.
func (s *Service) CreateUser(ctx context.Context, email, password string) (User, error) {
	if len(password) < 8 {
		return User{}, ErrUserPasswordTooShort
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return User{}, errors.Wrap(err, "[service] failed to create argon2id hash")
	}

	userID, err := base58.RandString(22) // 58**22 > 2**128
	if err != nil {
		return User{}, errors.Wrap(err, "[service] failed to generated random base58 string")
	}
	row, err := s.repo.InsertUser(ctx, store.AddUser{
		UserID:       userID,
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		return User{}, errors.Wrap(err, "[service] s.store.InsertUser failed")
	}

	return userFromRow(row), nil
}

// GetUser returns a single User with the given userID.
func (s *Service) GetUser(ctx context.Context, userID string) (User, error) {
	row, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, errors.Wrapf(err, "[service] s.store.GetUser(ctx, userID=%q) failed", userID)
	}

	return userFromRow(row), nil
}

// VerifyUserPassword accepts an email and password and returns a User entity
// if a user with the given email exists and the user's password matches.
//
// If the email is found but the password did not match an
// ErrUserWrongPassword is returned.
//
// If the email is not found an ErrUserNotFound is returned. You should not
// tell clients that their account could not be found otherwise
// VerifyUserPassword can be used to find valid emails. Instead return an
// authorized response to clients.
func (s *Service) VerifyUserPassword(ctx context.Context, email, password string) (User, error) {
	row, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			_, _ = argon2id.ComparePasswordAndHash(password, "dummyhash")
			return User{}, ErrUserNotFound
		}

		return User{}, errors.Wrap(err,
			"[service] s.store.VerifyUserPassword failed")
	}

	match, err := argon2id.ComparePasswordAndHash(password, row.PasswordHash)
	if err != nil {
		return User{}, errors.Wrap(err,
			"[service] failed to compare password and hash using argon2id")
	}
	if match {
		return userFromRow(row), nil
	}

	return User{}, ErrUserWrongPassword
}

func userFromRow(row store.User) User {
	return User{
		ID:        row.UserID,
		Email:     row.Email,
		CreatedAt: ISOTime(row.CreatedAt),
	}
}
