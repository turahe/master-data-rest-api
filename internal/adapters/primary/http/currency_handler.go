package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// CreateCurrency handles POST /api/v1/currencies
func (h *CurrencyHTTPHandler) CreateCurrency(c *fiber.Ctx) error {
	var req CreateCurrencyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	currency, err := h.currencyService.CreateCurrency(c.Context(), req.Name, req.Code, req.DecimalPlaces)
	if err != nil {
		return response.InternalServerError(c, "Failed to create currency: "+err.Error())
	}

	// Set optional symbol
	if req.Symbol != nil {
		currency.SetSymbol(*req.Symbol)
		if err := h.currencyService.UpdateCurrency(c.Context(), currency); err != nil {
			return response.InternalServerError(c, "Failed to update currency: "+err.Error())
		}
	}

	return response.Created(c, currency, "Currency created successfully")
}

// GetCurrencyByID handles GET /api/v1/currencies/:id
func (h *CurrencyHTTPHandler) GetCurrencyByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid currency ID: "+err.Error())
	}

	currency, err := h.currencyService.GetCurrencyByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Currency not found: "+err.Error())
	}

	return response.Success(c, currency, "Currency retrieved successfully")
}

// GetAllCurrencies handles GET /api/v1/currencies
func (h *CurrencyHTTPHandler) GetAllCurrencies(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	currencies, err := h.currencyService.GetAllCurrencies(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve currencies: "+err.Error())
	}

	return response.Success(c, currencies, "Currencies retrieved successfully")
}

// GetActiveCurrencies handles GET /api/v1/currencies/active
func (h *CurrencyHTTPHandler) GetActiveCurrencies(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	currencies, err := h.currencyService.GetActiveCurrencies(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve active currencies: "+err.Error())
	}

	return response.Success(c, currencies, "Active currencies retrieved successfully")
}

// SearchCurrencies handles GET /api/v1/currencies/search
func (h *CurrencyHTTPHandler) SearchCurrencies(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.BadRequest(c, "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	currencies, err := h.currencyService.SearchCurrencies(c.Context(), query, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to search currencies: "+err.Error())
	}

	return response.Success(c, currencies, "Currencies found")
}

// UpdateCurrency handles PUT /api/v1/currencies/:id
func (h *CurrencyHTTPHandler) UpdateCurrency(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid currency ID: "+err.Error())
	}

	var req UpdateCurrencyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	currency, err := h.currencyService.GetCurrencyByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Currency not found: "+err.Error())
	}

	// Update fields
	if req.Name != nil {
		currency.SetName(*req.Name)
	}
	if req.Code != nil {
		currency.SetCode(*req.Code)
	}
	if req.Symbol != nil {
		currency.SetSymbol(*req.Symbol)
	}
	if req.DecimalPlaces != nil {
		currency.SetDecimalPlaces(*req.DecimalPlaces)
	}

	if err := h.currencyService.UpdateCurrency(c.Context(), currency); err != nil {
		return response.InternalServerError(c, "Failed to update currency: "+err.Error())
	}

	return response.Success(c, currency, "Currency updated successfully")
}

// ActivateCurrency handles POST /api/v1/currencies/:id/activate
func (h *CurrencyHTTPHandler) ActivateCurrency(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid currency ID: "+err.Error())
	}

	if err := h.currencyService.ActivateCurrency(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to activate currency: "+err.Error())
	}

	return response.Success(c, nil, "Currency activated successfully")
}

// DeactivateCurrency handles POST /api/v1/currencies/:id/deactivate
func (h *CurrencyHTTPHandler) DeactivateCurrency(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid currency ID: "+err.Error())
	}

	if err := h.currencyService.DeactivateCurrency(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to deactivate currency: "+err.Error())
	}

	return response.Success(c, nil, "Currency deactivated successfully")
}

// DeleteCurrency handles DELETE /api/v1/currencies/:id
func (h *CurrencyHTTPHandler) DeleteCurrency(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid currency ID: "+err.Error())
	}

	if err := h.currencyService.DeleteCurrency(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete currency: "+err.Error())
	}

	return response.Success(c, nil, "Currency deleted successfully")
}

// Request/Response DTOs

type CreateCurrencyRequest struct {
	Name          string  `json:"name" validate:"required"`
	Code          string  `json:"code" validate:"required"`
	Symbol        *string `json:"symbol,omitempty"`
	DecimalPlaces int     `json:"decimal_places" validate:"min=0"`
}

type UpdateCurrencyRequest struct {
	Name          *string `json:"name,omitempty"`
	Code          *string `json:"code,omitempty"`
	Symbol        *string `json:"symbol,omitempty"`
	DecimalPlaces *int    `json:"decimal_places,omitempty"`
}
