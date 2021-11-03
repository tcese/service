package service

// Configurations exported
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	RepMode  string // Select mock over db
	LogFile  string // Name of the log file
}

// ServerConfigurations exported
type ServerConfig struct {
	Port       int
	Production bool // Is it in Debug or in Production
}

// DatabaseConfigurations exported
type DatabaseConfig struct {
	Server   string
	Schema   string
	User     string
	Password string
}
