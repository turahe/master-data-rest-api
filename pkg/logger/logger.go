package logger

import (
	"log"
	"os"
)

// Logger provides a simple logging interface
type Logger struct {
	*log.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.Printf("[DEBUG] "+format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.Printf("[WARN] "+format, v...)
} 