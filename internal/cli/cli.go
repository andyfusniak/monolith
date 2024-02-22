package cli

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/andyfusniak/monolith/internal/store/sqlite3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

type AppKey string

type App struct {
	version   string
	gitCommit string
	db        *sql.DB
	mg        *migrate.Migrate
	stdout    io.Writer
	stderr    io.Writer
}

type Option func(*App)

// NewApp creates a new CLI application.
func NewApp(options ...Option) *App {
	a := &App{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	for _, o := range options {
		o(a)
	}
	return a
}

// WithVersion option to set the cli version.
func WithVersion(s string) Option {
	return func(a *App) {
		a.version = s
	}
}

// WithGitCommit option to set the git commit hash.
func WithGitCommit(s string) Option {
	return func(a *App) {
		a.gitCommit = s
	}
}

// WithStdOut option to set default output stream.
func WithStdOut(w io.Writer) Option {
	return func(a *App) {
		a.stdout = w
	}
}

// WithStdErr option to set default error stream.
func WithStdErr(w io.Writer) Option {
	return func(a *App) {
		a.stderr = w
	}
}

// NewCmdInfo show the compile options used for the sqlite3
// driver.
func NewCmdInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info shows the sqlite3 driver compile options",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			// database
			db, err := sql.Open(sqlite3.DriverName, ":memory:")
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open sqlite3 in-memory client\n%+v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			const query = "select sqlite_version()"
			var sqliteVersion string
			if err := db.QueryRowContext(context.Background(), query).Scan(&sqliteVersion); err != nil {
				fmt.Fprintf(os.Stderr, "failed to select sqlite_version()\n%+v\n", err)
				os.Exit(1)
			}

			var builtWith string
			if sqlite3.DriverName == "sqlite" {
				builtWith = "pure go driver"
			} else {
				builtWith = "cgo driver"
			}

			fmt.Fprintf(app.stdout, "Running PRAGMA compile_options; (binary built with sqlite version %s %s)\n", sqliteVersion, builtWith)

			const pragma = "PRAGMA compile_options"
			rows, err := db.QueryContext(context.Background(), pragma)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to query context\n%+v\n", err)
				os.Exit(1)
			}
			defer rows.Close()

			for rows.Next() {
				var s string
				if err := rows.Scan(&s); err != nil {
					fmt.Fprintf(os.Stderr, "failed to scan row\n%+v\n", err)
				}
				fmt.Fprintf(os.Stdout, "%s\n", s)
			}
			if rows.Err() != nil {
				fmt.Fprintf(os.Stderr, "rows.Next err\n%+v\n", err)
				os.Exit(1)
			}

			return nil
		},
	}
	return cmd
}

// Version returns the cli application version.
func (a *App) Version() string {
	db, err := sql.Open(sqlite3.DriverName, ":memory:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open sqlite3 in-memory client\n%+v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	const query = "select sqlite_version()"
	var sqliteVersion string
	if err := db.QueryRowContext(context.Background(), query).Scan(&sqliteVersion); err != nil {
		fmt.Fprintf(os.Stderr, "failed to select sqlite_version()\n%+v\n", err)
		os.Exit(1)
	}

	var builtWith string
	if sqlite3.DriverName == "sqlite" {
		builtWith = "pure go driver"
	} else {
		builtWith = "cgo driver"
	}

	return fmt.Sprintf("%s (built with sqlite version %s %s)",
		a.version, sqliteVersion, builtWith)
}
