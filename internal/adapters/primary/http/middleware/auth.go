package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// APIKeyValidator interface for API key validation
type APIKeyValidator interface {
	ValidateAPIKey(ctx context.Context, key string) (*entities.APIKey, error)
}

// APIKeyAuth creates a middleware for API key authentication using database validation
func APIKeyAuth(apiKeyService APIKeyValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for API key in header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "API key is required")
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(c, "Invalid authorization format. Use 'Bearer <api-key>'")
		}

		// Extract the API key
		providedKey := strings.TrimPrefix(authHeader, "Bearer ")
		if providedKey == "" {
			return response.Unauthorized(c, "API key cannot be empty")
		}

		// Validate the API key against database
		apiKey, err := apiKeyService.ValidateAPIKey(c.Context(), providedKey)
		if err != nil {
			return response.Unauthorized(c, "Failed to validate API key")
		}

		if apiKey == nil {
			return response.Unauthorized(c, "Invalid API key")
		}

		// Store API key info in context for use in handlers
		c.Locals("api_key", apiKey)
		c.Locals("api_key_id", apiKey.ID.String())
		c.Locals("api_key_name", apiKey.Name)

		// API key is valid, continue to the next handler
		return c.Next()
	}
}

// APIKeyAuthWithHeader creates a middleware for API key authentication using custom header and database validation
func APIKeyAuthWithHeader(apiKeyService APIKeyValidator, headerName string) fiber.Handler {
	if headerName == "" {
		headerName = "X-API-Key"
	}

	return func(c *fiber.Ctx) error {
		// Check for API key in custom header
		providedKey := c.Get(headerName)
		if providedKey == "" {
			return response.Unauthorized(c, "API key is required in "+headerName+" header")
		}

		// Validate the API key against database
		apiKey, err := apiKeyService.ValidateAPIKey(c.Context(), providedKey)
		if err != nil {
			return response.Unauthorized(c, "Failed to validate API key")
		}

		if apiKey == nil {
			return response.Unauthorized(c, "Invalid API key")
		}

		// Store API key info in context for use in handlers
		c.Locals("api_key", apiKey)
		c.Locals("api_key_id", apiKey.ID.String())
		c.Locals("api_key_name", apiKey.Name)

		// API key is valid, continue to the next handler
		return c.Next()
	}
}

// OptionalAPIKeyAuth creates a middleware for optional API key authentication
// It validates API keys when provided but doesn't require them
func OptionalAPIKeyAuth(apiKeyService APIKeyValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for API key in header
		authHeader := c.Get("Authorization")

		// If no authorization header, continue without authentication
		if authHeader == "" {
			// Check for X-API-Key header as alternative
			apiKeyHeader := c.Get("X-API-Key")
			if apiKeyHeader == "" {
				// No authentication provided, continue as anonymous user
				c.Locals("api_key", nil)
				c.Locals("api_key_id", "")
				c.Locals("api_key_name", "anonymous")
				return c.Next()
			}
			// Use X-API-Key header
			authHeader = "Bearer " + apiKeyHeader
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue as anonymous user
			c.Locals("api_key", nil)
			c.Locals("api_key_id", "")
			c.Locals("api_key_name", "anonymous")
			return c.Next()
		}

		// Extract the API key
		providedKey := strings.TrimPrefix(authHeader, "Bearer ")
		if providedKey == "" {
			// Empty API key, continue as anonymous user
			c.Locals("api_key", nil)
			c.Locals("api_key_id", "")
			c.Locals("api_key_name", "anonymous")
			return c.Next()
		}

		// Validate the API key against database
		apiKey, err := apiKeyService.ValidateAPIKey(c.Context(), providedKey)
		if err != nil {
			// Validation failed, continue as anonymous user
			c.Locals("api_key", nil)
			c.Locals("api_key_id", "")
			c.Locals("api_key_name", "anonymous")
			return c.Next()
		}

		if apiKey == nil {
			// Invalid API key, continue as anonymous user
			c.Locals("api_key", nil)
			c.Locals("api_key_id", "")
			c.Locals("api_key_name", "anonymous")
			return c.Next()
		}

		// Store API key info in context for use in handlers
		c.Locals("api_key", apiKey)
		c.Locals("api_key_id", apiKey.ID.String())
		c.Locals("api_key_name", apiKey.Name)

		// API key is valid, continue to the next handler
		return c.Next()
	}
}
