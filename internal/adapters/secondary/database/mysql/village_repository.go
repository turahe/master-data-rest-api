package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// VillageRepository implements the VillageRepository interface for MySQL
type VillageRepository struct {
	db *sql.DB
}

// NewVillageRepository creates a new village repository
func NewVillageRepository(db *sql.DB) *VillageRepository {
	return &VillageRepository{db: db}
}

// Create creates a new village
func (r *VillageRepository) Create(village *entities.Village) error {
	query := `
		INSERT INTO tm_villages (district_id, name, code, postal_code, latitude, longitude, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	village.CreatedAt = now
	village.UpdatedAt = now

	result, err := r.db.Exec(query,
		village.DistrictID,
		village.Name,
		village.Code,
		village.PostalCode,
		village.Latitude,
		village.Longitude,
		village.CreatedAt,
		village.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create village: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	village.ID = uint(id)
	return nil
}

// GetByID retrieves a village by ID
func (r *VillageRepository) GetByID(id uint) (*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages WHERE id = ?
	`

	village := &entities.Village{}
	err := r.db.QueryRow(query, id).Scan(
		&village.ID,
		&village.DistrictID,
		&village.Name,
		&village.Code,
		&village.PostalCode,
		&village.Latitude,
		&village.Longitude,
		&village.CreatedAt,
		&village.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get village by id: %w", err)
	}

	return village, nil
}

// GetByName retrieves villages by name
func (r *VillageRepository) GetByName(name string) ([]*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get villages by name: %w", err)
	}
	defer rows.Close()

	var villages []*entities.Village
	for rows.Next() {
		village := &entities.Village{}
		err := rows.Scan(
			&village.ID,
			&village.DistrictID,
			&village.Name,
			&village.Code,
			&village.PostalCode,
			&village.Latitude,
			&village.Longitude,
			&village.CreatedAt,
			&village.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan village: %w", err)
		}
		villages = append(villages, village)
	}

	return villages, nil
}

// GetByCode retrieves a village by code
func (r *VillageRepository) GetByCode(code string) (*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages WHERE code = ?
	`

	village := &entities.Village{}
	err := r.db.QueryRow(query, code).Scan(
		&village.ID,
		&village.DistrictID,
		&village.Name,
		&village.Code,
		&village.PostalCode,
		&village.Latitude,
		&village.Longitude,
		&village.CreatedAt,
		&village.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get village by code: %w", err)
	}

	return village, nil
}

// GetByDistrictID retrieves villages by district ID
func (r *VillageRepository) GetByDistrictID(districtID uint) ([]*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages WHERE district_id = ?
	`

	rows, err := r.db.Query(query, districtID)
	if err != nil {
		return nil, fmt.Errorf("failed to get villages by district id: %w", err)
	}
	defer rows.Close()

	var villages []*entities.Village
	for rows.Next() {
		village := &entities.Village{}
		err := rows.Scan(
			&village.ID,
			&village.DistrictID,
			&village.Name,
			&village.Code,
			&village.PostalCode,
			&village.Latitude,
			&village.Longitude,
			&village.CreatedAt,
			&village.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan village: %w", err)
		}
		villages = append(villages, village)
	}

	return villages, nil
}

// GetByDistrictCode retrieves villages by district code
func (r *VillageRepository) GetByDistrictCode(districtCode string) ([]*entities.Village, error) {
	query := `
		SELECT v.id, v.district_id, v.name, v.code, v.postal_code, v.latitude, v.longitude, v.created_at, v.updated_at
		FROM tm_villages v
		JOIN tm_districts d ON v.district_id = d.id
		WHERE d.code = ?
	`

	rows, err := r.db.Query(query, districtCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get villages by district code: %w", err)
	}
	defer rows.Close()

	var villages []*entities.Village
	for rows.Next() {
		village := &entities.Village{}
		err := rows.Scan(
			&village.ID,
			&village.DistrictID,
			&village.Name,
			&village.Code,
			&village.PostalCode,
			&village.Latitude,
			&village.Longitude,
			&village.CreatedAt,
			&village.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan village: %w", err)
		}
		villages = append(villages, village)
	}

	return villages, nil
}

// GetByPostalCode retrieves villages by postal code
func (r *VillageRepository) GetByPostalCode(postalCode string) ([]*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages WHERE postal_code = ?
	`

	rows, err := r.db.Query(query, postalCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get villages by postal code: %w", err)
	}
	defer rows.Close()

	var villages []*entities.Village
	for rows.Next() {
		village := &entities.Village{}
		err := rows.Scan(
			&village.ID,
			&village.DistrictID,
			&village.Name,
			&village.Code,
			&village.PostalCode,
			&village.Latitude,
			&village.Longitude,
			&village.CreatedAt,
			&village.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan village: %w", err)
		}
		villages = append(villages, village)
	}

	return villages, nil
}

// GetAll retrieves all villages
func (r *VillageRepository) GetAll() ([]*entities.Village, error) {
	query := `
		SELECT id, district_id, name, code, postal_code, latitude, longitude, created_at, updated_at
		FROM tm_villages ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all villages: %w", err)
	}
	defer rows.Close()

	var villages []*entities.Village
	for rows.Next() {
		village := &entities.Village{}
		err := rows.Scan(
			&village.ID,
			&village.DistrictID,
			&village.Name,
			&village.Code,
			&village.PostalCode,
			&village.Latitude,
			&village.Longitude,
			&village.CreatedAt,
			&village.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan village: %w", err)
		}
		villages = append(villages, village)
	}

	return villages, nil
}

// Update updates a village
func (r *VillageRepository) Update(village *entities.Village) error {
	query := `
		UPDATE tm_villages 
		SET district_id = ?, name = ?, code = ?, postal_code = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?
	`

	village.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		village.DistrictID,
		village.Name,
		village.Code,
		village.PostalCode,
		village.Latitude,
		village.Longitude,
		village.UpdatedAt,
		village.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update village: %w", err)
	}

	return nil
}

// Delete deletes a village by ID
func (r *VillageRepository) Delete(id uint) error {
	query := `DELETE FROM tm_villages WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete village: %w", err)
	}

	return nil
}

// DeleteAll deletes all villages
func (r *VillageRepository) DeleteAll() error {
	query := `DELETE FROM tm_villages`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all villages: %w", err)
	}

	return nil
}

// Exists checks if a village exists by ID
func (r *VillageRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_villages WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if village exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of villages
func (r *VillageRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_villages`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count villages: %w", err)
	}

	return count, nil
}

// CountByDistrict returns the number of villages in a district
func (r *VillageRepository) CountByDistrict(districtID uint) (int64, error) {
	query := `SELECT COUNT(*) FROM tm_villages WHERE district_id = ?`

	var count int64
	err := r.db.QueryRow(query, districtID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count villages by district: %w", err)
	}

	return count, nil
}
