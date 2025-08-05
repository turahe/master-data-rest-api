package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// District represents a district in the master data system
type District struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CityID    uuid.UUID `json:"city_id" gorm:"type:char(36);not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Code      *string   `json:"code" gorm:"type:varchar(50)"`
	Latitude  *string   `json:"latitude" gorm:"type:varchar(20)"`
	Longitude *string   `json:"longitude" gorm:"type:varchar(20)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (d *District) TableName() string {
	return "tm_districts"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (d *District) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// NewDistrict creates a new District instance
func NewDistrict(cityID uuid.UUID, name string) *District {
	return &District{
		CityID: cityID,
		Name:   name,
	}
}

// SetCode sets the district code
func (d *District) SetCode(code string) {
	d.Code = &code
	d.UpdatedAt = time.Now()
}

// SetCoordinates sets the geographical coordinates
func (d *District) SetCoordinates(latitude, longitude string) {
	d.Latitude = &latitude
	d.Longitude = &longitude
	d.UpdatedAt = time.Now()
}

// Update updates the district information
func (d *District) Update(name string) {
	d.Name = name
	d.UpdatedAt = time.Now()
}
