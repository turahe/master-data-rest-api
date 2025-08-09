package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// GeodirectoryRepository implements the GeodirectoryRepository interface using pgx
type GeodirectoryRepository struct {
	pool *pgxpool.Pool
}

// NewGeodirectoryRepository creates a new GeodirectoryRepository instance
func NewGeodirectoryRepository(pool *pgxpool.Pool) *GeodirectoryRepository {
	return &GeodirectoryRepository{
		pool: pool,
	}
}

// Create creates a new geodirectory in the database using nested set model
func (r *GeodirectoryRepository) Create(ctx context.Context, geodirectory *entities.Geodirectory) error {
	return r.InsertNode(ctx, geodirectory, geodirectory.ParentID)
}

// GetByID retrieves a geodirectory by its ID
func (r *GeodirectoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE id = $1`

	var geodirectory entities.Geodirectory
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&geodirectory.ID, &geodirectory.Name, &geodirectory.Type, &geodirectory.Code,
		&geodirectory.PostalCode, &geodirectory.Longitude, &geodirectory.Latitude,
		&geodirectory.RecordLeft, &geodirectory.RecordRight, &geodirectory.RecordOrdering, &geodirectory.RecordDepth,
		&geodirectory.ParentID, &geodirectory.CreatedAt, &geodirectory.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("geodirectory not found")
		}
		return nil, err
	}

	return &geodirectory, nil
}

// GetAll retrieves all geodirectories with optional pagination
func (r *GeodirectoryRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		ORDER BY record_ordering, name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// Update updates an existing geodirectory
func (r *GeodirectoryRepository) Update(ctx context.Context, geodirectory *entities.Geodirectory) error {
	geodirectory.UpdatedAt = time.Now()

	query := `
		UPDATE tm_geodirectories SET
			name = $2, type = $3, code = $4, postal_code = $5, longitude = $6, latitude = $7,
			record_left = $8, record_right = $9, record_ordering = $10, record_depth = $11, parent_id = $12, updated_at = $13
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		geodirectory.ID, geodirectory.Name, geodirectory.Type, geodirectory.Code,
		geodirectory.PostalCode, geodirectory.Longitude, geodirectory.Latitude,
		geodirectory.RecordLeft, geodirectory.RecordRight, geodirectory.RecordOrdering, geodirectory.RecordDepth,
		geodirectory.ParentID, geodirectory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("geodirectory not found")
	}

	return nil
}

// Delete deletes a geodirectory by ID using nested set model
func (r *GeodirectoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteNode(ctx, id)
}

// Count returns the total number of geodirectories
func (r *GeodirectoryRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_geodirectories"

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Search searches geodirectories by name
func (r *GeodirectoryRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Geodirectory, error) {
	searchQuery := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE name ILIKE $1 OR code ILIKE $1 OR postal_code ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	rows, err := r.pool.Query(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetByName retrieves a geodirectory by name
func (r *GeodirectoryRepository) GetByName(ctx context.Context, name string) (*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE name = $1`

	var geodirectory entities.Geodirectory
	row := r.pool.QueryRow(ctx, query, name)

	err := row.Scan(
		&geodirectory.ID, &geodirectory.Name, &geodirectory.Type, &geodirectory.Code,
		&geodirectory.PostalCode, &geodirectory.Longitude, &geodirectory.Latitude,
		&geodirectory.RecordLeft, &geodirectory.RecordRight, &geodirectory.RecordOrdering, &geodirectory.RecordDepth,
		&geodirectory.ParentID, &geodirectory.CreatedAt, &geodirectory.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("geodirectory not found")
		}
		return nil, err
	}

	return &geodirectory, nil
}

// GetByCode retrieves a geodirectory by code
func (r *GeodirectoryRepository) GetByCode(ctx context.Context, code string) (*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE code = $1`

	var geodirectory entities.Geodirectory
	row := r.pool.QueryRow(ctx, query, code)

	err := row.Scan(
		&geodirectory.ID, &geodirectory.Name, &geodirectory.Type, &geodirectory.Code,
		&geodirectory.PostalCode, &geodirectory.Longitude, &geodirectory.Latitude,
		&geodirectory.RecordLeft, &geodirectory.RecordRight, &geodirectory.RecordOrdering, &geodirectory.RecordDepth,
		&geodirectory.ParentID, &geodirectory.CreatedAt, &geodirectory.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("geodirectory not found")
		}
		return nil, err
	}

	return &geodirectory, nil
}

// GetByPostalCode retrieves geodirectories by postal code
func (r *GeodirectoryRepository) GetByPostalCode(ctx context.Context, postalCode string) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE postal_code = $1
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query, postalCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetByType retrieves geodirectories by type
func (r *GeodirectoryRepository) GetByType(ctx context.Context, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE type = $1
		ORDER BY name
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, geoType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetCountries retrieves all countries
func (r *GeodirectoryRepository) GetCountries(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetByType(ctx, entities.GeoTypeCountry, limit, offset)
}

// GetProvinces retrieves all provinces
func (r *GeodirectoryRepository) GetProvinces(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetByType(ctx, entities.GeoTypeProvince, limit, offset)
}

// GetCities retrieves all cities
func (r *GeodirectoryRepository) GetCities(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetByType(ctx, entities.GeoTypeCity, limit, offset)
}

// GetDistricts retrieves all districts
func (r *GeodirectoryRepository) GetDistricts(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetByType(ctx, entities.GeoTypeDistrict, limit, offset)
}

// GetVillages retrieves all villages
func (r *GeodirectoryRepository) GetVillages(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetByType(ctx, entities.GeoTypeVillage, limit, offset)
}

// GetChildren retrieves children of a geodirectory
func (r *GeodirectoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE parent_id = $1
		ORDER BY record_ordering, name
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetChildrenByType retrieves children of a specific type
func (r *GeodirectoryRepository) GetChildrenByType(ctx context.Context, parentID uuid.UUID, geoType entities.GeoType, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE parent_id = $1 AND type = $2
		ORDER BY record_ordering, name
		LIMIT $3 OFFSET $4`

	rows, err := r.pool.Query(ctx, query, parentID, geoType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetCountryByCode retrieves a country by its code
func (r *GeodirectoryRepository) GetCountryByCode(ctx context.Context, code string) (*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE code = $1 AND type = 'COUNTRY'`

	var geodirectory entities.Geodirectory
	row := r.pool.QueryRow(ctx, query, code)

	err := row.Scan(
		&geodirectory.ID, &geodirectory.Name, &geodirectory.Type, &geodirectory.Code,
		&geodirectory.PostalCode, &geodirectory.Longitude, &geodirectory.Latitude,
		&geodirectory.RecordLeft, &geodirectory.RecordRight, &geodirectory.RecordOrdering, &geodirectory.RecordDepth,
		&geodirectory.ParentID, &geodirectory.CreatedAt, &geodirectory.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, err
	}

	return &geodirectory, nil
}

// GetProvincesByCountry retrieves provinces of a country
func (r *GeodirectoryRepository) GetProvincesByCountry(ctx context.Context, countryID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetChildrenByType(ctx, countryID, entities.GeoTypeProvince, limit, offset)
}

// GetCitiesByProvince retrieves cities of a province
func (r *GeodirectoryRepository) GetCitiesByProvince(ctx context.Context, provinceID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetChildrenByType(ctx, provinceID, entities.GeoTypeCity, limit, offset)
}

// GetDistrictsByCity retrieves districts of a city
func (r *GeodirectoryRepository) GetDistrictsByCity(ctx context.Context, cityID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetChildrenByType(ctx, cityID, entities.GeoTypeDistrict, limit, offset)
}

// GetVillagesByDistrict retrieves villages of a district
func (r *GeodirectoryRepository) GetVillagesByDistrict(ctx context.Context, districtID uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	return r.GetChildrenByType(ctx, districtID, entities.GeoTypeVillage, limit, offset)
}

// CountByType returns the count of geodirectories by type
func (r *GeodirectoryRepository) CountByType(ctx context.Context, geoType entities.GeoType) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_geodirectories WHERE type = $1"

	var count int64
	err := r.pool.QueryRow(ctx, query, geoType).Scan(&count)
	return count, err
}

// CountChildren returns the count of children for a geodirectory
func (r *GeodirectoryRepository) CountChildren(ctx context.Context, parentID uuid.UUID) (int64, error) {
	query := "SELECT COUNT(*) FROM tm_geodirectories WHERE parent_id = $1"

	var count int64
	err := r.pool.QueryRow(ctx, query, parentID).Scan(&count)
	return count, err
}

// HasChildren checks if a geodirectory has children
func (r *GeodirectoryRepository) HasChildren(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.CountChildren(ctx, id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetParent retrieves the parent of a geodirectory
func (r *GeodirectoryRepository) GetParent(ctx context.Context, id uuid.UUID) (*entities.Geodirectory, error) {
	query := `
		SELECT p.id, p.name, p.type, p.code, p.postal_code, p.longitude, p.latitude,
			   p.record_left, p.record_right, p.record_ordering, p.parent_id, p.created_at, p.updated_at
		FROM tm_geodirectories c
		JOIN tm_geodirectories p ON c.parent_id = p.id
		WHERE c.id = $1`

	var parent entities.Geodirectory
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&parent.ID, &parent.Name, &parent.Type, &parent.Code,
		&parent.PostalCode, &parent.Longitude, &parent.Latitude,
		&parent.RecordLeft, &parent.RecordRight, &parent.RecordOrdering, &parent.RecordDepth,
		&parent.ParentID, &parent.CreatedAt, &parent.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("parent not found")
		}
		return nil, err
	}

	return &parent, nil
}

// GetAncestors retrieves all ancestors of a geodirectory using nested set model
func (r *GeodirectoryRepository) GetAncestors(ctx context.Context, id uuid.UUID) ([]*entities.Geodirectory, error) {
	query := `
		SELECT p.id, p.name, p.type, p.code, p.postal_code, p.longitude, p.latitude,
			   p.record_left, p.record_right, p.record_ordering, p.parent_id, p.created_at, p.updated_at
		FROM tm_geodirectories n, tm_geodirectories p
		WHERE n.id = $1 
		  AND n.record_left BETWEEN p.record_left AND p.record_right
		  AND p.id != n.id
		ORDER BY p.record_left`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetDescendants retrieves all descendants of a geodirectory using nested set model
func (r *GeodirectoryRepository) GetDescendants(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT c.id, c.name, c.type, c.code, c.postal_code, c.longitude, c.latitude,
			   c.record_left, c.record_right, c.record_ordering, c.parent_id, c.created_at, c.updated_at
		FROM tm_geodirectories p, tm_geodirectories c
		WHERE p.id = $1 
		  AND c.record_left BETWEEN p.record_left AND p.record_right
		  AND c.id != p.id
		ORDER BY c.record_left
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetSiblings retrieves siblings of a geodirectory (same parent)
func (r *GeodirectoryRepository) GetSiblings(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT s.id, s.name, s.type, s.code, s.postal_code, s.longitude, s.latitude,
			   s.record_left, s.record_right, s.record_ordering, s.parent_id, s.created_at, s.updated_at
		FROM tm_geodirectories n
		JOIN tm_geodirectories s ON n.parent_id = s.parent_id
		WHERE n.id = $1 AND s.id != $1
		ORDER BY s.record_ordering, s.name
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetByNestedSetRange retrieves geodirectories within a nested set range
func (r *GeodirectoryRepository) GetByNestedSetRange(ctx context.Context, left, right int, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE record_left >= $1 AND record_right <= $2
		ORDER BY record_left
		LIMIT $3 OFFSET $4`

	rows, err := r.pool.Query(ctx, query, left, right, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// UpdateNestedSetValues updates the nested set values for a geodirectory
func (r *GeodirectoryRepository) UpdateNestedSetValues(ctx context.Context, id uuid.UUID, left, right, ordering int) error {
	query := `
		UPDATE tm_geodirectories 
		SET record_left = $2, record_right = $3, record_ordering = $4, updated_at = $5
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, left, right, ordering, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("geodirectory not found")
	}

	return nil
}

// InsertNode inserts a new node into the nested set tree
func (r *GeodirectoryRepository) InsertNode(ctx context.Context, geodirectory *entities.Geodirectory, parentID *uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var rightValue int

	if parentID == nil {
		// Insert as root node
		// Find the maximum right value
		err = tx.QueryRow(ctx, "SELECT COALESCE(MAX(record_right), 0) FROM tm_geodirectories").Scan(&rightValue)
		if err != nil {
			return err
		}
		rightValue += 1
	} else {
		// Insert as child of parent
		// Get parent's right value
		err = tx.QueryRow(ctx, "SELECT record_right FROM tm_geodirectories WHERE id = $1", parentID).Scan(&rightValue)
		if err != nil {
			return err
		}

		// Make space for new node by updating existing nodes
		_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_right = record_right + 2 WHERE record_right >= $1", rightValue)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = record_left + 2 WHERE record_left >= $1", rightValue)
		if err != nil {
			return err
		}
	}

	// Set nested set values for new node
	geodirectory.RecordLeft = &rightValue
	rightValue += 1
	geodirectory.RecordRight = &rightValue
	geodirectory.ParentID = parentID

	// Insert the new node
	geodirectory.GenerateID()
	geodirectory.CreatedAt = time.Now()
	geodirectory.UpdatedAt = time.Now()

	query := `
		INSERT INTO tm_geodirectories (
			id, name, type, code, postal_code, longitude, latitude,
			record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`

	_, err = tx.Exec(ctx, query,
		geodirectory.ID, geodirectory.Name, geodirectory.Type, geodirectory.Code,
		geodirectory.PostalCode, geodirectory.Longitude, geodirectory.Latitude,
		geodirectory.RecordLeft, geodirectory.RecordRight, geodirectory.RecordOrdering, geodirectory.RecordDepth,
		geodirectory.ParentID, geodirectory.CreatedAt, geodirectory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteNode deletes a node and all its descendants from the nested set tree
func (r *GeodirectoryRepository) DeleteNode(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get the node's left and right values
	var left, right int
	err = tx.QueryRow(ctx, "SELECT record_left, record_right FROM tm_geodirectories WHERE id = $1", id).Scan(&left, &right)
	if err != nil {
		return err
	}

	// Delete the node and all its descendants
	_, err = tx.Exec(ctx, "DELETE FROM tm_geodirectories WHERE record_left >= $1 AND record_right <= $2", left, right)
	if err != nil {
		return err
	}

	// Close the gap in the nested set
	width := right - left + 1
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = record_left - $1 WHERE record_left > $2", width, right)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_right = record_right - $1 WHERE record_right > $2", width, right)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// MoveNode moves a node to a new parent in the nested set tree
func (r *GeodirectoryRepository) MoveNode(ctx context.Context, nodeID, newParentID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get the node's current left and right values
	var nodeLeft, nodeRight int
	err = tx.QueryRow(ctx, "SELECT record_left, record_right FROM tm_geodirectories WHERE id = $1", nodeID).Scan(&nodeLeft, &nodeRight)
	if err != nil {
		return err
	}

	// Get the new parent's right value
	var parentRight int
	err = tx.QueryRow(ctx, "SELECT record_right FROM tm_geodirectories WHERE id = $1", newParentID).Scan(&parentRight)
	if err != nil {
		return err
	}

	// Calculate the width of the subtree being moved
	width := nodeRight - nodeLeft + 1

	// Temporarily set moved nodes to negative values to avoid conflicts
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = 0 - record_left, record_right = 0 - record_right WHERE record_left >= $1 AND record_right <= $2", nodeLeft, nodeRight)
	if err != nil {
		return err
	}

	// Close the gap left by the moved subtree
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = record_left - $1 WHERE record_left > $2", width, nodeRight)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_right = record_right - $1 WHERE record_right > $2", width, nodeRight)
	if err != nil {
		return err
	}

	// Adjust parent right if it was affected by the gap closure
	if parentRight > nodeRight {
		parentRight -= width
	}

	// Make space at the new location
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_right = record_right + $1 WHERE record_right >= $2", width, parentRight)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = record_left + $1 WHERE record_left >= $2", width, parentRight)
	if err != nil {
		return err
	}

	// Calculate the offset for the moved subtree
	offset := parentRight - nodeLeft

	// Move the subtree to its new location and update parent_id
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = 0 - record_left + $1, record_right = 0 - record_right + $1, parent_id = CASE WHEN id = $2 THEN $3 ELSE parent_id END WHERE record_left <= 0", offset, nodeID, newParentID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RebuildNestedSet rebuilds the entire nested set structure from parent-child relationships
func (r *GeodirectoryRepository) RebuildNestedSet(ctx context.Context) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Reset all nested set values
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = NULL, record_right = NULL, record_ordering = NULL")
	if err != nil {
		return err
	}

	// Get all root nodes (nodes with no parent)
	rows, err := tx.Query(ctx, "SELECT id FROM tm_geodirectories WHERE parent_id IS NULL ORDER BY name")
	if err != nil {
		return err
	}
	defer rows.Close()

	var rootIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return err
		}
		rootIDs = append(rootIDs, id)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Rebuild nested set for each root
	counter := 1
	for _, rootID := range rootIDs {
		counter, err = r.rebuildNestedSetForNode(ctx, tx, rootID, counter)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// rebuildNestedSetForNode recursively rebuilds nested set values for a node and its descendants
func (r *GeodirectoryRepository) rebuildNestedSetForNode(ctx context.Context, tx pgx.Tx, nodeID uuid.UUID, counter int) (int, error) {
	leftValue := counter
	counter++

	// Get children of this node
	rows, err := tx.Query(ctx, "SELECT id FROM tm_geodirectories WHERE parent_id = $1 ORDER BY name", nodeID)
	if err != nil {
		return counter, err
	}
	defer rows.Close()

	var childIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return counter, err
		}
		childIDs = append(childIDs, id)
	}

	if err := rows.Err(); err != nil {
		return counter, err
	}

	// Process each child
	ordering := 1
	for _, childID := range childIDs {
		counter, err = r.rebuildNestedSetForNode(ctx, tx, childID, counter)
		if err != nil {
			return counter, err
		}
		ordering++
	}

	rightValue := counter
	counter++

	// Update this node's values
	_, err = tx.Exec(ctx, "UPDATE tm_geodirectories SET record_left = $1, record_right = $2, record_ordering = $3 WHERE id = $4", leftValue, rightValue, ordering, nodeID)
	if err != nil {
		return counter, err
	}

	return counter, nil
}

// GetByCoordinates retrieves geodirectories by geographic coordinates (placeholder implementation)
func (r *GeodirectoryRepository) GetByCoordinates(ctx context.Context, latitude, longitude string, radius float64) ([]*entities.Geodirectory, error) {
	// This would require PostGIS extension for proper geographic queries
	// For now, return a simple string-based match
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE latitude = $1 AND longitude = $2
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query, latitude, longitude)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// GetNearby retrieves geodirectories near a given geodirectory (placeholder implementation)
func (r *GeodirectoryRepository) GetNearby(ctx context.Context, id uuid.UUID, radius float64, limit int) ([]*entities.Geodirectory, error) {
	// This would require PostGIS extension for proper geographic queries
	// For now, return siblings and nearby nodes in the hierarchy
	return r.GetSiblings(ctx, id, limit, 0)
}

func (r *GeodirectoryRepository) GetRoots(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE parent_id IS NULL
		ORDER BY record_ordering, name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

func (r *GeodirectoryRepository) GetLeaves(ctx context.Context, limit, offset int) ([]*entities.Geodirectory, error) {
	query := `
		SELECT id, name, type, code, postal_code, longitude, latitude,
			   record_left, record_right, record_ordering, record_depth, parent_id, created_at, updated_at
		FROM tm_geodirectories
		WHERE record_right - record_left = 1
		ORDER BY record_ordering, name
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanGeodirectories(rows)
}

// scanGeodirectories is a helper method to scan rows into geodirectory entities
func (r *GeodirectoryRepository) scanGeodirectories(rows pgx.Rows) ([]*entities.Geodirectory, error) {
	var geodirectories []*entities.Geodirectory

	for rows.Next() {
		var geodirectory entities.Geodirectory
		err := rows.Scan(
			&geodirectory.ID, &geodirectory.Name, &geodirectory.Type, &geodirectory.Code,
			&geodirectory.PostalCode, &geodirectory.Longitude, &geodirectory.Latitude,
			&geodirectory.RecordLeft, &geodirectory.RecordRight, &geodirectory.RecordOrdering, &geodirectory.RecordDepth,
			&geodirectory.ParentID, &geodirectory.CreatedAt, &geodirectory.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		geodirectories = append(geodirectories, &geodirectory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return geodirectories, nil
}

// Truncate removes all geodirectory records efficiently using TRUNCATE
func (r *GeodirectoryRepository) Truncate(ctx context.Context) error {
	query := `TRUNCATE TABLE tm_geodirectories RESTART IDENTITY CASCADE`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate geodirectories table: %w", err)
	}
	return nil
}
