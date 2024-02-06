package schema

import (
	"embed"
)

// Migrations used to embed the SQL init files.
//
//go:embed migrations
var Migrations embed.FS
