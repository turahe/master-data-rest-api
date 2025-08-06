package configs

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver       string
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	Charset      string
	ParseTime    string
	Loc          string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "Master Data REST API"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Driver:       getEnv("DB_DRIVER", "postgres"),
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME", "master_data"),
			User:         getEnv("DB_USER", "appuser"),
			Password:     getEnv("DB_PASSWORD", "apppassword"),
			Charset:      getEnv("DB_CHARSET", "utf8"),
			ParseTime:    getEnv("DB_PARSE_TIME", "true"),
			Loc:          getEnv("DB_LOC", "Local"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Server: ServerConfig{
			Host: getEnv("APP_HOST", "localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
