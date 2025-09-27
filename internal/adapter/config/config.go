package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port     string
	LogLevel string
}

// NewConfig creates a new configuration instance with defaults
func NewConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// GetPort returns the port with colon prefix for fiber
func (c *Config) GetPort() string {
	return ":" + c.Port
}

// GetPortInt returns the port as integer
func (c *Config) GetPortInt() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return port
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}