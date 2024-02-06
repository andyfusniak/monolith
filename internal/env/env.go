package env

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Config application configuration.
type Config struct {
	DBFilepath string
	App        AppConfig
	errors     []string
	warnings   []string
	errFatal   bool
}

type AppConfig struct {
	Port      string
	LogLevel  string
	CWD       string
	BaseURL   string
	IsDevMode bool
}

// HasWarnings returns true if there are any warnings.
func (c *Config) HasWarnings() bool {
	return len(c.warnings) > 0
}

// HasErrors return true if there are any errors.
func (c *Config) HasErrors() bool {
	return len(c.errors) > 0
}

// IsFatalErr returns true if any of the errors should terminate the service.
func (c *Config) IsFatalErr() bool {
	return c.errFatal
}

// Warnings returns configuration warnings. Warnings might prevent the
// application server functioning completely.
func (c *Config) Warnings() []string {
	return c.warnings
}

// Errors returns a slice of error messages that occured whilst reading
// the environment configuration.
func (c *Config) Errors() []string {
	return c.errors
}

// EnvToConfig reads the environment variables and creates a single Config.
func EnvToConfig() (*Config, error) {
	var cfg Config

	// current working directory
	// used to locate templates relative to the running binary
	var err error
	cfg.App.CWD, err = cwd()
	if err != nil {
		return nil, err
	}

	// PORT (optional)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg.App.Port = port

	// LOG_LEVEL (optional)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	if !isValidLogLevel(logLevel) {
		cfg.errors = append(cfg.errors, fmt.Sprintf("LOG_LEVEL %s is not valid", logLevel))
		cfg.errFatal = true
	}
	cfg.App.LogLevel = logLevel

	// DBFilepath full path to the sqlite3 database file.
	dbfilepath, found := os.LookupEnv("DB_FILEPATH")
	if !found {
		cfg.errors = append(cfg.errors, "DB_FILEPATH environment variable not set")
		cfg.errFatal = true
	}
	cfg.DBFilepath = dbfilepath

	return &cfg, nil
}

func cwd() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "[router] failed to getcwd")
	}
	return d, nil
}

func isValidLogLevel(v string) bool {
	switch v {
	case "panic", "fatal", "error", "warn", "info", "debug", "trace":
		return true
	default:
		return false
	}
}
