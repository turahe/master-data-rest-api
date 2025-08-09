package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

// RequestLoggerMiddleware creates a custom request logger using logrus for Fiber
func RequestLoggerMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate request duration
		latency := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()

		// Create log entry with structured fields
		entry := log.WithFields(logrus.Fields{
			"status":     status,
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"latency":    latency,
			"user_agent": c.Get("User-Agent"),
			"bytes_in":   len(c.Body()),
			"bytes_out":  len(c.Response().Body()),
		})

		// Add API key info if available
		if apiKeyID := c.Locals("api_key_id"); apiKeyID != nil {
			entry = entry.WithField("api_key_id", apiKeyID)
		}
		if apiKeyName := c.Locals("api_key_name"); apiKeyName != nil {
			entry = entry.WithField("api_key_name", apiKeyName)
		}

		// Log at appropriate level based on status code
		switch {
		case status >= 500:
			entry.Error("Server Error")
		case status >= 400:
			entry.Warn("Client Error")
		case status >= 300:
			entry.Info("Redirection")
		default:
			entry.Info("Request")
		}

		return err
	}
}

// ErrorLoggerMiddleware logs errors that occur during request processing
func ErrorLoggerMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		// Log any errors that occurred
		if err != nil {
			entry := log.WithFields(logrus.Fields{
				"method":     c.Method(),
				"path":       c.Path(),
				"ip":         c.IP(),
				"user_agent": c.Get("User-Agent"),
			})

			// Check if it's a Fiber error for better categorization
			if fiberErr, ok := err.(*fiber.Error); ok {
				entry = entry.WithField("error_code", fiberErr.Code)
				switch {
				case fiberErr.Code >= 500:
					entry.WithError(err).Error("Server error")
				case fiberErr.Code >= 400:
					entry.WithError(err).Warn("Client error")
				default:
					entry.WithError(err).Info("Request error")
				}
			} else {
				entry.WithError(err).Error("Request processing error")
			}
		}

		return err
	}
}
