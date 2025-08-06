package gorm

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"gorm.io/gorm"
)

// ProvinceRepository implements the ProvinceRepository interface using GORM
type ProvinceRepository struct {
	db *gorm.DB
}

// NewProvinceRepository creates a new GORM province repository
func NewProvinceRepository(db *gorm.DB) repositories.ProvinceRepository {
	return &ProvinceRepository{
		db: db,
	}
}

// Create creates a new province in the database
func (r *ProvinceRepository) Create(province *entities.Province) error {
	result := r.db.Create(province)
	if result.Error != nil {
		return fmt.Errorf("failed to create province: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a province by ID from the database
func (r *ProvinceRepository) GetByID(id uuid.UUID) (*entities.Province, error) {
	var province entities.Province
	result := r.db.Where("id = ?", id).First(&province)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("province not found")
		}
		return nil, fmt.Errorf("failed to get province by ID: %w", result.Error)
	}
	return &province, nil
}

// GetByCode retrieves a province by code from the database
func (r *ProvinceRepository) GetByCode(code string) (*entities.Province, error) {
	var province entities.Province
	result := r.db.Where("code = ?", code).First(&province)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("province not found")
		}
		return nil, fmt.Errorf("failed to get province by code: %w", result.Error)
	}
	return &province, nil
}

// GetByName retrieves a province by name from the database
func (r *ProvinceRepository) GetByName(name string) (*entities.Province, error) {
	var province entities.Province
	result := r.db.Where("name = ?", name).First(&province)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("province not found")
		}
		return nil, fmt.Errorf("failed to get province by name: %w", result.Error)
	}
	return &province, nil
}

// GetByCountryID retrieves provinces by country ID from the database
func (r *ProvinceRepository) GetByCountryID(countryID uuid.UUID) ([]*entities.Province, error) {
	var provinces []entities.Province
	result := r.db.Where("country_id = ?", countryID).Find(&provinces)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get provinces by country ID: %w", result.Error)
	}

	// Convert to slice of pointers
	var provincePtrs []*entities.Province
	for i := range provinces {
		provincePtrs = append(provincePtrs, &provinces[i])
	}
	return provincePtrs, nil
}

// GetByCountryCode retrieves provinces by country code from the database
func (r *ProvinceRepository) GetByCountryCode(countryCode string) ([]*entities.Province, error) {
	var provinces []entities.Province
	result := r.db.Joins("JOIN tm_countries ON tm_provinces.country_id = tm_countries.id").
		Where("tm_countries.country_code = ?", countryCode).
		Find(&provinces)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get provinces by country code: %w", result.Error)
	}

	// Convert to slice of pointers
	var provincePtrs []*entities.Province
	for i := range provinces {
		provincePtrs = append(provincePtrs, &provinces[i])
	}
	return provincePtrs, nil
}

// GetAll retrieves all provinces from the database
func (r *ProvinceRepository) GetAll() ([]*entities.Province, error) {
	var provinces []entities.Province
	result := r.db.Find(&provinces)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all provinces: %w", result.Error)
	}

	// Convert to slice of pointers
	var provincePtrs []*entities.Province
	for i := range provinces {
		provincePtrs = append(provincePtrs, &provinces[i])
	}
	return provincePtrs, nil
}

// Update updates a province in the database
func (r *ProvinceRepository) Update(province *entities.Province) error {
	result := r.db.Save(province)
	if result.Error != nil {
		return fmt.Errorf("failed to update province: %w", result.Error)
	}
	return nil
}

// Delete deletes a province by ID from the database
func (r *ProvinceRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entities.Province{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete province: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("province not found")
	}
	return nil
}

// DeleteAll deletes all provinces from the database
func (r *ProvinceRepository) DeleteAll() error {
	result := r.db.Delete(&entities.Province{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete all provinces: %w", result.Error)
	}
	return nil
}

// Exists checks if a province exists by ID
func (r *ProvinceRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.Model(&entities.Province{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check province existence: %w", result.Error)
	}
	return count > 0, nil
}

// Count returns the total number of provinces
func (r *ProvinceRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Province{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count provinces: %w", result.Error)
	}
	return count, nil
}

// CountByCountry returns the number of provinces in a country
func (r *ProvinceRepository) CountByCountry(countryID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.Model(&entities.Province{}).Where("country_id = ?", countryID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count provinces by country: %w", result.Error)
	}
	return count, nil
} 