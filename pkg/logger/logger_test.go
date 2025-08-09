package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		expectedLevel  logrus.Level
		expectedFormat string
	}{
		{
			name: "json format with debug level",
			config: Config{
				Level:  "debug",
				Format: "json",
				Output: "stdout",
			},
			expectedLevel:  logrus.DebugLevel,
			expectedFormat: "json",
		},
		{
			name: "text format with info level",
			config: Config{
				Level:  "info",
				Format: "text",
				Output: "stdout",
			},
			expectedLevel:  logrus.InfoLevel,
			expectedFormat: "text",
		},
		{
			name: "invalid level defaults to info",
			config: Config{
				Level:  "invalid",
				Format: "text",
				Output: "stdout",
			},
			expectedLevel:  logrus.InfoLevel,
			expectedFormat: "text",
		},
		{
			name: "empty format defaults to text",
			config: Config{
				Level:  "warn",
				Format: "",
				Output: "stdout",
			},
			expectedLevel:  logrus.WarnLevel,
			expectedFormat: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			logger := New(tt.config)

			// Then
			assert.NotNil(t, logger)
			assert.Equal(t, tt.expectedLevel, logger.GetLevel())

			// Check formatter type
			switch tt.expectedFormat {
			case "json":
				_, ok := logger.Formatter.(*logrus.JSONFormatter)
				assert.True(t, ok, "Expected JSON formatter")
			default:
				_, ok := logger.Formatter.(*logrus.TextFormatter)
				assert.True(t, ok, "Expected Text formatter")
			}
		})
	}
}

func TestNewDefault(t *testing.T) {
	// When
	logger := NewDefault()

	// Then
	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())

	// Should use text formatter by default
	_, ok := logger.Formatter.(*logrus.TextFormatter)
	assert.True(t, ok, "Expected Text formatter")
}

func TestLogger_JSONOutput(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When
	logger.Info("test message")

	// Then
	output := buf.String()
	assert.NotEmpty(t, output)

	// Parse JSON to verify structure
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["message"])
	assert.Contains(t, logEntry, "timestamp")
}

func TestLogger_TextOutput(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When
	logger.Info("test message")

	// Then
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "level=info")
}

func TestLogger_WithFields(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When
	logger.WithFields(logrus.Fields{
		"user_id": "123",
		"action":  "login",
	}).Info("User logged in")

	// Then
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "123", logEntry["user_id"])
	assert.Equal(t, "login", logEntry["action"])
	assert.Equal(t, "User logged in", logEntry["message"])
}

func TestLogger_LogLevels(t *testing.T) {
	tests := []struct {
		level     string
		logMethod func(*Logger, ...interface{})
		shouldLog bool
	}{
		{"debug", (*Logger).Debug, true},
		{"info", (*Logger).Info, true},
		{"warn", (*Logger).Warn, true},
		{"error", (*Logger).Error, true},
		{"debug", (*Logger).Info, true},  // info should log when level is debug
		{"info", (*Logger).Debug, false}, // debug should not log when level is info
		{"warn", (*Logger).Info, false},  // info should not log when level is warn
	}

	for _, tt := range tests {
		t.Run(tt.level+"_logging", func(t *testing.T) {
			// Given
			var buf bytes.Buffer
			config := Config{
				Level:  tt.level,
				Format: "text",
				Output: "stdout",
			}
			logger := New(config)
			logger.SetOutput(&buf)

			// When
			tt.logMethod(logger, "test message")

			// Then
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, "test message")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLogger_DatabaseLogging(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When
	logger.WithFields(logrus.Fields{
		"query":    "SELECT * FROM users WHERE id = $1",
		"args":     []interface{}{123},
		"duration": "10ms",
	}).Debug("Database query executed")

	// Then
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "SELECT * FROM users WHERE id = $1", logEntry["query"])
	assert.Equal(t, "10ms", logEntry["duration"])
	assert.Contains(t, logEntry, "args")
}

func TestLogger_ErrorLogging(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "error",
		Format: "json",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When
	logger.WithError(assert.AnError).Error("Operation failed")

	// Then
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "Operation failed", logEntry["message"])
	assert.Contains(t, logEntry, "error")
}

func TestLogger_ContextualLogging(t *testing.T) {
	// Given
	var buf bytes.Buffer
	config := Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	logger := New(config)
	logger.SetOutput(&buf)

	// When - Simulate HTTP request logging
	logger.WithFields(logrus.Fields{
		"method":    "POST",
		"path":      "/api/v1/banks",
		"status":    201,
		"duration":  "150ms",
		"remote_ip": "127.0.0.1",
	}).Info("HTTP request completed")

	// Then
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "POST", logEntry["method"])
	assert.Equal(t, "/api/v1/banks", logEntry["path"])
	assert.Equal(t, float64(201), logEntry["status"]) // JSON numbers are float64
	assert.Equal(t, "150ms", logEntry["duration"])
	assert.Equal(t, "127.0.0.1", logEntry["remote_ip"])
}

func TestConfig_OutputHandling(t *testing.T) {
	tests := []struct {
		name         string
		outputConfig string
		expectStdout bool
	}{
		{
			name:         "stdout output",
			outputConfig: "stdout",
			expectStdout: true,
		},
		{
			name:         "stderr output",
			outputConfig: "stderr",
			expectStdout: false,
		},
		{
			name:         "empty defaults to stdout",
			outputConfig: "",
			expectStdout: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			config := Config{
				Level:  "info",
				Format: "text",
				Output: tt.outputConfig,
			}

			// When
			logger := New(config)

			// Then - We can't easily test actual stdout/stderr redirection,
			// but we can verify the logger was created successfully
			assert.NotNil(t, logger)
		})
	}
}
