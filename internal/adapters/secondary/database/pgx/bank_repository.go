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

// BankRepository implements the BankRepository interface using pgx
type BankRepository struct {
	pool *pgxpool.Pool
}

// NewBankRepository creates a new BankRepository instance
func NewBankRepository(pool *pgxpool.Pool) *BankRepository {
	return &BankRepository{
		pool: pool,
	}
}

// Create creates a new bank in the database
func (r *BankRepository) Create(ctx context.Context, bank *entities.Bank) error {
	bank.GenerateID()
	bank.CreatedAt = time.Now()
	bank.UpdatedAt = time.Now()

	query := `
		INSERT INTO tm_banks (id, name, alias, company, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		bank.ID, bank.Name, bank.Alias, bank.Company, bank.Code,
		bank.CreatedAt, bank.UpdatedAt,
	)

	return err
}

// GetByID retrieves a bank by its ID
func (r *BankRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE id = $1`

	var bank entities.Bank
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&bank.ID, &bank.Name, &bank.Alias, &bank.Company, &bank.Code,
		&bank.CreatedAt, &bank.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, err
	}

	return &bank, nil
}

// GetAll retrieves all banks with pagination
func (r *BankRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBanks(rows)
}

// Update updates an existing bank
func (r *BankRepository) Update(ctx context.Context, bank *entities.Bank) error {
	bank.UpdatedAt = time.Now()

	query := `
		UPDATE tm_banks SET
			name = $2, alias = $3, company = $4, code = $5, updated_at = $6
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		bank.ID, bank.Name, bank.Alias, bank.Company, bank.Code, bank.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("bank not found")
	}

	return nil
}

// Delete deletes a bank by ID
func (r *BankRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM tm_banks WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("bank not found")
	}

	return nil
}

// Count returns the total number of banks
func (r *BankRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_banks"

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Search searches banks by name, alias, company, or code
func (r *BankRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error) {
	searchQuery := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE name ILIKE $1 OR alias ILIKE $1 OR company ILIKE $1 OR code ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBanks(rows)
}

// GetByName retrieves a bank by name
func (r *BankRepository) GetByName(ctx context.Context, name string) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE name = $1`

	var bank entities.Bank
	row := r.pool.QueryRow(ctx, query, name)

	err := row.Scan(
		&bank.ID, &bank.Name, &bank.Alias, &bank.Company, &bank.Code,
		&bank.CreatedAt, &bank.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, err
	}

	return &bank, nil
}

// GetByCode retrieves a bank by code
func (r *BankRepository) GetByCode(ctx context.Context, code string) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE code = $1`

	var bank entities.Bank
	row := r.pool.QueryRow(ctx, query, code)

	err := row.Scan(
		&bank.ID, &bank.Name, &bank.Alias, &bank.Company, &bank.Code,
		&bank.CreatedAt, &bank.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, err
	}

	return &bank, nil
}

// GetByAlias retrieves a bank by alias
func (r *BankRepository) GetByAlias(ctx context.Context, alias string) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE alias = $1`

	var bank entities.Bank
	row := r.pool.QueryRow(ctx, query, alias)

	err := row.Scan(
		&bank.ID, &bank.Name, &bank.Alias, &bank.Company, &bank.Code,
		&bank.CreatedAt, &bank.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, err
	}

	return &bank, nil
}

// GetByCompany retrieves banks by company
func (r *BankRepository) GetByCompany(ctx context.Context, company string, limit, offset int) ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks
		WHERE company = $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, company, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBanks(rows)
}

// ExistsByCode checks if a bank exists by code
func (r *BankRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_banks WHERE code = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, code).Scan(&exists)
	return exists, err
}

// ExistsByName checks if a bank exists by name
func (r *BankRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tm_banks WHERE name = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	return exists, err
}

// scanBanks is a helper method to scan rows into bank entities
func (r *BankRepository) scanBanks(rows pgx.Rows) ([]*entities.Bank, error) {
	var banks []*entities.Bank

	for rows.Next() {
		var bank entities.Bank
		err := rows.Scan(
			&bank.ID, &bank.Name, &bank.Alias, &bank.Company, &bank.Code,
			&bank.CreatedAt, &bank.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banks = append(banks, &bank)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

// Truncate removes all bank records efficiently using TRUNCATE
func (r *BankRepository) Truncate(ctx context.Context) error {
	query := `TRUNCATE TABLE tm_banks RESTART IDENTITY CASCADE`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate banks table: %w", err)
	}
	return nil
}
