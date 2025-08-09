package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// LanguageService implements business logic for language operations
type LanguageService struct {
	languageRepo repositories.LanguageRepository
}

// NewLanguageService creates a new LanguageService instance
func NewLanguageService(languageRepo repositories.LanguageRepository) *LanguageService {
	return &LanguageService{
		languageRepo: languageRepo,
	}
}

// CreateLanguage creates a new language with validation
func (s *LanguageService) CreateLanguage(ctx context.Context, name, code string) (*entities.Language, error) {
	// Validate required fields
	if name == "" {
		return nil, fmt.Errorf("language name is required")
	}
	if code == "" {
		return nil, fmt.Errorf("language code is required")
	}
	if len(code) > 10 {
		return nil, fmt.Errorf("language code must be 10 characters or less")
	}

	// Check if language with the same code already exists
	exists, err := s.languageRepo.ExistsByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to check language code existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("language with code '%s' already exists", code)
	}

	// Check if language with the same name already exists
	exists, err = s.languageRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check language name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("language with name '%s' already exists", name)
	}

	language := entities.NewLanguage(name, code)

	if err := s.languageRepo.Create(ctx, language); err != nil {
		return nil, fmt.Errorf("failed to create language: %w", err)
	}

	return language, nil
}

// GetLanguageByID retrieves a language by ID
func (s *LanguageService) GetLanguageByID(ctx context.Context, id uuid.UUID) (*entities.Language, error) {
	return s.languageRepo.GetByID(ctx, id)
}

// GetLanguageByCode retrieves a language by code
func (s *LanguageService) GetLanguageByCode(ctx context.Context, code string) (*entities.Language, error) {
	return s.languageRepo.GetByCode(ctx, code)
}

// GetLanguageByName retrieves a language by name
func (s *LanguageService) GetLanguageByName(ctx context.Context, name string) (*entities.Language, error) {
	return s.languageRepo.GetByName(ctx, name)
}

// GetAllLanguages retrieves all languages with pagination
func (s *LanguageService) GetAllLanguages(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	return s.languageRepo.GetAll(ctx, limit, offset)
}

// GetActiveLanguages retrieves all active languages
func (s *LanguageService) GetActiveLanguages(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	return s.languageRepo.GetActive(ctx, limit, offset)
}

// GetInactiveLanguages retrieves all inactive languages
func (s *LanguageService) GetInactiveLanguages(ctx context.Context, limit, offset int) ([]*entities.Language, error) {
	return s.languageRepo.GetInactive(ctx, limit, offset)
}

// SearchLanguages searches languages by query
func (s *LanguageService) SearchLanguages(ctx context.Context, query string, limit, offset int) ([]*entities.Language, error) {
	return s.languageRepo.Search(ctx, query, limit, offset)
}

// UpdateLanguage updates an existing language
func (s *LanguageService) UpdateLanguage(ctx context.Context, language *entities.Language) error {
	if !language.IsValid() {
		return fmt.Errorf("invalid language data")
	}

	// Check if updating to a code that already exists (but not for the same language)
	existingLanguage, err := s.languageRepo.GetByCode(ctx, language.Code)
	if err == nil && existingLanguage.ID != language.ID {
		return fmt.Errorf("language with code '%s' already exists", language.Code)
	}

	// Check if updating to a name that already exists (but not for the same language)
	existingLanguage, err = s.languageRepo.GetByName(ctx, language.Name)
	if err == nil && existingLanguage.ID != language.ID {
		return fmt.Errorf("language with name '%s' already exists", language.Name)
	}

	return s.languageRepo.Update(ctx, language)
}

// ActivateLanguage activates a language
func (s *LanguageService) ActivateLanguage(ctx context.Context, id uuid.UUID) error {
	return s.languageRepo.Activate(ctx, id)
}

// DeactivateLanguage deactivates a language
func (s *LanguageService) DeactivateLanguage(ctx context.Context, id uuid.UUID) error {
	return s.languageRepo.Deactivate(ctx, id)
}

// DeleteLanguage deletes a language by ID
func (s *LanguageService) DeleteLanguage(ctx context.Context, id uuid.UUID) error {
	return s.languageRepo.Delete(ctx, id)
}

// CountLanguages returns the total number of languages
func (s *LanguageService) CountLanguages(ctx context.Context) (int64, error) {
	return s.languageRepo.Count(ctx)
}

// ValidateLanguage validates a language entity
func (s *LanguageService) ValidateLanguage(language *entities.Language) error {
	if language == nil {
		return fmt.Errorf("language cannot be nil")
	}

	if !language.IsValid() {
		return fmt.Errorf("language validation failed: name and code are required, code must be 10 characters or less")
	}

	return nil
}
