package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// BankService implements business logic for bank operations
type BankService struct {
	bankRepo repositories.BankRepository
}

// NewBankService creates a new BankService instance
func NewBankService(bankRepo repositories.BankRepository) *BankService {
	return &BankService{
		bankRepo: bankRepo,
	}
}

// CreateBank creates a new bank with validation
func (s *BankService) CreateBank(ctx context.Context, name, alias, company, code string) (*entities.Bank, error) {
	// Validate required fields
	if name == "" {
		return nil, fmt.Errorf("bank name is required")
	}
	if code == "" {
		return nil, fmt.Errorf("bank code is required")
	}

	// Check if bank with the same code already exists
	exists, err := s.bankRepo.ExistsByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to check bank code existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("bank with code '%s' already exists", code)
	}

	// Check if bank with the same name already exists
	exists, err = s.bankRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check bank name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("bank with name '%s' already exists", name)
	}

	bank := entities.NewBank(name, alias, company, code)

	if err := s.bankRepo.Create(ctx, bank); err != nil {
		return nil, fmt.Errorf("failed to create bank: %w", err)
	}

	return bank, nil
}

// GetBankByID retrieves a bank by ID
func (s *BankService) GetBankByID(ctx context.Context, id uuid.UUID) (*entities.Bank, error) {
	return s.bankRepo.GetByID(ctx, id)
}

// GetBankByCode retrieves a bank by code
func (s *BankService) GetBankByCode(ctx context.Context, code string) (*entities.Bank, error) {
	return s.bankRepo.GetByCode(ctx, code)
}

// GetBankByName retrieves a bank by name
func (s *BankService) GetBankByName(ctx context.Context, name string) (*entities.Bank, error) {
	return s.bankRepo.GetByName(ctx, name)
}

// GetBankByAlias retrieves a bank by alias
func (s *BankService) GetBankByAlias(ctx context.Context, alias string) (*entities.Bank, error) {
	return s.bankRepo.GetByAlias(ctx, alias)
}

// GetAllBanks retrieves all banks with pagination
func (s *BankService) GetAllBanks(ctx context.Context, limit, offset int) ([]*entities.Bank, error) {
	return s.bankRepo.GetAll(ctx, limit, offset)
}

// SearchBanks searches banks by query
func (s *BankService) SearchBanks(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error) {
	return s.bankRepo.Search(ctx, query, limit, offset)
}

// GetBanksByCompany retrieves banks by company
func (s *BankService) GetBanksByCompany(ctx context.Context, company string, limit, offset int) ([]*entities.Bank, error) {
	return s.bankRepo.GetByCompany(ctx, company, limit, offset)
}

// UpdateBank updates an existing bank
func (s *BankService) UpdateBank(ctx context.Context, bank *entities.Bank) error {
	if !bank.IsValid() {
		return fmt.Errorf("invalid bank data")
	}

	// Check if updating to a code that already exists (but not for the same bank)
	existingBank, err := s.bankRepo.GetByCode(ctx, bank.Code)
	if err == nil && existingBank.ID != bank.ID {
		return fmt.Errorf("bank with code '%s' already exists", bank.Code)
	}

	// Check if updating to a name that already exists (but not for the same bank)
	existingBank, err = s.bankRepo.GetByName(ctx, bank.Name)
	if err == nil && existingBank.ID != bank.ID {
		return fmt.Errorf("bank with name '%s' already exists", bank.Name)
	}

	return s.bankRepo.Update(ctx, bank)
}

// DeleteBank deletes a bank by ID
func (s *BankService) DeleteBank(ctx context.Context, id uuid.UUID) error {
	return s.bankRepo.Delete(ctx, id)
}

// CountBanks returns the total number of banks
func (s *BankService) CountBanks(ctx context.Context) (int64, error) {
	return s.bankRepo.Count(ctx)
}

// ValidateBank validates a bank entity
func (s *BankService) ValidateBank(bank *entities.Bank) error {
	if bank == nil {
		return fmt.Errorf("bank cannot be nil")
	}

	if !bank.IsValid() {
		return fmt.Errorf("bank validation failed: name and code are required")
	}

	return nil
}
