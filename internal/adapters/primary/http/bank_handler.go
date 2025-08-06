package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// BankHTTPHandler handles HTTP requests for bank operations
type BankHTTPHandler struct {
	bankService *services.BankService
}

// NewBankHTTPHandler creates a new BankHTTPHandler instance
func NewBankHTTPHandler(bankService *services.BankService) *BankHTTPHandler {
	return &BankHTTPHandler{
		bankService: bankService,
	}
}

// CreateBank handles POST /api/v1/banks
// @Summary Create a new bank
// @Description Create a new bank with the provided information
// @Tags banks
// @Accept json
// @Produce json
// @Param request body object true "Bank information"
// @Success 201 {object} response.Response{data=entities.Bank} "Bank created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks [post]
func (h *BankHTTPHandler) CreateBank(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Alias   string `json:"alias" binding:"required"`
		Company string `json:"company" binding:"required"`
		Code    string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	bank, err := h.bankService.CreateBank(req.Name, req.Alias, req.Company, req.Code)
	if err != nil {
		response.InternalServerError(c, "Failed to create bank: "+err.Error())
		return
	}

	response.Created(c, bank, "Bank created successfully")
}

// GetBankByID handles GET /api/v1/banks/:id
// @Summary Get bank by ID
// @Description Get a bank by its UUID
// @Tags banks
// @Accept json
// @Produce json
// @Param id path string true "Bank ID" format(uuid)
// @Success 200 {object} response.Response{data=entities.Bank} "Bank found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/{id} [get]
func (h *BankHTTPHandler) GetBankByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid bank ID: "+err.Error())
		return
	}

	bank, err := h.bankService.GetBankByID(id)
	if err != nil {
		response.NotFound(c, "Bank not found: "+err.Error())
		return
	}

	response.Success(c, bank, "Bank retrieved successfully")
}

// GetBankByCode handles GET /api/v1/banks/code/:code
// @Summary Get bank by code
// @Description Get a bank by its code
// @Tags banks
// @Accept json
// @Produce json
// @Param code path string true "Bank code"
// @Success 200 {object} response.Response{data=entities.Bank} "Bank found"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/code/{code} [get]
func (h *BankHTTPHandler) GetBankByCode(c *gin.Context) {
	code := c.Param("code")
	bank, err := h.bankService.GetBankByCode(code)
	if err != nil {
		response.NotFound(c, "Bank not found: "+err.Error())
		return
	}

	response.Success(c, bank, "Bank retrieved successfully")
}

// GetBankByName handles GET /api/v1/banks/name/:name
// @Summary Get banks by name
// @Description Get banks by name (can return multiple)
// @Tags banks
// @Accept json
// @Produce json
// @Param name path string true "Bank name"
// @Success 200 {object} response.Response{data=[]entities.Bank} "Banks found"
// @Failure 404 {object} response.Response "Banks not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/name/{name} [get]
func (h *BankHTTPHandler) GetBankByName(c *gin.Context) {
	name := c.Param("name")
	banks, err := h.bankService.GetBankByName(name)
	if err != nil {
		response.NotFound(c, "Banks not found: "+err.Error())
		return
	}

	response.Success(c, banks, "Banks retrieved successfully")
}

// GetBankByAlias handles GET /api/v1/banks/alias/:alias
// @Summary Get banks by alias
// @Description Get banks by alias (can return multiple)
// @Tags banks
// @Accept json
// @Produce json
// @Param alias path string true "Bank alias"
// @Success 200 {object} response.Response{data=[]entities.Bank} "Banks found"
// @Failure 404 {object} response.Response "Banks not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/alias/{alias} [get]
func (h *BankHTTPHandler) GetBankByAlias(c *gin.Context) {
	alias := c.Param("alias")
	banks, err := h.bankService.GetBankByAlias(alias)
	if err != nil {
		response.NotFound(c, "Banks not found: "+err.Error())
		return
	}

	response.Success(c, banks, "Banks retrieved successfully")
}

// GetBanksByCompany handles GET /api/v1/banks/company/:company
// @Summary Get banks by company
// @Description Get all banks for a specific company
// @Tags banks
// @Accept json
// @Produce json
// @Param company path string true "Company name"
// @Success 200 {object} response.Response{data=[]entities.Bank} "Banks found"
// @Failure 404 {object} response.Response "Banks not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/company/{company} [get]
func (h *BankHTTPHandler) GetBanksByCompany(c *gin.Context) {
	company := c.Param("company")
	banks, err := h.bankService.GetBanksByCompany(company)
	if err != nil {
		response.NotFound(c, "Banks not found: "+err.Error())
		return
	}

	response.Success(c, banks, "Banks retrieved successfully")
}

// GetAllBanks handles GET /api/v1/banks
// @Summary Get all banks
// @Description Get all banks with optional pagination
// @Tags banks
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]entities.Bank} "Banks retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks [get]
func (h *BankHTTPHandler) GetAllBanks(c *gin.Context) {
	banks, err := h.bankService.GetAllBanks()
	if err != nil {
		response.InternalServerError(c, "Failed to get banks: "+err.Error())
		return
	}

	response.Success(c, banks, "Banks retrieved successfully")
}

// UpdateBank handles PUT /api/v1/banks/:id
// @Summary Update a bank
// @Description Update a bank with the provided information
// @Tags banks
// @Accept json
// @Produce json
// @Param id path string true "Bank ID" format(uuid)
// @Param request body object true "Bank information"
// @Success 200 {object} response.Response{data=entities.Bank} "Bank updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/{id} [put]
func (h *BankHTTPHandler) UpdateBank(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid bank ID: "+err.Error())
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Alias   string `json:"alias" binding:"required"`
		Company string `json:"company" binding:"required"`
		Code    string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Get existing bank
	bank, err := h.bankService.GetBankByID(id)
	if err != nil {
		response.NotFound(c, "Bank not found: "+err.Error())
		return
	}

	// Update fields
	bank.Update(req.Name, req.Alias, req.Company)
	bank.Code = req.Code

	if err := h.bankService.UpdateBank(bank); err != nil {
		response.InternalServerError(c, "Failed to update bank: "+err.Error())
		return
	}

	response.Success(c, bank, "Bank updated successfully")
}

// DeleteBank handles DELETE /api/v1/banks/:id
// @Summary Delete a bank
// @Description Delete a bank by its UUID
// @Tags banks
// @Accept json
// @Produce json
// @Param id path string true "Bank ID" format(uuid)
// @Success 200 {object} response.Response "Bank deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Bank not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/{id} [delete]
func (h *BankHTTPHandler) DeleteBank(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid bank ID: "+err.Error())
		return
	}

	if err := h.bankService.DeleteBank(id); err != nil {
		response.NotFound(c, "Bank not found: "+err.Error())
		return
	}

	response.Success(c, nil, "Bank deleted successfully")
}

// GetBankCount handles GET /api/v1/banks/count
// @Summary Get bank count
// @Description Get the total number of banks
// @Tags banks
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object} "Bank count retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/banks/count [get]
func (h *BankHTTPHandler) GetBankCount(c *gin.Context) {
	count, err := h.bankService.GetBankCount()
	if err != nil {
		response.InternalServerError(c, "Failed to get bank count: "+err.Error())
		return
	}

	response.Success(c, gin.H{"count": count}, "Bank count retrieved successfully")
} 