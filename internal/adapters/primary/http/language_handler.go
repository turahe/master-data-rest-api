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
