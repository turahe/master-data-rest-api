package search

import (
	"fmt"
	"log"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// MeilisearchClient wraps the Meilisearch client with our configuration
type MeilisearchClient struct {
	client meilisearch.ServiceManager
	config Config
}

// Config holds Meilisearch configuration
type Config struct {
	Host   string
	APIKey string
}

// NewMeilisearchClient creates a new Meilisearch client instance
func NewMeilisearchClient(config Config) *MeilisearchClient {
	var client meilisearch.ServiceManager
	if config.APIKey != "" {
		client = meilisearch.New(config.Host, meilisearch.WithAPIKey(config.APIKey))
	} else {
		client = meilisearch.New(config.Host)
	}

	return &MeilisearchClient{
		client: client,
		config: config,
	}
}

// GetClient returns the underlying Meilisearch client
func (m *MeilisearchClient) GetClient() meilisearch.ServiceManager {
	return m.client
}

// HealthCheck verifies the connection to Meilisearch
func (m *MeilisearchClient) HealthCheck() error {
	health, err := m.client.Health()
	if err != nil {
		return fmt.Errorf("meilisearch health check failed: %w", err)
	}

	if health.Status != "available" {
		return fmt.Errorf("meilisearch is not available, status: %s", health.Status)
	}

	log.Printf("Meilisearch health check passed - Status: %s", health.Status)
	return nil
}

// GetVersion returns Meilisearch server version
func (m *MeilisearchClient) GetVersion() (*meilisearch.Version, error) {
	return m.client.Version()
}

// CreateIndex creates a new index with the given name and primary key
func (m *MeilisearchClient) CreateIndex(indexName, primaryKey string) error {
	_, err := m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexName,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}

	log.Printf("Created Meilisearch index: %s", indexName)
	return nil
}

// GetIndex returns an index instance
func (m *MeilisearchClient) GetIndex(indexName string) meilisearch.IndexManager {
	return m.client.Index(indexName)
}

// DeleteIndex deletes an index
func (m *MeilisearchClient) DeleteIndex(indexName string) error {
	_, err := m.client.DeleteIndex(indexName)
	if err != nil {
		return fmt.Errorf("failed to delete index %s: %w", indexName, err)
	}

	log.Printf("Deleted Meilisearch index: %s", indexName)
	return nil
}

// ListIndexes returns all available indexes
func (m *MeilisearchClient) ListIndexes() ([]*meilisearch.IndexResult, error) {
	result, err := m.client.ListIndexes(&meilisearch.IndexesQuery{})
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}

	return result.Results, nil
}
