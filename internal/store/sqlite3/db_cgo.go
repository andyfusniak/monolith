//go:build cgo

package sqlite3

import (
	"database/sql"

	gosqlite3 "github.com/mattn/go-sqlite3"
)

const DriverName = "monolith_sqlite3"

func init() {
	sql.Register(DriverName,
		&gosqlite3.SQLiteDriver{
			ConnectHook: func(conn *gosqlite3.SQLiteConn) error {
				_, err := conn.Exec(`
				PRAGMA busy_timeout       = 10000;
				PRAGMA journal_mode       = WAL;
				PRAGMA journal_size_limit = 200000000;
				PRAGMA synchronous        = NORMAL;
				PRAGMA foreign_keys       = ON;
				PRAGMA temp_store         = MEMORY;
				PRAGMA cache_size         = -16000;
			`, nil)

				return err
			},
		},
	)
}

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open(DriverName, dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
