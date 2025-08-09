package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger to provide application-specific logging
type Logger struct {
	*logrus.Logger
}

// Config holds logger configuration
type Config struct {
	Level  string // debug, info, warn, error
	Format string // json, text
	Output string // stdout, stderr, file path
}

// New creates a new logger instance with configuration
func New(config Config) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	switch strings.ToLower(config.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// Set output
	var output io.Writer
	switch strings.ToLower(config.Output) {
	case "stderr":
		output = os.Stderr
	case "stdout", "":
		output = os.Stdout
	default:
		// If it's not stdout/stderr, treat as file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fallback to stdout if file can't be opened
			output = os.Stdout
			logger.WithError(err).Warn("Failed to open log file, falling back to stdout")
		} else {
			output = file
		}
	}
	logger.SetOutput(output)

	// Add caller information for better debugging
	logger.SetReportCaller(true)

	return &Logger{Logger: logger}
}

// NewDefault creates a logger with default configuration
func NewDefault() *Logger {
	return New(Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	})
}

// NewFromEnv creates a logger configured from environment variables
func NewFromEnv() *Logger {
	return New(Config{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "text"),
		Output: getEnv("LOG_OUTPUT", "stdout"),
	})
}

// WithField creates a new entry with a single field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields creates a new entry with multiple fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError creates a new entry with an error field
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithComponent creates a new entry with a component field
func (l *Logger) WithComponent(component string) *logrus.Entry {
	return l.Logger.WithField("component", component)
}

// WithRequest creates a new entry with request-related fields
func (l *Logger) WithRequest(method, url, userAgent string) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{
		"method":     method,
		"url":        url,
		"user_agent": userAgent,
	})
}

// WithDatabase creates a new entry with database-related fields
func (l *Logger) WithDatabase(operation, table string, duration interface{}) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{
		"db_operation": operation,
		"db_table":     table,
		"duration":     duration,
	})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
