package services

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// CityService implements business logic for city operations
type CityService struct {
	cityRepo repositories.CityRepository
}

// NewCityService creates a new CityService instance
func NewCityService(cityRepo repositories.CityRepository) *CityService {
	return &CityService{
		cityRepo: cityRepo,
	}
}

// CreateCity creates a new city
func (s *CityService) CreateCity(provinceID uuid.UUID, name string) (*entities.City, error) {
	city := entities.NewCity(provinceID, name)

	if err := s.cityRepo.Create(city); err != nil {
		return nil, err
	}

	return city, nil
}

// GetCityByID retrieves a city by ID
func (s *CityService) GetCityByID(id uuid.UUID) (*entities.City, error) {
	return s.cityRepo.GetByID(id)
}

// GetCityByName retrieves a city by name
func (s *CityService) GetCityByName(name string) (*entities.City, error) {
	return s.cityRepo.GetByName(name)
}

// GetCityByCode retrieves a city by code
func (s *CityService) GetCityByCode(code string) (*entities.City, error) {
	return s.cityRepo.GetByCode(code)
}

// GetCitiesByProvinceID retrieves cities by province ID
func (s *CityService) GetCitiesByProvinceID(provinceID uuid.UUID) ([]*entities.City, error) {
	return s.cityRepo.GetByProvinceID(provinceID)
}

// GetCitiesByProvinceCode retrieves cities by province code
func (s *CityService) GetCitiesByProvinceCode(provinceCode string) ([]*entities.City, error) {
	return s.cityRepo.GetByProvinceCode(provinceCode)
}

// GetAllCities retrieves all cities
func (s *CityService) GetAllCities() ([]*entities.City, error) {
	return s.cityRepo.GetAll()
}

// UpdateCity updates a city
func (s *CityService) UpdateCity(city *entities.City) error {
	return s.cityRepo.Update(city)
}

// DeleteCity deletes a city by ID
func (s *CityService) DeleteCity(id uuid.UUID) error {
	return s.cityRepo.Delete(id)
}

// CityExists checks if a city exists by ID
func (s *CityService) CityExists(id uuid.UUID) (bool, error) {
	return s.cityRepo.Exists(id)
}

// GetCityCount returns the total number of cities
func (s *CityService) GetCityCount() (int64, error) {
	return s.cityRepo.Count()
}

// GetCityCountByProvince returns the number of cities in a province
func (s *CityService) GetCityCountByProvince(provinceID uuid.UUID) (int64, error) {
	return s.cityRepo.CountByProvince(provinceID)
}
