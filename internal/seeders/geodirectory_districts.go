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

// seedDistricts seeds district data from districts/kec-*.json files
func (gs *GeodirectorySeeder) seedDistricts(ctx context.Context, geoDir string) error {
	districtsDir := filepath.Join(geoDir, "districts")
	gs.logger.WithField("directory", districtsDir).Info("Seeding districts")

	files, err := filepath.Glob(filepath.Join(districtsDir, "kec-*.json"))
	if err != nil {
		return fmt.Errorf("failed to list district files: %w", err)
	}

	totalSuccess := 0
	totalErrors := 0

	for _, filename := range files {
		// Extract province and city codes from filename (e.g., kec-11-01.json -> 11, 01)
		basename := filepath.Base(filename)
		parts := strings.Split(strings.TrimSuffix(basename, ".json"), "-")
		if len(parts) != 3 {
			gs.logger.WithField("filename", basename).Warn("Invalid district filename format")
			continue
		}

		provinceCode := parts[1]
		cityCode := parts[2]
		parentCityCode := provinceCode + cityCode

		// Get parent city
		parent, err := gs.repo.GetByCode(ctx, parentCityCode)
		if err != nil {
			gs.logger.WithError(err).WithField("city_code", parentCityCode).Warn("Failed to find parent city for districts")
			continue
		}

		file, err := os.Open(filename)
		if err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to open district file")
			continue
		}

		var districts map[string]string
		if err := json.NewDecoder(file).Decode(&districts); err != nil {
			gs.logger.WithError(err).WithField("file", filename).Warn("Failed to decode district JSON")
			file.Close()
			continue
		}
		file.Close()

		successCount := 0
		errorCount := 0
		orderingCounter := 1

		for code, name := range districts {
			fullCode := parentCityCode + code // Combine parent city and district codes

			// Check if district already exists
			existing, getErr := gs.repo.GetByCode(ctx, fullCode)
			if getErr != nil && !isNotFoundError(getErr) {
				gs.logger.WithError(getErr).WithField("code", fullCode).Warn("Failed to check existing district")
				errorCount++
				continue
			}

			if existing != nil {
				// Update existing district
				existing.Name = name
				existing.SetDepth(existing.GetDepthForType())
				existing.SetOrderingID(orderingCounter)
				if updateErr := gs.repo.Update(ctx, existing); updateErr != nil {
					gs.logger.WithError(updateErr).WithField("code", fullCode).Error("Failed to update district")
					errorCount++
					continue
				}
			} else {
				// Create new district
				district := entities.NewGeodirectory(name, entities.GeoTypeDistrict)
				district.SetCode(fullCode)
				district.SetParent(parent.ID)
				district.SetDepth(district.GetDepthForType())
				district.SetOrderingID(orderingCounter)

				if createErr := gs.repo.Create(ctx, district); createErr != nil {
					gs.logger.WithError(createErr).WithField("code", fullCode).Error("Failed to create district")
					errorCount++
					continue
				}
			}
			successCount++
			orderingCounter++
		}

		totalSuccess += successCount
		totalErrors += errorCount

		// Log progress every 50 files
		if len(files) > 50 && (totalSuccess+totalErrors)%500 == 0 {
			gs.logger.WithFields(map[string]interface{}{
				"processed":  totalSuccess + totalErrors,
				"successful": totalSuccess,
				"errors":     totalErrors,
			}).Info("District seeding progress")
		}
	}

	gs.logger.WithFields(map[string]interface{}{
		"total_processed": totalSuccess + totalErrors,
		"successful":      totalSuccess,
		"errors":          totalErrors,
	}).Info("District seeding completed")

	return nil
}
