package gorm

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
	"gorm.io/gorm"
)

// BankRepository implements the BankRepository interface using GORM
type BankRepository struct {
	db *gorm.DB
}

// NewBankRepository creates a new GORM bank repository
func NewBankRepository(db *gorm.DB) repositories.BankRepository {
	return &BankRepository{
		db: db,
	}
}

// Create creates a new bank in the database
func (r *BankRepository) Create(bank *entities.Bank) error {
	result := r.db.Create(bank)
	if result.Error != nil {
		return fmt.Errorf("failed to create bank: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a bank by ID from the database
func (r *BankRepository) GetByID(id uuid.UUID) (*entities.Bank, error) {
	var bank entities.Bank
	result := r.db.Where("id = ?", id).First(&bank)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, fmt.Errorf("failed to get bank by ID: %w", result.Error)
	}
	return &bank, nil
}

// GetByCode retrieves a bank by code from the database
func (r *BankRepository) GetByCode(code string) (*entities.Bank, error) {
	var bank entities.Bank
	result := r.db.Where("code = ?", code).First(&bank)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bank not found")
		}
		return nil, fmt.Errorf("failed to get bank by code: %w", result.Error)
	}
	return &bank, nil
}

// GetByName retrieves banks by name from the database
func (r *BankRepository) GetByName(name string) ([]*entities.Bank, error) {
	var banks []entities.Bank
	result := r.db.Where("name LIKE ?", "%"+name+"%").Order("name ASC").Find(&banks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get banks by name: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	bankPtrs := make([]*entities.Bank, len(banks))
	for i := range banks {
		bankPtrs[i] = &banks[i]
	}

	return bankPtrs, nil
}

// GetByAlias retrieves banks by alias from the database
func (r *BankRepository) GetByAlias(alias string) ([]*entities.Bank, error) {
	var banks []entities.Bank
	result := r.db.Where("alias LIKE ?", "%"+alias+"%").Order("name ASC").Find(&banks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get banks by alias: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	bankPtrs := make([]*entities.Bank, len(banks))
	for i := range banks {
		bankPtrs[i] = &banks[i]
	}

	return bankPtrs, nil
}

// GetByCompany retrieves banks by company from the database
func (r *BankRepository) GetByCompany(company string) ([]*entities.Bank, error) {
	var banks []entities.Bank
	result := r.db.Where("company LIKE ?", "%"+company+"%").Order("name ASC").Find(&banks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get banks by company: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	bankPtrs := make([]*entities.Bank, len(banks))
	for i := range banks {
		bankPtrs[i] = &banks[i]
	}

	return bankPtrs, nil
}

// GetAll retrieves all banks from the database
func (r *BankRepository) GetAll() ([]*entities.Bank, error) {
	var banks []entities.Bank
	result := r.db.Order("name ASC").Find(&banks)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all banks: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	bankPtrs := make([]*entities.Bank, len(banks))
	for i := range banks {
		bankPtrs[i] = &banks[i]
	}

	return bankPtrs, nil
}

// Update updates an existing bank in the database
func (r *BankRepository) Update(bank *entities.Bank) error {
	result := r.db.Save(bank)
	if result.Error != nil {
		return fmt.Errorf("failed to update bank: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("bank not found")
	}

	return nil
}

// Delete deletes a bank by ID from the database
func (r *BankRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&entities.Bank{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete bank: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("bank not found")
	}

	return nil
}

// DeleteAll deletes all banks from the database
func (r *BankRepository) DeleteAll() error {
	result := r.db.Delete(&entities.Bank{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete all banks: %w", result.Error)
	}
	return nil
}

// Exists checks if a bank exists by ID in the database
func (r *BankRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	result := r.db.Model(&entities.Bank{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check bank existence: %w", result.Error)
	}

	return count > 0, nil
}

// Count returns the total number of banks in the database
func (r *BankRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&entities.Bank{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count banks: %w", result.Error)
	}

	return count, nil
}
