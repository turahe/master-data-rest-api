package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CountryRepository implements the Country repository interface for MySQL
type CountryRepository struct {
	db *sql.DB
}

// NewCountryRepository creates a new CountryRepository instance
func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

// Create creates a new country
func (r *CountryRepository) Create(country *entities.Country) error {
	query := `
		INSERT INTO tm_countries (
			capital, citizenship, country_code, currency_name, currency_code,
			currency_sub_unit, currency_symbol, full_name, iso_3166_2, iso_3166_3,
			name, region_code, sub_region_code, eea, calling_code, flag,
			latitude, longitude, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		country.Capital, country.Citizenship, country.CountryCode, country.CurrencyName,
		country.CurrencyCode, country.CurrencySubUnit, country.CurrencySymbol,
		country.FullName, country.ISO31662, country.ISO31663, country.Name,
		country.RegionCode, country.SubRegionCode, country.EEA, country.CallingCode,
		country.Flag, country.Latitude, country.Longitude, country.CreatedAt, country.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create country: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	country.ID = uint(id)
	return nil
}

// GetByID retrieves a country by ID
func (r *CountryRepository) GetByID(id uint) (*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE id = ?`

	var country entities.Country
	err := r.db.QueryRow(query, id).Scan(
		&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
		&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
		&country.CurrencySymbol, &country.FullName, &country.ISO31662,
		&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
		&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
		&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get country by ID: %w", err)
	}

	return &country, nil
}

// GetByCode retrieves a country by country code
func (r *CountryRepository) GetByCode(code string) (*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE country_code = ?`

	var country entities.Country
	err := r.db.QueryRow(query, code).Scan(
		&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
		&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
		&country.CurrencySymbol, &country.FullName, &country.ISO31662,
		&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
		&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
		&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get country by code: %w", err)
	}

	return &country, nil
}

// GetByISO31662 retrieves a country by ISO 3166-2 code
func (r *CountryRepository) GetByISO31662(iso31662 string) (*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE iso_3166_2 = ?`

	var country entities.Country
	err := r.db.QueryRow(query, iso31662).Scan(
		&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
		&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
		&country.CurrencySymbol, &country.FullName, &country.ISO31662,
		&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
		&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
		&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get country by ISO 3166-2: %w", err)
	}

	return &country, nil
}

// GetByISO31663 retrieves a country by ISO 3166-3 code
func (r *CountryRepository) GetByISO31663(iso31663 string) (*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE iso_3166_3 = ?`

	var country entities.Country
	err := r.db.QueryRow(query, iso31663).Scan(
		&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
		&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
		&country.CurrencySymbol, &country.FullName, &country.ISO31662,
		&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
		&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
		&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get country by ISO 3166-3: %w", err)
	}

	return &country, nil
}

// GetByName retrieves a country by name
func (r *CountryRepository) GetByName(name string) (*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE name = ?`

	var country entities.Country
	err := r.db.QueryRow(query, name).Scan(
		&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
		&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
		&country.CurrencySymbol, &country.FullName, &country.ISO31662,
		&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
		&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
		&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get country by name: %w", err)
	}

	return &country, nil
}

// GetAll retrieves all countries
func (r *CountryRepository) GetAll() ([]*entities.Country, error) {
	query := `SELECT * FROM tm_countries ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all countries: %w", err)
	}
	defer rows.Close()

	var countries []*entities.Country
	for rows.Next() {
		var country entities.Country
		err := rows.Scan(
			&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
			&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
			&country.CurrencySymbol, &country.FullName, &country.ISO31662,
			&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
			&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
			&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// GetByRegion retrieves countries by region code
func (r *CountryRepository) GetByRegion(regionCode string) ([]*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE region_code = ? ORDER BY name`

	rows, err := r.db.Query(query, regionCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries by region: %w", err)
	}
	defer rows.Close()

	var countries []*entities.Country
	for rows.Next() {
		var country entities.Country
		err := rows.Scan(
			&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
			&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
			&country.CurrencySymbol, &country.FullName, &country.ISO31662,
			&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
			&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
			&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// GetBySubRegion retrieves countries by sub-region code
func (r *CountryRepository) GetBySubRegion(subRegionCode string) ([]*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE sub_region_code = ? ORDER BY name`

	rows, err := r.db.Query(query, subRegionCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries by sub-region: %w", err)
	}
	defer rows.Close()

	var countries []*entities.Country
	for rows.Next() {
		var country entities.Country
		err := rows.Scan(
			&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
			&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
			&country.CurrencySymbol, &country.FullName, &country.ISO31662,
			&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
			&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
			&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// GetEEACountries retrieves all EEA countries
func (r *CountryRepository) GetEEACountries() ([]*entities.Country, error) {
	query := `SELECT * FROM tm_countries WHERE eea = TRUE ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get EEA countries: %w", err)
	}
	defer rows.Close()

	var countries []*entities.Country
	for rows.Next() {
		var country entities.Country
		err := rows.Scan(
			&country.ID, &country.Capital, &country.Citizenship, &country.CountryCode,
			&country.CurrencyName, &country.CurrencyCode, &country.CurrencySubUnit,
			&country.CurrencySymbol, &country.FullName, &country.ISO31662,
			&country.ISO31663, &country.Name, &country.RegionCode, &country.SubRegionCode,
			&country.EEA, &country.CallingCode, &country.Flag, &country.Latitude,
			&country.Longitude, &country.CreatedAt, &country.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// Update updates a country
func (r *CountryRepository) Update(country *entities.Country) error {
	country.UpdatedAt = time.Now()

	query := `
		UPDATE tm_countries SET
			capital = ?, citizenship = ?, country_code = ?, currency_name = ?,
			currency_code = ?, currency_sub_unit = ?, currency_symbol = ?,
			full_name = ?, iso_3166_2 = ?, iso_3166_3 = ?, name = ?,
			region_code = ?, sub_region_code = ?, eea = ?, calling_code = ?,
			flag = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		country.Capital, country.Citizenship, country.CountryCode, country.CurrencyName,
		country.CurrencyCode, country.CurrencySubUnit, country.CurrencySymbol,
		country.FullName, country.ISO31662, country.ISO31663, country.Name,
		country.RegionCode, country.SubRegionCode, country.EEA, country.CallingCode,
		country.Flag, country.Latitude, country.Longitude, country.UpdatedAt, country.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update country: %w", err)
	}

	return nil
}

// Delete deletes a country by ID
func (r *CountryRepository) Delete(id uint) error {
	query := `DELETE FROM tm_countries WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete country: %w", err)
	}

	return nil
}

// Exists checks if a country exists by ID
func (r *CountryRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_countries WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check country existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of countries
func (r *CountryRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_countries`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count countries: %w", err)
	}

	return count, nil
}
