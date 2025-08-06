package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// BankRepository defines the interface for bank data operations
type BankRepository interface {
	Create(bank *entities.Bank) error
	GetByID(id uuid.UUID) (*entities.Bank, error)
	GetByCode(code string) (*entities.Bank, error)
	GetByName(name string) ([]*entities.Bank, error)
	GetByAlias(alias string) ([]*entities.Bank, error)
	GetByCompany(company string) ([]*entities.Bank, error)
	GetAll() ([]*entities.Bank, error)
	Update(bank *entities.Bank) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
}
