package seeders

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

// GeodirectorySeeder handles seeding geodirectory data
type GeodirectorySeeder struct {
	repo   *pgx.GeodirectoryRepository
	logger *logger.Logger
}

// NewGeodirectorySeeder creates a new geodirectory seeder
func NewGeodirectorySeeder(repo *pgx.GeodirectoryRepository, logger *logger.Logger) *GeodirectorySeeder {
	return &GeodirectorySeeder{
		repo:   repo,
		logger: logger,
	}
}

// Name returns the seeder name
func (gs *GeodirectorySeeder) Name() string {
	return "geodirectories"
}

// Seed seeds geodirectory data from CSV and JSON files
func (gs *GeodirectorySeeder) Seed(ctx context.Context, dataDir string) error {
	geoDir := filepath.Join(dataDir, "geodirectories")
	gs.logger.WithField("directory", geoDir).Info("Starting geodirectories seeding")

	// Seed in hierarchical order: countries -> provinces -> cities -> districts -> villages
	if err := gs.seedCountries(ctx, geoDir); err != nil {
		return fmt.Errorf("failed to seed countries: %w", err)
	}

	if err := gs.seedProvinces(ctx, geoDir); err != nil {
		return fmt.Errorf("failed to seed provinces: %w", err)
	}

	if err := gs.seedCities(ctx, geoDir); err != nil {
		return fmt.Errorf("failed to seed cities: %w", err)
	}

	if err := gs.seedDistricts(ctx, geoDir); err != nil {
		return fmt.Errorf("failed to seed districts: %w", err)
	}

	if err := gs.seedVillages(ctx, geoDir); err != nil {
		return fmt.Errorf("failed to seed villages: %w", err)
	}

	gs.logger.Info("Geodirectories seeding completed successfully")
	return nil
}

// seedCountries seeds country data from countries.csv
func (gs *GeodirectorySeeder) seedCountries(ctx context.Context, geoDir string) error {
	countryFile := filepath.Join(geoDir, "countries.csv")
	gs.logger.WithField("file", countryFile).Info("Seeding countries")

	file, err := os.Open(countryFile)
	if err != nil {
		return fmt.Errorf("failed to open country file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		gs.logger.Warn("No country records found in CSV file")
		return nil
	}

	// Skip header row
	records = records[1:]

	successCount := 0
	errorCount := 0
	orderingCounter := 1

	for _, record := range records {
		if len(record) < 4 {
			gs.logger.WithField("record", record).Warn("Skipping incomplete country record")
			errorCount++
			continue
		}

		code := strings.TrimSpace(record[0])
		if code == "" {
			continue // Skip empty lines
		}

		latStr := strings.TrimSpace(record[1])
		lonStr := strings.TrimSpace(record[2])
		name := strings.TrimSpace(record[3])

		// Parse coordinates if provided
		var latitude, longitude *float64
		if latStr != "" {
			if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
				latitude = &lat
			}
		}
		if lonStr != "" {
			if lon, err := strconv.ParseFloat(lonStr, 64); err == nil {
				longitude = &lon
			}
		}

		// Check if country already exists
		existing, getErr := gs.repo.GetByCode(ctx, code)
		if getErr != nil && !isNotFoundError(getErr) {
			gs.logger.WithError(getErr).WithField("code", code).Warn("Failed to check existing country")
			errorCount++
			continue
		}

		if existing != nil {
			// Update existing country
			existing.Name = name
			if latitude != nil && longitude != nil {
				existing.SetCoordinates(strconv.FormatFloat(*latitude, 'f', -1, 64), strconv.FormatFloat(*longitude, 'f', -1, 64))
			}
			existing.SetDepth(2) // Countries are depth 2
			existing.SetOrderingID(orderingCounter)
			if updateErr := gs.repo.Update(ctx, existing); updateErr != nil {
				gs.logger.WithError(updateErr).WithField("code", code).Error("Failed to update country")
				errorCount++
				continue
			}
			gs.logger.WithField("code", code).Debug("Updated existing country")
		} else {
			// Create new country
			country := entities.NewGeodirectory(name, entities.GeoTypeCountry)
			country.SetCode(code)
			if latitude != nil && longitude != nil {
				country.SetCoordinates(strconv.FormatFloat(*latitude, 'f', -1, 64), strconv.FormatFloat(*longitude, 'f', -1, 64))
			}
			country.SetDepth(2) // Countries are depth 2
			country.SetOrderingID(orderingCounter)

			if createErr := gs.repo.Create(ctx, country); createErr != nil {
				gs.logger.WithError(createErr).WithField("code", code).Error("Failed to create country")
				errorCount++
				continue
			}
			gs.logger.WithField("code", code).Debug("Created new country")
		}
		successCount++
		orderingCounter++

		// Log progress every 50 records
		if (successCount+errorCount)%50 == 0 {
			gs.logger.WithFields(map[string]interface{}{
				"processed": successCount + errorCount,
				"success":   successCount,
				"errors":    errorCount,
			}).Info("Country seeding progress")
		}
	}

	gs.logger.WithFields(map[string]interface{}{
		"total_processed": successCount + errorCount,
		"successful":      successCount,
		"errors":          errorCount,
	}).Info("Country seeding completed")

	return nil
}

// seedProvinces seeds province data from provinces/provinsi.json as children of Indonesia
func (gs *GeodirectorySeeder) seedProvinces(ctx context.Context, geoDir string) error {
	provinceFile := filepath.Join(geoDir, "provinces", "provinsi.json")
	gs.logger.WithField("file", provinceFile).Info("Seeding provinces")

	// First, find Indonesia country to set as parent
	indonesia, err := gs.repo.GetByCode(ctx, "ID")
	if err != nil {
		return fmt.Errorf("failed to get Indonesia country (code: ID): %w", err)
	}

	file, err := os.Open(provinceFile)
	if err != nil {
		return fmt.Errorf("failed to open province file: %w", err)
	}
	defer file.Close()

	var provinces map[string]string
	if err := json.NewDecoder(file).Decode(&provinces); err != nil {
		return fmt.Errorf("failed to decode province JSON: %w", err)
	}

	successCount := 0
	errorCount := 0
	orderingCounter := 1

	for code, name := range provinces {
		// Check if province already exists
		existing, getErr := gs.repo.GetByCode(ctx, code)
		if getErr != nil && !isNotFoundError(getErr) {
			gs.logger.WithError(getErr).WithField("code", code).Warn("Failed to check existing province")
			errorCount++
			continue
		}

		if existing != nil {
			// Update existing province
			existing.Name = name
			existing.SetParent(indonesia.ID)
			existing.SetDepth(existing.GetDepthForType())
			existing.SetOrderingID(orderingCounter)
			if updateErr := gs.repo.Update(ctx, existing); updateErr != nil {
				gs.logger.WithError(updateErr).WithField("code", code).Error("Failed to update province")
				errorCount++
				continue
			}
			gs.logger.WithField("code", code).Debug("Updated existing province")
		} else {
			// Create new province
			province := entities.NewGeodirectory(name, entities.GeoTypeProvince)
			province.SetCode(code)
			province.SetParent(indonesia.ID)
			province.SetDepth(province.GetDepthForType())
			province.SetOrderingID(orderingCounter)

			if createErr := gs.repo.Create(ctx, province); createErr != nil {
				gs.logger.WithError(createErr).WithField("code", code).Error("Failed to create province")
				errorCount++
				continue
			}
			gs.logger.WithField("code", code).Debug("Created new province")
		}
		successCount++
		orderingCounter++

		// Log progress every 10 records
		if (successCount+errorCount)%10 == 0 {
			gs.logger.WithFields(map[string]interface{}{
				"processed": successCount + errorCount,
				"success":   successCount,
				"errors":    errorCount,
			}).Info("Province seeding progress")
		}
	}

	gs.logger.WithFields(map[string]interface{}{
		"total_processed": successCount + errorCount,
		"successful":      successCount,
		"errors":          errorCount,
	}).Info("Province seeding completed")

	return nil
}

// Clear removes all geodirectory data
func (gs *GeodirectorySeeder) Clear(ctx context.Context) error {
	gs.logger.Info("Clearing geodirectory data using TRUNCATE")

	if err := gs.repo.Truncate(ctx); err != nil {
		return fmt.Errorf("failed to truncate geodirectories table: %w", err)
	}

	gs.logger.Info("Geodirectories table truncated successfully")
	return nil
}
