package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/http"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/gorm"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/response"

	// Swagger imports
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/turahe/master-data-rest-api/docs"
)

// @title           Master Data REST API
// @version         1.0
// @description     A REST API for managing master data including countries, provinces, cities, districts, villages, banks, currencies, and languages.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	config := configs.Load()

	// Convert port string to int
	port, err := strconv.Atoi(config.Database.Port)
	if err != nil {
		log.Fatal("Invalid database port:", err)
	}

	// Initialize database connection manager
	dbManager, err := database.NewGORMConnectionManager(&database.Config{
		Driver:    database.DriverType(config.Database.Driver),
		Host:      config.Database.Host,
		Port:      port,
		Username:  config.Database.User,
		Password:  config.Database.Password,
		Database:  config.Database.Name,
		Charset:   config.Database.Charset,
		ParseTime: config.Database.ParseTime == "true",
		Loc:       config.Database.Loc,
		SSLMode:   config.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database schema
	err = dbManager.AutoMigrate(
		&entities.Country{},
		&entities.Province{},
		&entities.City{},
		&entities.District{},
		&entities.Village{},
		&entities.Bank{},
		&entities.Currency{},
		&entities.Language{},
	)
	if err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	// Initialize repositories
	countryRepo := gorm.NewCountryRepository(dbManager.GetDB())
	provinceRepo := gorm.NewProvinceRepository(dbManager.GetDB())
	cityRepo := gorm.NewCityRepository(dbManager.GetDB())
	districtRepo := gorm.NewDistrictRepository(dbManager.GetDB())
	villageRepo := gorm.NewVillageRepository(dbManager.GetDB())
	bankRepo := gorm.NewBankRepository(dbManager.GetDB())
	currencyRepo := gorm.NewCurrencyRepository(dbManager.GetDB())
	languageRepo := gorm.NewLanguageRepository(dbManager.GetDB())

	// Initialize services
	countryService := services.NewCountryService(countryRepo)
	provinceService := services.NewProvinceService(provinceRepo)
	cityService := services.NewCityService(cityRepo)
	districtService := services.NewDistrictService(districtRepo)
	villageService := services.NewVillageService(villageRepo)
	bankService := services.NewBankService(bankRepo)
	currencyService := services.NewCurrencyService(currencyRepo)
	languageService := services.NewLanguageService(languageRepo)

	// Initialize handlers
	countryHandler := http.NewCountryHTTPHandler(countryService)
	provinceHandler := http.NewProvinceHTTPHandler(provinceService)
	cityHandler := http.NewCityHTTPHandler(cityService)
	districtHandler := http.NewDistrictHTTPHandler(districtService)
	villageHandler := http.NewVillageHTTPHandler(villageService)
	bankHandler := http.NewBankHTTPHandler(bankService)
	currencyHandler := http.NewCurrencyHTTPHandler(currencyService)
	languageHandler := http.NewLanguageHTTPHandler(languageService)

	// Setup router
	router := setupRouter(
		countryHandler,
		provinceHandler,
		cityHandler,
		districtHandler,
		villageHandler,
		bankHandler,
		currencyHandler,
		languageHandler,
	)

	// Start server
	log.Printf("Server starting on port %s", config.Server.Port)
	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(
	countryHandler *http.CountryHTTPHandler,
	provinceHandler *http.ProvinceHTTPHandler,
	cityHandler *http.CityHTTPHandler,
	districtHandler *http.DistrictHTTPHandler,
	villageHandler *http.VillageHTTPHandler,
	bankHandler *http.BankHTTPHandler,
	currencyHandler *http.CurrencyHTTPHandler,
	languageHandler *http.LanguageHTTPHandler,
) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"time":   "now",
		}, "Server is running")
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Country routes
		countries := api.Group("/countries")
		{
			countries.POST("/", countryHandler.CreateCountry)
			countries.GET("/", countryHandler.GetAllCountries)
			countries.GET("/:id", countryHandler.GetCountryByID)
			countries.GET("/code/:code", countryHandler.GetCountryByCode)
			countries.GET("/iso3166-2/:iso", countryHandler.GetCountryByISO31662)
			countries.GET("/iso3166-3/:iso", countryHandler.GetCountryByISO31663)
			countries.GET("/name/:name", countryHandler.GetCountryByName)
			countries.GET("/region/:region", countryHandler.GetCountriesByRegion)
			countries.GET("/subregion/:subregion", countryHandler.GetCountriesBySubRegion)
			countries.GET("/eea", countryHandler.GetEEACountries)
			countries.PUT("/:id", countryHandler.UpdateCountry)
			countries.DELETE("/:id", countryHandler.DeleteCountry)
			countries.GET("/count", countryHandler.GetCountryCount)
		}

		// Province routes
		provinces := api.Group("/provinces")
		{
			provinces.POST("/", provinceHandler.CreateProvince)
			provinces.GET("/", provinceHandler.GetAllProvinces)
			provinces.GET("/:id", provinceHandler.GetProvinceByID)
			provinces.GET("/code/:code", provinceHandler.GetProvinceByCode)
			provinces.GET("/name/:name", provinceHandler.GetProvinceByName)
			provinces.GET("/country/:country_id", provinceHandler.GetProvincesByCountry)
			provinces.PUT("/:id", provinceHandler.UpdateProvince)
			provinces.DELETE("/:id", provinceHandler.DeleteProvince)
			provinces.GET("/count", provinceHandler.GetProvinceCount)
		}

		// City routes
		cities := api.Group("/cities")
		{
			cities.POST("/", cityHandler.CreateCity)
			cities.GET("/", cityHandler.GetAllCities)
			cities.GET("/:id", cityHandler.GetCityByID)
			cities.GET("/code/:code", cityHandler.GetCityByCode)
			cities.GET("/name/:name", cityHandler.GetCityByName)
			cities.GET("/province/:province_id", cityHandler.GetCitiesByProvince)
			cities.PUT("/:id", cityHandler.UpdateCity)
			cities.DELETE("/:id", cityHandler.DeleteCity)
			cities.GET("/count", cityHandler.GetCityCount)
		}

		// Bank routes
		banks := api.Group("/banks")
		{
			banks.POST("/", bankHandler.CreateBank)
			banks.GET("/", bankHandler.GetAllBanks)
			banks.GET("/:id", bankHandler.GetBankByID)
			banks.GET("/code/:code", bankHandler.GetBankByCode)
			banks.GET("/name/:name", bankHandler.GetBankByName)
			banks.GET("/alias/:alias", bankHandler.GetBankByAlias)
			banks.GET("/company/:company", bankHandler.GetBanksByCompany)
			banks.PUT("/:id", bankHandler.UpdateBank)
			banks.DELETE("/:id", bankHandler.DeleteBank)
			banks.GET("/count", bankHandler.GetBankCount)
		}

		// Master data routes will be added here
		// Example: districts, villages, currencies, languages, etc.
	}

	return router
}
