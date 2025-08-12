package middleware

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// MockAPIKeyService is a mock implementation of APIKeyValidator
type MockAPIKeyService struct {
	mock.Mock
}

func (m *MockAPIKeyService) ValidateAPIKey(ctx context.Context, key string) (*entities.APIKey, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.APIKey), args.Error(1)
}

func TestOptionalAPIKeyAuth_NoAuthHeader(t *testing.T) {
	// Setup
	app := fiber.New()
	mockService := new(MockAPIKeyService)

	app.Use(OptionalAPIKeyAuth(mockService))

	app.Get("/test", func(c *fiber.Ctx) error {
		apiKey := c.Locals("api_key")
		apiKeyID := c.Locals("api_key_id")
		apiKeyName := c.Locals("api_key_name")

		assert.Nil(t, apiKey)
		assert.Equal(t, "", apiKeyID)
		assert.Equal(t, "anonymous", apiKeyName)

		return c.SendString("success")
	})

	// Test
	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)

	// Assertions
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertNotCalled(t, "ValidateAPIKey")
}

func TestOptionalAPIKeyAuth_ValidAPIKey(t *testing.T) {
	// Setup
	app := fiber.New()
	mockService := new(MockAPIKeyService)

	// Mock API key
	testID := uuid.New()
	mockAPIKey := &entities.APIKey{
		ID:   testID,
		Name: "Test Key",
	}

	mockService.On("ValidateAPIKey", mock.Anything, "valid-key").Return(mockAPIKey, nil)

	app.Use(OptionalAPIKeyAuth(mockService))

	app.Get("/test", func(c *fiber.Ctx) error {
		apiKey := c.Locals("api_key")
		apiKeyID := c.Locals("api_key_id")
		apiKeyName := c.Locals("api_key_name")

		assert.Equal(t, mockAPIKey, apiKey)
		assert.Equal(t, testID.String(), apiKeyID)
		assert.Equal(t, "Test Key", apiKeyName)

		return c.SendString("success")
	})

	// Test
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-key")
	resp, _ := app.Test(req)

	// Assertions
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestOptionalAPIKeyAuth_InvalidAPIKey(t *testing.T) {
	// Setup
	app := fiber.New()
	mockService := new(MockAPIKeyService)

	// Mock API key validation failure
	mockService.On("ValidateAPIKey", mock.Anything, "invalid-key").Return(nil, assert.AnError)

	app.Use(OptionalAPIKeyAuth(mockService))

	app.Get("/test", func(c *fiber.Ctx) error {
		apiKey := c.Locals("api_key")
		apiKeyID := c.Locals("api_key_id")
		apiKeyName := c.Locals("api_key_name")

		assert.Nil(t, apiKey)
		assert.Equal(t, "", apiKeyID)
		assert.Equal(t, "anonymous", apiKeyName)

		return c.SendString("success")
	})

	// Test
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-key")
	resp, _ := app.Test(req)

	// Assertions
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestOptionalAPIKeyAuth_XAPIKeyHeader(t *testing.T) {
	// Setup
	app := fiber.New()
	mockService := new(MockAPIKeyService)

	// Mock API key
	testID := uuid.New()
	mockAPIKey := &entities.APIKey{
		ID:   testID,
		Name: "Test Key",
	}

	mockService.On("ValidateAPIKey", mock.Anything, "valid-key").Return(mockAPIKey, nil)

	app.Use(OptionalAPIKeyAuth(mockService))

	app.Get("/test", func(c *fiber.Ctx) error {
		apiKey := c.Locals("api_key")
		apiKeyID := c.Locals("api_key_id")
		apiKeyName := c.Locals("api_key_name")

		assert.Equal(t, mockAPIKey, apiKey)
		assert.Equal(t, testID.String(), apiKeyID)
		assert.Equal(t, "Test Key", apiKeyName)

		return c.SendString("success")
	})

	// Test
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-key")
	resp, _ := app.Test(req)

	// Assertions
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}
