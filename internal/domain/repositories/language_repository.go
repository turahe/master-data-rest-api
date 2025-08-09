package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// LanguageRepository defines the interface for language data operations
type LanguageRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, language *entities.Language) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Language, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Language, error)
	Update(ctx context.Context, language *entities.Language) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Language, error)
	GetByName(ctx context.Context, name string) (*entities.Language, error)
	GetByCode(ctx context.Context, code string) (*entities.Language, error)

	// Status operations
	GetActive(ctx context.Context, limit, offset int) ([]*entities.Language, error)
	GetInactive(ctx context.Context, limit, offset int) ([]*entities.Language, error)
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error

	// Validation operations
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
