package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// APIKeyRepository defines the interface for API key data operations
type APIKeyRepository interface {
	// Create creates a new API key
	Create(ctx context.Context, apiKey *entities.APIKey) error

	// GetByKey retrieves an API key by its key value
	GetByKey(ctx context.Context, key string) (*entities.APIKey, error)

	// GetByID retrieves an API key by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.APIKey, error)

	// GetAll retrieves all API keys with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*entities.APIKey, error)

	// Update updates an existing API key
	Update(ctx context.Context, apiKey *entities.APIKey) error

	// Delete soft deletes an API key by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateLastUsed updates the last used timestamp for an API key
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error

	// Activate activates an API key
	Activate(ctx context.Context, id uuid.UUID) error

	// Deactivate deactivates an API key
	Deactivate(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of API keys
	Count(ctx context.Context) (int64, error)

	// Search searches API keys
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.APIKey, error)

	// ValidateKey validates an API key and returns it if valid
	ValidateKey(ctx context.Context, key string) (*entities.APIKey, error)
}
