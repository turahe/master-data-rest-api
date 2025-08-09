package http

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// APIKeyHTTPHandler handles HTTP requests for API key operations
type APIKeyHTTPHandler struct {
	apiKeyService *services.APIKeyService
}

// NewAPIKeyHTTPHandler creates a new APIKeyHTTPHandler instance
func NewAPIKeyHTTPHandler(apiKeyService *services.APIKeyService) *APIKeyHTTPHandler {
	return &APIKeyHTTPHandler{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKey handles POST /api/v1/api-keys
// @Summary Create a new API key
// @Description Create a new API key with the provided information
// @Tags api-keys
// @Accept json
// @Produce json
// @Param request body CreateAPIKeyRequest true "API key information"
// @Success 201 {object} response.Response "API key created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys [post]
func (h *APIKeyHTTPHandler) CreateAPIKey(c *fiber.Ctx) error {
	var req CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return response.BadRequest(c, "Invalid expires_at format. Use ISO 8601 format (e.g., 2023-12-31T23:59:59Z)")
		}
		expiresAt = &parsed
	}

	apiKey, err := h.apiKeyService.CreateAPIKey(context.Background(), req.Name, req.Description, expiresAt)
	if err != nil {
		return response.InternalServerError(c, "Failed to create API key: "+err.Error())
	}

	return response.Created(c, apiKey, "API key created successfully")
}

// GetAllAPIKeys handles GET /api/v1/api-keys
// @Summary Get all API keys
// @Description Get all API keys with pagination
// @Tags api-keys
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "API keys retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys [get]
func (h *APIKeyHTTPHandler) GetAllAPIKeys(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	apiKeys, err := h.apiKeyService.GetAllAPIKeys(context.Background(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve API keys: "+err.Error())
	}

	return response.Success(c, apiKeys, "API keys retrieved successfully")
}

// GetAPIKeyByID handles GET /api/v1/api-keys/:id
// @Summary Get API key by ID
// @Description Get an API key by its UUID
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID (UUID)"
// @Success 200 {object} response.Response "API key retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "API key not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys/{id} [get]
func (h *APIKeyHTTPHandler) GetAPIKeyByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid API key ID: "+err.Error())
	}

	apiKey, err := h.apiKeyService.GetAPIKeyByID(context.Background(), id)
	if err != nil {
		return response.NotFound(c, "API key not found: "+err.Error())
	}

	return response.Success(c, apiKey, "API key retrieved successfully")
}

// UpdateAPIKey handles PUT /api/v1/api-keys/:id
// @Summary Update an API key
// @Description Update an existing API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID (UUID)"
// @Param request body UpdateAPIKeyRequest true "Updated API key information"
// @Success 200 {object} response.Response "API key updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "API key not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys/{id} [put]
func (h *APIKeyHTTPHandler) UpdateAPIKey(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid API key ID: "+err.Error())
	}

	var req UpdateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	apiKey, err := h.apiKeyService.GetAPIKeyByID(context.Background(), id)
	if err != nil {
		return response.NotFound(c, "API key not found: "+err.Error())
	}

	// Update fields
	if req.Name != nil {
		apiKey.Name = *req.Name
	}
	if req.Description != nil {
		apiKey.SetDescription(*req.Description)
	}
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			apiKey.ExpiresAt = nil
		} else {
			parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return response.BadRequest(c, "Invalid expires_at format. Use ISO 8601 format")
			}
			apiKey.SetExpiration(parsed)
		}
	}

	if err := h.apiKeyService.UpdateAPIKey(context.Background(), apiKey); err != nil {
		return response.InternalServerError(c, "Failed to update API key: "+err.Error())
	}

	return response.Success(c, apiKey, "API key updated successfully")
}

// ActivateAPIKey handles POST /api/v1/api-keys/:id/activate
// @Summary Activate an API key
// @Description Activate a deactivated API key
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID (UUID)"
// @Success 200 {object} response.Response "API key activated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "API key not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys/{id}/activate [post]
func (h *APIKeyHTTPHandler) ActivateAPIKey(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid API key ID: "+err.Error())
	}

	if err := h.apiKeyService.ActivateAPIKey(context.Background(), id); err != nil {
		return response.InternalServerError(c, "Failed to activate API key: "+err.Error())
	}

	return response.Success(c, nil, "API key activated successfully")
}

// DeactivateAPIKey handles POST /api/v1/api-keys/:id/deactivate
// @Summary Deactivate an API key
// @Description Deactivate an active API key
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID (UUID)"
// @Success 200 {object} response.Response "API key deactivated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "API key not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys/{id}/deactivate [post]
func (h *APIKeyHTTPHandler) DeactivateAPIKey(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid API key ID: "+err.Error())
	}

	if err := h.apiKeyService.DeactivateAPIKey(context.Background(), id); err != nil {
		return response.InternalServerError(c, "Failed to deactivate API key: "+err.Error())
	}

	return response.Success(c, nil, "API key deactivated successfully")
}

// DeleteAPIKey handles DELETE /api/v1/api-keys/:id
// @Summary Delete an API key
// @Description Soft delete an API key
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID (UUID)"
// @Success 200 {object} response.Response "API key deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "API key not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/api-keys/{id} [delete]
func (h *APIKeyHTTPHandler) DeleteAPIKey(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid API key ID: "+err.Error())
	}

	if err := h.apiKeyService.DeleteAPIKey(context.Background(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete API key: "+err.Error())
	}

	return response.Success(c, nil, "API key deleted successfully")
}

// Request/Response DTOs

type CreateAPIKeyRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"` // ISO 8601 format
}

type UpdateAPIKeyRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty"` // ISO 8601 format
}
