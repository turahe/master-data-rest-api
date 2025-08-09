package configs

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
	Logging  LoggingConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name    string
	Env     string
	Version string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver       string
	Host         string
	Port         string
	Name         string
	User         string
	Username     string // Alias for User for backward compatibility
	Password     string
	Charset      string
	ParseTime    string
	Loc          string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	// Logging configuration
	LogLevel     string
	LogQueries   bool
	LogSlowQuery time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// Future auth configuration can be added here
	// APIKey moved to database management
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "Master Data REST API"),
			Env:     getEnv("APP_ENV", "development"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Driver:       getEnv("DB_DRIVER", "postgres"),
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME", "master_data"),
			User:         getEnv("DB_USER", "appuser"),
			Username:     getEnv("DB_USER", "appuser"), // Set both for compatibility
			Password:     getEnv("DB_PASSWORD", "apppassword"),
			Charset:      getEnv("DB_CHARSET", "utf8"),
			ParseTime:    getEnv("DB_PARSE_TIME", "true"),
			Loc:          getEnv("DB_LOC", "Local"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
			// Logging configuration
			LogLevel:     getEnv("DB_LOG_LEVEL", "info"),
			LogQueries:   getEnvAsBool("DB_LOG_QUERIES", true),
			LogSlowQuery: getEnvAsDuration("DB_LOG_SLOW_QUERY", 100*time.Millisecond),
		},
		Server: ServerConfig{
			Host: getEnv("APP_HOST", "localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Auth: AuthConfig{
			// API keys are now managed in database
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
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

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}
