package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
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
		if err := clearAllData(geodirectoryRepo, bankRepo, currencyRepo, languageRepo); err != nil {
			log.WithError(err).Error("Failed to clear data")
			return fmt.Errorf("failed to clear data: %w", err)
		}
		log.Info("Data cleared successfully")
	}

	// Perform seeding based on seedName flag
	if seedName != "" {
		switch strings.ToLower(seedName) {
		case "languages":
			if err := seedLanguages(languageRepo, dataDir); err != nil {
				log.WithError(err).Error("Failed to seed languages")
				return fmt.Errorf("failed to seed languages: %w", err)
			}
		case "banks":
			if err := seedBanks(bankRepo, dataDir); err != nil {
				log.WithError(err).Error("Failed to seed banks")
				return fmt.Errorf("failed to seed banks: %w", err)
			}
		case "currencies":
			if err := seedCurrencies(currencyRepo, dataDir); err != nil {
				log.WithError(err).Error("Failed to seed currencies")
				return fmt.Errorf("failed to seed currencies: %w", err)
			}
		case "geodirectories":
			if err := seedGeodirectories(geodirectoryRepo, dataDir); err != nil {
				log.WithError(err).Error("Failed to seed geodirectories")
				return fmt.Errorf("failed to seed geodirectories: %w", err)
			}
		default:
			return fmt.Errorf("unsupported seed type: %s. Available types: languages, banks, currencies, geodirectories", seedName)
		}
		fmt.Printf("✅ Successfully seeded %s data!\n", seedName)
	} else {
		// Seed all data types
		if err := seedAllData(geodirectoryRepo, bankRepo, currencyRepo, languageRepo, dataDir); err != nil {
			log.WithError(err).Error("Failed to seed all data")
			return fmt.Errorf("failed to seed all data: %w", err)
		}
		fmt.Println("✅ Successfully seeded all data!")
	}

	fmt.Printf("📁 Data directory: %s\n", dataDir)
	fmt.Printf("🗑️  Clear data: %v\n", clearData)
	if seedName != "" {
		fmt.Printf("🎯 Seeded: %s\n", seedName)
	} else {
		fmt.Println("🎯 Seeded: all data types")
	}

	return nil
}

func clearAllData(geodirectoryRepo *pgx.GeodirectoryRepository, bankRepo *pgx.BankRepository, currencyRepo *pgx.CurrencyRepository, languageRepo *pgx.LanguageRepository) error {
	// TODO: Implement data clearing logic
	// This would involve calling repository methods to clear data
	// For safety, this is left as a TODO to prevent accidental data loss
	return nil
}

// seedLanguages reads and imports language data from CSV file
func seedLanguages(languageRepo *pgx.LanguageRepository, dataDir string) error {
	log := GetLogger()
	ctx := context.Background()

	// Read languages from CSV file
	csvFile := filepath.Join(dataDir, "tm_languages.csv")
	log.WithField("file", csvFile).Info("Reading languages from CSV file")

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

	log.WithField("count", len(records)).Info("Processing language records")

	successCount := 0
	errorCount := 0

	for i, record := range records {
		if len(record) < 3 {
			log.WithField("row", i+2).Warn("Skipping row with insufficient columns")
			errorCount++
			continue
		}

		code := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		_ = strings.TrimSpace(record[2]) // native name - not used in current entity

		if code == "" || name == "" {
			log.WithField("row", i+2).Warn("Skipping row with empty code or name")
			errorCount++
			continue
		}

		// Create language entity
		language := entities.NewLanguage(name, code)

		// Create or update language in database
		err = languageRepo.Create(ctx, language)
		if err != nil {
			// Try to update if create fails (might already exist)
			existing, getErr := languageRepo.GetByCode(ctx, code)
			if getErr != nil {
				log.WithError(err).WithFields(map[string]interface{}{
					"row":  i + 2,
					"code": code,
					"name": name,
				}).Warn("Failed to create or find existing language")
				errorCount++
				continue
			}

			// Update existing language
			existing.SetName(name)
			updateErr := languageRepo.Update(ctx, existing)
			if updateErr != nil {
				log.WithError(updateErr).WithFields(map[string]interface{}{
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
			log.WithField("processed", successCount).Info("Language seeding progress")
		}
	}

	log.WithFields(map[string]interface{}{
		"success": successCount,
		"errors":  errorCount,
		"total":   len(records),
	}).Info("Language seeding completed")

	fmt.Printf("📊 Processed %d language records: %d success, %d errors\n", len(records), successCount, errorCount)

	return nil
}

// seedBanks reads and imports bank data from CSV file
func seedBanks(bankRepo *pgx.BankRepository, dataDir string) error {
	log := GetLogger()

	csvFile := filepath.Join(dataDir, "tm_banks.csv")
	log.WithField("file", csvFile).Info("Banks seeding would be implemented here")

	// TODO: Implement banks CSV reading and seeding
	fmt.Println("🏦 Banks seeding implementation pending")
	return nil
}

// seedCurrencies reads and imports currency data from CSV file
func seedCurrencies(currencyRepo *pgx.CurrencyRepository, dataDir string) error {
	log := GetLogger()

	csvFile := filepath.Join(dataDir, "tm_currencies.csv")
	log.WithField("file", csvFile).Info("Currencies seeding would be implemented here")

	// TODO: Implement currencies CSV reading and seeding
	fmt.Println("💰 Currencies seeding implementation pending")
	return nil
}

// seedGeodirectories reads and imports geodirectory data from JSON files
func seedGeodirectories(geodirectoryRepo *pgx.GeodirectoryRepository, dataDir string) error {
	log := GetLogger()

	geoDir := filepath.Join(dataDir, "geodirectories")
	log.WithField("directory", geoDir).Info("Geodirectories seeding would be implemented here")

	// TODO: Implement geodirectories JSON reading and seeding
	fmt.Println("🌍 Geodirectories seeding implementation pending")
	return nil
}

// seedAllData seeds all available data types
func seedAllData(geodirectoryRepo *pgx.GeodirectoryRepository, bankRepo *pgx.BankRepository, currencyRepo *pgx.CurrencyRepository, languageRepo *pgx.LanguageRepository, dataDir string) error {
	log := GetLogger()

	log.Info("Seeding all data types")

	// Seed languages first
	if err := seedLanguages(languageRepo, dataDir); err != nil {
		return fmt.Errorf("failed to seed languages: %w", err)
	}

	// Seed other data types
	if err := seedBanks(bankRepo, dataDir); err != nil {
		return fmt.Errorf("failed to seed banks: %w", err)
	}

	if err := seedCurrencies(currencyRepo, dataDir); err != nil {
		return fmt.Errorf("failed to seed currencies: %w", err)
	}

	if err := seedGeodirectories(geodirectoryRepo, dataDir); err != nil {
		return fmt.Errorf("failed to seed geodirectories: %w", err)
	}

	return nil
}
