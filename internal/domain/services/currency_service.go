package services

import (
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

// CreateCurrency creates a new currency
func (s *CurrencyService) CreateCurrency(isoCode string, subunitToUnit int, symbolFirst bool, smallestDenomination int) (*entities.Currency, error) {
	currency := entities.NewCurrency(isoCode, subunitToUnit, symbolFirst, smallestDenomination)

	if err := s.currencyRepo.Create(currency); err != nil {
		return nil, err
	}

	return currency, nil
}

// GetCurrencyByID retrieves a currency by ID
func (s *CurrencyService) GetCurrencyByID(id uint) (*entities.Currency, error) {
	return s.currencyRepo.GetByID(id)
}

// GetCurrencyByISOCode retrieves a currency by ISO code
func (s *CurrencyService) GetCurrencyByISOCode(isoCode string) (*entities.Currency, error) {
	return s.currencyRepo.GetByISOCode(isoCode)
}

// GetCurrencyByName retrieves a currency by name
func (s *CurrencyService) GetCurrencyByName(name string) (*entities.Currency, error) {
	return s.currencyRepo.GetByName(name)
}

// GetCurrenciesBySymbol retrieves currencies by symbol
func (s *CurrencyService) GetCurrenciesBySymbol(symbol string) ([]*entities.Currency, error) {
	return s.currencyRepo.GetBySymbol(symbol)
}

// GetCurrencyByISONumeric retrieves a currency by ISO numeric code
func (s *CurrencyService) GetCurrencyByISONumeric(isoNumeric string) (*entities.Currency, error) {
	return s.currencyRepo.GetByISONumeric(isoNumeric)
}

// GetAllCurrencies retrieves all currencies
func (s *CurrencyService) GetAllCurrencies() ([]*entities.Currency, error) {
	return s.currencyRepo.GetAll()
}

// GetCurrenciesByPriority retrieves currencies by priority
func (s *CurrencyService) GetCurrenciesByPriority(priority int) ([]*entities.Currency, error) {
	return s.currencyRepo.GetByPriority(priority)
}

// UpdateCurrency updates a currency
func (s *CurrencyService) UpdateCurrency(currency *entities.Currency) error {
	return s.currencyRepo.Update(currency)
}

// DeleteCurrency deletes a currency by ID
func (s *CurrencyService) DeleteCurrency(id uint) error {
	return s.currencyRepo.Delete(id)
}

// CurrencyExists checks if a currency exists by ID
func (s *CurrencyService) CurrencyExists(id uint) (bool, error) {
	return s.currencyRepo.Exists(id)
}

// GetCurrencyCount returns the total number of currencies
func (s *CurrencyService) GetCurrencyCount() (int64, error) {
	return s.currencyRepo.Count()
}
