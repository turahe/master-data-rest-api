package services

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// VillageService implements business logic for village operations
type VillageService struct {
	villageRepo repositories.VillageRepository
}

// NewVillageService creates a new VillageService instance
func NewVillageService(villageRepo repositories.VillageRepository) *VillageService {
	return &VillageService{
		villageRepo: villageRepo,
	}
}

// CreateVillage creates a new village
func (s *VillageService) CreateVillage(districtID uuid.UUID, name string) (*entities.Village, error) {
	village := entities.NewVillage(districtID, name)

	if err := s.villageRepo.Create(village); err != nil {
		return nil, err
	}

	return village, nil
}

// GetVillageByID retrieves a village by ID
func (s *VillageService) GetVillageByID(id uuid.UUID) (*entities.Village, error) {
	return s.villageRepo.GetByID(id)
}

// GetVillageByName retrieves a village by name
func (s *VillageService) GetVillageByName(name string) (*entities.Village, error) {
	return s.villageRepo.GetByName(name)
}

// GetVillageByCode retrieves a village by code
func (s *VillageService) GetVillageByCode(code string) (*entities.Village, error) {
	return s.villageRepo.GetByCode(code)
}

// GetVillagesByDistrictID retrieves villages by district ID
func (s *VillageService) GetVillagesByDistrictID(districtID uuid.UUID) ([]*entities.Village, error) {
	return s.villageRepo.GetByDistrictID(districtID)
}

// GetVillagesByDistrictCode retrieves villages by district code
func (s *VillageService) GetVillagesByDistrictCode(districtCode string) ([]*entities.Village, error) {
	return s.villageRepo.GetByDistrictCode(districtCode)
}

// GetVillagesByPostalCode retrieves villages by postal code
func (s *VillageService) GetVillagesByPostalCode(postalCode string) ([]*entities.Village, error) {
	return s.villageRepo.GetByPostalCode(postalCode)
}

// GetAllVillages retrieves all villages
func (s *VillageService) GetAllVillages() ([]*entities.Village, error) {
	return s.villageRepo.GetAll()
}

// UpdateVillage updates a village
func (s *VillageService) UpdateVillage(village *entities.Village) error {
	return s.villageRepo.Update(village)
}

// DeleteVillage deletes a village by ID
func (s *VillageService) DeleteVillage(id uuid.UUID) error {
	return s.villageRepo.Delete(id)
}

// VillageExists checks if a village exists by ID
func (s *VillageService) VillageExists(id uuid.UUID) (bool, error) {
	return s.villageRepo.Exists(id)
}

// GetVillageCount returns the total number of villages
func (s *VillageService) GetVillageCount() (int64, error) {
	return s.villageRepo.Count()
}

// GetVillageCountByDistrict returns the number of villages in a district
func (s *VillageService) GetVillageCountByDistrict(districtID uuid.UUID) (int64, error) {
	return s.villageRepo.CountByDistrict(districtID)
}
