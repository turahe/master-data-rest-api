package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Currency represents a currency in the master data system
type Currency struct {
	ID                   uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Priority             *int      `json:"priority" gorm:"type:int"`
	ISOCode              string    `json:"iso_code" gorm:"type:varchar(3);unique;not null"`
	Name                 *string   `json:"name" gorm:"type:varchar(255)"`
	Symbol               *string   `json:"symbol" gorm:"type:varchar(10)"`
	DisambiguateSymbol   *string   `json:"disambiguate_symbol" gorm:"type:varchar(10)"`
	AlternateSymbols     *string   `json:"alternate_symbols" gorm:"type:text"`
	Subunit              *string   `json:"subunit" gorm:"type:varchar(255)"`
	SubunitToUnit        int       `json:"subunit_to_unit" gorm:"not null"`
	SymbolFirst          bool      `json:"symbol_first" gorm:"default:false"`
	HTMLEntity           *string   `json:"html_entity" gorm:"type:varchar(50)"`
	DecimalMark          *string   `json:"decimal_mark" gorm:"type:varchar(10)"`
	ThousandsSeparator   *string   `json:"thousands_separator" gorm:"type:varchar(10)"`
	ISONumeric           *string   `json:"iso_numeric" gorm:"type:varchar(10)"`
	SmallestDenomination int       `json:"smallest_denomination" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (c *Currency) TableName() string {
	return "tm_currencies"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (c *Currency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// NewCurrency creates a new Currency instance
func NewCurrency(isoCode string, subunitToUnit int, symbolFirst bool, smallestDenomination int) *Currency {
	return &Currency{
		ISOCode:              isoCode,
		SubunitToUnit:        subunitToUnit,
		SymbolFirst:          symbolFirst,
		SmallestDenomination: smallestDenomination,
	}
}

// SetPriority sets the priority of the currency
func (c *Currency) SetPriority(priority int) {
	c.Priority = &priority
	c.UpdatedAt = time.Now()
}

// SetName sets the name of the currency
func (c *Currency) SetName(name string) {
	c.Name = &name
	c.UpdatedAt = time.Now()
}

// SetSymbol sets the symbol of the currency
func (c *Currency) SetSymbol(symbol string) {
	c.Symbol = &symbol
	c.UpdatedAt = time.Now()
}

// SetDisambiguateSymbol sets the disambiguate symbol
func (c *Currency) SetDisambiguateSymbol(symbol string) {
	c.DisambiguateSymbol = &symbol
	c.UpdatedAt = time.Now()
}

// SetAlternateSymbols sets the alternate symbols
func (c *Currency) SetAlternateSymbols(symbols string) {
	c.AlternateSymbols = &symbols
	c.UpdatedAt = time.Now()
}

// SetSubunit sets the subunit name
func (c *Currency) SetSubunit(subunit string) {
	c.Subunit = &subunit
	c.UpdatedAt = time.Now()
}

// SetHTMLEntity sets the HTML entity
func (c *Currency) SetHTMLEntity(entity string) {
	c.HTMLEntity = &entity
	c.UpdatedAt = time.Now()
}

// SetDecimalMark sets the decimal mark
func (c *Currency) SetDecimalMark(mark string) {
	c.DecimalMark = &mark
	c.UpdatedAt = time.Now()
}

// SetThousandsSeparator sets the thousands separator
func (c *Currency) SetThousandsSeparator(separator string) {
	c.ThousandsSeparator = &separator
	c.UpdatedAt = time.Now()
}

// SetISONumeric sets the ISO numeric code
func (c *Currency) SetISONumeric(code string) {
	c.ISONumeric = &code
	c.UpdatedAt = time.Now()
}

// GetDisplaySymbol returns the display symbol (disambiguate symbol if available, otherwise symbol)
func (c *Currency) GetDisplaySymbol() string {
	if c.DisambiguateSymbol != nil && *c.DisambiguateSymbol != "" {
		return *c.DisambiguateSymbol
	}
	if c.Symbol != nil && *c.Symbol != "" {
		return *c.Symbol
	}
	return c.ISOCode
}
