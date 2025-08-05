package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bank represents a bank in the master data system
type Bank struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Alias     string    `json:"alias" gorm:"type:varchar(255);not null"`
	Company   string    `json:"company" gorm:"type:varchar(255);not null"`
	Code      string    `json:"code" gorm:"type:varchar(50);unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (b *Bank) TableName() string {
	return "tm_banks"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (b *Bank) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// NewBank creates a new Bank instance
func NewBank(name, alias, company, code string) *Bank {
	return &Bank{
		Name:    name,
		Alias:   alias,
		Company: company,
		Code:    code,
	}
}

// Update updates the bank information
func (b *Bank) Update(name, alias, company string) {
	b.Name = name
	b.Alias = alias
	b.Company = company
	b.UpdatedAt = time.Now()
}

// GetDisplayName returns the display name (alias if available, otherwise name)
func (b *Bank) GetDisplayName() string {
	if b.Alias != "" {
		return b.Alias
	}
	return b.Name
}
