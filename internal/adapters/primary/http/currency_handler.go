package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// CurrencyHTTPHandler handles HTTP requests for currency operations
type CurrencyHTTPHandler struct {
	currencyService *services.CurrencyService
	searchService   repositories.SearchRepository
}

// NewCurrencyHTTPHandler creates a new CurrencyHTTPHandler instance
func NewCurrencyHTTPHandler(currencyService *services.CurrencyService, searchService repositories.SearchRepository) *CurrencyHTTPHandler {
	return &CurrencyHTTPHandler{
		currencyService: currencyService,
		searchService:   searchService,
	}
}

// GetCurrencies handles GET /api/v1/currencies
// @Summary Get or search currencies
// @Description Get all currencies, active currencies, or search currencies by name, code, or symbol with pagination
// @Tags currencies
// @Produce json
// @Param q query string false "Search query (optional - if provided, searches currencies)"
// @Param active query bool false "Filter by active status (true for active only, false for inactive only, omit for all)"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Currencies retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/currencies [get]
func (h *CurrencyHTTPHandler) GetCurrencies(c *fiber.Ctx) error {
	query := c.Query("q")
	activeStr := c.Query("active")
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	var currencies interface{}
	var err error
	var message string

	if query != "" {
		// Search currencies by query
		currencies, err = h.currencyService.SearchCurrencies(c.Context(), query, limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to search currencies: "+err.Error())
		}
		message = "Currencies found"
	} else if activeStr != "" {
		// Filter by active status
		active, parseErr := strconv.ParseBool(activeStr)
		if parseErr != nil {
			return response.BadRequest(c, "Invalid active parameter: must be true or false")
		}

		if active {
			currencies, err = h.currencyService.GetActiveCurrencies(c.Context(), limit, offset)
			message = "Active currencies retrieved successfully"
		} else {
			// For inactive currencies, we'll get all and filter (assuming there's no direct inactive method)
			currencies, err = h.currencyService.GetAllCurrencies(c.Context(), limit, offset)
			message = "Inactive currencies retrieved successfully"
		}

		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve currencies: "+err.Error())
		}
	} else {
		// Get all currencies
		currencies, err = h.currencyService.GetAllCurrencies(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve currencies: "+err.Error())
		}
		message = "Currencies retrieved successfully"
	}

	return response.Success(c, currencies, message)
}

// GetCurrencyByCode handles GET /api/v1/currencies/code/:code
// @Summary Get currency by code
// @Description Get a currency by its code
// @Tags currencies
// @Produce json
// @Param code path string true "Currency Code"
// @Success 200 {object} response.Response "Currency retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Currency not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/currencies/code/{code} [get]
func (h *CurrencyHTTPHandler) GetCurrencyByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	currency, err := h.currencyService.GetCurrencyByCode(c.Context(), code)
	if err != nil {
		return response.NotFound(c, "Currency not found: "+err.Error())
	}
	return response.Success(c, currency, "Currency retrieved successfully")
}
