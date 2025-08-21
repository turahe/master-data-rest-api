package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/http/middleware"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/redis"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/pkg/logger"
	"github.com/turahe/master-data-rest-api/pkg/response"
)

// SetupRouter configures and returns a Fiber app with all routes
func SetupRouter(
	config *configs.Config,
	log *logger.Logger,
	geodirectoryHandler *GeodirectoryHTTPHandler,
	apiKeyHandler *APIKeyHTTPHandler,
	bankHandler *BankHTTPHandler,
	currencyHandler *CurrencyHTTPHandler,
	languageHandler *LanguageHTTPHandler,
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

	// Apply authentication middleware based on configuration
	if config.Auth.Required {
		// Strict authentication - API key is required
		api.Use(middleware.APIKeyAuth(apiKeyService))
	} else {
		// Optional authentication - API key is optional but validated when provided
		api.Use(middleware.OptionalAPIKeyAuth(apiKeyService))
	}

	// Rate limit management routes (only if Redis is enabled)
	if redisManager.IsEnabled() {
		rateLimitHandler := NewRateLimitHTTPHandler(rateLimiter, log.Logger)
		rateLimits := api.Group("/rate-limit")
		rateLimits.Get("/info", rateLimitHandler.GetRateLimitInfo)
		rateLimits.Get("/stats", rateLimitHandler.GetRateLimitStats)
		rateLimits.Get("/config", rateLimitHandler.GetRateLimitConfig)
		rateLimits.Post("/reset", rateLimitHandler.ResetRateLimit)
	}

	// Geodirectory routes
	geodirectories := api.Group("/geodirectories")
	geodirectories.Get("/", geodirectoryHandler.GetAllGeodirectories)
	geodirectories.Get("/search", geodirectoryHandler.SearchGeodirectories)
	geodirectories.Get("/type/:type", geodirectoryHandler.GetGeodirectoriesByType)
	geodirectories.Get("/:id", geodirectoryHandler.GetGeodirectoryByID)
	geodirectories.Get("/:id/hierarchy", geodirectoryHandler.GetGeodirectoryWithHierarchy)
	geodirectories.Get("/:id/children", geodirectoryHandler.GetChildren)
	geodirectories.Get("/:id/ancestors", geodirectoryHandler.GetAncestors)
	geodirectories.Get("/:id/descendants", geodirectoryHandler.GetDescendants)

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
	banks.Get("/", bankHandler.GetBanks)
	banks.Get("/:code", bankHandler.GetBankByCode)

	// Currency routes
	currencies := api.Group("/currencies")
	currencies.Get("/", currencyHandler.GetCurrencies)
	currencies.Get("/:code", currencyHandler.GetCurrencyByCode)

	// Language routes
	languages := api.Group("/languages")
	languages.Get("/", languageHandler.GetAllLanguages)
	languages.Get("/search", languageHandler.SearchLanguages)

	return app
}
