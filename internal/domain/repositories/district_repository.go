package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// DistrictRepository defines the interface for district data operations
type DistrictRepository interface {
	Create(district *entities.District) error
	GetByID(id uuid.UUID) (*entities.District, error)
	GetByName(name string) (*entities.District, error)
	GetByCode(code string) (*entities.District, error)
	GetByCityID(cityID uuid.UUID) ([]*entities.District, error)
	GetByCityCode(cityCode string) ([]*entities.District, error)
	GetAll() ([]*entities.District, error)
	Update(district *entities.District) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
	CountByCity(cityID uuid.UUID) (int64, error)
}
