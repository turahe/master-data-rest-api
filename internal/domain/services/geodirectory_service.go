package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// GeodirectoryService implements business logic for geodirectory operations
type GeodirectoryService struct {
	geodirectoryRepo repositories.GeodirectoryRepository
}

// NewGeodirectoryService creates a new GeodirectoryService instance
func NewGeodirectoryService(geodirectoryRepo repositories.GeodirectoryRepository) *GeodirectoryService {
	return &GeodirectoryService{
		geodirectoryRepo: geodirectoryRepo,
	}
}

// CreateGeodirectory creates a new geodirectory with validation
func (s *GeodirectoryService) CreateGeodirectory(
	ctx context.Context,
	name string,
	geoType entities.GeoType,
	parentID *uuid.UUID,
) (*entities.Geodirectory, error) {
	// Create the geodirectory entity
	geodirectory := entities.NewGeodirectory(name, geoType)
	geodirectory.ParentID = parentID

	// Validate the geodirectory type
	if !geodirectory.ValidateType() {
		return nil, fmt.Errorf("invalid geodirectory type: %s", geoType)
	}

	// Validate parent-child relationship if parent is specified
	if parentID != nil {
		parent, err := s.geodirectoryRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("parent geodirectory not found: %w", err)
		}

		if !geodirectory.CanHaveParentType(parent.Type) {
			return nil, fmt.Errorf("geodirectory type %s cannot have parent type %s", geoType, parent.Type)
		}
	}

	// Create the geodirectory using nested set model
	if err := s.geodirectoryRepo.Create(ctx, geodirectory); err != nil {
		return nil, fmt.Errorf("failed to create geodirectory: %w", err)
	}

	return geodirectory, nil
}

// GetGeodirectoryByID retrieves a geodirectory by ID
func (s *GeodirectoryService) GetGeodirectoryByID(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetByID(ctx, id)
}

// GetGeodirectoryByCode retrieves a geodirectory by code
func (s *GeodirectoryService) GetGeodirectoryByCode(ctx context.Context, code string) (*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetByCode(ctx, code)
}

// GetGeodirectoryByName retrieves a geodirectory by name
func (s *GeodirectoryService) GetGeodirectoryByName(ctx context.Context, name string) (*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetByName(ctx, name)
}

// GetAllGeodirectories retrieves all geodirectories with pagination
func (s *GeodirectoryService) GetAllGeodirectories(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetAll(ctx, limit, offset)
}

// SearchGeodirectories searches geodirectories by query
func (s *GeodirectoryService) SearchGeodirectories(ctx context.Context, query string, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.Search(ctx, query, limit, offset)
}

// UpdateGeodirectory updates an existing geodirectory
func (s *GeodirectoryService) UpdateGeodirectory(ctx context.Context, geodirectory *entities.Geodirectory) error {
	// Validate the geodirectory type
	if !geodirectory.ValidateType() {
		return fmt.Errorf("invalid geodirectory type: %s", geodirectory.Type)
	}

	return s.geodirectoryRepo.Update(ctx, geodirectory)
}

// DeleteGeodirectory deletes a geodirectory and all its descendants
func (s *GeodirectoryService) DeleteGeodirectory(ctx context.Context, id uuid.UUID) error {
	return s.geodirectoryRepo.Delete(ctx, id)
}

// CountGeodirectories returns the total number of geodirectories
func (s *GeodirectoryService) CountGeodirectories(ctx context.Context) (int64, error) {
	return s.geodirectoryRepo.Count(ctx)
}

// GetGeodirectoriesByType retrieves geodirectories by type
func (s *GeodirectoryService) GetGeodirectoriesByType(ctx context.Context, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetByType(ctx, geoType, limit, offset)
}

// GetCountries retrieves all countries
func (s *GeodirectoryService) GetCountries(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetCountries(ctx, limit, offset)
}

// GetProvinces retrieves all provinces
func (s *GeodirectoryService) GetProvinces(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetProvinces(ctx, limit, offset)
}

// GetCities retrieves all cities
func (s *GeodirectoryService) GetCities(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetCities(ctx, limit, offset)
}

// GetDistricts retrieves all districts
func (s *GeodirectoryService) GetDistricts(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetDistricts(ctx, limit, offset)
}

// GetVillages retrieves all villages
func (s *GeodirectoryService) GetVillages(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetVillages(ctx, limit, offset)
}

// Hierarchical operations

// GetChildren retrieves children of a geodirectory
func (s *GeodirectoryService) GetChildren(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetChildren(ctx, parentID, limit, offset)
}

// GetChildrenByType retrieves children of a specific type
func (s *GeodirectoryService) GetChildrenByType(ctx context.Context, parentID uuid.UUID, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetChildrenByType(ctx, parentID, geoType, limit, offset)
}

// GetParent retrieves the parent of a geodirectory
func (s *GeodirectoryService) GetParent(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetParent(ctx, id)
}

// GetAncestors retrieves all ancestors of a geodirectory
func (s *GeodirectoryService) GetAncestors(ctx context.Context, id uuid.UUID) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetAncestors(ctx, id)
}

// GetDescendants retrieves all descendants of a geodirectory
func (s *GeodirectoryService) GetDescendants(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetDescendants(ctx, id, limit, offset)
}

// GetSiblings retrieves siblings of a geodirectory
func (s *GeodirectoryService) GetSiblings(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetSiblings(ctx, id, limit, offset)
}

// GetRoots retrieves all root geodirectories
func (s *GeodirectoryService) GetRoots(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetRoots(ctx, limit, offset)
}

// GetLeaves retrieves all leaf geodirectories
func (s *GeodirectoryService) GetLeaves(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetLeaves(ctx, limit, offset)
}

// Country-specific operations (for backward compatibility)

// GetCountryByCode retrieves a country by its code
func (s *GeodirectoryService) GetCountryByCode(ctx context.Context, code string) (*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetCountryByCode(ctx, code)
}

// GetProvincesByCountry retrieves provinces of a country
func (s *GeodirectoryService) GetProvincesByCountry(ctx context.Context, countryID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetProvincesByCountry(ctx, countryID, limit, offset)
}

// GetCitiesByProvince retrieves cities of a province
func (s *GeodirectoryService) GetCitiesByProvince(ctx context.Context, provinceID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetCitiesByProvince(ctx, provinceID, limit, offset)
}

// GetDistrictsByCity retrieves districts of a city
func (s *GeodirectoryService) GetDistrictsByCity(ctx context.Context, cityID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetDistrictsByCity(ctx, cityID, limit, offset)
}

// GetVillagesByDistrict retrieves villages of a district
func (s *GeodirectoryService) GetVillagesByDistrict(ctx context.Context, districtID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return s.geodirectoryRepo.GetVillagesByDistrict(ctx, districtID, limit, offset)
}

// Administrative operations

// CountByType returns the count of geodirectories by type
func (s *GeodirectoryService) CountByType(ctx context.Context, geoType entities.GeoType) (int64, error) {
	return s.geodirectoryRepo.CountByType(ctx, geoType)
}

// CountChildren returns the count of children for a geodirectory
func (s *GeodirectoryService) CountChildren(ctx context.Context, parentID uuid.UUID) (int64, error) {
	return s.geodirectoryRepo.CountChildren(ctx, parentID)
}

// HasChildren checks if a geodirectory has children
func (s *GeodirectoryService) HasChildren(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.geodirectoryRepo.HasChildren(ctx, id)
}

// Advanced nested set operations

// MoveGeodirectory moves a geodirectory to a new parent
func (s *GeodirectoryService) MoveGeodirectory(ctx context.Context, nodeID, newParentID uuid.UUID) error {
	// Validate that the new parent exists
	parent, err := s.geodirectoryRepo.GetByID(ctx, newParentID)
	if err != nil {
		return fmt.Errorf("new parent not found: %w", err)
	}

	// Get the node being moved
	node, err := s.geodirectoryRepo.GetByID(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("node to move not found: %w", err)
	}

	// Validate parent-child relationship
	if !node.CanHaveParentType(parent.Type) {
		return fmt.Errorf("geodirectory type %s cannot have parent type %s", node.Type, parent.Type)
	}

	// Check if the new parent is not a descendant of the node being moved
	descendants, err := s.geodirectoryRepo.GetDescendants(ctx, nodeID, 1000, 0) // Get reasonable limit
	if err != nil {
		return fmt.Errorf("failed to check descendants: %w", err)
	}

	for _, descendant := range descendants {
		if descendant.ID == newParentID {
			return fmt.Errorf("cannot move node to its own descendant")
		}
	}

	// Perform the move operation
	return s.geodirectoryRepo.MoveNode(ctx, nodeID, newParentID)
}

// RebuildNestedSet rebuilds the nested set structure
func (s *GeodirectoryService) RebuildNestedSet(ctx context.Context) error {
	return s.geodirectoryRepo.RebuildNestedSet(ctx)
}

// GetGeodirectoryWithHierarchy retrieves a geodirectory with its parent and children
func (s *GeodirectoryService) GetGeodirectoryWithHierarchy(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error) {
	// Get the main geodirectory
	geodirectory, err := s.geodirectoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get parent if exists
	if geodirectory.ParentID != nil {
		parent, err := s.geodirectoryRepo.GetByID(ctx, *geodirectory.ParentID)
		if err == nil {
			geodirectory.Parent = parent
		}
	}

	// Get children
	children, err := s.geodirectoryRepo.GetChildren(ctx, id, 100, 0) // Reasonable limit
	if err == nil {
		geodirectory.Children = children
	}

	return geodirectory, nil
}

// ValidateHierarchy validates the hierarchical structure of geodirectories
func (s *GeodirectoryService) ValidateHierarchy(ctx context.Context) ([]string, error) {
	var errors []string

	// Get all geodirectories
	geodirectories, err := s.geodirectoryRepo.GetAll(ctx, 10000, 0) // Large limit to get all
	if err != nil {
		return nil, fmt.Errorf("failed to get geodirectories: %w", err)
	}

	// Validate each geodirectory
	for _, geo := range geodirectories {
		// Check if parent-child relationship is valid
		if geo.ParentID != nil {
			parent, err := s.geodirectoryRepo.GetByID(ctx, *geo.ParentID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Geodirectory %s has invalid parent ID %s", geo.Name, geo.ParentID.String()))
				continue
			}

			if !geo.CanHaveParentType(parent.Type) {
				errors = append(errors, fmt.Sprintf("Geodirectory %s (type %s) cannot have parent %s (type %s)", geo.Name, geo.Type, parent.Name, parent.Type))
			}
		}

		// Validate nested set values
		if geo.RecordLeft != nil && geo.RecordRight != nil {
			if *geo.RecordLeft >= *geo.RecordRight {
				errors = append(errors, fmt.Sprintf("Geodirectory %s has invalid nested set values: left=%d, right=%d", geo.Name, *geo.RecordLeft, *geo.RecordRight))
			}
		}
	}

	return errors, nil
}
