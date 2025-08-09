package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// CurrencyService implements business logic for currency operations
type CurrencyService struct {
	currencyRepo repositories.CurrencyRepository
}

// NewCurrencyService creates a new CurrencyService instance
func NewCurrencyService(currencyRepo repositories.CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
	}
}

// CreateCurrency creates a new currency with validation
func (s *CurrencyService) CreateCurrency(ctx context.Context, name, code string, decimalPlaces int) (*entities.Currency, error) {
	// Validate required fields
	if name == "" {
		return nil, fmt.Errorf("currency name is required")
	}
	if code == "" {
		return nil, fmt.Errorf("currency code is required")
	}
	if len(code) > 3 {
		return nil, fmt.Errorf("currency code must be 3 characters or less")
	}
	if decimalPlaces < 0 {
		return nil, fmt.Errorf("decimal places cannot be negative")
	}

	// Check if currency with the same code already exists
	exists, err := s.currencyRepo.ExistsByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to check currency code existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("currency with code '%s' already exists", code)
	}

	// Check if currency with the same name already exists
	exists, err = s.currencyRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check currency name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("currency with name '%s' already exists", name)
	}

	currency := entities.NewCurrency(name, code, decimalPlaces)

	if err := s.currencyRepo.Create(ctx, currency); err != nil {
		return nil, fmt.Errorf("failed to create currency: %w", err)
	}

	return currency, nil
}

// GetCurrencyByID retrieves a currency by ID
func (s *CurrencyService) GetCurrencyByID(ctx context.Context, id uuid.UUID) (*entities.Currency, error) {
	return s.currencyRepo.GetByID(ctx, id)
}

// GetCurrencyByCode retrieves a currency by code
func (s *CurrencyService) GetCurrencyByCode(ctx context.Context, code string) (*entities.Currency, error) {
	return s.currencyRepo.GetByCode(ctx, code)
}

// GetCurrencyByName retrieves a currency by name
func (s *CurrencyService) GetCurrencyByName(ctx context.Context, name string) (*entities.Currency, error) {
	return s.currencyRepo.GetByName(ctx, name)
}

// GetCurrenciesBySymbol retrieves currencies by symbol
func (s *CurrencyService) GetCurrenciesBySymbol(ctx context.Context, symbol string) ([]*entities.Currency, error) {
	return s.currencyRepo.GetBySymbol(ctx, symbol)
}

// GetAllCurrencies retrieves all currencies with pagination
func (s *CurrencyService) GetAllCurrencies(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	return s.currencyRepo.GetAll(ctx, limit, offset)
}

// GetActiveCurrencies retrieves all active currencies
func (s *CurrencyService) GetActiveCurrencies(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	return s.currencyRepo.GetActive(ctx, limit, offset)
}

// GetInactiveCurrencies retrieves all inactive currencies
func (s *CurrencyService) GetInactiveCurrencies(ctx context.Context, limit, offset int) ([]*entities.Currency, error) {
	return s.currencyRepo.GetInactive(ctx, limit, offset)
}

// SearchCurrencies searches currencies by query
func (s *CurrencyService) SearchCurrencies(ctx context.Context, query string, limit, offset int) ([]*entities.Currency, error) {
	return s.currencyRepo.Search(ctx, query, limit, offset)
}

// UpdateCurrency updates an existing currency
func (s *CurrencyService) UpdateCurrency(ctx context.Context, currency *entities.Currency) error {
	if !currency.IsValid() {
		return fmt.Errorf("invalid currency data")
	}

	// Check if updating to a code that already exists (but not for the same currency)
	existingCurrency, err := s.currencyRepo.GetByCode(ctx, currency.Code)
	if err == nil && existingCurrency.ID != currency.ID {
		return fmt.Errorf("currency with code '%s' already exists", currency.Code)
	}

	// Check if updating to a name that already exists (but not for the same currency)
	existingCurrency, err = s.currencyRepo.GetByName(ctx, currency.Name)
	if err == nil && existingCurrency.ID != currency.ID {
		return fmt.Errorf("currency with name '%s' already exists", currency.Name)
	}

	return s.currencyRepo.Update(ctx, currency)
}

// ActivateCurrency activates a currency
func (s *CurrencyService) ActivateCurrency(ctx context.Context, id uuid.UUID) error {
	return s.currencyRepo.Activate(ctx, id)
}

// DeactivateCurrency deactivates a currency
func (s *CurrencyService) DeactivateCurrency(ctx context.Context, id uuid.UUID) error {
	return s.currencyRepo.Deactivate(ctx, id)
}

// DeleteCurrency deletes a currency by ID
func (s *CurrencyService) DeleteCurrency(ctx context.Context, id uuid.UUID) error {
	return s.currencyRepo.Delete(ctx, id)
}

// CountCurrencies returns the total number of currencies
func (s *CurrencyService) CountCurrencies(ctx context.Context) (int64, error) {
	return s.currencyRepo.Count(ctx)
}

// ValidateCurrency validates a currency entity
func (s *CurrencyService) ValidateCurrency(currency *entities.Currency) error {
	if currency == nil {
		return fmt.Errorf("currency cannot be nil")
	}

	if !currency.IsValid() {
		return fmt.Errorf("currency validation failed: name and code are required, code must be 3 characters or less, decimal places cannot be negative")
	}

	return nil
}
