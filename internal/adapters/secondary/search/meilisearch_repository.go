package search

import (
	"context"
	"fmt"
	"log"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// Index names constants
const (
	BanksIndex          = "banks"
	CurrenciesIndex     = "currencies"
	LanguagesIndex      = "languages"
	GeodirectoriesIndex = "geodirectories"
)

// MeilisearchRepository implements SearchRepository using Meilisearch
type MeilisearchRepository struct {
	client *MeilisearchClient
}

// NewMeilisearchRepository creates a new Meilisearch repository instance
func NewMeilisearchRepository(client *MeilisearchClient) repositories.SearchRepository {
	return &MeilisearchRepository{
		client: client,
	}
}

// CreateIndexes creates all required indexes with proper configuration
func (r *MeilisearchRepository) CreateIndexes() error {
	indexes := []struct {
		name       string
		primaryKey string
	}{
		{BanksIndex, "id"},
		{CurrenciesIndex, "id"},
		{LanguagesIndex, "id"},
		{GeodirectoriesIndex, "id"},
	}

	for _, idx := range indexes {
		if err := r.client.CreateIndex(idx.name, idx.primaryKey); err != nil {
			// Index might already exist, log but continue
			log.Printf("Index creation warning for %s: %v", idx.name, err)
		}
	}

	// Configure search settings for each index
	if err := r.configureIndexSettings(); err != nil {
		return fmt.Errorf("failed to configure index settings: %w", err)
	}

	return nil
}

// configureIndexSettings sets up searchable attributes and other settings
func (r *MeilisearchRepository) configureIndexSettings() error {
	// Banks index settings
	banksIndex := r.client.GetIndex(BanksIndex)
	searchableAttrs := []string{"name", "alias", "company", "code"}
	_, err := banksIndex.UpdateSearchableAttributes(&searchableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update banks searchable attributes: %w", err)
	}

	filterableAttrs := []interface{}{"code", "company"}
	_, err = banksIndex.UpdateFilterableAttributes(&filterableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update banks filterable attributes: %w", err)
	}

	// Currencies index settings
	currenciesIndex := r.client.GetIndex(CurrenciesIndex)
	currencySearchableAttrs := []string{"name", "code", "symbol"}
	_, err = currenciesIndex.UpdateSearchableAttributes(&currencySearchableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update currencies searchable attributes: %w", err)
	}

	currencyFilterableAttrs := []interface{}{"code", "is_active", "decimal_places"}
	_, err = currenciesIndex.UpdateFilterableAttributes(&currencyFilterableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update currencies filterable attributes: %w", err)
	}

	// Languages index settings
	languagesIndex := r.client.GetIndex(LanguagesIndex)
	languageSearchableAttrs := []string{"name", "code"}
	_, err = languagesIndex.UpdateSearchableAttributes(&languageSearchableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update languages searchable attributes: %w", err)
	}

	languageFilterableAttrs := []interface{}{"code", "is_active"}
	_, err = languagesIndex.UpdateFilterableAttributes(&languageFilterableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update languages filterable attributes: %w", err)
	}

	// Geodirectories index settings
	geoIndex := r.client.GetIndex(GeodirectoriesIndex)
	geoSearchableAttrs := []string{"name", "code", "type", "postal_code"}
	_, err = geoIndex.UpdateSearchableAttributes(&geoSearchableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update geodirectories searchable attributes: %w", err)
	}

	geoFilterableAttrs := []interface{}{"type", "code", "parent_id"}
	_, err = geoIndex.UpdateFilterableAttributes(&geoFilterableAttrs)
	if err != nil {
		return fmt.Errorf("failed to update geodirectories filterable attributes: %w", err)
	}

	log.Println("Meilisearch index settings configured successfully")
	return nil
}

// IndexBank adds or updates a bank in the search index
func (r *MeilisearchRepository) IndexBank(ctx context.Context, bank *entities.Bank) error {
	index := r.client.GetIndex(BanksIndex)
	_, err := index.AddDocuments([]interface{}{bank}, nil)
	if err != nil {
		return fmt.Errorf("failed to index bank %s: %w", bank.ID, err)
	}
	return nil
}

// IndexCurrency adds or updates a currency in the search index
func (r *MeilisearchRepository) IndexCurrency(ctx context.Context, currency *entities.Currency) error {
	index := r.client.GetIndex(CurrenciesIndex)
	_, err := index.AddDocuments([]interface{}{currency}, nil)
	if err != nil {
		return fmt.Errorf("failed to index currency %s: %w", currency.ID, err)
	}
	return nil
}

// IndexLanguage adds or updates a language in the search index
func (r *MeilisearchRepository) IndexLanguage(ctx context.Context, language *entities.Language) error {
	index := r.client.GetIndex(LanguagesIndex)
	_, err := index.AddDocuments([]interface{}{language}, nil)
	if err != nil {
		return fmt.Errorf("failed to index language %s: %w", language.ID, err)
	}
	return nil
}

// IndexGeodirectory adds or updates a geodirectory in the search index
func (r *MeilisearchRepository) IndexGeodirectory(ctx context.Context, geodirectory *entities.Geodirectory) error {
	index := r.client.GetIndex(GeodirectoriesIndex)
	_, err := index.AddDocuments([]interface{}{geodirectory}, nil)
	if err != nil {
		return fmt.Errorf("failed to index geodirectory %s: %w", geodirectory.ID, err)
	}
	return nil
}

// SearchBanks performs a search query on banks index
func (r *MeilisearchRepository) SearchBanks(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error) {
	index := r.client.GetIndex(BanksIndex)
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	result, err := index.Search(query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search banks: %w", err)
	}

	var banks []*entities.Bank
	for _, hit := range result.Hits {
		bank := &entities.Bank{}
		if err := mapToStruct(hit, bank); err != nil {
			log.Printf("Failed to map search hit to bank: %v", err)
			continue
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// SearchCurrencies performs a search query on currencies index
func (r *MeilisearchRepository) SearchCurrencies(ctx context.Context, query string, limit, offset int) ([]*entities.Currency, error) {
	index := r.client.GetIndex(CurrenciesIndex)
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	result, err := index.Search(query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search currencies: %w", err)
	}

	var currencies []*entities.Currency
	for _, hit := range result.Hits {
		currency := &entities.Currency{}
		if err := mapToStruct(hit, currency); err != nil {
			log.Printf("Failed to map search hit to currency: %v", err)
			continue
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// SearchLanguages performs a search query on languages index
func (r *MeilisearchRepository) SearchLanguages(ctx context.Context, query string, limit, offset int) ([]*entities.Language, error) {
	index := r.client.GetIndex(LanguagesIndex)
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	result, err := index.Search(query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search languages: %w", err)
	}

	var languages []*entities.Language
	for _, hit := range result.Hits {
		language := &entities.Language{}
		if err := mapToStruct(hit, language); err != nil {
			log.Printf("Failed to map search hit to language: %v", err)
			continue
		}
		languages = append(languages, language)
	}

	return languages, nil
}

// SearchGeodirectories performs a search query on geodirectories index
func (r *MeilisearchRepository) SearchGeodirectories(ctx context.Context, query string, limit, offset int) ([]*entities.Geodirectory, error) {
	index := r.client.GetIndex(GeodirectoriesIndex)
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	result, err := index.Search(query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search geodirectories: %w", err)
	}

	var geodirectories []*entities.Geodirectory
	for _, hit := range result.Hits {
		geodirectory := &entities.Geodirectory{}
		if err := mapToStruct(hit, geodirectory); err != nil {
			log.Printf("Failed to map search hit to geodirectory: %v", err)
			continue
		}
		geodirectories = append(geodirectories, geodirectory)
	}

	return geodirectories, nil
}

// IndexAllBanks bulk indexes all banks
func (r *MeilisearchRepository) IndexAllBanks(ctx context.Context, banks []*entities.Bank) error {
	if len(banks) == 0 {
		return nil
	}

	index := r.client.GetIndex(BanksIndex)
	documents := make([]interface{}, len(banks))
	for i, bank := range banks {
		documents[i] = bank
	}

	_, err := index.AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to bulk index banks: %w", err)
	}

	log.Printf("Successfully indexed %d banks", len(banks))
	return nil
}

// IndexAllCurrencies bulk indexes all currencies
func (r *MeilisearchRepository) IndexAllCurrencies(ctx context.Context, currencies []*entities.Currency) error {
	if len(currencies) == 0 {
		return nil
	}

	index := r.client.GetIndex(CurrenciesIndex)
	documents := make([]interface{}, len(currencies))
	for i, currency := range currencies {
		documents[i] = currency
	}

	_, err := index.AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to bulk index currencies: %w", err)
	}

	log.Printf("Successfully indexed %d currencies", len(currencies))
	return nil
}

// IndexAllLanguages bulk indexes all languages
func (r *MeilisearchRepository) IndexAllLanguages(ctx context.Context, languages []*entities.Language) error {
	if len(languages) == 0 {
		return nil
	}

	index := r.client.GetIndex(LanguagesIndex)
	documents := make([]interface{}, len(languages))
	for i, language := range languages {
		documents[i] = language
	}

	_, err := index.AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to bulk index languages: %w", err)
	}

	log.Printf("Successfully indexed %d languages", len(languages))
	return nil
}

// IndexAllGeodirectories bulk indexes all geodirectories
func (r *MeilisearchRepository) IndexAllGeodirectories(ctx context.Context, geodirectories []*entities.Geodirectory) error {
	if len(geodirectories) == 0 {
		return nil
	}

	index := r.client.GetIndex(GeodirectoriesIndex)
	documents := make([]interface{}, len(geodirectories))
	for i, geodirectory := range geodirectories {
		documents[i] = geodirectory
	}

	_, err := index.AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to bulk index geodirectories: %w", err)
	}

	log.Printf("Successfully indexed %d geodirectories", len(geodirectories))
	return nil
}

// DeleteBankFromIndex removes a bank from the search index
func (r *MeilisearchRepository) DeleteBankFromIndex(ctx context.Context, bankID string) error {
	index := r.client.GetIndex(BanksIndex)
	_, err := index.DeleteDocument(bankID)
	if err != nil {
		return fmt.Errorf("failed to delete bank %s from index: %w", bankID, err)
	}
	return nil
}

// DeleteCurrencyFromIndex removes a currency from the search index
func (r *MeilisearchRepository) DeleteCurrencyFromIndex(ctx context.Context, currencyID string) error {
	index := r.client.GetIndex(CurrenciesIndex)
	_, err := index.DeleteDocument(currencyID)
	if err != nil {
		return fmt.Errorf("failed to delete currency %s from index: %w", currencyID, err)
	}
	return nil
}

// DeleteLanguageFromIndex removes a language from the search index
func (r *MeilisearchRepository) DeleteLanguageFromIndex(ctx context.Context, languageID string) error {
	index := r.client.GetIndex(LanguagesIndex)
	_, err := index.DeleteDocument(languageID)
	if err != nil {
		return fmt.Errorf("failed to delete language %s from index: %w", languageID, err)
	}
	return nil
}

// DeleteGeodirectoryFromIndex removes a geodirectory from the search index
func (r *MeilisearchRepository) DeleteGeodirectoryFromIndex(ctx context.Context, geodirectoryID string) error {
	index := r.client.GetIndex(GeodirectoriesIndex)
	_, err := index.DeleteDocument(geodirectoryID)
	if err != nil {
		return fmt.Errorf("failed to delete geodirectory %s from index: %w", geodirectoryID, err)
	}
	return nil
}

// ReindexAll reindexes all data from the database
func (r *MeilisearchRepository) ReindexAll(ctx context.Context) error {
	// This would typically fetch all data from database repositories
	// and reindex everything. Implementation depends on how you want to
	// coordinate with your existing repositories.
	log.Println("Full reindex operation initiated")
	return nil
}

// GetSearchStats returns statistics about search indexes
func (r *MeilisearchRepository) GetSearchStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	indexes, err := r.client.ListIndexes()
	if err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}

	for _, index := range indexes {
		indexStats, err := index.GetStats()
		if err != nil {
			log.Printf("Failed to get stats for index %s: %v", index.UID, err)
			continue
		}
		stats[index.UID] = indexStats
	}

	return stats, nil
}
