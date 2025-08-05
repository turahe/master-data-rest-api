package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/turahe/master-data-rest-api/configs"
	"github.com/turahe/master-data-rest-api/internal/adapters/primary/application"
	httphandler "github.com/turahe/master-data-rest-api/internal/adapters/primary/http"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/gorm"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	config := configs.Load()

	// Initialize database connection
	dbManager, err := initDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	// Run GORM auto-migrations
	if err := dbManager.AutoMigrate(
		&entities.User{},
		&entities.Country{},
		&entities.Province{},
		&entities.City{},
		&entities.District{},
		&entities.Village{},
		&entities.Bank{},
		&entities.Currency{},
		&entities.Language{},
	); err != nil {
		log.Fatalf("Failed to run auto-migrations: %v", err)
	}

	// Get GORM DB instance
	gormDB := dbManager.GetDB()

	// Initialize repositories
	userRepo := gorm.NewUserRepository(gormDB)

	// Initialize domain services
	userService := services.NewUserService(userRepo)

	// Initialize application services
	userAppService := application.NewUserApplicationService(userService)

	// Initialize HTTP handlers
	userHandler := httphandler.NewUserHTTPHandler(userAppService)

	// Setup router
	router := setupRouter(userHandler)

	// Get server configuration
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// initDatabase initializes the GORM database connection
func initDatabase(config *configs.Config) (*database.GORMConnectionManager, error) {
	// Create database factory
	dbFactory := database.NewFactory()

	// Create GORM connection manager
	dbManager, err := dbFactory.CreateGORMConnectionManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM connection manager: %w", err)
	}

	// Set log level based on environment
	if config.App.Env == "production" {
		dbManager.SetLogLevel(1) // Error level
	} else {
		dbManager.SetLogLevel(4) // Info level
	}

	return dbManager, nil
}

// setupRouter configures the HTTP router with middleware and routes
func setupRouter(userHandler *httphandler.UserHTTPHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/email", userHandler.GetUserByEmail)
			users.PUT("/:id", userHandler.UpdateUser)
			users.PATCH("/:id/activate", userHandler.ActivateUser)
			users.PATCH("/:id/deactivate", userHandler.DeactivateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}
