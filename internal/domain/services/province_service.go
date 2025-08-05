package services

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// ProvinceService implements business logic for province operations
type ProvinceService struct {
	provinceRepo repositories.ProvinceRepository
}

// NewProvinceService creates a new ProvinceService instance
func NewProvinceService(provinceRepo repositories.ProvinceRepository) *ProvinceService {
	return &ProvinceService{
		provinceRepo: provinceRepo,
	}
}

// CreateProvince creates a new province
func (s *ProvinceService) CreateProvince(countryID uuid.UUID, name string) (*entities.Province, error) {
	province := entities.NewProvince(countryID, name)

	if err := s.provinceRepo.Create(province); err != nil {
		return nil, err
	}

	return province, nil
}

// GetProvinceByID retrieves a province by ID
func (s *ProvinceService) GetProvinceByID(id uuid.UUID) (*entities.Province, error) {
	return s.provinceRepo.GetByID(id)
}

// GetProvinceByName retrieves a province by name
func (s *ProvinceService) GetProvinceByName(name string) (*entities.Province, error) {
	return s.provinceRepo.GetByName(name)
}

// GetProvinceByCode retrieves a province by code
func (s *ProvinceService) GetProvinceByCode(code string) (*entities.Province, error) {
	return s.provinceRepo.GetByCode(code)
}

// GetProvincesByCountryID retrieves provinces by country ID
func (s *ProvinceService) GetProvincesByCountryID(countryID uuid.UUID) ([]*entities.Province, error) {
	return s.provinceRepo.GetByCountryID(countryID)
}

// GetProvincesByCountryCode retrieves provinces by country code
func (s *ProvinceService) GetProvincesByCountryCode(countryCode string) ([]*entities.Province, error) {
	return s.provinceRepo.GetByCountryCode(countryCode)
}

// GetAllProvinces retrieves all provinces
func (s *ProvinceService) GetAllProvinces() ([]*entities.Province, error) {
	return s.provinceRepo.GetAll()
}

// UpdateProvince updates a province
func (s *ProvinceService) UpdateProvince(province *entities.Province) error {
	return s.provinceRepo.Update(province)
}

// DeleteProvince deletes a province by ID
func (s *ProvinceService) DeleteProvince(id uuid.UUID) error {
	return s.provinceRepo.Delete(id)
}

// ProvinceExists checks if a province exists by ID
func (s *ProvinceService) ProvinceExists(id uuid.UUID) (bool, error) {
	return s.provinceRepo.Exists(id)
}

// GetProvinceCount returns the total number of provinces
func (s *ProvinceService) GetProvinceCount() (int64, error) {
	return s.provinceRepo.Count()
}

// GetProvinceCountByCountry returns the number of provinces in a country
func (s *ProvinceService) GetProvinceCountByCountry(countryID uuid.UUID) (int64, error) {
	return s.provinceRepo.CountByCountry(countryID)
}
