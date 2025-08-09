package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// LanguageHTTPHandler handles HTTP requests for language operations
type LanguageHTTPHandler struct {
	languageService *services.LanguageService
	searchService   repositories.SearchRepository
}

// NewLanguageHTTPHandler creates a new LanguageHTTPHandler instance
func NewLanguageHTTPHandler(languageService *services.LanguageService, searchService repositories.SearchRepository) *LanguageHTTPHandler {
	return &LanguageHTTPHandler{
		languageService: languageService,
		searchService:   searchService,
	}
}

// CreateLanguage handles POST /api/v1/languages
// @Summary Create a new language
// @Description Create a new language with the provided information
// @Tags languages
// @Accept json
// @Produce json
// @Param language body CreateLanguageRequest true "Language information"
// @Success 201 {object} response.Response "Language created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages [post]
func (h *LanguageHTTPHandler) CreateLanguage(c *fiber.Ctx) error {
	var req CreateLanguageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	language, err := h.languageService.CreateLanguage(c.Context(), req.Name, req.Code)
	if err != nil {
		return response.InternalServerError(c, "Failed to create language: "+err.Error())
	}

	return response.Created(c, language, "Language created successfully")
}

// GetLanguageByID handles GET /api/v1/languages/:id
// @Summary Get language by ID
// @Description Get a language by its UUID
// @Tags languages
// @Produce json
// @Param id path string true "Language ID (UUID)"
// @Success 200 {object} response.Response "Language retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Language not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/{id} [get]
func (h *LanguageHTTPHandler) GetLanguageByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid language ID: "+err.Error())
	}

	language, err := h.languageService.GetLanguageByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Language not found: "+err.Error())
	}

	return response.Success(c, language, "Language retrieved successfully")
}

// GetAllLanguages handles GET /api/v1/languages
// @Summary Get all languages
// @Description Get all languages with pagination
// @Tags languages
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Languages retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages [get]
func (h *LanguageHTTPHandler) GetAllLanguages(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	languages, err := h.languageService.GetAllLanguages(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve languages: "+err.Error())
	}

	return response.Success(c, languages, "Languages retrieved successfully")
}

// GetActiveLanguages handles GET /api/v1/languages/active
// @Summary Get active languages
// @Description Get all active languages with pagination
// @Tags languages
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Active languages retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/active [get]
func (h *LanguageHTTPHandler) GetActiveLanguages(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	languages, err := h.languageService.GetActiveLanguages(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve active languages: "+err.Error())
	}

	return response.Success(c, languages, "Active languages retrieved successfully")
}

// SearchLanguages handles GET /api/v1/languages/search
// @Summary Search languages
// @Description Search languages by name or code
// @Tags languages
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Languages found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/search [get]
func (h *LanguageHTTPHandler) SearchLanguages(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.BadRequest(c, "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	languages, err := h.languageService.SearchLanguages(c.Context(), query, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to search languages: "+err.Error())
	}

	return response.Success(c, languages, "Languages found")
}

// UpdateLanguage handles PUT /api/v1/languages/:id
// @Summary Update a language
// @Description Update an existing language
// @Tags languages
// @Accept json
// @Produce json
// @Param id path string true "Language ID (UUID)"
// @Param language body UpdateLanguageRequest true "Updated language information"
// @Success 200 {object} response.Response "Language updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Language not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/{id} [put]
func (h *LanguageHTTPHandler) UpdateLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid language ID: "+err.Error())
	}

	var req UpdateLanguageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	language, err := h.languageService.GetLanguageByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Language not found: "+err.Error())
	}

	// Update fields
	if req.Name != nil {
		language.SetName(*req.Name)
	}
	if req.Code != nil {
		language.SetCode(*req.Code)
	}

	if err := h.languageService.UpdateLanguage(c.Context(), language); err != nil {
		return response.InternalServerError(c, "Failed to update language: "+err.Error())
	}

	return response.Success(c, language, "Language updated successfully")
}

// ActivateLanguage handles POST /api/v1/languages/:id/activate
// @Summary Activate a language
// @Description Activate a deactivated language
// @Tags languages
// @Produce json
// @Param id path string true "Language ID (UUID)"
// @Success 200 {object} response.Response "Language activated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Language not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/{id}/activate [post]
func (h *LanguageHTTPHandler) ActivateLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid language ID: "+err.Error())
	}

	if err := h.languageService.ActivateLanguage(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to activate language: "+err.Error())
	}

	return response.Success(c, nil, "Language activated successfully")
}

// DeactivateLanguage handles POST /api/v1/languages/:id/deactivate
// @Summary Deactivate a language
// @Description Deactivate an active language
// @Tags languages
// @Produce json
// @Param id path string true "Language ID (UUID)"
// @Success 200 {object} response.Response "Language deactivated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Language not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/{id}/deactivate [post]
func (h *LanguageHTTPHandler) DeactivateLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid language ID: "+err.Error())
	}

	if err := h.languageService.DeactivateLanguage(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to deactivate language: "+err.Error())
	}

	return response.Success(c, nil, "Language deactivated successfully")
}

// DeleteLanguage handles DELETE /api/v1/languages/:id
// @Summary Delete a language
// @Description Delete a language by ID
// @Tags languages
// @Produce json
// @Param id path string true "Language ID (UUID)"
// @Success 200 {object} response.Response "Language deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Language not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/languages/{id} [delete]
func (h *LanguageHTTPHandler) DeleteLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid language ID: "+err.Error())
	}

	if err := h.languageService.DeleteLanguage(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete language: "+err.Error())
	}

	return response.Success(c, nil, "Language deleted successfully")
}

// Request/Response DTOs

type CreateLanguageRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type UpdateLanguageRequest struct {
	Name *string `json:"name,omitempty"`
	Code *string `json:"code,omitempty"`
}
