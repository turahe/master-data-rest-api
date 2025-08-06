package http

import (
	"github.com/gin-gonic/gin"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
)

// CurrencyHTTPHandler handles HTTP requests for currency operations
type CurrencyHTTPHandler struct {
	currencyService *services.CurrencyService
}

// NewCurrencyHTTPHandler creates a new CurrencyHTTPHandler instance
func NewCurrencyHTTPHandler(currencyService *services.CurrencyService) *CurrencyHTTPHandler {
	return &CurrencyHTTPHandler{
		currencyService: currencyService,
	}
}

// Placeholder methods - to be implemented
func (h *CurrencyHTTPHandler) CreateCurrency(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) GetAllCurrencies(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) GetCurrencyByID(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) GetCurrencyByCode(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) GetCurrencyByName(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) UpdateCurrency(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) DeleteCurrency(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
}

func (h *CurrencyHTTPHandler) GetCurrencyCount(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet"})
} 