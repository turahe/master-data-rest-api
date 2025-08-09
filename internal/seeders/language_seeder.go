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

// LanguageSeeder handles seeding language data
type LanguageSeeder struct {
	repo   *pgx.LanguageRepository
	logger *logger.Logger
}

// NewLanguageSeeder creates a new language seeder
func NewLanguageSeeder(repo *pgx.LanguageRepository, logger *logger.Logger) *LanguageSeeder {
	return &LanguageSeeder{
		repo:   repo,
		logger: logger,
	}
}

// Name returns the seeder name
func (ls *LanguageSeeder) Name() string {
	return "languages"
}

// Seed seeds language data from CSV file
func (ls *LanguageSeeder) Seed(ctx context.Context, dataDir string) error {
	csvFile := filepath.Join(dataDir, "tm_languages.csv")
	ls.logger.WithField("file", csvFile).Info("Reading languages from CSV file")

	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open languages CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found in CSV file")
	}

	// Skip header row
	if len(records) > 1 {
		records = records[1:]
	}

	ls.logger.WithField("count", len(records)).Info("Processing language records")

	successCount := 0
	errorCount := 0

	for i, record := range records {
		if len(record) < 3 {
			ls.logger.WithField("row", i+2).Warn("Skipping row with insufficient columns")
			errorCount++
			continue
		}

		code := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		_ = strings.TrimSpace(record[2]) // native name - not used in current entity

		if code == "" || name == "" {
			ls.logger.WithField("row", i+2).Warn("Skipping row with empty code or name")
			errorCount++
			continue
		}

		// Create language entity
		language := entities.NewLanguage(name, code)

		// Check if language already exists
		existing, getErr := ls.repo.GetByCode(ctx, code)
		if getErr != nil {
			// Language doesn't exist, create new one
			err = ls.repo.Create(ctx, language)
			if err != nil {
				ls.logger.WithError(err).WithFields(map[string]interface{}{
					"row":  i + 2,
					"code": code,
					"name": name,
				}).Warn("Failed to create language")
				errorCount++
				continue
			}
		} else {
			// Language exists, update it
			existing.SetName(name)
			updateErr := ls.repo.Update(ctx, existing)
			if updateErr != nil {
				ls.logger.WithError(updateErr).WithFields(map[string]interface{}{
					"row":  i + 2,
					"code": code,
					"name": name,
				}).Warn("Failed to update existing language")
				errorCount++
				continue
			}
		}

		successCount++
		if successCount%50 == 0 {
			ls.logger.WithField("processed", successCount).Info("Language seeding progress")
		}
	}

	ls.logger.WithFields(map[string]interface{}{
		"success": successCount,
		"errors":  errorCount,
		"total":   len(records),
	}).Info("Language seeding completed")

	fmt.Printf("📊 Processed %d language records: %d success, %d errors\n", len(records), successCount, errorCount)

	return nil
}

// Clear removes all language data
func (ls *LanguageSeeder) Clear(ctx context.Context) error {
	// Get all languages first to count them
	languages, err := ls.repo.GetAll(ctx, 10000, 0) // Get up to 10k records with 0 offset
	if err != nil {
		return fmt.Errorf("failed to get languages for clearing: %w", err)
	}

	ls.logger.WithField("count", len(languages)).Info("Clearing existing languages")

	for _, language := range languages {
		if err := ls.repo.Delete(ctx, language.ID); err != nil {
			ls.logger.WithError(err).WithField("id", language.ID).Warn("Failed to delete language")
			// Continue with others even if one fails
		}
	}

	ls.logger.WithField("count", len(languages)).Info("Languages cleared successfully")
	return nil
}
