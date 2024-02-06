	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/andyfusniak/monolith/internal/app"
	"github.com/andyfusniak/monolith/internal/env"
	"github.com/andyfusniak/monolith/internal/store/sqlite3"
	"github.com/andyfusniak/monolith/service"
	"github.com/spf13/cobra"
)

const (
	defaultMaxOpenConns int = 120
	defaultMaxIdleConns int = 20
)

// NewCmdServer creates a new server command. This command starts the web service.
func NewCmdServer(version, gitcommit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "start the web service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// environment
			cfg, err := env.EnvToConfig()
			if err != nil {
				return err
			}
			if cfg.HasWarnings() {
				for _, warning := range cfg.Warnings() {
					log.Warnf("[main] warning %s", warning)
				}
			}
			if cfg.HasErrors() {
				for _, e := range cfg.Errors() {
					log.Errorf("[main] %s", e)
				}
				if cfg.IsFatalErr() {
					log.Errorf("[main] exiting as configuration failed")
					os.Exit(1)
				}
			}

			// set up logging
			defer func() {
				log.Infof("[main] goodbye from monolith version %s (%s)", version, gitcommit)
			}()
			initLogging(cfg.App.LogLevel)
			log.Infof("[main] hello from monolith version %s (%s) %s for %s %s",
				version, gitcommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)

			// database connection
			// one read-only with high concurrency
			// one read-write for non-concurrent queries
			rw, err := sqlite3.OpenDB(cfg.DBFilepath)
			if err != nil {
				return err
			}
			defer rw.Close()
			rw.SetMaxOpenConns(1)
			rw.SetMaxIdleConns(1)
			rw.SetConnMaxIdleTime(5 * time.Minute)

			ro, err := sqlite3.OpenDB(cfg.DBFilepath)
			if err != nil {
				return err
			}
			defer ro.Close()
			ro.SetMaxOpenConns(defaultMaxOpenConns)
			ro.SetMaxIdleConns(defaultMaxIdleConns)
			ro.SetConnMaxIdleTime(5 * time.Minute)

			// store and service
			store := sqlite3.NewStore(ro, rw)
			svc := service.New(service.WithStore(store))

			// HTTP application server
			app, err := app.New(cfg.App, app.WithService(svc))
			if err != nil {
				return err
			}
			if err := app.Start(context.Background()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func initLogging(logLevel string) {
	// Output logs with colour
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Log debug level severity or above.
	logrusLevel := logLevelToLogrusLevel(logLevel)
	log.SetLevel(logrusLevel)
}

func logLevelToLogrusLevel(v string) log.Level {
	switch v {
	case "panic":
		return log.PanicLevel
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.DebugLevel
	}
}
