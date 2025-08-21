package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
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

// Request/Response DTOs
