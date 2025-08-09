package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/search"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Manage search indexes",
	Long:  `Manage Meilisearch indexes and search functionality`,
}

var searchInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize search indexes",
	Long:  `Initialize all search indexes in Meilisearch with proper configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return initializeSearchIndexes()
	},
}

var searchReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all data",
	Long:  `Reindex all data from the database into Meilisearch search indexes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return reindexAllData()
	},
}

var searchStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show search statistics",
	Long:  `Display statistics about search indexes and performance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return showSearchStats()
	},
}

var searchHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check search service health",
	Long:  `Check the health status of the Meilisearch service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return checkSearchHealth()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchInitCmd)
	searchCmd.AddCommand(searchReindexCmd)
	searchCmd.AddCommand(searchStatsCmd)
	searchCmd.AddCommand(searchHealthCmd)
}

func initializeSearchIndexes() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Initializing search indexes...")

	// Create Meilisearch client
	meilisearchClient := search.NewMeilisearchClient(search.Config{
		Host:   config.Meilisearch.Host,
		APIKey: config.Meilisearch.APIKey,
	})

	// Create search repository
	searchRepo := search.NewMeilisearchRepository(meilisearchClient)

	// Create indexes
	if err := searchRepo.CreateIndexes(); err != nil {
		log.WithError(err).Error("Failed to create search indexes")
		return fmt.Errorf("failed to create search indexes: %w", err)
	}

	log.Info("✅ Search indexes initialized successfully!")
	fmt.Println("✅ Search indexes initialized successfully!")
	return nil
}

func reindexAllData() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Starting full reindex operation...")
	fmt.Println("🔄 Starting full reindex operation...")

	// Create Meilisearch client
	meilisearchClient := search.NewMeilisearchClient(search.Config{
		Host:   config.Meilisearch.Host,
		APIKey: config.Meilisearch.APIKey,
	})

	// Create search repository
	searchRepo := search.NewMeilisearchRepository(meilisearchClient)

	// Create database connection
	dbConnection := database.NewPgxConnection(config.Database)
	if err := dbConnection.Connect(); err != nil {
		log.WithError(err).Error("Failed to connect to database")
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbConnection.Close()

	// TODO: Implement actual data fetching and indexing
	// This would involve:
	// 1. Fetching all banks, currencies, languages, geodirectories from database
	// 2. Bulk indexing them into Meilisearch
	// For now, just show that the operation would complete

	ctx := context.Background()

	// Example: Index sample data (this should be replaced with actual database queries)
	log.Info("Reindexing banks...")
	fmt.Println("📊 Reindexing banks...")

	log.Info("Reindexing currencies...")
	fmt.Println("💰 Reindexing currencies...")

	log.Info("Reindexing languages...")
	fmt.Println("🗣️  Reindexing languages...")

	log.Info("Reindexing geodirectories...")
	fmt.Println("🌍 Reindexing geodirectories...")

	// Trigger reindex
	if err := searchRepo.ReindexAll(ctx); err != nil {
		log.WithError(err).Error("Failed to reindex data")
		return fmt.Errorf("failed to reindex data: %w", err)
	}

	log.Info("✅ Full reindex completed successfully!")
	fmt.Println("✅ Full reindex completed successfully!")
	return nil
}

func showSearchStats() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Fetching search statistics...")

	// Create Meilisearch client
	meilisearchClient := search.NewMeilisearchClient(search.Config{
		Host:   config.Meilisearch.Host,
		APIKey: config.Meilisearch.APIKey,
	})

	// Create search repository
	searchRepo := search.NewMeilisearchRepository(meilisearchClient)

	// Get statistics
	stats, err := searchRepo.GetSearchStats()
	if err != nil {
		log.WithError(err).Error("Failed to get search statistics")
		return fmt.Errorf("failed to get search statistics: %w", err)
	}

	fmt.Println("📊 Search Index Statistics")
	fmt.Println("==========================")

	for indexName, indexStats := range stats {
		fmt.Printf("\n📄 Index: %s\n", indexName)
		fmt.Printf("Stats: %+v\n", indexStats)
	}

	return nil
}

func checkSearchHealth() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Checking search service health...")

	// Create Meilisearch client
	meilisearchClient := search.NewMeilisearchClient(search.Config{
		Host:   config.Meilisearch.Host,
		APIKey: config.Meilisearch.APIKey,
	})

	// Check health
	if err := meilisearchClient.HealthCheck(); err != nil {
		log.WithError(err).Error("Search service health check failed")
		fmt.Printf("❌ Search service health check failed: %v\n", err)
		return err
	}

	// Get version info
	version, err := meilisearchClient.GetVersion()
	if err != nil {
		log.WithError(err).Warn("Failed to get version info")
	} else {
		fmt.Printf("✅ Meilisearch is healthy!\n")
		fmt.Printf("📋 Version: %s\n", version.PkgVersion)
		fmt.Printf("🔗 Host: %s\n", config.Meilisearch.Host)
	}

	return nil
}
