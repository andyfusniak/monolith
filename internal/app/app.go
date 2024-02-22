package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/andyfusniak/monolith/internal/env"
	"github.com/andyfusniak/monolith/internal/handler"

	log "github.com/sirupsen/logrus"

	"github.com/andyfusniak/monolith/service"
)

type App struct {
	svc     *service.Service
	cfg     env.AppConfig
	router  *http.ServeMux
	handler *handler.Handler
}

// Option is a function that configures an App.
type Option func(a *App)

// New creates a new app server.
func New(cfg env.AppConfig, opts ...Option) (*App, error) {
	app := &App{
		cfg: cfg,
	}
	for _, o := range opts {
		o(app)
	}

	// application handlers
	app.handler = handler.New(app.svc)

	// routing
	app.router = app.v1Routes()

	return app, nil
}

// WithService sets the service for the app.
func WithService(svc *service.Service) Option {
	return func(a *App) {
		a.svc = svc
	}
}

// Start the app server.
func (a *App) Start(ctx context.Context) error {
	// HTTP Service
	srv := http.Server{
		Addr:    "0.0.0.0:" + a.cfg.Port,
		Handler: a.router,
	}

	// signals
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		switch sig := <-sigint; sig {
		case syscall.SIGINT:
			log.Info("[main] received signal SIGINT")
		case syscall.SIGTERM:
			log.Info("[main] received signal SIGTERM")
		default:
			log.Info("[main] received unexpected signal", "signal", sig)
		}

		log.Info("[main] gracefully shutting down the server...")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Infof("[main] HTTP server Shutdown: %+v", err)
		}
		log.Info("[main] HTTP server shutdown complete")
		close(idleConnsClosed)
	}()

	log.Infof("[main] server listening on HTTP port %s", a.cfg.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("[main] error: %+v", err)
		return err
	}
	<-idleConnsClosed

	return nil
}
