package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// CurrencyRepository implements the CurrencyRepository interface for MySQL
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository creates a new currency repository
func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

// Create creates a new currency
func (r *CurrencyRepository) Create(currency *entities.Currency) error {
	query := `
		INSERT INTO tm_currencies (priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	currency.CreatedAt = now
	currency.UpdatedAt = now

	result, err := r.db.Exec(query,
		currency.Priority,
		currency.ISOCode,
		currency.Name,
		currency.Symbol,
		currency.DisambiguateSymbol,
		currency.AlternateSymbols,
		currency.Subunit,
		currency.SubunitToUnit,
		currency.SymbolFirst,
		currency.HTMLEntity,
		currency.DecimalMark,
		currency.ThousandsSeparator,
		currency.ISONumeric,
		currency.SmallestDenomination,
		currency.CreatedAt,
		currency.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create currency: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	currency.ID = uint(id)
	return nil
}

// GetByID retrieves a currency by ID
func (r *CurrencyRepository) GetByID(id uint) (*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE id = ?
	`

	currency := &entities.Currency{}
	err := r.db.QueryRow(query, id).Scan(
		&currency.ID,
		&currency.Priority,
		&currency.ISOCode,
		&currency.Name,
		&currency.Symbol,
		&currency.DisambiguateSymbol,
		&currency.AlternateSymbols,
		&currency.Subunit,
		&currency.SubunitToUnit,
		&currency.SymbolFirst,
		&currency.HTMLEntity,
		&currency.DecimalMark,
		&currency.ThousandsSeparator,
		&currency.ISONumeric,
		&currency.SmallestDenomination,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency by id: %w", err)
	}

	return currency, nil
}

// GetByISOCode retrieves a currency by ISO code
func (r *CurrencyRepository) GetByISOCode(isoCode string) (*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE iso_code = ?
	`

	currency := &entities.Currency{}
	err := r.db.QueryRow(query, isoCode).Scan(
		&currency.ID,
		&currency.Priority,
		&currency.ISOCode,
		&currency.Name,
		&currency.Symbol,
		&currency.DisambiguateSymbol,
		&currency.AlternateSymbols,
		&currency.Subunit,
		&currency.SubunitToUnit,
		&currency.SymbolFirst,
		&currency.HTMLEntity,
		&currency.DecimalMark,
		&currency.ThousandsSeparator,
		&currency.ISONumeric,
		&currency.SmallestDenomination,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency by ISO code: %w", err)
	}

	return currency, nil
}

// GetByName retrieves currencies by name
func (r *CurrencyRepository) GetByName(name string) ([]*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE name LIKE ?
	`

	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies by name: %w", err)
	}
	defer rows.Close()

	var currencies []*entities.Currency
	for rows.Next() {
		currency := &entities.Currency{}
		err := rows.Scan(
			&currency.ID,
			&currency.Priority,
			&currency.ISOCode,
			&currency.Name,
			&currency.Symbol,
			&currency.DisambiguateSymbol,
			&currency.AlternateSymbols,
			&currency.Subunit,
			&currency.SubunitToUnit,
			&currency.SymbolFirst,
			&currency.HTMLEntity,
			&currency.DecimalMark,
			&currency.ThousandsSeparator,
			&currency.ISONumeric,
			&currency.SmallestDenomination,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// GetBySymbol retrieves currencies by symbol
func (r *CurrencyRepository) GetBySymbol(symbol string) ([]*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE symbol = ?
	`

	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies by symbol: %w", err)
	}
	defer rows.Close()

	var currencies []*entities.Currency
	for rows.Next() {
		currency := &entities.Currency{}
		err := rows.Scan(
			&currency.ID,
			&currency.Priority,
			&currency.ISOCode,
			&currency.Name,
			&currency.Symbol,
			&currency.DisambiguateSymbol,
			&currency.AlternateSymbols,
			&currency.Subunit,
			&currency.SubunitToUnit,
			&currency.SymbolFirst,
			&currency.HTMLEntity,
			&currency.DecimalMark,
			&currency.ThousandsSeparator,
			&currency.ISONumeric,
			&currency.SmallestDenomination,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// GetByISONumeric retrieves a currency by ISO numeric code
func (r *CurrencyRepository) GetByISONumeric(isoNumeric string) (*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE iso_numeric = ?
	`

	currency := &entities.Currency{}
	err := r.db.QueryRow(query, isoNumeric).Scan(
		&currency.ID,
		&currency.Priority,
		&currency.ISOCode,
		&currency.Name,
		&currency.Symbol,
		&currency.DisambiguateSymbol,
		&currency.AlternateSymbols,
		&currency.Subunit,
		&currency.SubunitToUnit,
		&currency.SymbolFirst,
		&currency.HTMLEntity,
		&currency.DecimalMark,
		&currency.ThousandsSeparator,
		&currency.ISONumeric,
		&currency.SmallestDenomination,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency by ISO numeric: %w", err)
	}

	return currency, nil
}

// GetAll retrieves all currencies
func (r *CurrencyRepository) GetAll() ([]*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies ORDER BY priority, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all currencies: %w", err)
	}
	defer rows.Close()

	var currencies []*entities.Currency
	for rows.Next() {
		currency := &entities.Currency{}
		err := rows.Scan(
			&currency.ID,
			&currency.Priority,
			&currency.ISOCode,
			&currency.Name,
			&currency.Symbol,
			&currency.DisambiguateSymbol,
			&currency.AlternateSymbols,
			&currency.Subunit,
			&currency.SubunitToUnit,
			&currency.SymbolFirst,
			&currency.HTMLEntity,
			&currency.DecimalMark,
			&currency.ThousandsSeparator,
			&currency.ISONumeric,
			&currency.SmallestDenomination,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// GetByPriority retrieves currencies by priority
func (r *CurrencyRepository) GetByPriority(priority int) ([]*entities.Currency, error) {
	query := `
		SELECT id, priority, iso_code, name, symbol, disambiguate_symbol, alternate_symbols, subunit, subunit_to_unit, symbol_first, html_entity, decimal_mark, thousands_separator, iso_numeric, smallest_denomination, created_at, updated_at
		FROM tm_currencies WHERE priority = ? ORDER BY name
	`

	rows, err := r.db.Query(query, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies by priority: %w", err)
	}
	defer rows.Close()

	var currencies []*entities.Currency
	for rows.Next() {
		currency := &entities.Currency{}
		err := rows.Scan(
			&currency.ID,
			&currency.Priority,
			&currency.ISOCode,
			&currency.Name,
			&currency.Symbol,
			&currency.DisambiguateSymbol,
			&currency.AlternateSymbols,
			&currency.Subunit,
			&currency.SubunitToUnit,
			&currency.SymbolFirst,
			&currency.HTMLEntity,
			&currency.DecimalMark,
			&currency.ThousandsSeparator,
			&currency.ISONumeric,
			&currency.SmallestDenomination,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan currency: %w", err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// Update updates a currency
func (r *CurrencyRepository) Update(currency *entities.Currency) error {
	query := `
		UPDATE tm_currencies 
		SET priority = ?, iso_code = ?, name = ?, symbol = ?, disambiguate_symbol = ?, alternate_symbols = ?, subunit = ?, subunit_to_unit = ?, symbol_first = ?, html_entity = ?, decimal_mark = ?, thousands_separator = ?, iso_numeric = ?, smallest_denomination = ?, updated_at = ?
		WHERE id = ?
	`

	currency.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		currency.Priority,
		currency.ISOCode,
		currency.Name,
		currency.Symbol,
		currency.DisambiguateSymbol,
		currency.AlternateSymbols,
		currency.Subunit,
		currency.SubunitToUnit,
		currency.SymbolFirst,
		currency.HTMLEntity,
		currency.DecimalMark,
		currency.ThousandsSeparator,
		currency.ISONumeric,
		currency.SmallestDenomination,
		currency.UpdatedAt,
		currency.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update currency: %w", err)
	}

	return nil
}

// Delete deletes a currency by ID
func (r *CurrencyRepository) Delete(id uint) error {
	query := `DELETE FROM tm_currencies WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete currency: %w", err)
	}

	return nil
}

// DeleteAll deletes all currencies
func (r *CurrencyRepository) DeleteAll() error {
	query := `DELETE FROM tm_currencies`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete all currencies: %w", err)
	}

	return nil
}

// Exists checks if a currency exists by ID
func (r *CurrencyRepository) Exists(id uint) (bool, error) {
	query := `SELECT COUNT(*) FROM tm_currencies WHERE id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if currency exists: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of currencies
func (r *CurrencyRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM tm_currencies`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count currencies: %w", err)
	}

	return count, nil
}
