package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// LanguageRepository implements the LanguageRepository interface using pgx
type LanguageRepository struct {
	pool *pgxpool.Pool
}

// NewLanguageRepository creates a new LanguageRepository instance
func NewLanguageRepository(pool *pgxpool.Pool) *LanguageRepository {
	return &LanguageRepository{
		pool: pool,
	}
}

// Create creates a new language in the database
func (r *LanguageRepository) Create(ctx context.Context, language *entities.Language) error {
	language.GenerateID()
	language.CreatedAt = time.Now()
	language.UpdatedAt = time.Now()

	query := `
		INSERT INTO tm_languages (id, name, code, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		language.ID, language.Name, language.Code, language.IsActive,
		language.CreatedAt, language.UpdatedAt,
	)

	return err
}

// GetByID retrieves a language by its ID
func (r *LanguageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE id = $1`

	var language entities.Language
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&language.ID, &language.Name, &language.Code, &language.IsActive,
		&language.CreatedAt, &language.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("language not found")
		}
		return nil, err
	}

	return &language, nil
}

// GetAll retrieves all languages with pagination
func (r *LanguageRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLanguages(rows)
}

// Update updates an existing language
func (r *LanguageRepository) Update(ctx context.Context, language *entities.Language) error {
	language.UpdatedAt = time.Now()

	query := `
		UPDATE tm_languages SET
			name = $2, code = $3, is_active = $4, updated_at = $5
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		language.ID, language.Name, language.Code, language.IsActive, language.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

// Delete deletes a language by ID
func (r *LanguageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM tm_languages WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

// Count returns the total number of languages
func (r *LanguageRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_languages"

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Search searches languages by name or code
func (r *LanguageRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Language, error) {
	searchQuery := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE name ILIKE $1 OR code ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLanguages(rows)
}

// GetByName retrieves a language by name
func (r *LanguageRepository) GetByName(ctx context.Context, name string) (*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE name = $1`

	var language entities.Language
	row := r.pool.QueryRow(ctx, query, name)

	err := row.Scan(
		&language.ID, &language.Name, &language.Code, &language.IsActive,
		&language.CreatedAt, &language.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("language not found")
		}
		return nil, err
	}

	return &language, nil
}

// GetByCode retrieves a language by code
func (r *LanguageRepository) GetByCode(ctx context.Context, code string) (*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE code = $1`

	var language entities.Language
	row := r.pool.QueryRow(ctx, query, code)

	err := row.Scan(
		&language.ID, &language.Name, &language.Code, &language.IsActive,
		&language.CreatedAt, &language.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("language not found")
		}
		return nil, err
	}

	return &language, nil
}

// GetActive retrieves all active languages
func (r *LanguageRepository) GetActive(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE is_active = true
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLanguages(rows)
}

// GetInactive retrieves all inactive languages
func (r *LanguageRepository) GetInactive(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM tm_languages
		WHERE is_active = false
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLanguages(rows)
}

// Activate activates a language
func (r *LanguageRepository) Activate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_languages SET
			is_active = true, updated_at = $2
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

// Deactivate deactivates a language
func (r *LanguageRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_languages SET
			is_active = false, updated_at = $2
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

// ExistsByCode checks if a language exists by code
func (r *LanguageRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_languages WHERE code = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, code).Scan(&exists)
	return exists, err
}

// ExistsByName checks if a language exists by name
func (r *LanguageRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_languages WHERE name = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	return exists, err
}

// scanLanguages is a helper method to scan rows into language entities
func (r *LanguageRepository) scanLanguages(rows pgx.Rows) ([]*entities.Language, error) {
	var languages []*entities.Language

	for rows.Next() {
		var language entities.Language
		err := rows.Scan(
			&language.ID, &language.Name, &language.Code, &language.IsActive,
			&language.CreatedAt, &language.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		languages = append(languages, &language)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}
