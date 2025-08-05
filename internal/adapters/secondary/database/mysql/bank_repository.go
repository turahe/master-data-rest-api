package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// BankRepository implements the BankRepository interface for MySQL
type BankRepository struct {
	db *sql.DB
}

// NewBankRepository creates a new bank repository
func NewBankRepository(db *sql.DB) *BankRepository {
	return &BankRepository{db: db}
}

// Create creates a new bank
func (r *BankRepository) Create(bank *entities.Bank) error {
	query := `
		INSERT INTO tm_banks (name, alias, company, code, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	bank.CreatedAt = now
	bank.UpdatedAt = now

	result, err := r.db.Exec(query,
		bank.Name,
		bank.Alias,
		bank.Company,
		bank.Code,
		bank.CreatedAt,
		bank.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create bank: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	bank.ID = uint(id)
	return nil
}

// GetByID retrieves a bank by ID
func (r *BankRepository) GetByID(id uint) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks WHERE id = ?
	`

	bank := &entities.Bank{}
	err := r.db.QueryRow(query, id).Scan(
		&bank.ID,
		&bank.Name,
		&bank.Alias,
		&bank.Company,
		&bank.Code,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank by id: %w", err)
	}

	return bank, nil
}

// GetByCode retrieves a bank by code
func (r *BankRepository) GetByCode(code string) (*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks WHERE code = ?
	`

	bank := &entities.Bank{}
	err := r.db.QueryRow(query, code).Scan(
		&bank.ID,
		&bank.Name,
		&bank.Alias,
		&bank.Company,
		&bank.Code,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank by code: %w", err)
	}

	return bank, nil
}

// GetByName retrieves banks by name
func (r *BankRepository) GetByName(name string) ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get banks by name: %w", err)
	}
	defer rows.Close()

	var banks []*entities.Bank
	for rows.Next() {
		bank := &entities.Bank{}
		err := rows.Scan(
			&bank.ID,
			&bank.Name,
			&bank.Alias,
			&bank.Company,
			&bank.Code,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank: %w", err)
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// GetByAlias retrieves banks by alias
func (r *BankRepository) GetByAlias(alias string) ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks WHERE alias LIKE ?
	`

	rows, err := r.db.Query(query, "%"+alias+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get banks by alias: %w", err)
	}
	defer rows.Close()

	var banks []*entities.Bank
	for rows.Next() {
		bank := &entities.Bank{}
		err := rows.Scan(
			&bank.ID,
			&bank.Name,
			&bank.Alias,
			&bank.Company,
			&bank.Code,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank: %w", err)
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// GetByCompany retrieves banks by company
func (r *BankRepository) GetByCompany(company string) ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks WHERE company LIKE ?
	`

	rows, err := r.db.Query(query, "%"+company+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get banks by company: %w", err)
	}
	defer rows.Close()

	var banks []*entities.Bank
	for rows.Next() {
		bank := &entities.Bank{}
		err := rows.Scan(
			&bank.ID,
			&bank.Name,
			&bank.Alias,
			&bank.Company,
			&bank.Code,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank: %w", err)
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// GetAll retrieves all banks
func (r *BankRepository) GetAll() ([]*entities.Bank, error) {
	query := `
		SELECT id, name, alias, company, code, created_at, updated_at
		FROM tm_banks ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all banks: %w", err)
	}
	defer rows.Close()

	var banks []*entities.Bank
	for rows.Next() {
		bank := &entities.Bank{}
		err := rows.Scan(
			&bank.ID,
			&bank.Name,
			&bank.Alias,
			&bank.Company,
			&bank.Code,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bank: %w", err)
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// Update updates a bank
func (r *BankRepository) Update(bank *entities.Bank) error {
	query := `
		UPDATE tm_banks 
		SET name = ?, alias = ?, company = ?, code = ?, updated_at = ?
		WHERE id = ?
	`

	bank.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		bank.Name,
		bank.Alias,
		bank.Company,
		bank.Code,
		bank.UpdatedAt,
		bank.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update bank: %w", err)
	}

	return nil
}

// Delete deletes a bank by ID
func (r *BankRepository) Delete(id uint) error {
	query := `DELETE FROM tm_banks WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bank: %w", err)
	}

	return nil
}

// DeleteAll deletes all banks
func (r *BankRepository) DeleteAll() error {
	query := `DELETE FROM tm_banks`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all banks: %w", err)
	}

	return nil
}

// Exists checks if a bank exists by ID
func (r *BankRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_banks WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if bank exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of banks
func (r *BankRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_banks`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count banks: %w", err)
	}

	return count, nil
}
