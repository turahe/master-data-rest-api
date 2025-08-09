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

// CurrencyRepository implements the CurrencyRepository interface using pgx
type CurrencyRepository struct {
	pool *pgxpool.Pool
}

// NewCurrencyRepository creates a new CurrencyRepository instance
func NewCurrencyRepository(pool *pgxpool.Pool) *CurrencyRepository {
	return &CurrencyRepository{
		pool: pool,
	}
}

// Create creates a new currency in the database
func (r *CurrencyRepository) Create(ctx context.Context, currency *entities.Currency) error {
	currency.GenerateID()
	currency.CreatedAt = time.Now()
	currency.UpdatedAt = time.Now()

	query := `
		INSERT INTO tm_currencies (id, name, code, symbol, decimal_places, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, query,
		currency.ID, currency.Name, currency.Code, currency.Symbol,
		currency.DecimalPlaces, currency.IsActive, currency.CreatedAt, currency.UpdatedAt,
	)

	return err
}

// GetByID retrieves a currency by its ID
func (r *CurrencyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE id = $1`

	var currency entities.Currency
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&currency.ID, &currency.Name, &currency.Code, &currency.Symbol,
		&currency.DecimalPlaces, &currency.IsActive, &currency.CreatedAt, &currency.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// GetAll retrieves all currencies with pagination
func (r *CurrencyRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCurrencies(rows)
}

// Update updates an existing currency
func (r *CurrencyRepository) Update(ctx context.Context, currency *entities.Currency) error {
	currency.UpdatedAt = time.Now()

	query := `
		UPDATE tm_currencies SET
			name = $2, code = $3, symbol = $4, decimal_places = $5, is_active = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		currency.ID, currency.Name, currency.Code, currency.Symbol,
		currency.DecimalPlaces, currency.IsActive, currency.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("currency not found")
	}

	return nil
}

// Delete deletes a currency by ID
func (r *CurrencyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM tm_currencies WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("currency not found")
	}

	return nil
}

// Count returns the total number of currencies
func (r *CurrencyRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_currencies"

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Search searches currencies by name, code, or symbol
func (r *CurrencyRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Currency, error) {
	searchQuery := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE name ILIKE $1 OR code ILIKE $1 OR symbol ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCurrencies(rows)
}

// GetByName retrieves a currency by name
func (r *CurrencyRepository) GetByName(ctx context.Context, name string) (*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE name = $1`

	var currency entities.Currency
	row := r.pool.QueryRow(ctx, query, name)

	err := row.Scan(
		&currency.ID, &currency.Name, &currency.Code, &currency.Symbol,
		&currency.DecimalPlaces, &currency.IsActive, &currency.CreatedAt, &currency.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// GetByCode retrieves a currency by code
func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE code = $1`

	var currency entities.Currency
	row := r.pool.QueryRow(ctx, query, code)

	err := row.Scan(
		&currency.ID, &currency.Name, &currency.Code, &currency.Symbol,
		&currency.DecimalPlaces, &currency.IsActive, &currency.CreatedAt, &currency.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// GetBySymbol retrieves currencies by symbol
func (r *CurrencyRepository) GetBySymbol(ctx context.Context, symbol string) ([]*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE symbol = $1
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCurrencies(rows)
}

// GetActive retrieves all active currencies
func (r *CurrencyRepository) GetActive(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE is_active = true
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCurrencies(rows)
}

// GetInactive retrieves all inactive currencies
func (r *CurrencyRepository) GetInactive(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	query := `
		SELECT id, name, code, symbol, decimal_places, is_active, created_at, updated_at
		FROM tm_currencies
		WHERE is_active = false
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCurrencies(rows)
}

// Activate activates a currency
func (r *CurrencyRepository) Activate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_currencies SET
			is_active = true, updated_at = $2
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("currency not found")
	}

	return nil
}

// Deactivate deactivates a currency
func (r *CurrencyRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tm_currencies SET
			is_active = false, updated_at = $2
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("currency not found")
	}

	return nil
}

// ExistsByCode checks if a currency exists by code
func (r *CurrencyRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_currencies WHERE code = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, code).Scan(&exists)
	return exists, err
}

// ExistsByName checks if a currency exists by name
func (r *CurrencyRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_currencies WHERE name = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	return exists, err
}

// scanCurrencies is a helper method to scan rows into currency entities
func (r *CurrencyRepository) scanCurrencies(rows pgx.Rows) ([]*entities.Currency, error) {
	var currencies []*entities.Currency

	for rows.Next() {
		var currency entities.Currency
		err := rows.Scan(
			&currency.ID, &currency.Name, &currency.Code, &currency.Symbol,
			&currency.DecimalPlaces, &currency.IsActive, &currency.CreatedAt, &currency.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &currency)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
