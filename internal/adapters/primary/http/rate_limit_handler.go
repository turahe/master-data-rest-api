package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/redis"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// RateLimitHTTPHandler handles HTTP requests for rate limit operations
type RateLimitHTTPHandler struct {
	rateLimiter *redis.RateLimiter
	log         *logrus.Logger
}

// NewRateLimitHTTPHandler creates a new RateLimitHTTPHandler instance
func NewRateLimitHTTPHandler(rateLimiter *redis.RateLimiter, log *logrus.Logger) *RateLimitHTTPHandler {
	return &RateLimitHTTPHandler{
		rateLimiter: rateLimiter,
		log:         log,
	}
}

// GetRateLimitInfo handles GET /api/v1/rate-limit/info
// @Summary Get rate limit information
// @Description Get current rate limit information for the requesting client
// @Tags rate-limit
// @Produce json
// @Param identifier query string false "Client identifier (IP, user ID, or API key)"
// @Success 200 {object} response.Response "Rate limit information retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/rate-limit/info [get]
func (h *RateLimitHTTPHandler) GetRateLimitInfo(c *fiber.Ctx) error {
	identifier := c.Query("identifier")
	if identifier == "" {
		identifier = c.IP()
	}

	// Get rate limit info for different types
	results := make(map[string]interface{})

	// API rate limit
	apiResult, err := h.rateLimiter.GetRateLimitInfo(c.Context(), identifier, redis.RateLimitConfig{
		Requests: 100,
		Window:   time.Minute,
		Key:      "api",
	})
	if err != nil {
		h.log.WithError(err).Error("Failed to get API rate limit info")
	} else {
		results["api"] = map[string]interface{}{
			"allowed":    apiResult.Allowed,
			"remaining":  apiResult.Remaining,
			"reset_time": apiResult.ResetTime,
			"limit":      100,
			"window":     "1m",
		}
	}

	// IP-based rate limit
	ipResult, err := h.rateLimiter.GetRateLimitInfo(c.Context(), identifier, redis.RateLimitConfig{
		Requests: 1000,
		Window:   time.Minute,
		Key:      "ip",
	})
	if err != nil {
		h.log.WithError(err).Error("Failed to get IP rate limit info")
	} else {
		results["ip"] = map[string]interface{}{
			"allowed":    ipResult.Allowed,
			"remaining":  ipResult.Remaining,
			"reset_time": ipResult.ResetTime,
			"limit":      1000,
			"window":     "1m",
		}
	}

	return response.Success(c, results, "Rate limit information retrieved successfully")
}

// GetRateLimitStats handles GET /api/v1/rate-limit/stats
// @Summary Get rate limit statistics
// @Description Get comprehensive rate limit statistics across all clients
// @Tags rate-limit
// @Produce json
// @Param key query string false "Rate limit key to get stats for (api, ip, user, apikey)"
// @Success 200 {object} response.Response "Rate limit statistics retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/rate-limit/stats [get]
func (h *RateLimitHTTPHandler) GetRateLimitStats(c *fiber.Ctx) error {
	key := c.Query("key", "api")

	// Get stats for the specified key
	stats, err := h.rateLimiter.GetRateLimitStats(c.Context(), redis.RateLimitConfig{
		Requests: 100,
		Window:   time.Minute,
		Key:      key,
	})

	if err != nil {
		h.log.WithError(err).Error("Failed to get rate limit stats")
		return response.InternalServerError(c, "Failed to get rate limit statistics")
	}

	return response.Success(c, stats, "Rate limit statistics retrieved successfully")
}

// ResetRateLimit handles POST /api/v1/rate-limit/reset
// @Summary Reset rate limit for a client
// @Description Reset rate limit for a specific client identifier
// @Tags rate-limit
// @Accept json
// @Produce json
// @Param identifier body string true "Client identifier to reset"
// @Param key body string false "Rate limit key to reset (default: api)"
// @Success 200 {object} response.Response "Rate limit reset successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/rate-limit/reset [post]
func (h *RateLimitHTTPHandler) ResetRateLimit(c *fiber.Ctx) error {
	var req struct {
		Identifier string `json:"identifier" validate:"required"`
		Key        string `json:"key"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if req.Key == "" {
		req.Key = "api"
	}

	// Reset the rate limit
	err := h.rateLimiter.ResetRateLimit(c.Context(), req.Identifier, redis.RateLimitConfig{
		Requests: 100,
		Window:   time.Minute,
		Key:      req.Key,
	})

	if err != nil {
		h.log.WithError(err).Error("Failed to reset rate limit")
		return response.InternalServerError(c, "Failed to reset rate limit")
	}

	h.log.WithFields(logrus.Fields{
		"identifier": req.Identifier,
		"key":        req.Key,
	}).Info("Rate limit reset")

	return response.Success(c, map[string]interface{}{
		"identifier": req.Identifier,
		"key":        req.Key,
		"reset_time": time.Now(),
	}, "Rate limit reset successfully")
}

// GetRateLimitConfig handles GET /api/v1/rate-limit/config
// @Summary Get rate limit configuration
// @Description Get current rate limit configuration for different endpoints
// @Tags rate-limit
// @Produce json
// @Success 200 {object} response.Response "Rate limit configuration retrieved successfully"
// @Router /api/v1/rate-limit/config [get]
func (h *RateLimitHTTPHandler) GetRateLimitConfig(c *fiber.Ctx) error {
	config := map[string]interface{}{
		"global": map[string]interface{}{
			"requests": 100,
			"window":   "1m",
			"key":      "api",
		},
		"tiered": map[string]interface{}{
			"geographic": map[string]interface{}{
				"endpoints": []string{"/api/v1/countries", "/api/v1/provinces"},
				"requests":  1000,
				"window":    "1m",
				"key":       "geo",
			},
			"financial": map[string]interface{}{
				"endpoints": []string{"/api/v1/banks", "/api/v1/currencies"},
				"requests":  500,
				"window":    "1m",
				"key":       "financial",
			},
			"languages": map[string]interface{}{
				"endpoints": []string{"/api/v1/languages"},
				"requests":  200,
				"window":    "1m",
				"key":       "languages",
			},
		},
		"identifiers": map[string]interface{}{
			"ip": map[string]interface{}{
				"requests": 1000,
				"window":   "1m",
				"key":      "ip",
			},
			"user": map[string]interface{}{
				"requests": 500,
				"window":   "1m",
				"key":      "user",
			},
			"apikey": map[string]interface{}{
				"requests": 2000,
				"window":   "1m",
				"key":      "apikey",
			},
		},
	}

	return response.Success(c, config, "Rate limit configuration retrieved successfully")
}
