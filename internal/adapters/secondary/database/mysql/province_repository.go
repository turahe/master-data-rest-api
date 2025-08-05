package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// ProvinceRepository implements the Province repository interface for MySQL
type ProvinceRepository struct {
	db *sql.DB
}

// NewProvinceRepository creates a new ProvinceRepository instance
func NewProvinceRepository(db *sql.DB) *ProvinceRepository {
	return &ProvinceRepository{db: db}
}

// Create creates a new province
func (r *ProvinceRepository) Create(province *entities.Province) error {
	query := `
		INSERT INTO tm_provinces (
			country_id, name, region, iso_3166_2, code, latitude, longitude,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		province.CountryID, province.Name, province.Region, province.ISO31662,
		province.Code, province.Latitude, province.Longitude,
		province.CreatedAt, province.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create province: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	province.ID = uint(id)
	return nil
}

// GetByID retrieves a province by ID
func (r *ProvinceRepository) GetByID(id uint) (*entities.Province, error) {
	query := `SELECT * FROM tm_provinces WHERE id = ?`

	var province entities.Province
	err := r.db.QueryRow(query, id).Scan(
		&province.ID, &province.CountryID, &province.Name, &province.Region,
		&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
		&province.CreatedAt, &province.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get province by ID: %w", err)
	}

	return &province, nil
}

// GetByName retrieves a province by name
func (r *ProvinceRepository) GetByName(name string) (*entities.Province, error) {
	query := `SELECT * FROM tm_provinces WHERE name = ?`

	var province entities.Province
	err := r.db.QueryRow(query, name).Scan(
		&province.ID, &province.CountryID, &province.Name, &province.Region,
		&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
		&province.CreatedAt, &province.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get province by name: %w", err)
	}

	return &province, nil
}

// GetByCode retrieves a province by code
func (r *ProvinceRepository) GetByCode(code string) (*entities.Province, error) {
	query := `SELECT * FROM tm_provinces WHERE code = ?`

	var province entities.Province
	err := r.db.QueryRow(query, code).Scan(
		&province.ID, &province.CountryID, &province.Name, &province.Region,
		&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
		&province.CreatedAt, &province.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get province by code: %w", err)
	}

	return &province, nil
}

// GetByCountryID retrieves provinces by country ID
func (r *ProvinceRepository) GetByCountryID(countryID uint) ([]*entities.Province, error) {
	query := `SELECT * FROM tm_provinces WHERE country_id = ? ORDER BY name`

	rows, err := r.db.Query(query, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces by country ID: %w", err)
	}
	defer rows.Close()

	var provinces []*entities.Province
	for rows.Next() {
		var province entities.Province
		err := rows.Scan(
			&province.ID, &province.CountryID, &province.Name, &province.Region,
			&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
			&province.CreatedAt, &province.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan province: %w", err)
		}
		provinces = append(provinces, &province)
	}

	return provinces, nil
}

// GetByCountryCode retrieves provinces by country code
func (r *ProvinceRepository) GetByCountryCode(countryCode string) ([]*entities.Province, error) {
	query := `
		SELECT p.* FROM tm_provinces p
		JOIN tm_countries c ON p.country_id = c.id
		WHERE c.country_code = ?
		ORDER BY p.name
	`

	rows, err := r.db.Query(query, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces by country code: %w", err)
	}
	defer rows.Close()

	var provinces []*entities.Province
	for rows.Next() {
		var province entities.Province
		err := rows.Scan(
			&province.ID, &province.CountryID, &province.Name, &province.Region,
			&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
			&province.CreatedAt, &province.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan province: %w", err)
		}
		provinces = append(provinces, &province)
	}

	return provinces, nil
}

// GetAll retrieves all provinces
func (r *ProvinceRepository) GetAll() ([]*entities.Province, error) {
	query := `SELECT * FROM tm_provinces ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all provinces: %w", err)
	}
	defer rows.Close()

	var provinces []*entities.Province
	for rows.Next() {
		var province entities.Province
		err := rows.Scan(
			&province.ID, &province.CountryID, &province.Name, &province.Region,
			&province.ISO31662, &province.Code, &province.Latitude, &province.Longitude,
			&province.CreatedAt, &province.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan province: %w", err)
		}
		provinces = append(provinces, &province)
	}

	return provinces, nil
}

// Update updates a province
func (r *ProvinceRepository) Update(province *entities.Province) error {
	province.UpdatedAt = time.Now()

	query := `
		UPDATE tm_provinces SET
			country_id = ?, name = ?, region = ?, iso_3166_2 = ?,
			code = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		province.CountryID, province.Name, province.Region, province.ISO31662,
		province.Code, province.Latitude, province.Longitude, province.UpdatedAt, province.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update province: %w", err)
	}

	return nil
}

// Delete deletes a province by ID
func (r *ProvinceRepository) Delete(id uint) error {
	query := `DELETE FROM tm_provinces WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete province: %w", err)
	}

	return nil
}

// Exists checks if a province exists by ID
func (r *ProvinceRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_provinces WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check province existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of provinces
func (r *ProvinceRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_provinces`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count provinces: %w", err)
	}

	return count, nil
}

// CountByCountry returns the number of provinces in a country
func (r *ProvinceRepository) CountByCountry(countryID uint) (int64, error) {
	query := `SELECT COUNT(*) FROM tm_provinces WHERE country_id = ?`

	var count int64
	err := r.db.QueryRow(query, countryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count provinces by country: %w", err)
	}

	return count, nil
}
