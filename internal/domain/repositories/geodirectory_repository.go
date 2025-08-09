package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// GeodirectoryRepository defines the interface for geodirectory data operations
type GeodirectoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, geodirectory *entities.Geodirectory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	Update(ctx context.Context, geodirectory *entities.Geodirectory) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Geodirectory, error)
	GetByName(ctx context.Context, name string) (*entities.Geodirectory, error)
	GetByCode(ctx context.Context, code string) (*entities.Geodirectory, error)
	GetByPostalCode(ctx context.Context, postalCode string) ([]*entities.Geodirectory, error)

	// Type-based queries
	GetByType(ctx context.Context, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error)
	GetCountries(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	GetProvinces(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	GetCities(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	GetDistricts(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	GetVillages(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)

	// Hierarchical operations
	GetChildren(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)
	GetChildrenByType(ctx context.Context, parentID uuid.UUID, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error)
	GetParent(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error)
	GetAncestors(ctx context.Context, id uuid.UUID) ([]*entities.Geodirectory, error)
	GetDescendants(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)
	GetSiblings(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)

	// Nested set model operations
	GetByNestedSetRange(ctx context.Context, left, right int, limit, offset int) ([]*entities.Geodirectory, error)
	UpdateNestedSetValues(ctx context.Context, id uuid.UUID, left, right, ordering int) error
	RebuildNestedSet(ctx context.Context) error
	MoveNode(ctx context.Context, nodeID, newParentID uuid.UUID) error

	// Geographic operations
	GetByCoordinates(ctx context.Context, latitude, longitude string, radius float64) ([]*entities.Geodirectory, error)
	GetNearby(ctx context.Context, id uuid.UUID, radius float64, limit int) ([]*entities.Geodirectory, error)

	// Country-specific operations (for backward compatibility)
	GetCountryByCode(ctx context.Context, code string) (*entities.Geodirectory, error)
	GetProvincesByCountry(ctx context.Context, countryID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)
	GetCitiesByProvince(ctx context.Context, provinceID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)
	GetDistrictsByCity(ctx context.Context, cityID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)
	GetVillagesByDistrict(ctx context.Context, districtID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error)

	// Administrative operations
	CountByType(ctx context.Context, geoType entities.GeoType) (int64, error)
	CountChildren(ctx context.Context, parentID uuid.UUID) (int64, error)
	HasChildren(ctx context.Context, id uuid.UUID) (bool, error)
	GetRoots(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
	GetLeaves(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error)
}
