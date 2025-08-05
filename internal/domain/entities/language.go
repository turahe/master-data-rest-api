package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Language represents a language in the master data system
type Language struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Code      string    `json:"code" gorm:"type:varchar(10);unique;not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Native    string    `json:"native" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (l *Language) TableName() string {
	return "tm_languages"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (l *Language) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// NewLanguage creates a new Language instance
func NewLanguage(code, name, native string) *Language {
	return &Language{
		Code:   code,
		Name:   name,
		Native: native,
	}
}

// Update updates the language information
func (l *Language) Update(name, native string) {
	l.Name = name
	l.Native = native
	l.UpdatedAt = time.Now()
}

// GetDisplayName returns the display name (native name if available, otherwise name)
func (l *Language) GetDisplayName() string {
	if l.Native != "" {
		return l.Native
	}
	return l.Name
}
