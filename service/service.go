package service

import (
	"time"

	"github.com/andyfusniak/monolith/internal/store"
	"github.com/andyfusniak/monolith/internal/store/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	store store.Repository
}

type Option func(*Service)

// New constructs a new service from the given Options. At a minimum New
// must be called using the WithSqlite3 configurator since the service
// requires a functional store to persist state.
func New(opts ...Option) *Service {
	service := &Service{}
	for _, o := range opts {
		o(service)
	}
	return service
}

// WithSqlite3 configures the service with sqlite3 as its store.
func WithSqlite3(st *sqlite3.Store) Option {
	return func(s *Service) {
		s.store = st
	}
}

const jsonTime = "2006-01-02T15:04:05.000Z07:00" // .000Z = keep trailing zeros

// ISOTime custom type to allow for JSON microsecond formating.
type ISOTime time.Time

// MarshalJSON provides microsecond formating
func (t ISOTime) MarshalJSON() ([]byte, error) {
	vt := time.Time(t)
	vt = vt.UTC().Round(time.Millisecond)
	return []byte(vt.Format(`"` + jsonTime + `"`)), nil
}
