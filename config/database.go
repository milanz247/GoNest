package config

// DatabaseConfig holds MySQL connection settings. The connection itself is
// established later in bootstrap once GORM is wired in.
type DatabaseConfig struct {
	Driver    string
	Host      string
	Port      int
	Name      string
	User      string
	Password  string
	Charset   string
	ParseTime bool
	Location  string
}

// String masks the password so the config can be safely logged.
func (d DatabaseConfig) String() string {
	return "DatabaseConfig{Driver: " + d.Driver + ", Host: " + d.Host + ", Name: " + d.Name + ", User: " + d.User + ", Password: ****}"
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:    getEnv("DB_DRIVER", "mysql"),
		Host:      getEnv("DB_HOST", "127.0.0.1"),
		Port:      getEnvInt("DB_PORT", 3306),
		Name:      getEnv("DB_NAME", "gomvc"),
		User:      getEnv("DB_USER", "root"),
		Password:  getEnv("DB_PASSWORD", ""),
		Charset:   getEnv("DB_CHARSET", "utf8mb4"),
		ParseTime: getEnvBool("DB_PARSE_TIME", true),
		Location:  getEnv("DB_LOCATION", "Local"),
	}
}
