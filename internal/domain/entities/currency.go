package entities

import (
	"time"

	"github.com/google/uuid"
)

// Currency represents a currency entity
type Currency struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Code          string    `json:"code" db:"code"`
	Symbol        *string   `json:"symbol,omitempty" db:"symbol"`
	DecimalPlaces int       `json:"decimal_places" db:"decimal_places"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the Currency entity
func (c *Currency) TableName() string {
	return "tm_currencies"
}

// GenerateID generates a new UUID for the currency if not set
func (c *Currency) GenerateID() {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
}

// NewCurrency creates a new Currency instance
func NewCurrency(name, code string, decimalPlaces int) *Currency {
	return &Currency{
		ID:            uuid.New(),
		Name:          name,
		Code:          code,
		DecimalPlaces: decimalPlaces,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// SetName sets the name of the currency
func (c *Currency) SetName(name string) {
	c.Name = name
	c.UpdatedAt = time.Now()
}

// SetCode sets the code of the currency
func (c *Currency) SetCode(code string) {
	c.Code = code
	c.UpdatedAt = time.Now()
}

// SetSymbol sets the symbol of the currency
func (c *Currency) SetSymbol(symbol string) {
	c.Symbol = &symbol
	c.UpdatedAt = time.Now()
}

// SetDecimalPlaces sets the decimal places of the currency
func (c *Currency) SetDecimalPlaces(places int) {
	c.DecimalPlaces = places
	c.UpdatedAt = time.Now()
}

// Activate activates the currency
func (c *Currency) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the currency
func (c *Currency) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// IsValid validates the currency entity
func (c *Currency) IsValid() bool {
	return c.Name != "" && c.Code != "" && len(c.Code) <= 3 && c.DecimalPlaces >= 0
}

// GetDisplayName returns the display name with symbol if available
func (c *Currency) GetDisplayName() string {
	if c.Symbol != nil && *c.Symbol != "" {
		return c.Name + " (" + *c.Symbol + ")"
	}
	return c.Name
}
