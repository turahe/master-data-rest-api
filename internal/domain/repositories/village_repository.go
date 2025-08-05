package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// VillageRepository defines the interface for village data operations
type VillageRepository interface {
	Create(village *entities.Village) error
	GetByID(id uuid.UUID) (*entities.Village, error)
	GetByName(name string) (*entities.Village, error)
	GetByCode(code string) (*entities.Village, error)
	GetByDistrictID(districtID uuid.UUID) ([]*entities.Village, error)
	GetByDistrictCode(districtCode string) ([]*entities.Village, error)
	GetByPostalCode(postalCode string) ([]*entities.Village, error)
	GetAll() ([]*entities.Village, error)
	Update(village *entities.Village) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
	CountByDistrict(districtID uuid.UUID) (int64, error)
}
