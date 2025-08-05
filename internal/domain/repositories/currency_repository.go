package repositories

import "github.com/turahe/master-data-rest-api/internal/domain/entities"

// CurrencyRepository defines the interface for currency data operations
type CurrencyRepository interface {
	Create(currency *entities.Currency) error
	GetByID(id uint) (*entities.Currency, error)
	GetByISOCode(isoCode string) (*entities.Currency, error)
	GetByName(name string) (*entities.Currency, error)
	GetBySymbol(symbol string) ([]*entities.Currency, error)
	GetByISONumeric(isoNumeric string) (*entities.Currency, error)
	GetAll() ([]*entities.Currency, error)
	GetByPriority(priority int) ([]*entities.Currency, error)
	Update(currency *entities.Currency) error
	Delete(id uint) error
	DeleteAll() error
	Exists(id uint) (bool, error)
	Count() (int64, error)
}
