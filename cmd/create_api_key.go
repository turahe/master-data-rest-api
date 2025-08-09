package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database/pgx"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
)

var (
	keyName        string
	keyDescription string
	keyExpires     string
)

// createAPIKeyCmd represents the create-api-key command
var createAPIKeyCmd = &cobra.Command{
	Use:   "create-api-key",
	Short: "Create a new API key for authentication",
	Long: `Create a new API key for authentication to the Master Data REST API.

This command generates a secure API key that can be used to authenticate
requests to the API endpoints. The key can optionally have an expiration
date and a custom description.

Examples:
  # Create a basic API key
  master-data-api create-api-key

  # Create an API key with custom name and description
  master-data-api create-api-key --name "Production Key" --description "API key for production environment"

  # Create an API key with expiration (ISO 8601 format)
  master-data-api create-api-key --expires "2024-12-31T23:59:59Z"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}

		return createAPIKey()
	},
}

func init() {
	rootCmd.AddCommand(createAPIKeyCmd)

	// API key creation flags
	createAPIKeyCmd.Flags().StringVarP(&keyName, "name", "n", "Default API Key", "name for the API key")
	createAPIKeyCmd.Flags().StringVarP(&keyDescription, "description", "d", "API key for accessing Master Data REST API", "description for the API key")
	createAPIKeyCmd.Flags().StringVarP(&keyExpires, "expires", "e", "", "expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
}

func createAPIKey() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Starting API key creation utility")

	// Initialize database connection
	log.Info("Connecting to database")
	dbConnection := database.NewPgxConnectionWithLogger(config.Database, log)
	if err := dbConnection.Connect(); err != nil {
		log.WithError(err).Error("Failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbConnection.Close()

	log.Info("Database connected successfully")

	// Run migrations to ensure API key table exists
	log.Info("Running database migrations")
	migrator := database.NewMigrator(config.Database)
	if err := migrator.RunMigrations("migrations"); err != nil {
		log.WithError(err).Error("Failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Info("Database migrations completed successfully")

	// Initialize repository and service
	apiKeyRepo := pgx.NewAPIKeyRepository(dbConnection.GetPool())
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)

	// Parse expiration date if provided
	var expiresAt *time.Time
	if keyExpires != "" {
		parsed, err := time.Parse(time.RFC3339, keyExpires)
		if err != nil {
			log.WithError(err).Error("Invalid expiration date format")
			return fmt.Errorf("invalid expiration date format. Use ISO 8601 format (e.g., 2024-12-31T23:59:59Z): %w", err)
		}
		expiresAt = &parsed
	}

	// Create API key
	log.WithFields(map[string]interface{}{
		"name":        keyName,
		"description": keyDescription,
		"expires_at":  expiresAt,
	}).Info("Creating API key")

	apiKey, err := apiKeyService.CreateAPIKey(context.Background(), keyName, keyDescription, expiresAt)
	if err != nil {
		log.WithError(err).Error("Failed to create API key")
		return fmt.Errorf("failed to create API key: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"api_key_id":   apiKey.ID.String(),
		"api_key_name": apiKey.Name,
		"is_active":    apiKey.IsActive,
		"created_at":   apiKey.CreatedAt,
		"expires_at":   apiKey.ExpiresAt,
	}).Info("API key created successfully")

	// Display success information
	fmt.Println("✅ API Key created successfully!")
	fmt.Printf("📝 Name: %s\n", apiKey.Name)
	if apiKey.Description != nil {
		fmt.Printf("📄 Description: %s\n", *apiKey.Description)
	}
	fmt.Printf("🔑 API Key: %s\n", apiKey.Key)
	fmt.Printf("🆔 ID: %s\n", apiKey.ID.String())
	fmt.Printf("✅ Active: %v\n", apiKey.IsActive)
	fmt.Printf("📅 Created: %s\n", apiKey.CreatedAt.Format("2006-01-02 15:04:05"))
	if apiKey.ExpiresAt != nil {
		fmt.Printf("⏰ Expires: %s\n", apiKey.ExpiresAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("⏰ Expires: Never")
	}

	fmt.Println("\n🚀 You can now use this API key to authenticate requests:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:%s/api/v1/countries\n", apiKey.Key, config.Server.Port)

	return nil
}
