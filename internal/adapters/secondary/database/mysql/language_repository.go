package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// LanguageRepository implements the LanguageRepository interface for MySQL
type LanguageRepository struct {
	db *sql.DB
}

// NewLanguageRepository creates a new language repository
func NewLanguageRepository(db *sql.DB) *LanguageRepository {
	return &LanguageRepository{db: db}
}

// Create creates a new language
func (r *LanguageRepository) Create(language *entities.Language) error {
	query := `
		INSERT INTO tm_languages (code, name, native, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	language.CreatedAt = now
	language.UpdatedAt = now

	result, err := r.db.Exec(query,
		language.Code,
		language.Name,
		language.Native,
		language.CreatedAt,
		language.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create language: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	language.ID = uint(id)
	return nil
}

// GetByID retrieves a language by ID
func (r *LanguageRepository) GetByID(id uint) (*entities.Language, error) {
	query := `
		SELECT id, code, name, native, created_at, updated_at
		FROM tm_languages WHERE id = ?
	`

	language := &entities.Language{}
	err := r.db.QueryRow(query, id).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.Native,
		&language.CreatedAt,
		&language.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get language by id: %w", err)
	}

	return language, nil
}

// GetByCode retrieves a language by code
func (r *LanguageRepository) GetByCode(code string) (*entities.Language, error) {
	query := `
		SELECT id, code, name, native, created_at, updated_at
		FROM tm_languages WHERE code = ?
	`

	language := &entities.Language{}
	err := r.db.QueryRow(query, code).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.Native,
		&language.CreatedAt,
		&language.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get language by code: %w", err)
	}

	return language, nil
}

// GetByName retrieves languages by name
func (r *LanguageRepository) GetByName(name string) ([]*entities.Language, error) {
	query := `
		SELECT id, code, name, native, created_at, updated_at
		FROM tm_languages WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get languages by name: %w", err)
	}
	defer rows.Close()

	var languages []*entities.Language
	for rows.Next() {
		language := &entities.Language{}
		err := rows.Scan(
			&language.ID,
			&language.Code,
			&language.Name,
			&language.Native,
			&language.CreatedAt,
			&language.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, language)
	}

	return languages, nil
}

// GetByNative retrieves languages by native name
func (r *LanguageRepository) GetByNative(native string) ([]*entities.Language, error) {
	query := `
		SELECT id, code, name, native, created_at, updated_at
		FROM tm_languages WHERE native LIKE ?
	`

	rows, err := r.db.Query(query, "%"+native+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get languages by native: %w", err)
	}
	defer rows.Close()

	var languages []*entities.Language
	for rows.Next() {
		language := &entities.Language{}
		err := rows.Scan(
			&language.ID,
			&language.Code,
			&language.Name,
			&language.Native,
			&language.CreatedAt,
			&language.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, language)
	}

	return languages, nil
}

// GetAll retrieves all languages
func (r *LanguageRepository) GetAll() ([]*entities.Language, error) {
	query := `
		SELECT id, code, name, native, created_at, updated_at
		FROM tm_languages ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all languages: %w", err)
	}
	defer rows.Close()

	var languages []*entities.Language
	for rows.Next() {
		language := &entities.Language{}
		err := rows.Scan(
			&language.ID,
			&language.Code,
			&language.Name,
			&language.Native,
			&language.CreatedAt,
			&language.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, language)
	}

	return languages, nil
}

// Update updates a language
func (r *LanguageRepository) Update(language *entities.Language) error {
	query := `
		UPDATE tm_languages 
		SET code = ?, name = ?, native = ?, updated_at = ?
		WHERE id = ?
	`

	language.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		language.Code,
		language.Name,
		language.Native,
		language.UpdatedAt,
		language.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update language: %w", err)
	}

	return nil
}

// Delete deletes a language by ID
func (r *LanguageRepository) Delete(id uint) error {
	query := `DELETE FROM tm_languages WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete language: %w", err)
	}

	return nil
}

// DeleteAll deletes all languages
func (r *LanguageRepository) DeleteAll() error {
	query := `DELETE FROM tm_languages`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all languages: %w", err)
	}

	return nil
}

// Exists checks if a language exists by ID
func (r *LanguageRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_languages WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if language exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of languages
func (r *LanguageRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_languages`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count languages: %w", err)
	}

	return count, nil
}
