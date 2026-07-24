// Package bootstrap wires together configuration, logging, routing and (in
// later phases) the database into a ready-to-run framework.Application.
package bootstrap

import (
	"fmt"

	"myapp/config"
	"myapp/framework"
)

// Boot loads configuration and constructs the Application. It is the single
// entry point main.go calls to assemble the app.
func Boot() (*framework.Application, error) {
	cfg, err := config.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("bootstrap: %w", err)
	}

	logger := NewLogger(cfg)
	router := NewRouter()

	app := framework.New(cfg, logger, router)

	return app, nil
}
