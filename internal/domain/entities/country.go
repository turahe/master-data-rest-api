package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Country represents a country in the master data system
type Country struct {
	ID              uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Capital         *string   `json:"capital" gorm:"type:varchar(255)"`
	Citizenship     *string   `json:"citizenship" gorm:"type:varchar(255)"`
	CountryCode     string    `json:"country_code" gorm:"type:varchar(3);unique;not null"`
	CurrencyName    *string   `json:"currency_name" gorm:"type:varchar(255)"`
	CurrencyCode    *string   `json:"currency_code" gorm:"type:varchar(3)"`
	CurrencySubUnit *string   `json:"currency_sub_unit" gorm:"type:varchar(255)"`
	CurrencySymbol  *string   `json:"currency_symbol" gorm:"type:varchar(10)"`
	FullName        *string   `json:"full_name" gorm:"type:varchar(255)"`
	ISO31662        string    `json:"iso_3166_2" gorm:"type:varchar(2);unique;not null"`
	ISO31663        string    `json:"iso_3166_3" gorm:"type:varchar(3);unique;not null"`
	Name            string    `json:"name" gorm:"type:varchar(255);not null"`
	RegionCode      *string   `json:"region_code" gorm:"type:varchar(10)"`
	SubRegionCode   *string   `json:"sub_region_code" gorm:"type:varchar(10)"`
	EEA             bool      `json:"eea" gorm:"default:false"`
	CallingCode     string    `json:"calling_code" gorm:"type:varchar(10);not null"`
	Flag            *string   `json:"flag" gorm:"type:varchar(10)"`
	Latitude        *string   `json:"latitude" gorm:"type:varchar(20)"`
	Longitude       *string   `json:"longitude" gorm:"type:varchar(20)"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (c *Country) TableName() string {
	return "tm_countries"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (c *Country) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// NewCountry creates a new Country instance
func NewCountry(
	countryCode string,
	iso31662 string,
	iso31663 string,
	name string,
	eea bool,
	callingCode string,
) *Country {
	return &Country{
		CountryCode: countryCode,
		ISO31662:    iso31662,
		ISO31663:    iso31663,
		Name:        name,
		EEA:         eea,
		CallingCode: callingCode,
	}
}

// SetCapital sets the capital of the country
func (c *Country) SetCapital(capital string) {
	c.Capital = &capital
	c.UpdatedAt = time.Now()
}

// SetCitizenship sets the citizenship of the country
func (c *Country) SetCitizenship(citizenship string) {
	c.Citizenship = &citizenship
	c.UpdatedAt = time.Now()
}

// SetCurrencyInfo sets the currency information
func (c *Country) SetCurrencyInfo(name, code, subUnit, symbol string) {
	c.CurrencyName = &name
	c.CurrencyCode = &code
	c.CurrencySubUnit = &subUnit
	c.CurrencySymbol = &symbol
	c.UpdatedAt = time.Now()
}

// SetFullName sets the full name of the country
func (c *Country) SetFullName(fullName string) {
	c.FullName = &fullName
	c.UpdatedAt = time.Now()
}

// SetRegionInfo sets the region information
func (c *Country) SetRegionInfo(regionCode, subRegionCode string) {
	c.RegionCode = &regionCode
	c.SubRegionCode = &subRegionCode
	c.UpdatedAt = time.Now()
}

// SetFlag sets the flag code
func (c *Country) SetFlag(flag string) {
	c.Flag = &flag
	c.UpdatedAt = time.Now()
}

// SetCoordinates sets the geographical coordinates
func (c *Country) SetCoordinates(latitude, longitude string) {
	c.Latitude = &latitude
	c.Longitude = &longitude
	c.UpdatedAt = time.Now()
}

// IsEEA returns whether the country is part of the European Economic Area
func (c *Country) IsEEA() bool {
	return c.EEA
}

// GetDisplayName returns the display name (full name if available, otherwise name)
func (c *Country) GetDisplayName() string {
	if c.FullName != nil && *c.FullName != "" {
		return *c.FullName
	}
	return c.Name
}
