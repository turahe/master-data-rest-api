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

// seedCities seeds city and regency data from cities/kab-*.json files
// Names with "KAB " prefix are treated as regencies, "KOTA " prefix as cities
func (gs *GeodirectorySeeder) seedCities(ctx context.Context, geoDir string) error {
	citiesDir := filepath.Join(geoDir, "cities")
	gs.logger.WithField("directory", citiesDir).Info("Seeding cities and regencies")

	files, err := filepath.Glob(filepath.Join(citiesDir, "kab-*.json"))
	if err != nil {
		return fmt.Errorf("failed to list city files: %w", err)
	}

	totalSuccess := 0
	totalErrors := 0

	for _, filename := range files {
		// Extract province code from filename (e.g., kab-11.json -> 11)
		basename := filepath.Base(filename)
		provinceCode := strings.TrimPrefix(strings.TrimSuffix(basename, ".json"), "kab-")

		// Get parent province
		parent, err := gs.repo.GetByCode(ctx, provinceCode)
		if err != nil {
			gs.logger.WithError(err).WithField("province_code", provinceCode).Warn("Failed to find parent province for cities/regencies")
			continue
		}

		file, err := os.Open(filename)
		if err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to open city file")
			continue
		}

		var cities map[string]string
		if err := json.NewDecoder(file).Decode(&cities); err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to decode city JSON")
			file.Close()
			continue
		}
		file.Close()

		successCount := 0
		errorCount := 0
		orderingCounter := 1

		for code, name := range cities {
			fullCode := provinceCode + code // Combine province and city codes

			// Check if city already exists
			existing, getErr := gs.repo.GetByCode(ctx, fullCode)
			if getErr != nil && !isNotFoundError(getErr) {
				gs.logger.WithError(getErr).WithField("code", fullCode).Warn("Failed to check existing city")
				errorCount++
				continue
			}

			// Determine the geo type based on name prefix
			var geoType entities.GeoType
			if strings.HasPrefix(strings.ToUpper(name), "KAB ") {
				geoType = entities.GeoTypeRegency
			} else if strings.HasPrefix(strings.ToUpper(name), "KOTA ") {
				geoType = entities.GeoTypeCity
			} else {
				// Default to city if no prefix
				geoType = entities.GeoTypeCity
			}

			if existing != nil {
				// Update existing city/regency
				existing.Name = name
				existing.Type = geoType
				existing.SetDepth(existing.GetDepthForType())
				existing.SetOrderingID(orderingCounter)
				if updateErr := gs.repo.Update(ctx, existing); updateErr != nil {
					gs.logger.WithError(updateErr).WithField("code", fullCode).Error("Failed to update city/regency")
					errorCount++
					continue
				}
			} else {
				// Create new city/regency
				geodirectory := entities.NewGeodirectory(name, geoType)
				geodirectory.SetCode(fullCode)
				geodirectory.SetParent(parent.ID)
				geodirectory.SetDepth(geodirectory.GetDepthForType())
				geodirectory.SetOrderingID(orderingCounter)

				if createErr := gs.repo.Create(ctx, geodirectory); createErr != nil {
					gs.logger.WithError(createErr).WithField("code", fullCode).Error("Failed to create city/regency")
					errorCount++
					continue
				}
			}
			successCount++
			orderingCounter++
		}

		totalSuccess += successCount
		totalErrors += errorCount

		gs.logger.WithFields(map[string]interface{}{
			"province_code": provinceCode,
			"processed":     successCount + errorCount,
			"successful":    successCount,
			"errors":        errorCount,
		}).Info("City/regency seeding progress for province")
	}

	gs.logger.WithFields(map[string]interface{}{
		"total_processed": totalSuccess + totalErrors,
		"successful":      totalSuccess,
		"errors":          totalErrors,
	}).Info("City/regency seeding completed")

	return nil
}
