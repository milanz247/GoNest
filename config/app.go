package config

import "fmt"

// AppConfig holds application-level (non-database) settings.
type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	Host  string
	Port  int
}

// Addr returns the host:port pair suitable for http.Server.Addr.
func (a AppConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// IsProduction reports whether the app is running in a production-like env.
func (a AppConfig) IsProduction() bool {
	return a.Env == "production"
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Name:  getEnv("APP_NAME", "GoMVC"),
		Env:   getEnv("APP_ENV", "local"),
		Debug: getEnvBool("APP_DEBUG", true),
		Host:  getEnv("APP_HOST", "127.0.0.1"),
		Port:  getEnvInt("APP_PORT", 8080),
	}
}
