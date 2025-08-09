package repositories

import (
	"context"

	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// SearchRepository defines the interface for search operations
type SearchRepository interface {
	// Index operations
	IndexBank(ctx context.Context, bank *entities.Bank) error
	IndexCurrency(ctx context.Context, currency *entities.Currency) error
	IndexLanguage(ctx context.Context, language *entities.Language) error
	IndexGeodirectory(ctx context.Context, geodirectory *entities.Geodirectory) error

	// Search operations
	SearchBanks(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error)
	SearchCurrencies(ctx context.Context, query string, limit, offset int) ([]*entities.Currency, error)
	SearchLanguages(ctx context.Context, query string, limit, offset int) ([]*entities.Language, error)
	SearchGeodirectories(ctx context.Context, query string, limit, offset int) ([]*entities.Geodirectory, error)

	// Bulk operations
	IndexAllBanks(ctx context.Context, banks []*entities.Bank) error
	IndexAllCurrencies(ctx context.Context, currencies []*entities.Currency) error
	IndexAllLanguages(ctx context.Context, languages []*entities.Language) error
	IndexAllGeodirectories(ctx context.Context, geodirectories []*entities.Geodirectory) error

	// Delete operations
	DeleteBankFromIndex(ctx context.Context, bankID string) error
	DeleteCurrencyFromIndex(ctx context.Context, currencyID string) error
	DeleteLanguageFromIndex(ctx context.Context, languageID string) error
	DeleteGeodirectoryFromIndex(ctx context.Context, geodirectoryID string) error

	// Index management
	CreateIndexes() error
	ReindexAll(ctx context.Context) error
	GetSearchStats() (map[string]interface{}, error)
}

// SearchResult represents a generic search result with metadata
type SearchResult struct {
	Query          string      `json:"query"`
	Results        interface{} `json:"results"`
	TotalHits      int64       `json:"total_hits"`
	ProcessingTime int64       `json:"processing_time_ms"`
	Limit          int         `json:"limit"`
	Offset         int         `json:"offset"`
}

// SearchFilters represents search filters for advanced queries
type SearchFilters struct {
	IsActive   *bool             `json:"is_active,omitempty"`
	Type       string            `json:"type,omitempty"`
	Country    string            `json:"country,omitempty"`
	Attributes []string          `json:"attributes,omitempty"`
	Sort       []string          `json:"sort,omitempty"`
	Facets     map[string]string `json:"facets,omitempty"`
}
