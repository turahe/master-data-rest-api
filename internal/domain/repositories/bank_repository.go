package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// BankRepository defines the interface for bank data operations
type BankRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, bank *entities.Bank) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Bank, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Bank, error)
	Update(ctx context.Context, bank *entities.Bank) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error)
	GetByName(ctx context.Context, name string) (*entities.Bank, error)
	GetByCode(ctx context.Context, code string) (*entities.Bank, error)
	GetByAlias(ctx context.Context, alias string) (*entities.Bank, error)
	GetByCompany(ctx context.Context, company string, limit, offset int) ([]*entities.Bank, error)

	// Validation operations
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
