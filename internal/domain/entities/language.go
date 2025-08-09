package entities

import (
	"time"

	"github.com/google/uuid"
)

// Language represents a language entity
type Language struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the Language entity
func (l *Language) TableName() string {
	return "tm_languages"
}

// GenerateID generates a new UUID for the language if not set
func (l *Language) GenerateID() {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
}

// NewLanguage creates a new Language instance
func NewLanguage(name, code string) *Language {
	return &Language{
		ID:        uuid.New(),
		Name:      name,
		Code:      code,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetName sets the name of the language
func (l *Language) SetName(name string) {
	l.Name = name
	l.UpdatedAt = time.Now()
}

// SetCode sets the code of the language
func (l *Language) SetCode(code string) {
	l.Code = code
	l.UpdatedAt = time.Now()
}

// Activate activates the language
func (l *Language) Activate() {
	l.IsActive = true
	l.UpdatedAt = time.Now()
}

// Deactivate deactivates the language
func (l *Language) Deactivate() {
	l.IsActive = false
	l.UpdatedAt = time.Now()
}

// IsValid validates the language entity
func (l *Language) IsValid() bool {
	return l.Name != "" && l.Code != "" && len(l.Code) <= 10
}
