package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CityRepository implements the CityRepository interface for MySQL
type CityRepository struct {
	db *sql.DB
}

// NewCityRepository creates a new city repository
func NewCityRepository(db *sql.DB) *CityRepository {
	return &CityRepository{db: db}
}

// Create creates a new city
func (r *CityRepository) Create(city *entities.City) error {
	query := `
		INSERT INTO tm_cities (province_id, name, type, code, latitude, longitude, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	city.CreatedAt = now
	city.UpdatedAt = now

	result, err := r.db.Exec(query,
		city.ProvinceID,
		city.Name,
		city.Type,
		city.Code,
		city.Latitude,
		city.Longitude,
		city.CreatedAt,
		city.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create city: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	city.ID = uint(id)
	return nil
}

// GetByID retrieves a city by ID
func (r *CityRepository) GetByID(id uint) (*entities.City, error) {
	query := `
		SELECT id, province_id, name, type, code, latitude, longitude, created_at, updated_at
		FROM tm_cities WHERE id = ?
	`

	city := &entities.City{}
	err := r.db.QueryRow(query, id).Scan(
		&city.ID,
		&city.ProvinceID,
		&city.Name,
		&city.Type,
		&city.Code,
		&city.Latitude,
		&city.Longitude,
		&city.CreatedAt,
		&city.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get city by id: %w", err)
	}

	return city, nil
}

// GetByName retrieves cities by name
func (r *CityRepository) GetByName(name string) ([]*entities.City, error) {
	query := `
		SELECT id, province_id, name, type, code, latitude, longitude, created_at, updated_at
		FROM tm_cities WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get cities by name: %w", err)
	}
	defer rows.Close()

	var cities []*entities.City
	for rows.Next() {
		city := &entities.City{}
		err := rows.Scan(
			&city.ID,
			&city.ProvinceID,
			&city.Name,
			&city.Type,
			&city.Code,
			&city.Latitude,
			&city.Longitude,
			&city.CreatedAt,
			&city.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// GetByCode retrieves a city by code
func (r *CityRepository) GetByCode(code string) (*entities.City, error) {
	query := `
		SELECT id, province_id, name, type, code, latitude, longitude, created_at, updated_at
		FROM tm_cities WHERE code = ?
	`

	city := &entities.City{}
	err := r.db.QueryRow(query, code).Scan(
		&city.ID,
		&city.ProvinceID,
		&city.Name,
		&city.Type,
		&city.Code,
		&city.Latitude,
		&city.Longitude,
		&city.CreatedAt,
		&city.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get city by code: %w", err)
	}

	return city, nil
}

// GetByProvinceID retrieves cities by province ID
func (r *CityRepository) GetByProvinceID(provinceID uint) ([]*entities.City, error) {
	query := `
		SELECT id, province_id, name, type, code, latitude, longitude, created_at, updated_at
		FROM tm_cities WHERE province_id = ?
	`

	rows, err := r.db.Query(query, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities by province id: %w", err)
	}
	defer rows.Close()

	var cities []*entities.City
	for rows.Next() {
		city := &entities.City{}
		err := rows.Scan(
			&city.ID,
			&city.ProvinceID,
			&city.Name,
			&city.Type,
			&city.Code,
			&city.Latitude,
			&city.Longitude,
			&city.CreatedAt,
			&city.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// GetByProvinceCode retrieves cities by province code
func (r *CityRepository) GetByProvinceCode(provinceCode string) ([]*entities.City, error) {
	query := `
		SELECT c.id, c.province_id, c.name, c.type, c.code, c.latitude, c.longitude, c.created_at, c.updated_at
		FROM tm_cities c
		JOIN tm_provinces p ON c.province_id = p.id
		WHERE p.code = ?
	`

	rows, err := r.db.Query(query, provinceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities by province code: %w", err)
	}
	defer rows.Close()

	var cities []*entities.City
	for rows.Next() {
		city := &entities.City{}
		err := rows.Scan(
			&city.ID,
			&city.ProvinceID,
			&city.Name,
			&city.Type,
			&city.Code,
			&city.Latitude,
			&city.Longitude,
			&city.CreatedAt,
			&city.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// GetAll retrieves all cities
func (r *CityRepository) GetAll() ([]*entities.City, error) {
	query := `
		SELECT id, province_id, name, type, code, latitude, longitude, created_at, updated_at
		FROM tm_cities ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all cities: %w", err)
	}
	defer rows.Close()

	var cities []*entities.City
	for rows.Next() {
		city := &entities.City{}
		err := rows.Scan(
			&city.ID,
			&city.ProvinceID,
			&city.Name,
			&city.Type,
			&city.Code,
			&city.Latitude,
			&city.Longitude,
			&city.CreatedAt,
			&city.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// Update updates a city
func (r *CityRepository) Update(city *entities.City) error {
	query := `
		UPDATE tm_cities 
		SET province_id = ?, name = ?, type = ?, code = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?
	`

	city.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		city.ProvinceID,
		city.Name,
		city.Type,
		city.Code,
		city.Latitude,
		city.Longitude,
		city.UpdatedAt,
		city.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update city: %w", err)
	}

	return nil
}

// Delete deletes a city by ID
func (r *CityRepository) Delete(id uint) error {
	query := `DELETE FROM tm_cities WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete city: %w", err)
	}

	return nil
}

// DeleteAll deletes all cities
func (r *CityRepository) DeleteAll() error {
	query := `DELETE FROM tm_cities`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all cities: %w", err)
	}

	return nil
}

// Exists checks if a city exists by ID
func (r *CityRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_cities WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if city exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of cities
func (r *CityRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_cities`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cities: %w", err)
	}

	return count, nil
}

// CountByProvince returns the number of cities in a province
func (r *CityRepository) CountByProvince(provinceID uint) (int64, error) {
	query := `SELECT COUNT(*) FROM tm_cities WHERE province_id = ?`

	var count int64
	err := r.db.QueryRow(query, provinceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cities by province: %w", err)
	}

	return count, nil
}
