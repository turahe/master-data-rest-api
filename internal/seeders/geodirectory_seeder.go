package seeders

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
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
	gs.logger.WithField("directory", geoDir).Info("Geodirectories seeding")

	// TODO: Implement geodirectories JSON reading and seeding
	// This would involve:
	// 1. Reading provinces from geodirectories/provinces/provinsi.json
	// 2. Reading cities from geodirectories/cities/*.json
	// 3. Reading districts from geodirectories/districts/*.json
	// 4. Reading villages from geodirectories/villages/*.json
	// 5. Creating nested set relationships

	fmt.Println("🌍 Geodirectories seeding implementation pending")
	gs.logger.Info("Geodirectories seeding completed (implementation pending)")

	return nil
}

// Clear removes all geodirectory data
func (gs *GeodirectorySeeder) Clear(ctx context.Context) error {
	// TODO: Implement geodirectory clearing
	// This would involve clearing all geodirectory records in proper order
	// to maintain referential integrity

	gs.logger.Info("Geodirectories clearing completed (implementation pending)")
	return nil
}
