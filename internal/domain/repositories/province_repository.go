package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// ProvinceRepository defines the interface for province data operations
type ProvinceRepository interface {
	Create(province *entities.Province) error
	GetByID(id uuid.UUID) (*entities.Province, error)
	GetByName(name string) (*entities.Province, error)
	GetByCode(code string) (*entities.Province, error)
	GetByCountryID(countryID uuid.UUID) ([]*entities.Province, error)
	GetByCountryCode(countryCode string) ([]*entities.Province, error)
	GetAll() ([]*entities.Province, error)
	Update(province *entities.Province) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
	CountByCountry(countryID uuid.UUID) (int64, error)
}
