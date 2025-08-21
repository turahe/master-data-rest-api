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
// @description Enter your API key in the format: Bearer YOUR_API_KEY (optional when AUTH_REQUIRED=false)
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
	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/http"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/search"
	"github.com/turahe/master-data-rest-api/internal/domain/services"

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
	app := http.SetupRouter(config, log, geodirectoryHandler, apiKeyHandler, bankHandler, currencyHandler, languageHandler, apiKeyService, geodirectoryService)

	// Start server
	port := ":" + config.Server.Port
	log.WithField("port", config.Server.Port).Info("Starting HTTP server")

	if err := app.Listen(port); err != nil {
		log.WithError(err).Fatal("Failed to start HTTP server")
		return err
	}

	return nil
}
