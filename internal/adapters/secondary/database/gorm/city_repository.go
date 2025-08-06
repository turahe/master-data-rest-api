package gorm

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"gorm.io/gorm"
)

// CityRepository implements the CityRepository interface using GORM
type CityRepository struct {
	db *gorm.DB
}

// NewCityRepository creates a new GORM city repository
func NewCityRepository(db *gorm.DB) repositories.CityRepository {
	return &CityRepository{
		db: db,
	}
}

// Create creates a new city in the database
func (r *CityRepository) Create(city *entities.City) error {
	result := r.db.Create(city)
	if result.Error != nil {
		return fmt.Errorf("failed to create city: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a city by ID from the database
func (r *CityRepository) GetByID(id uuid.UUID) (*entities.City, error) {
	var city entities.City
	result := r.db.Where("id = ?", id).First(&city)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("city not found")
		}
		return nil, fmt.Errorf("failed to get city by ID: %w", result.Error)
	}
	return &city, nil
}

// GetByCode retrieves a city by code from the database
func (r *CityRepository) GetByCode(code string) (*entities.City, error) {
	var city entities.City
	result := r.db.Where("code = ?", code).First(&city)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("city not found")
		}
		return nil, fmt.Errorf("failed to get city by code: %w", result.Error)
	}
	return &city, nil
}

// GetByName retrieves a city by name from the database
func (r *CityRepository) GetByName(name string) (*entities.City, error) {
	var city entities.City
	result := r.db.Where("name = ?", name).First(&city)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("city not found")
		}
		return nil, fmt.Errorf("failed to get city by name: %w", result.Error)
	}
	return &city, nil
}

// GetByProvinceID retrieves cities by province ID from the database
func (r *CityRepository) GetByProvinceID(provinceID uuid.UUID) ([]*entities.City, error) {
	var cities []entities.City
	result := r.db.Where("province_id = ?", provinceID).Find(&cities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get cities by province ID: %w", result.Error)
	}

	// Convert to slice of pointers
	var cityPtrs []*entities.City
	for i := range cities {
		cityPtrs = append(cityPtrs, &cities[i])
	}
	return cityPtrs, nil
}

// GetAll retrieves all cities from the database
func (r *CityRepository) GetAll() ([]*entities.City, error) {
	var cities []entities.City
	result := r.db.Find(&cities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all cities: %w", result.Error)
	}

	// Convert to slice of pointers
	var cityPtrs []*entities.City
	for i := range cities {
		cityPtrs = append(cityPtrs, &cities[i])
	}
	return cityPtrs, nil
}

// Update updates a city in the database
func (r *CityRepository) Update(city *entities.City) error {
	result := r.db.Save(city)
	if result.Error != nil {
		return fmt.Errorf("failed to update city: %w", result.Error)
	}
	return nil
}

// Delete deletes a city by ID from the database
func (r *CityRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entities.City{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete city: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("city not found")
	}
	return nil
}

// DeleteAll deletes all cities from the database
func (r *CityRepository) DeleteAll() error {
	result := r.db.Delete(&entities.City{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete all cities: %w", result.Error)
	}
	return nil
}

// Exists checks if a city exists by ID
func (r *CityRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.Model(&entities.City{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check city existence: %w", result.Error)
	}
	return count > 0, nil
}

// Count returns the total number of cities
func (r *CityRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.City{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count cities: %w", result.Error)
	}
	return count, nil
}

// CountByProvince returns the number of cities in a province
func (r *CityRepository) CountByProvince(provinceID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.Model(&entities.City{}).Where("province_id = ?", provinceID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count cities by province: %w", result.Error)
	}
	return count, nil
}

// GetByProvinceCode retrieves cities by province code from the database
func (r *CityRepository) GetByProvinceCode(provinceCode string) ([]*entities.City, error) {
	var cities []entities.City
	result := r.db.Joins("JOIN tm_provinces ON tm_cities.province_id = tm_provinces.id").
		Where("tm_provinces.code = ?", provinceCode).
		Find(&cities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get cities by province code: %w", result.Error)
	}

	// Convert to slice of pointers
	var cityPtrs []*entities.City
	for i := range cities {
		cityPtrs = append(cityPtrs, &cities[i])
	}
	return cityPtrs, nil
} 