package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/seeders"
)

var (
	dataDir   string
	clearData bool
	seedOnly  bool
	seedName  string
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

  # Seed specific data type
  master-data-api seed --name languages
  master-data-api seed --name banks
  master-data-api seed --name currencies
  master-data-api seed --name geodirectories

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
	seedCmd.Flags().StringVarP(&seedName, "name", "n", "", "seed specific data type (languages, banks, currencies, geodirectories)")
}

func runSeeder() error {
	config := GetConfig()
	log := GetLogger()
	ctx := context.Background()

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
	bankRepo := pgx.NewBankRepository(dbConnection.GetPool())
	currencyRepo := pgx.NewCurrencyRepository(dbConnection.GetPool())
	languageRepo := pgx.NewLanguageRepository(dbConnection.GetPool())

	// Initialize seeder manager
	seederManager := seeders.NewSeederManager(
		geodirectoryRepo,
		bankRepo,
		currencyRepo,
		languageRepo,
		log,
	)

	log.WithFields(map[string]interface{}{
		"data_dir":   dataDir,
		"clear_data": clearData,
		"seed_only":  seedOnly,
		"seed_name":  seedName,
	}).Info("Seeder configuration loaded")

	log.Info("Successfully created repositories using pgx:")
	log.WithField("repositories", []string{
		"Geodirectory", "Bank", "Currency", "Language",
	}).Info("All repositories initialized with pgx driver")

	// Clear data if requested
	if clearData && !seedOnly {
		log.Info("Clearing existing data")
		if err := seederManager.Clear(ctx, seedName); err != nil {
			log.WithError(err).Error("Failed to clear data")
			return fmt.Errorf("failed to clear data: %w", err)
		}
		log.Info("Data cleared successfully")
	}

	// Perform seeding
	log.Info("Starting seeding process")
	if err := seederManager.Seed(ctx, dataDir, seedName); err != nil {
		log.WithError(err).Error("Failed to seed data")
		return fmt.Errorf("failed to seed data: %w", err)
	}

	fmt.Printf("📁 Data directory: %s\n", dataDir)
	fmt.Printf("🗑️  Clear data: %v\n", clearData)
	if seedName != "" {
		fmt.Printf("🎯 Seeded: %s\n", seedName)
	} else {
		fmt.Println("🎯 Seeded: all data types")
	}

	fmt.Println("✅ Seeding completed successfully!")
	return nil
}
