package seeders

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

// CurrencySeeder handles seeding currency data
type CurrencySeeder struct {
	repo   *pgx.CurrencyRepository
	logger *logger.Logger
}

// NewCurrencySeeder creates a new currency seeder
func NewCurrencySeeder(repo *pgx.CurrencyRepository, logger *logger.Logger) *CurrencySeeder {
	return &CurrencySeeder{
		repo:   repo,
		logger: logger,
	}
}

// Name returns the seeder name
func (cs *CurrencySeeder) Name() string {
	return "currencies"
}

// Seed seeds currency data from CSV file
func (cs *CurrencySeeder) Seed(ctx context.Context, dataDir string) error {
	csvFile := filepath.Join(dataDir, "tm_currencies.csv")
	cs.logger.WithField("file", csvFile).Info("Starting currencies seeding")

	// Open CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open currencies CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read currencies CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("currencies CSV file must contain at least a header and one data row")
	}

	// Skip header and process records
	successCount := 0
	errorCount := 0

	fmt.Printf("💰 Processing %d currency records...\n", len(records)-1)

	for i, record := range records[1:] { // Skip header
		if len(record) < 3 {
			cs.logger.WithField("row", i+2).Warn("Currency record has insufficient columns")
			errorCount++
			continue
		}

		// Parse CSV fields: priority,iso_code,name,symbol,disambiguate_symbol,alternate_symbols,subunit,subunit_to_unit,symbol_first,html_entity,decimal_mark,thousands_separator,iso_numeric,smallest_denomination
		isoCode := strings.TrimSpace(record[1])
		name := strings.TrimSpace(record[2])
		_ = strings.TrimSpace(record[3]) // symbol - not used in constructor

		// Create currency entity (using iso_code as the code field, default 2 decimal places)
		currency := entities.NewCurrency(name, isoCode, 2)

		// Check if currency already exists by code
		existing, err := cs.repo.GetByCode(ctx, isoCode)
		if err == nil && existing != nil {
			// Update existing currency
			existing.SetName(name)
			// Note: Symbol and decimal places cannot be updated via setters

			err = cs.repo.Update(ctx, existing)
			if err != nil {
				cs.logger.WithError(err).WithFields(map[string]interface{}{
					"row":      i + 2,
					"iso_code": isoCode,
					"name":     name,
				}).Warn("Failed to update existing currency")
				errorCount++
				continue
			}
		} else {
			// Create new currency
			err = cs.repo.Create(ctx, currency)
			if err != nil {
				cs.logger.WithError(err).WithFields(map[string]interface{}{
					"row":      i + 2,
					"iso_code": isoCode,
					"name":     name,
				}).Warn("Failed to create currency")
				errorCount++
				continue
			}
		}

		successCount++

		// Log progress every 50 records
		if (i+1)%50 == 0 {
			cs.logger.WithFields(map[string]interface{}{
				"processed": i + 1,
				"total":     len(records) - 1,
				"success":   successCount,
				"errors":    errorCount,
			}).Info("Currencies seeding progress")
		}
	}

	cs.logger.WithFields(map[string]interface{}{
		"total_processed": len(records) - 1,
		"successful":      successCount,
		"errors":          errorCount,
	}).Info("Currencies seeding completed")

	fmt.Printf("✅ Currencies seeding completed: %d successful, %d errors\n", successCount, errorCount)
	return nil
}

// Clear removes all currency data
func (cs *CurrencySeeder) Clear(ctx context.Context) error {
	// Get all currencies first to count them
	currencies, err := cs.repo.GetAll(ctx, 10000, 0) // Get up to 10k records with 0 offset
	if err != nil {
		return fmt.Errorf("failed to get currencies for clearing: %w", err)
	}

	cs.logger.WithField("count", len(currencies)).Info("Clearing existing currencies")

	for _, currency := range currencies {
		if err := cs.repo.Delete(ctx, currency.ID); err != nil {
			cs.logger.WithError(err).WithField("id", currency.ID).Warn("Failed to delete currency")
			// Continue with others even if one fails
		}
	}

	cs.logger.WithField("count", len(currencies)).Info("Currencies cleared successfully")
	return nil
}
