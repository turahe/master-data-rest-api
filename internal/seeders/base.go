package seeders

import (
	"context"
	"fmt"

	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/pkg/logger"
)

// Seeder defines the interface for all data seeders
type Seeder interface {
	Seed(ctx context.Context, dataDir string) error
	Clear(ctx context.Context) error
	Name() string
}

// SeederManager manages all seeders
type SeederManager struct {
	logger  *logger.Logger
	seeders map[string]Seeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(
	geodirectoryRepo *pgx.GeodirectoryRepository,
	bankRepo *pgx.BankRepository,
	currencyRepo *pgx.CurrencyRepository,
	languageRepo *pgx.LanguageRepository,
	logger *logger.Logger,
) *SeederManager {
	seeders := map[string]Seeder{
		"languages":      NewLanguageSeeder(languageRepo, logger),
		"banks":          NewBankSeeder(bankRepo, logger),
		"currencies":     NewCurrencySeeder(currencyRepo, logger),
		"geodirectories": NewGeodirectorySeeder(geodirectoryRepo, logger),
	}

	return &SeederManager{
		logger:  logger,
		seeders: seeders,
	}
}

// Seed seeds specific data type or all if name is empty
func (sm *SeederManager) Seed(ctx context.Context, dataDir string, name string) error {
	if name != "" {
		seeder, exists := sm.seeders[name]
		if !exists {
			return fmt.Errorf("unknown seeder '%s'. Available: languages, banks, currencies, geodirectories", name)
		}

		sm.logger.WithField("seeder", name).Info("Starting specific seeding")
		return seeder.Seed(ctx, dataDir)
	}

	// Seed all data types
	sm.logger.Info("Starting seeding for all data types")
	for name, seeder := range sm.seeders {
		sm.logger.WithField("seeder", name).Info("Starting seeding")
		if err := seeder.Seed(ctx, dataDir); err != nil {
			return fmt.Errorf("failed to seed %s: %w", name, err)
		}
		sm.logger.WithField("seeder", name).Info("Seeding completed")
	}

	return nil
}

// Clear clears specific data type or all if name is empty
func (sm *SeederManager) Clear(ctx context.Context, name string) error {
	if name != "" {
		seeder, exists := sm.seeders[name]
		if !exists {
			return fmt.Errorf("unknown seeder '%s'. Available: languages, banks, currencies, geodirectories", name)
		}

		sm.logger.WithField("seeder", name).Info("Starting specific clearing using TRUNCATE")
		return seeder.Clear(ctx)
	}

	// Clear all data types using TRUNCATE for efficient bulk deletion
	sm.logger.Info("Starting clearing for all data types using TRUNCATE")
	for name, seeder := range sm.seeders {
		sm.logger.WithField("seeder", name).Info("Starting clearing with TRUNCATE")
		if err := seeder.Clear(ctx); err != nil {
			return fmt.Errorf("failed to clear %s: %w", name, err)
		}
		sm.logger.WithField("seeder", name).Info("Clearing completed")
	}

	return nil
}

// GetAvailableSeeders returns list of available seeder names
func (sm *SeederManager) GetAvailableSeeders() []string {
	names := make([]string, 0, len(sm.seeders))
	for name := range sm.seeders {
		names = append(names, name)
	}
	return names
}
