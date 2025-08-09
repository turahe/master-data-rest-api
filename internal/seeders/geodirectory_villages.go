package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// seedVillages seeds village data from villages/keldesa-*.json files
func (gs *GeodirectorySeeder) seedVillages(ctx context.Context, geoDir string) error {
	villagesDir := filepath.Join(geoDir, "villages")
	gs.logger.WithField("directory", villagesDir).Info("Seeding villages")

	files, err := filepath.Glob(filepath.Join(villagesDir, "keldesa-*.json"))
	if err != nil {
		return fmt.Errorf("failed to list village files: %w", err)
	}

	gs.logger.WithField("file_count", len(files)).Info("Found village files")

	totalSuccess := 0
	totalErrors := 0

	for i, filename := range files {
		// Extract codes from filename (e.g., keldesa-15-02-010.json -> 15, 02, 010)
		basename := filepath.Base(filename)
		parts := strings.Split(strings.TrimSuffix(basename, ".json"), "-")
		if len(parts) != 4 {
			gs.logger.WithField("filename", basename).Warn("Invalid village filename format")
			continue
		}

		provinceCode := parts[1]
		cityCode := parts[2]
		districtCode := parts[3]
		parentDistrictCode := provinceCode + cityCode + districtCode

		// Get parent district
		parent, err := gs.repo.GetByCode(ctx, parentDistrictCode)
		if err != nil {
			gs.logger.WithError(err).WithField("district_code", parentDistrictCode).Warn("Failed to find parent district for villages")
			continue
		}

		file, err := os.Open(filename)
		if err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to open village file")
			continue
		}

		var villages map[string]string
		if err := json.NewDecoder(file).Decode(&villages); err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to decode village JSON")
			file.Close()
			continue
		}
		file.Close()

		successCount := 0
		errorCount := 0
		orderingCounter := 1

		for code, name := range villages {
			fullCode := parentDistrictCode + code // Combine parent district and village codes

			// Check if village already exists
			existing, getErr := gs.repo.GetByCode(ctx, fullCode)
			if getErr != nil && !isNotFoundError(getErr) {
				gs.logger.WithError(getErr).WithField("code", fullCode).Warn("Failed to check existing village")
				errorCount++
				continue
			}

			if existing != nil {
				// Update existing village
				existing.Name = name
				existing.SetDepth(existing.GetDepthForType())
				existing.SetOrderingID(orderingCounter)
				if updateErr := gs.repo.Update(ctx, existing); updateErr != nil {
					gs.logger.WithError(updateErr).WithField("code", fullCode).Error("Failed to update village")
					errorCount++
					continue
				}
			} else {
				// Create new village
				village := entities.NewGeodirectory(name, entities.GeoTypeVillage)
				village.SetCode(fullCode)
				village.SetParent(parent.ID)
				village.SetDepth(village.GetDepthForType())
				village.SetOrderingID(orderingCounter)

				if createErr := gs.repo.Create(ctx, village); createErr != nil {
					gs.logger.WithError(createErr).WithField("code", fullCode).Error("Failed to create village")
					errorCount++
					continue
				}
			}
			successCount++
			orderingCounter++
		}

		totalSuccess += successCount
		totalErrors += errorCount

		// Log progress every 100 files
		if (i+1)%100 == 0 || i == len(files)-1 {
			gs.logger.WithFields(map[string]interface{}{
				"files_processed":   i + 1,
				"total_files":       len(files),
				"records_processed": totalSuccess + totalErrors,
				"successful":        totalSuccess,
				"errors":            totalErrors,
			}).Info("Village seeding progress")
		}
	}

	gs.logger.WithFields(map[string]interface{}{
		"total_processed": totalSuccess + totalErrors,
		"successful":      totalSuccess,
		"errors":          totalErrors,
	}).Info("Village seeding completed")

	return nil
}
