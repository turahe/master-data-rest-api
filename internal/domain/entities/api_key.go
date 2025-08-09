package entities

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key entity
type APIKey struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Key         string     `json:"key" db:"key"`
	Description *string    `json:"description,omitempty" db:"description"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

// TableName returns the table name for the APIKey entity
func (APIKey) TableName() string {
	return "tm_api_keys"
}

// NewAPIKey creates a new APIKey instance
func NewAPIKey(name, key string) *APIKey {
	return &APIKey{
		ID:       uuid.New(),
		Name:     name,
		Key:      key,
		IsActive: true,
	}
}

// SetDescription sets the description for the API key
func (a *APIKey) SetDescription(description string) {
	a.Description = &description
}

// SetExpiration sets the expiration date for the API key
func (a *APIKey) SetExpiration(expiresAt time.Time) {
	a.ExpiresAt = &expiresAt
}

// IsExpired checks if the API key has expired
func (a *APIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsValid checks if the API key is valid (active and not expired)
func (a *APIKey) IsValid() bool {
	return a.IsActive && !a.IsExpired()
}

// UpdateLastUsed updates the last used timestamp
func (a *APIKey) UpdateLastUsed() {
	now := time.Now()
	a.LastUsedAt = &now
}

// Deactivate deactivates the API key
func (a *APIKey) Deactivate() {
	a.IsActive = false
}

// Activate activates the API key
func (a *APIKey) Activate() {
	a.IsActive = true
}
