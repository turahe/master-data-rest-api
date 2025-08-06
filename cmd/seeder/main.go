package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/mysql"
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

	// Create GORM connection for bank repository
	gormConnManager, err := dbFactory.CreateGORMConnectionManager(config)
	if err != nil {
		log.Fatalf("Failed to create GORM database connection: %v", err)
	}
	defer gormConnManager.Close()

	// Create repositories
	countryRepo := mysql.NewCountryRepository(connManager.GetDB())
	provinceRepo := mysql.NewProvinceRepository(connManager.GetDB())
	cityRepo := mysql.NewCityRepository(connManager.GetDB())
	districtRepo := mysql.NewDistrictRepository(connManager.GetDB())
	villageRepo := mysql.NewVillageRepository(connManager.GetDB())
	bankRepo := mysql.NewBankRepository(gormConnManager.GetDB())
	currencyRepo := mysql.NewCurrencyRepository(connManager.GetDB())
	languageRepo := mysql.NewLanguageRepository(connManager.GetDB())

	log.Printf("Successfully created repositories:")
	log.Printf("- Country repository: %T", countryRepo)
	log.Printf("- Province repository: %T", provinceRepo)
	log.Printf("- City repository: %T", cityRepo)
	log.Printf("- District repository: %T", districtRepo)
	log.Printf("- Village repository: %T", villageRepo)
	log.Printf("- Bank repository: %T (using GORM)", bankRepo)
	log.Printf("- Currency repository: %T", currencyRepo)
	log.Printf("- Language repository: %T", languageRepo)

	log.Printf("Bank repository now uses GORM instead of raw SQL!")
	log.Printf("Action: %s, Data directory: %s, Clear: %v", *action, *dataDir, *clear)
}
