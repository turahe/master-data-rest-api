package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Village represents a village in the master data system
type Village struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	DistrictID uuid.UUID `json:"district_id" gorm:"type:char(36);not null"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`
	Code       *string   `json:"code" gorm:"type:varchar(50)"`
	PostalCode *string   `json:"postal_code" gorm:"type:varchar(20)"`
	Latitude   *string   `json:"latitude" gorm:"type:varchar(20)"`
	Longitude  *string   `json:"longitude" gorm:"type:varchar(20)"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (v *Village) TableName() string {
	return "tm_villages"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (v *Village) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// NewVillage creates a new Village instance
func NewVillage(districtID uuid.UUID, name string) *Village {
	return &Village{
		DistrictID: districtID,
		Name:       name,
	}
}

// SetCode sets the village code
func (v *Village) SetCode(code string) {
	v.Code = &code
	v.UpdatedAt = time.Now()
}

// SetPostalCode sets the postal code
func (v *Village) SetPostalCode(postalCode string) {
	v.PostalCode = &postalCode
	v.UpdatedAt = time.Now()
}

// SetCoordinates sets the geographical coordinates
func (v *Village) SetCoordinates(latitude, longitude string) {
	v.Latitude = &latitude
	v.Longitude = &longitude
	v.UpdatedAt = time.Now()
}

// Update updates the village information
func (v *Village) Update(name string) {
	v.Name = name
	v.UpdatedAt = time.Now()
}
