package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/redis"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// RateLimiterConfig holds configuration for the rate limiter middleware
type RateLimiterConfig struct {
	// Rate limit configuration
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window for the limit
	Key      string        // Redis key prefix for this limit

	// Identifier configuration
	IdentifierFunc func(c *fiber.Ctx) string // Function to extract identifier from request

	// Response configuration
	LimitExceededHandler func(c *fiber.Ctx, result *redis.RateLimitResult) error // Custom handler for limit exceeded

	// Skip configuration
	SkipFunc func(c *fiber.Ctx) bool // Function to determine if rate limiting should be skipped

	// Headers configuration
	IncludeHeaders bool // Whether to include rate limit headers in response
}

// DefaultRateLimiterConfig returns a default rate limiter configuration
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Requests:       100,         // 100 requests
		Window:         time.Minute, // per minute
		Key:            "api",       // API rate limit
		IncludeHeaders: true,
		IdentifierFunc: func(c *fiber.Ctx) string {
			// Default: use IP address as identifier
			return c.IP()
		},
		LimitExceededHandler: func(c *fiber.Ctx, result *redis.RateLimitResult) error {
			// Default: return 429 Too Many Requests
			return response.TooManyRequests(c, map[string]interface{}{
				"retry_after": int(result.RetryAfter.Seconds()),
				"reset_time":  result.ResetTime,
			}, fmt.Sprintf("Rate limit exceeded. Try again in %d seconds", int(result.RetryAfter.Seconds())))
		},
		SkipFunc: func(c *fiber.Ctx) bool {
			// Skip rate limiting for health checks and static files
			return c.Path() == "/health" || c.Path() == "/metrics"
		},
	}
}

// RateLimiter creates a new rate limiter middleware
func RateLimiter(rateLimiter *redis.RateLimiter, config *RateLimiterConfig) fiber.Handler {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return func(c *fiber.Ctx) error {
		// Skip rate limiting if configured
		if config.SkipFunc != nil && config.SkipFunc(c) {
			return c.Next()
		}

		// Get identifier for this request
		identifier := config.IdentifierFunc(c)
		if identifier == "" {
			// If no identifier, allow the request
			return c.Next()
		}

		// Check rate limit
		result, err := rateLimiter.CheckRateLimit(c.Context(), identifier, redis.RateLimitConfig{
			Requests: config.Requests,
			Window:   config.Window,
			Key:      config.Key,
		})

		if err != nil {
			// Log error but allow request to proceed
			logrus.WithError(err).WithField("identifier", identifier).Error("Rate limit check failed")
			return c.Next()
		}

		// Add rate limit headers if configured
		if config.IncludeHeaders {
			c.Set("X-RateLimit-Limit", strconv.Itoa(config.Requests))
			c.Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
			c.Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))

			if !result.Allowed {
				c.Set("X-RateLimit-RetryAfter", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
			}
		}

		// Check if request is allowed
		if !result.Allowed {
			// Call custom handler or use default
			if config.LimitExceededHandler != nil {
				return config.LimitExceededHandler(c, result)
			}

			// Default response
			return response.TooManyRequests(c, map[string]interface{}{
				"retry_after": int(result.RetryAfter.Seconds()),
				"reset_time":  result.ResetTime,
			}, "Rate limit exceeded")
		}

		// Request allowed, continue
		return c.Next()
	}
}

// IPBasedRateLimiter creates a rate limiter that uses IP address as identifier
func IPBasedRateLimiter(rateLimiter *redis.RateLimiter, requests int, window time.Duration) fiber.Handler {
	config := DefaultRateLimiterConfig()
	config.Requests = requests
	config.Window = window
	config.Key = "ip"

	return RateLimiter(rateLimiter, config)
}

// UserBasedRateLimiter creates a rate limiter that uses user ID as identifier
func UserBasedRateLimiter(rateLimiter *redis.RateLimiter, requests int, window time.Duration) fiber.Handler {
	config := DefaultRateLimiterConfig()
	config.Requests = requests
	config.Window = window
	config.Key = "user"
	config.IdentifierFunc = func(c *fiber.Ctx) string {
		// Extract user ID from context or header
		// This is a placeholder - implement based on your auth system
		userID := c.Get("X-User-ID")
		if userID == "" {
			// Fallback to IP if no user ID
			return c.IP()
		}
		return userID
	}

	return RateLimiter(rateLimiter, config)
}

// APIKeyBasedRateLimiter creates a rate limiter that uses API key as identifier
func APIKeyBasedRateLimiter(rateLimiter *redis.RateLimiter, requests int, window time.Duration) fiber.Handler {
	config := DefaultRateLimiterConfig()
	config.Requests = requests
	config.Window = window
	config.Key = "apikey"
	config.IdentifierFunc = func(c *fiber.Ctx) string {
		// Extract API key from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return ""
		}

		// Remove "Bearer " prefix if present
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}

		return authHeader
	}

	return RateLimiter(rateLimiter, config)
}

// TieredRateLimiter creates different rate limits based on request path or other criteria
func TieredRateLimiter(rateLimiter *redis.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requests int
		var window time.Duration
		var key string

		// Determine rate limit based on path
		switch {
		case c.Path() == "/api/v1/countries" || c.Path() == "/api/v1/provinces":
			// Geographic data: higher limits
			requests = 1000
			window = time.Minute
			key = "geo"
		case c.Path() == "/api/v1/banks" || c.Path() == "/api/v1/currencies":
			// Financial data: medium limits
			requests = 500
			window = time.Minute
			key = "financial"
		case c.Path() == "/api/v1/languages":
			// Language data: standard limits
			requests = 200
			window = time.Minute
			key = "languages"
		default:
			// Default rate limit
			requests = 100
			window = time.Minute
			key = "default"
		}

		// Create temporary config for this request
		config := &RateLimiterConfig{
			Requests:       requests,
			Window:         window,
			Key:            key,
			IncludeHeaders: true,
			IdentifierFunc: func(c *fiber.Ctx) string {
				return c.IP()
			},
		}

		// Apply rate limiting
		return RateLimiter(rateLimiter, config)(c)
	}
}
