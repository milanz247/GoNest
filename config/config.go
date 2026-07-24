// Package config loads and validates typed application configuration from
// the environment (and an optional .env file).
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the root typed configuration for the application.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

// Load reads a .env file (if present), then builds and validates the typed
// Config from the process environment. Real environment variables always
// take precedence over values found in the .env file.
func Load(envPath string) (*Config, error) {
	if err := loadDotEnv(envPath); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg := &Config{
		App:      loadAppConfig(),
		Database: loadDatabaseConfig(),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.App.Name == "" {
		missing = append(missing, "APP_NAME")
	}
	if c.App.Host == "" {
		missing = append(missing, "APP_HOST")
	}
	if c.App.Port == 0 {
		missing = append(missing, "APP_PORT")
	}
	if c.Database.Driver == "" {
		missing = append(missing, "DB_DRIVER")
	}
	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Database.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

// loadDotEnv parses a simple KEY=VALUE .env file and applies each entry via
// os.Setenv, without overriding variables already present in the
// environment. Missing files are not an error since real env vars may be
// supplied by the host/deployment platform instead.
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
