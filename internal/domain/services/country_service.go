package services

import (
	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// CountryService implements business logic for country operations
type CountryService struct {
	countryRepo repositories.CountryRepository
}

// NewCountryService creates a new CountryService instance
func NewCountryService(countryRepo repositories.CountryRepository) *CountryService {
	return &CountryService{
		countryRepo: countryRepo,
	}
}

// CreateCountry creates a new country
func (s *CountryService) CreateCountry(
	countryCode string,
	iso31662 string,
	iso31663 string,
	name string,
	eea bool,
	callingCode string,
) (*entities.Country, error) {
	country := entities.NewCountry(countryCode, iso31662, iso31663, name, eea, callingCode)

	if err := s.countryRepo.Create(country); err != nil {
		return nil, err
	}

	return country, nil
}

// GetCountryByID retrieves a country by ID
func (s *CountryService) GetCountryByID(id uuid.UUID) (*entities.Country, error) {
	return s.countryRepo.GetByID(id)
}

// GetCountryByCode retrieves a country by country code
func (s *CountryService) GetCountryByCode(code string) (*entities.Country, error) {
	return s.countryRepo.GetByCode(code)
}

// GetCountryByISO31662 retrieves a country by ISO 3166-2 code
func (s *CountryService) GetCountryByISO31662(iso31662 string) (*entities.Country, error) {
	return s.countryRepo.GetByISO31662(iso31662)
}

// GetCountryByISO31663 retrieves a country by ISO 3166-3 code
func (s *CountryService) GetCountryByISO31663(iso31663 string) (*entities.Country, error) {
	return s.countryRepo.GetByISO31663(iso31663)
}

// GetCountryByName retrieves a country by name
func (s *CountryService) GetCountryByName(name string) (*entities.Country, error) {
	return s.countryRepo.GetByName(name)
}

// GetAllCountries retrieves all countries
func (s *CountryService) GetAllCountries() ([]*entities.Country, error) {
	return s.countryRepo.GetAll()
}

// GetCountriesByRegion retrieves countries by region code
func (s *CountryService) GetCountriesByRegion(regionCode string) ([]*entities.Country, error) {
	return s.countryRepo.GetByRegion(regionCode)
}

// GetCountriesBySubRegion retrieves countries by sub-region code
func (s *CountryService) GetCountriesBySubRegion(subRegionCode string) ([]*entities.Country, error) {
	return s.countryRepo.GetBySubRegion(subRegionCode)
}

// GetEEACountries retrieves all EEA countries
func (s *CountryService) GetEEACountries() ([]*entities.Country, error) {
	return s.countryRepo.GetEEACountries()
}

// UpdateCountry updates a country
func (s *CountryService) UpdateCountry(country *entities.Country) error {
	return s.countryRepo.Update(country)
}

// DeleteCountry deletes a country by ID
func (s *CountryService) DeleteCountry(id uuid.UUID) error {
	return s.countryRepo.Delete(id)
}

// CountryExists checks if a country exists by ID
func (s *CountryService) CountryExists(id uuid.UUID) (bool, error) {
	return s.countryRepo.Exists(id)
}

// GetCountryCount returns the total number of countries
func (s *CountryService) GetCountryCount() (int64, error) {
	return s.countryRepo.Count()
}
