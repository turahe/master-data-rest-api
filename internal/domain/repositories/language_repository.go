package repositories

import "github.com/turahe/master-data-rest-api/internal/domain/entities"

// LanguageRepository defines the interface for language data operations
type LanguageRepository interface {
	Create(language *entities.Language) error
	GetByID(id uint) (*entities.Language, error)
	GetByCode(code string) (*entities.Language, error)
	GetByName(name string) (*entities.Language, error)
	GetByNative(native string) (*entities.Language, error)
	GetAll() ([]*entities.Language, error)
	Update(language *entities.Language) error
	Delete(id uint) error
	DeleteAll() error
	Exists(id uint) (bool, error)
	Count() (int64, error)
}
