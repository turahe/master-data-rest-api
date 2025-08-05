package repositories

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CountryRepository defines the interface for country data operations
type CountryRepository interface {
	Create(country *entities.Country) error
	GetByID(id uuid.UUID) (*entities.Country, error)
	GetByCode(code string) (*entities.Country, error)
	GetByISO31662(iso31662 string) (*entities.Country, error)
	GetByISO31663(iso31663 string) (*entities.Country, error)
	GetByName(name string) (*entities.Country, error)
	GetAll() ([]*entities.Country, error)
	GetByRegion(regionCode string) ([]*entities.Country, error)
	GetBySubRegion(subRegionCode string) ([]*entities.Country, error)
	GetEEACountries() ([]*entities.Country, error)
	Update(country *entities.Country) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	Exists(id uuid.UUID) (bool, error)
	Count() (int64, error)
}
