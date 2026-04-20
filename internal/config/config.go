package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Path:         getEnv("DATABASE_PATH", "./mocker.db"),
			MaxOpenConns: getIntEnv("DATABASE_MAX_OPEN_CONNS", 10),
			MaxIdleConns: getIntEnv("DATABASE_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getDurationEnv("DATABASE_MAX_LIFETIME", 5*time.Minute),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetDatabaseDSN returns the SQLite DSN
func (c *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.Path)
}

// GetPort returns the server port as a string with colon prefix
func (c *ServerConfig) GetPort() string {
	return ":" + c.Port
}
