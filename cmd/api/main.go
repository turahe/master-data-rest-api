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
	_ "github.com/turahe/master-data-rest-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Initialize services
	countryService := services.NewCountryService(countryRepo)

	// Initialize handlers
	countryHandler := http.NewCountryHTTPHandler(countryService)

	// Setup router
	router := setupRouter(countryHandler)

	// Start server
	log.Printf("Server starting on port %s", config.Server.Port)
	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(countryHandler *http.CountryHTTPHandler) *gin.Engine {
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

		// Master data routes will be added here
		// Example: provinces, cities, etc.
	}

	return router
}
