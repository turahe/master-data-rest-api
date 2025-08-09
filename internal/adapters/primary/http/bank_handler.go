package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// CreateBank handles POST /api/v1/banks
// @Summary Create a new bank
// @Description Create a new bank with the provided information
// @Tags banks
// @Accept json
// @Produce json
// @Param bank body CreateBankRequest true "Bank information"
// @Success 201 {object} response.Response "Bank created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks [post]
func (h *BankHTTPHandler) CreateBank(c *fiber.Ctx) error {
	var req CreateBankRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	bank, err := h.bankService.CreateBank(c.Context(), req.Name, req.Alias, req.Company, req.Code)
	if err != nil {
		return response.InternalServerError(c, "Failed to create bank: "+err.Error())
	}

	// Index the new bank in Meilisearch
	if err := h.searchService.IndexBank(c.Context(), bank); err != nil {
		// Log error but don't fail the operation
		// TODO: Consider using a background job for indexing
		// log.Printf("Failed to index bank in search: %v", err)
	}

	return response.Created(c, bank, "Bank created successfully")
}

// GetBankByID handles GET /api/v1/banks/:id
// @Summary Get bank by ID
// @Description Get a bank by its UUID
// @Tags banks
// @Produce json
// @Param id path string true "Bank ID (UUID)"
// @Success 200 {object} response.Response "Bank retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks/{id} [get]
func (h *BankHTTPHandler) GetBankByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid bank ID: "+err.Error())
	}

	bank, err := h.bankService.GetBankByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Bank not found: "+err.Error())
	}

	return response.Success(c, bank, "Bank retrieved successfully")
}

// GetAllBanks handles GET /api/v1/banks
// @Summary Get all banks
// @Description Get all banks with pagination
// @Tags banks
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Banks retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks [get]
func (h *BankHTTPHandler) GetAllBanks(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	banks, err := h.bankService.GetAllBanks(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve banks: "+err.Error())
	}

	return response.Success(c, banks, "Banks retrieved successfully")
}

// SearchBanks handles GET /api/v1/banks/search
// @Summary Search banks
// @Description Search banks by name, alias, company, or code
// @Tags banks
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response "Banks found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks/search [get]
func (h *BankHTTPHandler) SearchBanks(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.BadRequest(c, "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Use Meilisearch for fast search functionality
	banks, err := h.searchService.SearchBanks(c.Context(), query, limit, offset)
	if err != nil {
		// Fallback to database search if Meilisearch fails
		banks, err = h.bankService.SearchBanks(c.Context(), query, limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to search banks: "+err.Error())
		}
	}

	return response.Success(c, banks, "Banks found")
}

// GetBankByCode handles GET /api/v1/banks/code/:code
// @Summary Get bank by code
// @Description Get a bank by its code
// @Tags banks
// @Produce json
// @Param code path string true "Bank Code"
// @Success 200 {object} response.Response "Bank retrieved successfully"
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

// UpdateBank handles PUT /api/v1/banks/:id
// @Summary Update a bank
// @Description Update an existing bank
// @Tags banks
// @Accept json
// @Produce json
// @Param id path string true "Bank ID (UUID)"
// @Param bank body UpdateBankRequest true "Updated bank information"
// @Success 200 {object} response.Response "Bank updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks/{id} [put]
func (h *BankHTTPHandler) UpdateBank(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid bank ID: "+err.Error())
	}

	var req UpdateBankRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	bank, err := h.bankService.GetBankByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Bank not found: "+err.Error())
	}

	// Update fields
	if req.Name != nil {
		bank.SetName(*req.Name)
	}
	if req.Alias != nil {
		bank.SetAlias(*req.Alias)
	}
	if req.Company != nil {
		bank.SetCompany(*req.Company)
	}
	if req.Code != nil {
		bank.SetCode(*req.Code)
	}

	if err := h.bankService.UpdateBank(c.Context(), bank); err != nil {
		return response.InternalServerError(c, "Failed to update bank: "+err.Error())
	}

	return response.Success(c, bank, "Bank updated successfully")
}

// DeleteBank handles DELETE /api/v1/banks/:id
// @Summary Delete a bank
// @Description Delete a bank by ID
// @Tags banks
// @Produce json
// @Param id path string true "Bank ID (UUID)"
// @Success 200 {object} response.Response "Bank deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/banks/{id} [delete]
func (h *BankHTTPHandler) DeleteBank(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid bank ID: "+err.Error())
	}

	if err := h.bankService.DeleteBank(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete bank: "+err.Error())
	}

	return response.Success(c, nil, "Bank deleted successfully")
}

// Request/Response DTOs

type CreateBankRequest struct {
	Name    string `json:"name" validate:"required"`
	Alias   string `json:"alias" validate:"required"`
	Company string `json:"company" validate:"required"`
	Code    string `json:"code" validate:"required"`
}

type UpdateBankRequest struct {
	Name    *string `json:"name,omitempty"`
	Alias   *string `json:"alias,omitempty"`
	Company *string `json:"company,omitempty"`
	Code    *string `json:"code,omitempty"`
}
