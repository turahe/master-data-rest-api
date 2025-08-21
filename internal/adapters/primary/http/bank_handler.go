package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// BankHTTPHandler handles HTTP requests for bank operations
type BankHTTPHandler struct {
	bankService   *services.BankService
	searchService repositories.SearchRepository
}

// NewBankHTTPHandler creates a new BankHTTPHandler instance
func NewBankHTTPHandler(bankService *services.BankService, searchService repositories.SearchRepository) *BankHTTPHandler {
	return &BankHTTPHandler{
		bankService:   bankService,
		searchService: searchService,
	}
}

// GetBanks handles GET /api/v1/banks
// @Summary Get or search banks
// @Description Get all banks or search banks by name, alias, company, or code with pagination
// @Tags banks
// @Produce json
// @Param q query string false "Search query (optional - if provided, searches banks; if not provided, gets all banks)"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Banks retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks [get]
func (h *BankHTTPHandler) GetBanks(c *fiber.Ctx) error {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	var banks interface{}
	var err error
	var message string

	if query != "" {
		// Search banks by query
		// Use Meilisearch for fast search functionality
		banks, err = h.searchService.SearchBanks(c.Context(), query, limit, offset)
		if err != nil {
			// Fallback to database search if Meilisearch fails
			banks, err = h.bankService.SearchBanks(c.Context(), query, limit, offset)
			if err != nil {
				return response.InternalServerError(c, "Failed to search banks: "+err.Error())
			}
		}
		message = "Banks found"
	} else {
		// Get all banks
		banks, err = h.bankService.GetAllBanks(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve banks: "+err.Error())
		}
		message = "Banks retrieved successfully"
	}

	return response.Success(c, banks, message)
}

// GetBankByCode handles GET /api/v1/banks/code/:code
// @Summary Get bank by code
// @Description Get a bank by its code
// @Tags banks
// @Produce json
// @Param code path string true "Bank Code"
// @Success 200 {object} response.Response "Bank retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks/code/{code} [get]
func (h *BankHTTPHandler) GetBankByCode(c *fiber.Ctx) error {
	code := c.Params("code")

	bank, err := h.bankService.GetBankByCode(c.Context(), code)
	if err != nil {
		return response.NotFound(c, "Bank not found: "+err.Error())
	}

	return response.Success(c, bank, "Bank retrieved successfully")
}

// Request/Response DTOs
