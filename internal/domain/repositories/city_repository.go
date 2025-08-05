package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CityRepository defines the interface for city data operations
type CityRepository interface {
	Create(city *entities.City) error
	GetByID(id uuid.UUID) (*entities.City, error)
	GetByName(name string) (*entities.City, error)
	GetByCode(code string) (*entities.City, error)
	GetByProvinceID(provinceID uuid.UUID) ([]*entities.City, error)
	GetByProvinceCode(provinceCode string) ([]*entities.City, error)
	GetAll() ([]*entities.City, error)
	Update(city *entities.City) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
	CountByProvince(provinceID uuid.UUID) (int64, error)
}
