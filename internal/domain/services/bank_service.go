package services

import (
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

// CreateBank creates a new bank
func (s *BankService) CreateBank(name, alias, company, code string) (*entities.Bank, error) {
	bank := entities.NewBank(name, alias, company, code)

	if err := s.bankRepo.Create(bank); err != nil {
		return nil, err
	}

	return bank, nil
}

// GetBankByID retrieves a bank by ID
func (s *BankService) GetBankByID(id uint) (*entities.Bank, error) {
	return s.bankRepo.GetByID(id)
}

// GetBankByCode retrieves a bank by code
func (s *BankService) GetBankByCode(code string) (*entities.Bank, error) {
	return s.bankRepo.GetByCode(code)
}

// GetBankByName retrieves a bank by name
func (s *BankService) GetBankByName(name string) (*entities.Bank, error) {
	return s.bankRepo.GetByName(name)
}

// GetBankByAlias retrieves a bank by alias
func (s *BankService) GetBankByAlias(alias string) (*entities.Bank, error) {
	return s.bankRepo.GetByAlias(alias)
}

// GetBanksByCompany retrieves banks by company
func (s *BankService) GetBanksByCompany(company string) ([]*entities.Bank, error) {
	return s.bankRepo.GetByCompany(company)
}

// GetAllBanks retrieves all banks
func (s *BankService) GetAllBanks() ([]*entities.Bank, error) {
	return s.bankRepo.GetAll()
}

// UpdateBank updates a bank
func (s *BankService) UpdateBank(bank *entities.Bank) error {
	return s.bankRepo.Update(bank)
}

// DeleteBank deletes a bank by ID
func (s *BankService) DeleteBank(id uint) error {
	return s.bankRepo.Delete(id)
}

// BankExists checks if a bank exists by ID
func (s *BankService) BankExists(id uint) (bool, error) {
	return s.bankRepo.Exists(id)
}

// GetBankCount returns the total number of banks
func (s *BankService) GetBankCount() (int64, error) {
	return s.bankRepo.Count()
}
