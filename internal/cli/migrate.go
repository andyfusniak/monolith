package cli

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/andyfusniak/monolith/internal/store/sqlite3"
	"github.com/andyfusniak/monolith/internal/store/sqlite3/schema"

	"github.com/golang-migrate/migrate/v4"

	driversqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/spf13/cobra"
)

// NewCmdMigrate migrate sub command.
func NewCmdMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "database migration",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			dbfile := os.Getenv("DB_FILEPATH")
			if dbfile == "" {
				fmt.Fprint(app.stderr, "DB_FILEPATH not set\n")
				os.Exit(1)
			}

			db, err := sqlite3.OpenDB(dbfile)
			if err != nil {
				fmt.Fprint(app.stderr, "failed to open sqlite3 database file - check DB_FILEPATH\n")
				os.Exit(1)
			}

			db.SetMaxOpenConns(1)
			db.SetMaxIdleConns(1)
			db.SetConnMaxLifetime(5 * time.Minute)

			driver, err := driversqlite3.WithInstance(db, &driversqlite3.Config{NoTxWrap: true})
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed with instance %+v\n", err)
				os.Exit(1)
			}

			source, err := httpfs.New(http.FS(schema.Migrations), "migrations")
			if err != nil {
				fmt.Fprintf(app.stderr, "%+v\n", err)
				os.Exit(1)
			}

			mg, err := migrate.NewWithInstance("https", source, "sqlite3", driver)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get new migrate instance %+v\n", err)
				os.Exit(1)
			}

			app.db = db
			app.mg = mg
		},
	}

	cmd.AddCommand(NewCmdMigrateUp())
	cmd.AddCommand(NewCmdMigrateDown())
	return cmd
}

// NewCmdMigrateUp migrate up brings up the database schema.
func NewCmdMigrateUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "apply all or N up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			if err := app.mg.Up(); err != nil {
				fmt.Fprintf(os.Stderr, "migrate up failed: %+v\n", err)
			}

			app.db.Close()

			return nil
		},
	}
	return cmd
}

// NewCmdMigrateDown migrate down brings down the database schema.
func NewCmdMigrateDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "apply all or N down migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			if err := app.mg.Down(); err != nil {
				fmt.Fprintf(os.Stderr, "migrate down failed: %+v\n", err)
			}

			app.db.Close()

			return nil
		},
	}
	return cmd
}
