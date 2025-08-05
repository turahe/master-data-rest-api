package services

import (
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

// CreateLanguage creates a new language
func (s *LanguageService) CreateLanguage(code, name, native string) (*entities.Language, error) {
	language := entities.NewLanguage(code, name, native)

	if err := s.languageRepo.Create(language); err != nil {
		return nil, err
	}

	return language, nil
}

// GetLanguageByID retrieves a language by ID
func (s *LanguageService) GetLanguageByID(id uint) (*entities.Language, error) {
	return s.languageRepo.GetByID(id)
}

// GetLanguageByCode retrieves a language by code
func (s *LanguageService) GetLanguageByCode(code string) (*entities.Language, error) {
	return s.languageRepo.GetByCode(code)
}

// GetLanguageByName retrieves a language by name
func (s *LanguageService) GetLanguageByName(name string) (*entities.Language, error) {
	return s.languageRepo.GetByName(name)
}

// GetLanguageByNative retrieves a language by native name
func (s *LanguageService) GetLanguageByNative(native string) (*entities.Language, error) {
	return s.languageRepo.GetByNative(native)
}

// GetAllLanguages retrieves all languages
func (s *LanguageService) GetAllLanguages() ([]*entities.Language, error) {
	return s.languageRepo.GetAll()
}

// UpdateLanguage updates a language
func (s *LanguageService) UpdateLanguage(language *entities.Language) error {
	return s.languageRepo.Update(language)
}

// DeleteLanguage deletes a language by ID
func (s *LanguageService) DeleteLanguage(id uint) error {
	return s.languageRepo.Delete(id)
}

// LanguageExists checks if a language exists by ID
func (s *LanguageService) LanguageExists(id uint) (bool, error) {
	return s.languageRepo.Exists(id)
}

// GetLanguageCount returns the total number of languages
func (s *LanguageService) GetLanguageCount() (int64, error) {
	return s.languageRepo.Count()
}
