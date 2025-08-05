package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// DistrictRepository implements the DistrictRepository interface for MySQL
type DistrictRepository struct {
	db *sql.DB
}

// NewDistrictRepository creates a new district repository
func NewDistrictRepository(db *sql.DB) *DistrictRepository {
	return &DistrictRepository{db: db}
}

// Create creates a new district
func (r *DistrictRepository) Create(district *entities.District) error {
	query := `
		INSERT INTO tm_districts (city_id, name, code, latitude, longitude, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	district.CreatedAt = now
	district.UpdatedAt = now

	result, err := r.db.Exec(query,
		district.CityID,
		district.Name,
		district.Code,
		district.Latitude,
		district.Longitude,
		district.CreatedAt,
		district.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create district: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	district.ID = uint(id)
	return nil
}

// GetByID retrieves a district by ID
func (r *DistrictRepository) GetByID(id uint) (*entities.District, error) {
	query := `
		SELECT id, city_id, name, code, latitude, longitude, created_at, updated_at
		FROM tm_districts WHERE id = ?
	`

	district := &entities.District{}
	err := r.db.QueryRow(query, id).Scan(
		&district.ID,
		&district.CityID,
		&district.Name,
		&district.Code,
		&district.Latitude,
		&district.Longitude,
		&district.CreatedAt,
		&district.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get district by id: %w", err)
	}

	return district, nil
}

// GetByName retrieves districts by name
func (r *DistrictRepository) GetByName(name string) ([]*entities.District, error) {
	query := `
		SELECT id, city_id, name, code, latitude, longitude, created_at, updated_at
		FROM tm_districts WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get districts by name: %w", err)
	}
	defer rows.Close()

	var districts []*entities.District
	for rows.Next() {
		district := &entities.District{}
		err := rows.Scan(
			&district.ID,
			&district.CityID,
			&district.Name,
			&district.Code,
			&district.Latitude,
			&district.Longitude,
			&district.CreatedAt,
			&district.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan district: %w", err)
		}
		districts = append(districts, district)
	}

	return districts, nil
}

// GetByCode retrieves a district by code
func (r *DistrictRepository) GetByCode(code string) (*entities.District, error) {
	query := `
		SELECT id, city_id, name, code, latitude, longitude, created_at, updated_at
		FROM tm_districts WHERE code = ?
	`

	district := &entities.District{}
	err := r.db.QueryRow(query, code).Scan(
		&district.ID,
		&district.CityID,
		&district.Name,
		&district.Code,
		&district.Latitude,
		&district.Longitude,
		&district.CreatedAt,
		&district.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get district by code: %w", err)
	}

	return district, nil
}

// GetByCityID retrieves districts by city ID
func (r *DistrictRepository) GetByCityID(cityID uint) ([]*entities.District, error) {
	query := `
		SELECT id, city_id, name, code, latitude, longitude, created_at, updated_at
		FROM tm_districts WHERE city_id = ?
	`

	rows, err := r.db.Query(query, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get districts by city id: %w", err)
	}
	defer rows.Close()

	var districts []*entities.District
	for rows.Next() {
		district := &entities.District{}
		err := rows.Scan(
			&district.ID,
			&district.CityID,
			&district.Name,
			&district.Code,
			&district.Latitude,
			&district.Longitude,
			&district.CreatedAt,
			&district.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan district: %w", err)
		}
		districts = append(districts, district)
	}

	return districts, nil
}

// GetByCityCode retrieves districts by city code
func (r *DistrictRepository) GetByCityCode(cityCode string) ([]*entities.District, error) {
	query := `
		SELECT d.id, d.city_id, d.name, d.code, d.latitude, d.longitude, d.created_at, d.updated_at
		FROM tm_districts d
		JOIN tm_cities c ON d.city_id = c.id
		WHERE c.code = ?
	`

	rows, err := r.db.Query(query, cityCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get districts by city code: %w", err)
	}
	defer rows.Close()

	var districts []*entities.District
	for rows.Next() {
		district := &entities.District{}
		err := rows.Scan(
			&district.ID,
			&district.CityID,
			&district.Name,
			&district.Code,
			&district.Latitude,
			&district.Longitude,
			&district.CreatedAt,
			&district.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan district: %w", err)
		}
		districts = append(districts, district)
	}

	return districts, nil
}

// GetAll retrieves all districts
func (r *DistrictRepository) GetAll() ([]*entities.District, error) {
	query := `
		SELECT id, city_id, name, code, latitude, longitude, created_at, updated_at
		FROM tm_districts ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all districts: %w", err)
	}
	defer rows.Close()

	var districts []*entities.District
	for rows.Next() {
		district := &entities.District{}
		err := rows.Scan(
			&district.ID,
			&district.CityID,
			&district.Name,
			&district.Code,
			&district.Latitude,
			&district.Longitude,
			&district.CreatedAt,
			&district.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan district: %w", err)
		}
		districts = append(districts, district)
	}

	return districts, nil
}

// Update updates a district
func (r *DistrictRepository) Update(district *entities.District) error {
	query := `
		UPDATE tm_districts 
		SET city_id = ?, name = ?, code = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?
	`

	district.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		district.CityID,
		district.Name,
		district.Code,
		district.Latitude,
		district.Longitude,
		district.UpdatedAt,
		district.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update district: %w", err)
	}

	return nil
}

// Delete deletes a district by ID
func (r *DistrictRepository) Delete(id uint) error {
	query := `DELETE FROM tm_districts WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete district: %w", err)
	}

	return nil
}

// DeleteAll deletes all districts
func (r *DistrictRepository) DeleteAll() error {
	query := `DELETE FROM tm_districts`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all districts: %w", err)
	}

	return nil
}

// Exists checks if a district exists by ID
func (r *DistrictRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_districts WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if district exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of districts
func (r *DistrictRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_districts`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count districts: %w", err)
	}

	return count, nil
}

// CountByCity returns the number of districts in a city
func (r *DistrictRepository) CountByCity(cityID uint) (int64, error) {
	query := `SELECT COUNT(*) FROM tm_districts WHERE city_id = ?`

	var count int64
	err := r.db.QueryRow(query, cityID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count districts by city: %w", err)
	}

	return count, nil
}
