package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// Seed seeds geodirectory data from JSON files
func (gs *GeodirectorySeeder) Seed(ctx context.Context, dataDir string) error {
	geoDir := filepath.Join(dataDir, "geodirectories")
	gs.logger.WithField("directory", geoDir).Info("Starting geodirectories seeding")

	// Seed in hierarchical order: provinces -> cities -> districts -> villages
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

// seedProvinces seeds province data from provinces/provinsi.json
func (gs *GeodirectorySeeder) seedProvinces(ctx context.Context, geoDir string) error {
	provinceFile := filepath.Join(geoDir, "provinces", "provinsi.json")
	gs.logger.WithField("file", provinceFile).Info("Seeding provinces")

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
	gs.logger.Info("Clearing geodirectory data")

	// Clear in reverse hierarchical order to maintain referential integrity
	// Villages -> Districts -> Cities -> Provinces

	if err := gs.clearByType(ctx, entities.GeoTypeVillage); err != nil {
		return fmt.Errorf("failed to clear villages: %w", err)
	}

	if err := gs.clearByType(ctx, entities.GeoTypeDistrict); err != nil {
		return fmt.Errorf("failed to clear districts: %w", err)
	}

	if err := gs.clearByType(ctx, entities.GeoTypeCity); err != nil {
		return fmt.Errorf("failed to clear cities: %w", err)
	}

	if err := gs.clearByType(ctx, entities.GeoTypeProvince); err != nil {
		return fmt.Errorf("failed to clear provinces: %w", err)
	}

	gs.logger.Info("Geodirectory data cleared successfully")
	return nil
}

// clearByType clears all geodirectories of a specific type
func (gs *GeodirectorySeeder) clearByType(ctx context.Context, geoType entities.GeoType) error {
	gs.logger.WithField("type", geoType).Info("Clearing geodirectories by type")

	// Get all geodirectories of this type
	geodirectories, err := gs.repo.GetByType(ctx, geoType, 10000, 0)
	if err != nil {
		if isNotFoundError(err) {
			gs.logger.WithField("type", geoType).Info("No geodirectories found to clear")
			return nil
		}
		return fmt.Errorf("failed to get geodirectories by type %s: %w", geoType, err)
	}

	if len(geodirectories) == 0 {
		gs.logger.WithField("type", geoType).Info("No geodirectories found to clear")
		return nil
	}

	deletedCount := 0
	for _, geo := range geodirectories {
		if err := gs.repo.Delete(ctx, geo.ID); err != nil {
			gs.logger.WithError(err).WithFields(map[string]interface{}{
				"id":   geo.ID,
				"type": geoType,
				"name": geo.Name,
			}).Warn("Failed to delete geodirectory during clear")
			continue
		}
		deletedCount++
	}

	gs.logger.WithFields(map[string]interface{}{
		"type":  geoType,
		"count": deletedCount,
	}).Info("Geodirectories cleared by type")

	return nil
}
