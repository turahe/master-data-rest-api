package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/mysql"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse command line flags
	var (
		dataDir = flag.String("data", "./data", "Directory containing JSON seed files")
		clear   = flag.Bool("clear", false, "Clear all seeded data before seeding")
		action  = flag.String("action", "seed", "Action to perform: seed, clear, or both")
	)
	flag.Parse()

	// Load configuration
	config := configs.Load()

	// Create database connection
	dbFactory := database.NewFactory()
	connManager, err := dbFactory.CreateConnectionManager(config)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}
	defer connManager.Close()

	// Run migrations
	if err := connManager.RunMigrations("./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	countryRepo := mysql.NewCountryRepository(connManager.GetDB())
	provinceRepo := mysql.NewProvinceRepository(connManager.GetDB())
	cityRepo := mysql.NewCityRepository(connManager.GetDB())
	districtRepo := mysql.NewDistrictRepository(connManager.GetDB())
	villageRepo := mysql.NewVillageRepository(connManager.GetDB())
	bankRepo := mysql.NewBankRepository(connManager.GetDB())
	currencyRepo := mysql.NewCurrencyRepository(connManager.GetDB())
	languageRepo := mysql.NewLanguageRepository(connManager.GetDB())

	// Create seeder service
	seederService := services.NewSeederService(
		countryRepo,
		provinceRepo,
		cityRepo,
		districtRepo,
		villageRepo,
		bankRepo,
		currencyRepo,
		languageRepo,
	)

	// Perform requested action
	switch *action {
	case "seed":
		if *clear {
			log.Println("Clearing existing data...")
			if err := seederService.ClearAll(); err != nil {
				log.Printf("Warning: Failed to clear data: %v", err)
			}
		}

		log.Println("Starting to seed data...")
		if err := seederService.SeedAll(*dataDir); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
		log.Println("Data seeding completed successfully!")

	case "clear":
		log.Println("Clearing all seeded data...")
		if err := seederService.ClearAll(); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		log.Println("Data clearing completed successfully!")

	case "both":
		log.Println("Clearing existing data...")
		if err := seederService.ClearAll(); err != nil {
			log.Printf("Warning: Failed to clear data: %v", err)
		}

		log.Println("Starting to seed data...")
		if err := seederService.SeedAll(*dataDir); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
		log.Println("Data seeding completed successfully!")

	default:
		log.Fatalf("Invalid action: %s. Use 'seed', 'clear', or 'both'", *action)
	}
}
