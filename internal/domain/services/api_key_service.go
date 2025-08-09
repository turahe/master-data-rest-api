package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// APIKeyService provides business logic for API key operations
type APIKeyService struct {
	apiKeyRepo repositories.APIKeyRepository
}

// NewAPIKeyService creates a new APIKeyService instance
func NewAPIKeyService(apiKeyRepo repositories.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// GenerateAPIKey generates a secure random API key
func (s *APIKeyService) GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateAPIKey creates a new API key
func (s *APIKeyService) CreateAPIKey(ctx context.Context, name, description string, expiresAt *time.Time) (*entities.APIKey, error) {
	key, err := s.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Check if key already exists (very unlikely but good to be safe)
	existingKey, err := s.apiKeyRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to check key existence: %w", err)
	}
	if existingKey != nil {
		// Regenerate key if it already exists
		return s.CreateAPIKey(ctx, name, description, expiresAt)
	}

	apiKey := entities.NewAPIKey(name, key)
	if description != "" {
		apiKey.SetDescription(description)
	}
	if expiresAt != nil {
		apiKey.SetExpiration(*expiresAt)
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKey, nil
}

// ValidateAPIKey validates an API key and returns the associated entity
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, key string) (*entities.APIKey, error) {
	return s.apiKeyRepo.ValidateKey(ctx, key)
}

// GetAPIKeyByID retrieves an API key by its ID
func (s *APIKeyService) GetAPIKeyByID(ctx context.Context, id uuid.UUID) (*entities.APIKey, error) {
	apiKey, err := s.apiKeyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	return apiKey, nil
}

// GetAllAPIKeys retrieves all API keys with pagination
func (s *APIKeyService) GetAllAPIKeys(ctx context.Context, limit, offset int) ([]*entities.APIKey, error) {
	apiKeys, err := s.apiKeyRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get API keys: %w", err)
	}
	return apiKeys, nil
}

// UpdateAPIKey updates an existing API key
func (s *APIKeyService) UpdateAPIKey(ctx context.Context, apiKey *entities.APIKey) error {
	if err := s.apiKeyRepo.Update(ctx, apiKey); err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}
	return nil
}

// ActivateAPIKey activates an API key
func (s *APIKeyService) ActivateAPIKey(ctx context.Context, id uuid.UUID) error {
	if err := s.apiKeyRepo.Activate(ctx, id); err != nil {
		return fmt.Errorf("failed to activate API key: %w", err)
	}
	return nil
}

// DeactivateAPIKey deactivates an API key
func (s *APIKeyService) DeactivateAPIKey(ctx context.Context, id uuid.UUID) error {
	if err := s.apiKeyRepo.Deactivate(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate API key: %w", err)
	}
	return nil
}

// DeleteAPIKey soft deletes an API key
func (s *APIKeyService) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	if err := s.apiKeyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

// CountAPIKeys returns the total number of API keys
func (s *APIKeyService) CountAPIKeys(ctx context.Context) (int64, error) {
	count, err := s.apiKeyRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count API keys: %w", err)
	}
	return count, nil
}

// SearchAPIKeys searches API keys by query
func (s *APIKeyService) SearchAPIKeys(ctx context.Context, query string, limit, offset int) ([]*entities.APIKey, error) {
	apiKeys, err := s.apiKeyRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search API keys: %w", err)
	}
	return apiKeys, nil
}

// ValidateAPIKeyByName validates an API key by checking if it exists and is valid
func (s *APIKeyService) ValidateAPIKeyByName(ctx context.Context, name string) (*entities.APIKey, error) {
	apiKeys, err := s.apiKeyRepo.Search(ctx, name, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to search API key by name: %w", err)
	}

	if len(apiKeys) == 0 {
		return nil, fmt.Errorf("API key not found")
	}

	apiKey := apiKeys[0]
	if !apiKey.IsValid() {
		return nil, fmt.Errorf("API key is not valid")
	}

	return apiKey, nil
}

// GetAPIKeyStatistics returns statistics about API keys
func (s *APIKeyService) GetAPIKeyStatistics(ctx context.Context) (map[string]interface{}, error) {
	total, err := s.CountAPIKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// For now, return basic statistics
	// In the future, you could add more detailed statistics
	stats := map[string]interface{}{
		"total_keys": total,
	}

	return stats, nil
}
