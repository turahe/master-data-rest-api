package logger

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
)

// PgxLogger adapts our Logger to work with pgx's tracelog.Logger interface
type PgxLogger struct {
	logger *Logger
}

// NewPgxLogger creates a new PgxLogger instance
func NewPgxLogger(logger *Logger) *PgxLogger {
	return &PgxLogger{
		logger: logger,
	}
}

// Log implements tracelog.Logger interface
func (l *PgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	entry := l.logger.WithFields(logrus.Fields(data))

	// Add database context
	entry = entry.WithField("component", "database")

	switch level {
	case tracelog.LogLevelTrace:
		entry.Trace(msg)
	case tracelog.LogLevelDebug:
		entry.Debug(msg)
	case tracelog.LogLevelInfo:
		entry.Info(msg)
	case tracelog.LogLevelWarn:
		entry.Warn(msg)
	case tracelog.LogLevelError:
		entry.Error(msg)
	default:
		entry.Info(msg)
	}
}

// GetPgxLogLevel converts string log level to pgx tracelog.LogLevel
func GetPgxLogLevel(level string) tracelog.LogLevel {
	switch level {
	case "trace":
		return tracelog.LogLevelTrace
	case "debug":
		return tracelog.LogLevelDebug
	case "info":
		return tracelog.LogLevelInfo
	case "warn", "warning":
		return tracelog.LogLevelWarn
	case "error":
		return tracelog.LogLevelError
	default:
		return tracelog.LogLevelInfo
	}
}

// QueryLogger provides structured logging for database queries
type QueryLogger struct {
	logger *Logger
}

// NewQueryLogger creates a new QueryLogger instance
func NewQueryLogger(logger *Logger) *QueryLogger {
	return &QueryLogger{
		logger: logger,
	}
}

// LogQuery logs a database query with timing and metadata
func (ql *QueryLogger) LogQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	entry := ql.logger.WithFields(logrus.Fields{
		"component":    "database",
		"query":        query,
		"args_count":   len(args),
		"duration_ms":  duration.Milliseconds(),
		"duration_str": duration.String(),
	})

	// Add query arguments in debug mode
	if ql.logger.Level >= logrus.DebugLevel {
		entry = entry.WithField("args", args)
	}

	if err != nil {
		entry.WithError(err).Error("Database query failed")
	} else {
		if duration > 1*time.Second {
			entry.Warn("Slow database query detected")
		} else if duration > 100*time.Millisecond {
			entry.Info("Database query executed")
		} else {
			entry.Debug("Database query executed")
		}
	}
}

// LogConnection logs database connection events
func (ql *QueryLogger) LogConnection(ctx context.Context, event string, data map[string]interface{}) {
	entry := ql.logger.WithFields(logrus.Fields{
		"component": "database",
		"event":     event,
	})

	if data != nil {
		entry = entry.WithFields(logrus.Fields(data))
	}

	entry.Info("Database connection event")
}

// LogTransaction logs database transaction events
func (ql *QueryLogger) LogTransaction(ctx context.Context, event string, duration time.Duration) {
	entry := ql.logger.WithFields(logrus.Fields{
		"component":    "database",
		"event":        event,
		"duration_ms":  duration.Milliseconds(),
		"duration_str": duration.String(),
	})

	switch event {
	case "begin":
		entry.Debug("Transaction started")
	case "commit":
		entry.Debug("Transaction committed")
	case "rollback":
		entry.Info("Transaction rolled back")
	default:
		entry.Debug("Transaction event")
	}
}
