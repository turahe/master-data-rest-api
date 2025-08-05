package gorm

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"gorm.io/gorm"
)

// CountryRepository implements the CountryRepository interface using GORM
type CountryRepository struct {
	db *gorm.DB
}

// NewCountryRepository creates a new GORM country repository
func NewCountryRepository(db *gorm.DB) repositories.CountryRepository {
	return &CountryRepository{
		db: db,
	}
}

// Create creates a new country in the database
func (r *CountryRepository) Create(country *entities.Country) error {
	result := r.db.Create(country)
	if result.Error != nil {
		return fmt.Errorf("failed to create country: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a country by ID from the database
func (r *CountryRepository) GetByID(id uuid.UUID) (*entities.Country, error) {
	var country entities.Country
	result := r.db.Where("id = ?", id).First(&country)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country by ID: %w", result.Error)
	}
	return &country, nil
}

// GetByCode retrieves a country by country code from the database
func (r *CountryRepository) GetByCode(code string) (*entities.Country, error) {
	var country entities.Country
	result := r.db.Where("country_code = ?", code).First(&country)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country by code: %w", result.Error)
	}
	return &country, nil
}

// GetByISO31662 retrieves a country by ISO 3166-2 code from the database
func (r *CountryRepository) GetByISO31662(iso31662 string) (*entities.Country, error) {
	var country entities.Country
	result := r.db.Where("iso_3166_2 = ?", iso31662).First(&country)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country by ISO 3166-2: %w", result.Error)
	}
	return &country, nil
}

// GetByISO31663 retrieves a country by ISO 3166-3 code from the database
func (r *CountryRepository) GetByISO31663(iso31663 string) (*entities.Country, error) {
	var country entities.Country
	result := r.db.Where("iso_3166_3 = ?", iso31663).First(&country)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country by ISO 3166-3: %w", result.Error)
	}
	return &country, nil
}

// GetByName retrieves a country by name from the database
func (r *CountryRepository) GetByName(name string) (*entities.Country, error) {
	var country entities.Country
	result := r.db.Where("name = ?", name).First(&country)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country by name: %w", result.Error)
	}
	return &country, nil
}

// GetAll retrieves all countries from the database
func (r *CountryRepository) GetAll() ([]*entities.Country, error) {
	var countries []entities.Country
	result := r.db.Order("name ASC").Find(&countries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all countries: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	countryPtrs := make([]*entities.Country, len(countries))
	for i := range countries {
		countryPtrs[i] = &countries[i]
	}

	return countryPtrs, nil
}

// GetByRegion retrieves countries by region code from the database
func (r *CountryRepository) GetByRegion(regionCode string) ([]*entities.Country, error) {
	var countries []entities.Country
	result := r.db.Where("region_code = ?", regionCode).Order("name ASC").Find(&countries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get countries by region: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	countryPtrs := make([]*entities.Country, len(countries))
	for i := range countries {
		countryPtrs[i] = &countries[i]
	}

	return countryPtrs, nil
}

// GetBySubRegion retrieves countries by sub-region code from the database
func (r *CountryRepository) GetBySubRegion(subRegionCode string) ([]*entities.Country, error) {
	var countries []entities.Country
	result := r.db.Where("sub_region_code = ?", subRegionCode).Order("name ASC").Find(&countries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get countries by sub-region: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	countryPtrs := make([]*entities.Country, len(countries))
	for i := range countries {
		countryPtrs[i] = &countries[i]
	}

	return countryPtrs, nil
}

// GetEEACountries retrieves all EEA countries from the database
func (r *CountryRepository) GetEEACountries() ([]*entities.Country, error) {
	var countries []entities.Country
	result := r.db.Where("eea = ?", true).Order("name ASC").Find(&countries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get EEA countries: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	countryPtrs := make([]*entities.Country, len(countries))
	for i := range countries {
		countryPtrs[i] = &countries[i]
	}

	return countryPtrs, nil
}

// Update updates an existing country in the database
func (r *CountryRepository) Update(country *entities.Country) error {
	result := r.db.Save(country)
	if result.Error != nil {
		return fmt.Errorf("failed to update country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

// Delete deletes a country by ID from the database
func (r *CountryRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entities.Country{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

// DeleteAll deletes all countries from the database
func (r *CountryRepository) DeleteAll() error {
	result := r.db.Delete(&entities.Country{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete all countries: %w", result.Error)
	}
	return nil
}

// Exists checks if a country exists by ID in the database
func (r *CountryRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.Model(&entities.Country{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check country existence: %w", result.Error)
	}

	return count > 0, nil
}

// Count returns the total number of countries in the database
func (r *CountryRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Country{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count countries: %w", result.Error)
	}

	return count, nil
}
