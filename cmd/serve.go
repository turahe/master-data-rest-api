// Package cmd provides command-line interface for the Master Data REST API
//
// @title Master Data REST API
// @version 1.0.0
// @description A comprehensive REST API for managing master data including geographical information, banks, currencies, and languages. Built with Go, Fiber, and PostgreSQL using Hexagonal Architecture.
// @termsOfService https://github.com/turahe/master-data-rest-api
//
// @contact.name API Support
// @contact.url https://github.com/turahe/master-data-rest-api/issues
// @contact.email support@turahe.com
//
// @license.name MIT
// @license.url https://github.com/turahe/master-data-rest-api/blob/main/LICENSE
//
// @host localhost:8080
// @BasePath /api/v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter your API key in the format: Bearer YOUR_API_KEY
//
// @tag.name geodirectories
// @tag.description Operations for managing geographical directories (countries, provinces, cities, districts, villages)
//
// @tag.name banks
// @tag.description Operations for managing bank information
//
// @tag.name currencies
// @tag.description Operations for managing currency information
//
// @tag.name languages
// @tag.description Operations for managing language information
//
// @tag.name api-keys
// @tag.description Operations for managing API keys and authentication
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/http"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/http/middleware"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/redis"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/search"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/logger"
	"github.com/turahe/master-data-rest-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	// Swagger imports
	_ "github.com/turahe/master-data-rest-api/docs"
)

var (
	serverHost string
	serverPort string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Master Data REST API server",
	Long: `Start the Master Data REST API server with full functionality including:
- RESTful endpoints for geographical data management
- Banks, currencies, and languages management
- API key authentication and management
- Comprehensive logging and monitoring
- Auto-generated Swagger documentation

The server supports configurable host, port, and database logging options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}

		return runServer(cmd)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Server-specific flags
	serveCmd.Flags().StringVarP(&serverHost, "host", "H", "", "server host (overrides config)")
	serveCmd.Flags().StringVarP(&serverPort, "port", "p", "", "server port (overrides config)")
	serveCmd.Flags().Bool("migrate", true, "run database migrations on startup")
}

func runServer(cmd *cobra.Command) error {
	config := GetConfig()
	log := GetLogger()

	log.WithField("app", config.App.Name).Info("Starting application")

	// Override config with flags if provided
	if serverHost != "" {
		config.Server.Host = serverHost
	}
	if serverPort != "" {
		config.Server.Port = serverPort
	}

	// Initialize database connection
	log.Info("Connecting to database")
	dbConnection := database.NewPgxConnectionWithLogger(config.Database, log)
	if err := dbConnection.Connect(); err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
		return err
	}
	defer dbConnection.Close()

	log.Info("Database connected successfully")

	// Run migrations by default (can be disabled with --migrate=false)
	if migrate, _ := cmd.Flags().GetBool("migrate"); migrate {
		log.Info("Running database migrations")
		migrator := database.NewMigrator(config.Database)
		if err := migrator.RunMigrations("migrations"); err != nil {
			log.WithError(err).Fatal("Failed to run migrations")
			return err
		}
		log.Info("Database migrations completed successfully")
	}

	// Initialize repositories
	log.Info("Initializing repositories")
	geodirectoryRepo := pgx.NewGeodirectoryRepository(dbConnection.GetPool())
	apiKeyRepo := pgx.NewAPIKeyRepository(dbConnection.GetPool())
	bankRepo := pgx.NewBankRepository(dbConnection.GetPool())
	currencyRepo := pgx.NewCurrencyRepository(dbConnection.GetPool())
	languageRepo := pgx.NewLanguageRepository(dbConnection.GetPool())

	// Initialize search service
	log.Info("Initializing search service")
	meilisearchClient := search.NewMeilisearchClient(search.Config{
		Host:   config.Meilisearch.Host,
		APIKey: config.Meilisearch.APIKey,
	})
	searchService := search.NewMeilisearchRepository(meilisearchClient)

	// Initialize services
	log.Info("Initializing services")
	geodirectoryService := services.NewGeodirectoryService(geodirectoryRepo)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	bankService := services.NewBankService(bankRepo)
	currencyService := services.NewCurrencyService(currencyRepo)
	languageService := services.NewLanguageService(languageRepo)

	// Initialize handlers
	log.Info("Initializing HTTP handlers")
	geodirectoryHandler := http.NewGeodirectoryHTTPHandler(geodirectoryService, searchService)
	apiKeyHandler := http.NewAPIKeyHTTPHandler(apiKeyService)
	bankHandler := http.NewBankHTTPHandler(bankService, searchService)
	currencyHandler := http.NewCurrencyHTTPHandler(currencyService, searchService)
	languageHandler := http.NewLanguageHTTPHandler(languageService, searchService)

	// Setup router
	app := setupRouter(config, log, geodirectoryHandler, apiKeyHandler, bankHandler, currencyHandler, languageHandler, apiKeyService, geodirectoryService)

	// Start server
	port := ":" + config.Server.Port
	log.WithField("port", config.Server.Port).Info("Starting HTTP server")

	if err := app.Listen(port); err != nil {
		log.WithError(err).Fatal("Failed to start HTTP server")
		return err
	}

	return nil
}

func setupRouter(
	config *configs.Config,
	log *logger.Logger,
	geodirectoryHandler *http.GeodirectoryHTTPHandler,
	apiKeyHandler *http.APIKeyHTTPHandler,
	bankHandler *http.BankHTTPHandler,
	currencyHandler *http.CurrencyHTTPHandler,
	languageHandler *http.LanguageHTTPHandler,
	apiKeyService *services.APIKeyService,
	geodirectoryService *services.GeodirectoryService,
) *fiber.App {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.WithError(err).Error("Request failed")
			return response.InternalServerError(c, "Internal server error")
		},
	})

	// Initialize Redis manager and rate limiter
	redisManager := redis.NewManager(&config.Redis, log.Logger)
	rateLimiter := redis.NewRateLimiter(redisManager.GetClient(), log.Logger)

	// Connect to Redis if enabled
	if config.Redis.Enabled {
		if err := redisManager.Connect(context.Background()); err != nil {
			log.WithError(err).Warn("Failed to connect to Redis, rate limiting will be disabled")
		}
	}

	// Add custom middleware
	app.Use(middleware.RequestLoggerMiddleware(log))
	app.Use(middleware.ErrorLoggerMiddleware(log))

	// Add rate limiting middleware if Redis is enabled
	if redisManager.IsEnabled() {
		app.Use(middleware.TieredRateLimiter(rateLimiter))
		log.Info("Rate limiting middleware enabled")
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": config.App.Name,
			"version": config.App.Version,
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/api/v1")

	// Apply API key authentication middleware
	api.Use(middleware.APIKeyAuth(apiKeyService))

	// Rate limit management routes (only if Redis is enabled)
	if redisManager.IsEnabled() {
		rateLimitHandler := http.NewRateLimitHTTPHandler(rateLimiter, log.Logger)
		rateLimits := api.Group("/rate-limit")
		rateLimits.Get("/info", rateLimitHandler.GetRateLimitInfo)
		rateLimits.Get("/stats", rateLimitHandler.GetRateLimitStats)
		rateLimits.Get("/config", rateLimitHandler.GetRateLimitConfig)
		rateLimits.Post("/reset", rateLimitHandler.ResetRateLimit)
	}

	// Geodirectory routes
	geodirectories := api.Group("/geodirectories")
	geodirectories.Post("/", geodirectoryHandler.CreateGeodirectory)
	geodirectories.Get("/", geodirectoryHandler.GetAllGeodirectories)
	geodirectories.Get("/search", geodirectoryHandler.SearchGeodirectories)
	geodirectories.Get("/type/:type", geodirectoryHandler.GetGeodirectoriesByType)
	geodirectories.Post("/rebuild", geodirectoryHandler.RebuildNestedSet)
	geodirectories.Get("/:id", geodirectoryHandler.GetGeodirectoryByID)
	geodirectories.Get("/:id/hierarchy", geodirectoryHandler.GetGeodirectoryWithHierarchy)
	geodirectories.Get("/:id/children", geodirectoryHandler.GetChildren)
	geodirectories.Get("/:id/ancestors", geodirectoryHandler.GetAncestors)
	geodirectories.Get("/:id/descendants", geodirectoryHandler.GetDescendants)
	geodirectories.Put("/:id", geodirectoryHandler.UpdateGeodirectory)
	geodirectories.Post("/:id/move", geodirectoryHandler.MoveGeodirectory)
	geodirectories.Delete("/:id", geodirectoryHandler.DeleteGeodirectory)

	// Backward compatibility routes for countries, provinces, cities, etc.
	countries := api.Group("/countries")
	countries.Get("/", func(c *fiber.Ctx) error {
		limit := 50
		offset := 0
		geodirectories, err := geodirectoryService.GetCountries(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve countries: "+err.Error())
		}
		return response.Success(c, geodirectories, "Countries retrieved successfully")
	})

	provinces := api.Group("/provinces")
	provinces.Get("/", func(c *fiber.Ctx) error {
		limit := 50
		offset := 0
		geodirectories, err := geodirectoryService.GetProvinces(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve provinces: "+err.Error())
		}
		return response.Success(c, geodirectories, "Provinces retrieved successfully")
	})

	cities := api.Group("/cities")
	cities.Get("/", func(c *fiber.Ctx) error {
		limit := 50
		offset := 0
		geodirectories, err := geodirectoryService.GetCities(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve cities: "+err.Error())
		}
		return response.Success(c, geodirectories, "Cities retrieved successfully")
	})

	districts := api.Group("/districts")
	districts.Get("/", func(c *fiber.Ctx) error {
		limit := 50
		offset := 0
		geodirectories, err := geodirectoryService.GetDistricts(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve districts: "+err.Error())
		}
		return response.Success(c, geodirectories, "Districts retrieved successfully")
	})

	villages := api.Group("/villages")
	villages.Get("/", func(c *fiber.Ctx) error {
		limit := 50
		offset := 0
		geodirectories, err := geodirectoryService.GetVillages(c.Context(), limit, offset)
		if err != nil {
			return response.InternalServerError(c, "Failed to retrieve villages: "+err.Error())
		}
		return response.Success(c, geodirectories, "Villages retrieved successfully")
	})

	// API key management routes
	apiKeys := api.Group("/api-keys")
	apiKeys.Post("/", apiKeyHandler.CreateAPIKey)
	apiKeys.Get("/", apiKeyHandler.GetAllAPIKeys)
	apiKeys.Get("/:id", apiKeyHandler.GetAPIKeyByID)
	apiKeys.Put("/:id", apiKeyHandler.UpdateAPIKey)
	apiKeys.Post("/:id/activate", apiKeyHandler.ActivateAPIKey)
	apiKeys.Post("/:id/deactivate", apiKeyHandler.DeactivateAPIKey)
	apiKeys.Delete("/:id", apiKeyHandler.DeleteAPIKey)

	// Bank routes
	banks := api.Group("/banks")
	banks.Post("/", bankHandler.CreateBank)
	banks.Get("/", bankHandler.GetAllBanks)
	banks.Get("/search", bankHandler.SearchBanks)
	banks.Get("/code/:code", bankHandler.GetBankByCode)
	banks.Get("/:id", bankHandler.GetBankByID)
	banks.Put("/:id", bankHandler.UpdateBank)
	banks.Delete("/:id", bankHandler.DeleteBank)

	// Currency routes
	currencies := api.Group("/currencies")
	currencies.Post("/", currencyHandler.CreateCurrency)
	currencies.Get("/", currencyHandler.GetAllCurrencies)
	currencies.Get("/active", currencyHandler.GetActiveCurrencies)
	currencies.Get("/search", currencyHandler.SearchCurrencies)
	currencies.Get("/:id", currencyHandler.GetCurrencyByID)
	currencies.Put("/:id", currencyHandler.UpdateCurrency)
	currencies.Post("/:id/activate", currencyHandler.ActivateCurrency)
	currencies.Post("/:id/deactivate", currencyHandler.DeactivateCurrency)
	currencies.Delete("/:id", currencyHandler.DeleteCurrency)

	// Language routes
	languages := api.Group("/languages")
	languages.Post("/", languageHandler.CreateLanguage)
	languages.Get("/", languageHandler.GetAllLanguages)
	languages.Get("/active", languageHandler.GetActiveLanguages)
	languages.Get("/search", languageHandler.SearchLanguages)
	languages.Get("/:id", languageHandler.GetLanguageByID)
	languages.Put("/:id", languageHandler.UpdateLanguage)
	languages.Post("/:id/activate", languageHandler.ActivateLanguage)
	languages.Post("/:id/deactivate", languageHandler.DeactivateLanguage)
	languages.Delete("/:id", languageHandler.DeleteLanguage)

	return app
}
