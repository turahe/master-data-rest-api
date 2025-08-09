package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// APIKeyAuth creates a middleware for API key authentication using database validation
func APIKeyAuth(apiKeyService *services.APIKeyService) fiber.Handler {
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
func APIKeyAuthWithHeader(apiKeyService *services.APIKeyService, headerName string) fiber.Handler {
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
