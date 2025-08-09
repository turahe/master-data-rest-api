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

// BankSeeder handles seeding bank data
type BankSeeder struct {
	repo   *pgx.BankRepository
	logger *logger.Logger
}

// NewBankSeeder creates a new bank seeder
func NewBankSeeder(repo *pgx.BankRepository, logger *logger.Logger) *BankSeeder {
	return &BankSeeder{
		repo:   repo,
		logger: logger,
	}
}

// Name returns the seeder name
func (bs *BankSeeder) Name() string {
	return "banks"
}

// Seed seeds bank data from CSV file
func (bs *BankSeeder) Seed(ctx context.Context, dataDir string) error {
	csvFile := filepath.Join(dataDir, "tm_banks.csv")
	bs.logger.WithField("file", csvFile).Info("Starting banks seeding")

	// Open CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open banks CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read banks CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("banks CSV file must contain at least a header and one data row")
	}

	// Skip header and process records
	successCount := 0
	errorCount := 0

	fmt.Printf("🏦 Processing %d bank records...\n", len(records)-1)

	for i, record := range records[1:] { // Skip header
		if len(record) < 4 {
			bs.logger.WithField("row", i+2).Warn("Bank record has insufficient columns")
			errorCount++
			continue
		}

		name := strings.TrimSpace(record[0])
		alias := strings.TrimSpace(record[1])
		company := strings.TrimSpace(record[2])
		code := strings.TrimSpace(record[3])

		// Create bank entity
		bank := entities.NewBank(name, alias, company, code)

		// Check if bank already exists by code
		existing, err := bs.repo.GetByCode(ctx, code)
		if err == nil && existing != nil {
			// Update existing bank
			existing.SetName(name)
			existing.SetAlias(alias)
			existing.SetCompany(company)

			err = bs.repo.Update(ctx, existing)
			if err != nil {
				bs.logger.WithError(err).WithFields(map[string]interface{}{
					"row":  i + 2,
					"code": code,
					"name": name,
				}).Warn("Failed to update existing bank")
				errorCount++
				continue
			}
		} else {
			// Create new bank
			err = bs.repo.Create(ctx, bank)
			if err != nil {
				bs.logger.WithError(err).WithFields(map[string]interface{}{
					"row":  i + 2,
					"code": code,
					"name": name,
				}).Warn("Failed to create bank")
				errorCount++
				continue
			}
		}

		successCount++

		// Log progress every 50 records
		if (i+1)%50 == 0 {
			bs.logger.WithFields(map[string]interface{}{
				"processed": i + 1,
				"total":     len(records) - 1,
				"success":   successCount,
				"errors":    errorCount,
			}).Info("Banks seeding progress")
		}
	}

	bs.logger.WithFields(map[string]interface{}{
		"total_processed": len(records) - 1,
		"successful":      successCount,
		"errors":          errorCount,
	}).Info("Banks seeding completed")

	fmt.Printf("✅ Banks seeding completed: %d successful, %d errors\n", successCount, errorCount)
	return nil
}

// Clear removes all bank data
func (bs *BankSeeder) Clear(ctx context.Context) error {
	bs.logger.Info("Clearing bank data using TRUNCATE")

	if err := bs.repo.Truncate(ctx); err != nil {
		return fmt.Errorf("failed to truncate banks table: %w", err)
	}

	bs.logger.Info("Banks table truncated successfully")
	return nil
}
