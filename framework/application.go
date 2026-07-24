// Package framework provides the reusable, application-agnostic core of the
// web framework: the Application lifecycle (start, signal handling,
// graceful shutdown). App-specific wiring lives in bootstrap.
package framework

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myapp/config"
)

// Application ties together configuration, logging and the HTTP server, and
// owns the process lifecycle (startup + graceful shutdown).
type Application struct {
	Config *config.Config
	Logger *log.Logger
	Server *http.Server

	// ShutdownTimeout bounds how long graceful shutdown may take before the
	// process is forced to exit.
	ShutdownTimeout time.Duration
}

// New builds an Application ready to Run.
func New(cfg *config.Config, logger *log.Logger, handler http.Handler) *Application {
	return &Application{
		Config: cfg,
		Logger: logger,
		Server: &http.Server{
			Addr:         cfg.App.Addr(),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		ShutdownTimeout: 10 * time.Second,
	}
}

// Run starts the HTTP server and blocks until a termination signal is
// received, at which point it gracefully shuts the server down.
func (a *Application) Run() error {
	serveErr := make(chan error, 1)

	go func() {
		a.Logger.Printf("%s starting on http://%s (env=%s)", a.Config.App.Name, a.Server.Addr, a.Config.App.Env)
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	case sig := <-quit:
		a.Logger.Printf("received signal %s, shutting down gracefully", sig)
		return a.Shutdown()
	}
}

// Shutdown gracefully stops the HTTP server within ShutdownTimeout.
func (a *Application) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.ShutdownTimeout)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	a.Logger.Println("shutdown complete")
	return nil
}
