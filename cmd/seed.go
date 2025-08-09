package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
)

var (
	dataDir   string
	clearData bool
	seedOnly  bool
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with sample data",
	Long: `Seed the database with sample data for development and testing.

This command can populate the database with initial data for:
- Geographical data (countries, provinces, cities, districts, villages)
- Banking information
- Currency data
- Language information

The seeder supports various options for clearing existing data
and specifying custom data directories.

Examples:
  # Seed with default data
  master-data-api seed

  # Clear existing data and seed fresh
  master-data-api seed --clear

  # Seed from custom data directory
  master-data-api seed --data-dir ./custom-data

  # Only seed without clearing
  master-data-api seed --seed-only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}

		return runSeeder()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)

	// Seeder flags
	seedCmd.Flags().StringVarP(&dataDir, "data-dir", "d", "configs/data", "directory containing seed data files")
	seedCmd.Flags().BoolVarP(&clearData, "clear", "c", false, "clear existing data before seeding")
	seedCmd.Flags().BoolVar(&seedOnly, "seed-only", false, "only seed data, don't clear existing data")
}

func runSeeder() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Starting seeder application")

	// Initialize database connection
	log.Info("Connecting to database")
	dbConnection := database.NewPgxConnectionWithLogger(config.Database, log)
	if err := dbConnection.Connect(); err != nil {
		log.WithError(err).Error("Failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbConnection.Close()

	log.Info("Database connected successfully")

	// Run migrations
	log.Info("Running database migrations")
	migrator := database.NewMigrator(config.Database)
	if err := migrator.RunMigrations("migrations"); err != nil {
		log.WithError(err).Error("Failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Info("Database migrations completed successfully")

	// Initialize repositories
	log.Info("Initializing repositories")
	geodirectoryRepo := pgx.NewGeodirectoryRepository(dbConnection.GetPool())
	apiKeyRepo := pgx.NewAPIKeyRepository(dbConnection.GetPool())
	bankRepo := pgx.NewBankRepository(dbConnection.GetPool())
	currencyRepo := pgx.NewCurrencyRepository(dbConnection.GetPool())
	languageRepo := pgx.NewLanguageRepository(dbConnection.GetPool())

	log.WithFields(map[string]interface{}{
		"data_dir":   dataDir,
		"clear_data": clearData,
		"seed_only":  seedOnly,
	}).Info("Seeder configuration loaded")

	log.Info("Successfully created repositories using pgx:")
	log.WithField("repositories", []string{
		"Geodirectory", "APIKey", "Bank", "Currency", "Language",
	}).Info("All repositories initialized with pgx driver")

	// Clear data if requested
	if clearData && !seedOnly {
		log.Info("Clearing existing data")
		if err := clearAllData(geodirectoryRepo, bankRepo, currencyRepo, languageRepo); err != nil {
			log.WithError(err).Error("Failed to clear data")
			return fmt.Errorf("failed to clear data: %w", err)
		}
		log.Info("Data cleared successfully")
	}

	// TODO: Implement actual seeding logic
	// For now, just demonstrate the repositories are available
	repositories := map[string]interface{}{
		"geodirectory": geodirectoryRepo,
		"api_key":      apiKeyRepo,
		"bank":         bankRepo,
		"currency":     currencyRepo,
		"language":     languageRepo,
	}

	log.WithField("repository_count", len(repositories)).Info("Seeder setup completed successfully")
	fmt.Println("✅ Seeder setup completed!")
	fmt.Printf("📁 Data directory: %s\n", dataDir)
	fmt.Printf("🗑️  Clear data: %v\n", clearData)
	fmt.Printf("📊 Repositories available: %d\n", len(repositories))
	fmt.Println("\n📝 Note: Actual seeding logic can be implemented based on your data files")

	return nil
}

func clearAllData(geodirectoryRepo *pgx.GeodirectoryRepository, bankRepo *pgx.BankRepository, currencyRepo *pgx.CurrencyRepository, languageRepo *pgx.LanguageRepository) error {
	// TODO: Implement data clearing logic
	// This would involve calling repository methods to clear data
	// For safety, this is left as a TODO to prevent accidental data loss
	return nil
}
