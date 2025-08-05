package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// City represents a city in the master data system
type City struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	ProvinceID uuid.UUID `json:"province_id" gorm:"type:char(36);not null"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`
	Type       *string   `json:"type" gorm:"type:varchar(100)"`
	Code       *string   `json:"code" gorm:"type:varchar(50)"`
	Latitude   *string   `json:"latitude" gorm:"type:varchar(20)"`
	Longitude  *string   `json:"longitude" gorm:"type:varchar(20)"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (c *City) TableName() string {
	return "tm_cities"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (c *City) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// NewCity creates a new City instance
func NewCity(provinceID uuid.UUID, name string) *City {
	return &City{
		ProvinceID: provinceID,
		Name:       name,
	}
}

// SetType sets the type of the city
func (c *City) SetType(cityType string) {
	c.Type = &cityType
	c.UpdatedAt = time.Now()
}

// SetCode sets the city code
func (c *City) SetCode(code string) {
	c.Code = &code
	c.UpdatedAt = time.Now()
}

// SetCoordinates sets the geographical coordinates
func (c *City) SetCoordinates(latitude, longitude string) {
	c.Latitude = &latitude
	c.Longitude = &longitude
	c.UpdatedAt = time.Now()
}

// Update updates the city information
func (c *City) Update(name string) {
	c.Name = name
	c.UpdatedAt = time.Now()
}
