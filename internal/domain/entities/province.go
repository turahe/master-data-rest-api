package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Province represents a province/state in the master data system
type Province struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CountryID uuid.UUID `json:"country_id" gorm:"type:char(36);not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Region    *string   `json:"region" gorm:"type:varchar(255)"`
	ISO31662  *string   `json:"iso_3166_2" gorm:"type:varchar(10)"`
	Code      *string   `json:"code" gorm:"type:varchar(50)"`
	Latitude  *string   `json:"latitude" gorm:"type:varchar(20)"`
	Longitude *string   `json:"longitude" gorm:"type:varchar(20)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (p *Province) TableName() string {
	return "tm_provinces"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (p *Province) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// NewProvince creates a new Province instance
func NewProvince(countryID uuid.UUID, name string) *Province {
	return &Province{
		CountryID: countryID,
		Name:      name,
	}
}

// SetRegion sets the region of the province
func (p *Province) SetRegion(region string) {
	p.Region = &region
	p.UpdatedAt = time.Now()
}

// SetISO31662 sets the ISO 3166-2 code
func (p *Province) SetISO31662(code string) {
	p.ISO31662 = &code
	p.UpdatedAt = time.Now()
}

// SetCode sets the province code
func (p *Province) SetCode(code string) {
	p.Code = &code
	p.UpdatedAt = time.Now()
}

// SetCoordinates sets the geographical coordinates
func (p *Province) SetCoordinates(latitude, longitude string) {
	p.Latitude = &latitude
	p.Longitude = &longitude
	p.UpdatedAt = time.Now()
}

// Update updates the province information
func (p *Province) Update(name string) {
	p.Name = name
	p.UpdatedAt = time.Now()
}
