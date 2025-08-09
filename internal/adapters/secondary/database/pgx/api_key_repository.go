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

// APIKeyRepository implements the APIKeyRepository interface using pgx
type APIKeyRepository struct {
	pool *pgxpool.Pool
}

// NewAPIKeyRepository creates a new APIKeyRepository instance
func NewAPIKeyRepository(pool *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{
		pool: pool,
	}
}

// Create creates a new API key in the database
func (r *APIKeyRepository) Create(ctx context.Context, apiKey *entities.APIKey) error {
	if apiKey.ID == uuid.Nil {
		apiKey.ID = uuid.New()
	}
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()

	query := `
		INSERT INTO tm_api_keys (
			id, name, key, description, is_active, expires_at, last_used_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	_, err := r.pool.Exec(ctx, query,
		apiKey.ID, apiKey.Name, apiKey.Key, apiKey.Description, apiKey.IsActive,
		apiKey.ExpiresAt, apiKey.LastUsedAt, apiKey.CreatedAt, apiKey.UpdatedAt,
	)

	return err
}

// GetByID retrieves an API key by its ID
func (r *APIKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.APIKey, error) {
	query := `
		SELECT id, name, key, description, is_active, expires_at, last_used_at, created_at, updated_at, deleted_at
		FROM tm_api_keys
		WHERE id = $1 AND deleted_at IS NULL`

	var apiKey entities.APIKey
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&apiKey.ID, &apiKey.Name, &apiKey.Key, &apiKey.Description, &apiKey.IsActive,
		&apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.CreatedAt, &apiKey.UpdatedAt, &apiKey.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

// GetByKey retrieves an API key by its key value
func (r *APIKeyRepository) GetByKey(ctx context.Context, key string) (*entities.APIKey, error) {
	query := `
		SELECT id, name, key, description, is_active, expires_at, last_used_at, created_at, updated_at, deleted_at
		FROM tm_api_keys
		WHERE key = $1 AND deleted_at IS NULL`

	var apiKey entities.APIKey
	row := r.pool.QueryRow(ctx, query, key)

	err := row.Scan(
		&apiKey.ID, &apiKey.Name, &apiKey.Key, &apiKey.Description, &apiKey.IsActive,
		&apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.CreatedAt, &apiKey.UpdatedAt, &apiKey.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for not found, not an error for validation
		}
		return nil, err
	}

	return &apiKey, nil
}

// GetAll retrieves all API keys with optional pagination
func (r *APIKeyRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.APIKey, error) {
	query := `
		SELECT id, name, key, description, is_active, expires_at, last_used_at, created_at, updated_at, deleted_at
		FROM tm_api_keys
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []*entities.APIKey
	for rows.Next() {
		var apiKey entities.APIKey
		err := rows.Scan(
			&apiKey.ID, &apiKey.Name, &apiKey.Key, &apiKey.Description, &apiKey.IsActive,
			&apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.CreatedAt, &apiKey.UpdatedAt, &apiKey.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// Update updates an existing API key
func (r *APIKeyRepository) Update(ctx context.Context, apiKey *entities.APIKey) error {
	apiKey.UpdatedAt = time.Now()

	query := `
		UPDATE tm_api_keys SET
			name = $2, description = $3, is_active = $4, expires_at = $5,
			last_used_at = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query,
		apiKey.ID, apiKey.Name, apiKey.Description, apiKey.IsActive,
		apiKey.ExpiresAt, apiKey.LastUsedAt, apiKey.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp for an API key
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE tm_api_keys SET
			last_used_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id, now)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// Delete soft deletes an API key by setting deleted_at timestamp
func (r *APIKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE tm_api_keys SET
			deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id, now)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// Activate activates an API key
func (r *APIKeyRepository) Activate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_api_keys SET
			is_active = true, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// Deactivate deactivates an API key
func (r *APIKeyRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_api_keys SET
			is_active = false, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// Count returns the total number of active API keys
func (r *APIKeyRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_api_keys WHERE deleted_at IS NULL"

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Search searches API keys by name
func (r *APIKeyRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.APIKey, error) {
	searchQuery := `
		SELECT id, name, key, description, is_active, expires_at, last_used_at, created_at, updated_at, deleted_at
		FROM tm_api_keys
		WHERE deleted_at IS NULL AND (name ILIKE $1 OR description ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []*entities.APIKey
	for rows.Next() {
		var apiKey entities.APIKey
		err := rows.Scan(
			&apiKey.ID, &apiKey.Name, &apiKey.Key, &apiKey.Description, &apiKey.IsActive,
			&apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.CreatedAt, &apiKey.UpdatedAt, &apiKey.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// ValidateKey validates an API key and returns it if valid
func (r *APIKeyRepository) ValidateKey(ctx context.Context, key string) (*entities.APIKey, error) {
	apiKey, err := r.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if apiKey == nil {
		return nil, nil // Key not found
	}

	// Check if the key is valid (active and not expired)
	if !apiKey.IsValid() {
		return nil, nil // Key is invalid
	}

	// Update last used timestamp
	apiKey.UpdateLastUsed()
	if err := r.UpdateLastUsed(ctx, apiKey.ID); err != nil {
		// Log the error but don't fail the validation
		// The key is still valid even if we can't update the timestamp
	}

	return apiKey, nil
}
