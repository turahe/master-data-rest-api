package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CurrencyRepository defines the interface for currency data operations
type CurrencyRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, currency *entities.Currency) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Currency, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Currency, error)
	Update(ctx context.Context, currency *entities.Currency) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Currency, error)
	GetByName(ctx context.Context, name string) (*entities.Currency, error)
	GetByCode(ctx context.Context, code string) (*entities.Currency, error)
	GetBySymbol(ctx context.Context, symbol string) ([]*entities.Currency, error)

	// Status operations
	GetActive(ctx context.Context, limit, offset int) ([]*entities.Currency, error)
	GetInactive(ctx context.Context, limit, offset int) ([]*entities.Currency, error)
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error

	// Validation operations
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
