package services

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// DistrictService implements business logic for district operations
type DistrictService struct {
	districtRepo repositories.DistrictRepository
}

// NewDistrictService creates a new DistrictService instance
func NewDistrictService(districtRepo repositories.DistrictRepository) *DistrictService {
	return &DistrictService{
		districtRepo: districtRepo,
	}
}

// CreateDistrict creates a new district
func (s *DistrictService) CreateDistrict(cityID uuid.UUID, name string) (*entities.District, error) {
	district := entities.NewDistrict(cityID, name)

	if err := s.districtRepo.Create(district); err != nil {
		return nil, err
	}

	return district, nil
}

// GetDistrictByID retrieves a district by ID
func (s *DistrictService) GetDistrictByID(id uuid.UUID) (*entities.District, error) {
	return s.districtRepo.GetByID(id)
}

// GetDistrictByName retrieves a district by name
func (s *DistrictService) GetDistrictByName(name string) (*entities.District, error) {
	return s.districtRepo.GetByName(name)
}

// GetDistrictByCode retrieves a district by code
func (s *DistrictService) GetDistrictByCode(code string) (*entities.District, error) {
	return s.districtRepo.GetByCode(code)
}

// GetDistrictsByCityID retrieves districts by city ID
func (s *DistrictService) GetDistrictsByCityID(cityID uuid.UUID) ([]*entities.District, error) {
	return s.districtRepo.GetByCityID(cityID)
}

// GetDistrictsByCityCode retrieves districts by city code
func (s *DistrictService) GetDistrictsByCityCode(cityCode string) ([]*entities.District, error) {
	return s.districtRepo.GetByCityCode(cityCode)
}

// GetAllDistricts retrieves all districts
func (s *DistrictService) GetAllDistricts() ([]*entities.District, error) {
	return s.districtRepo.GetAll()
}

// UpdateDistrict updates a district
func (s *DistrictService) UpdateDistrict(district *entities.District) error {
	return s.districtRepo.Update(district)
}

// DeleteDistrict deletes a district by ID
func (s *DistrictService) DeleteDistrict(id uuid.UUID) error {
	return s.districtRepo.Delete(id)
}

// DistrictExists checks if a district exists by ID
func (s *DistrictService) DistrictExists(id uuid.UUID) (bool, error) {
	return s.districtRepo.Exists(id)
}

// GetDistrictCount returns the total number of districts
func (s *DistrictService) GetDistrictCount() (int64, error) {
	return s.districtRepo.Count()
}

// GetDistrictCountByCity returns the number of districts in a city
func (s *DistrictService) GetDistrictCountByCity(cityID uuid.UUID) (int64, error) {
	return s.districtRepo.CountByCity(cityID)
}
