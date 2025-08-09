package entities

import (
	"time"

	"github.com/google/uuid"
)

// GeoType represents the type of geographical location
type GeoType string

const (
	GeoTypeContinent    GeoType = "CONTINENT"
	GeoTypeSubcontinent GeoType = "SUBCONTINENT"
	GeoTypeCountry      GeoType = "COUNTRY"
	GeoTypeState        GeoType = "STATE"
	GeoTypeProvince     GeoType = "PROVINCE"
	GeoTypeRegency      GeoType = "REGENCY"
	GeoTypeCity         GeoType = "CITY"
	GeoTypeDistrict     GeoType = "DISTRICT"
	GeoTypeSubdistrict  GeoType = "SUBDISTRICT"
	GeoTypeVillage      GeoType = "VILLAGE"
)

// Geodirectory represents a geographical location in a hierarchical structure
type Geodirectory struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Type           GeoType    `json:"type" db:"type"`
	Code           *string    `json:"code,omitempty" db:"code"`
	PostalCode     *string    `json:"postal_code,omitempty" db:"postal_code"`
	Longitude      *string    `json:"longitude,omitempty" db:"longitude"`
	Latitude       *string    `json:"latitude,omitempty" db:"latitude"`
	RecordLeft     *int       `json:"record_left,omitempty" db:"record_left"`
	RecordRight    *int       `json:"record_right,omitempty" db:"record_right"`
	RecordOrdering *int       `json:"record_ordering,omitempty" db:"record_ordering"`
	RecordDepth    *int       `json:"record_depth,omitempty" db:"record_depth"`
	ParentID       *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	// Relations (not stored in DB, populated when needed)
	Parent   *Geodirectory   `json:"parent,omitempty"`
	Children []*Geodirectory `json:"children,omitempty"`
}

// TableName returns the table name for the Geodirectory entity
func (g *Geodirectory) TableName() string {
	return "tm_geodirectories"
}

// GenerateID generates a new UUID for the geodirectory if not set
func (g *Geodirectory) GenerateID() {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
}

// NewGeodirectory creates a new Geodirectory instance
func NewGeodirectory(name string, geoType GeoType) *Geodirectory {
	return &Geodirectory{
		ID:        uuid.New(),
		Name:      name,
		Type:      geoType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetCode sets the code for the geodirectory
func (g *Geodirectory) SetCode(code string) {
	g.Code = &code
	g.UpdatedAt = time.Now()
}

// SetPostalCode sets the postal code for the geodirectory
func (g *Geodirectory) SetPostalCode(postalCode string) {
	g.PostalCode = &postalCode
	g.UpdatedAt = time.Now()
}

// SetCoordinates sets the latitude and longitude for the geodirectory
func (g *Geodirectory) SetCoordinates(latitude, longitude string) {
	g.Latitude = &latitude
	g.Longitude = &longitude
	g.UpdatedAt = time.Now()
}

// SetParent sets the parent geodirectory
func (g *Geodirectory) SetParent(parentID uuid.UUID) {
	g.ParentID = &parentID
	g.UpdatedAt = time.Now()
}

// SetNestedSetValues sets the nested set model values for hierarchical queries
func (g *Geodirectory) SetNestedSetValues(left, right, ordering int) {
	g.RecordLeft = &left
	g.RecordRight = &right
	g.RecordOrdering = &ordering
	g.UpdatedAt = time.Now()
}

// SetDepth sets the hierarchical depth for the geodirectory
func (g *Geodirectory) SetDepth(depth int) {
	g.RecordDepth = &depth
	g.UpdatedAt = time.Now()
}

// SetOrderingID sets the record_ordering as a sortable identifier
func (g *Geodirectory) SetOrderingID(orderingID int) {
	g.RecordOrdering = &orderingID
	g.UpdatedAt = time.Now()
}

// GetDepthForType returns the appropriate depth based on geodirectory type
func (g *Geodirectory) GetDepthForType() int {
	switch g.Type {
	case GeoTypeContinent:
		return 1
	case GeoTypeSubcontinent:
		return 2
	case GeoTypeCountry:
		return 3
	case GeoTypeState, GeoTypeProvince:
		return 4
	case GeoTypeRegency:
		return 5
	case GeoTypeCity:
		return 6
	case GeoTypeDistrict:
		return 7
	case GeoTypeSubdistrict:
		return 8
	case GeoTypeVillage:
		return 9
	default:
		return 0
	}
}

// IsLeaf checks if this geodirectory is a leaf node (has no children)
func (g *Geodirectory) IsLeaf() bool {
	if g.RecordLeft == nil || g.RecordRight == nil {
		return false
	}
	return *g.RecordRight-*g.RecordLeft == 1
}

// HasChildren checks if this geodirectory has children
func (g *Geodirectory) HasChildren() bool {
	return !g.IsLeaf()
}

// IsRoot checks if this geodirectory is a root node (has no parent)
func (g *Geodirectory) IsRoot() bool {
	return g.ParentID == nil
}

// GetHierarchyLevel returns the hierarchy level based on the type
func (g *Geodirectory) GetHierarchyLevel() int {
	switch g.Type {
	case GeoTypeContinent:
		return 1
	case GeoTypeSubcontinent:
		return 2
	case GeoTypeCountry:
		return 3
	case GeoTypeState, GeoTypeProvince:
		return 4
	case GeoTypeRegency:
		return 5
	case GeoTypeCity:
		return 6
	case GeoTypeDistrict:
		return 7
	case GeoTypeSubdistrict:
		return 8
	case GeoTypeVillage:
		return 9
	default:
		return 0
	}
}

// GetFullPath returns the full hierarchical path as a string
func (g *Geodirectory) GetFullPath() string {
	if g.Parent == nil {
		return g.Name
	}
	return g.Parent.GetFullPath() + " > " + g.Name
}

// ValidateType checks if the geodirectory type is valid
func (g *Geodirectory) ValidateType() bool {
	validTypes := []GeoType{
		GeoTypeContinent, GeoTypeSubcontinent, GeoTypeCountry,
		GeoTypeState, GeoTypeProvince, GeoTypeRegency,
		GeoTypeCity, GeoTypeDistrict, GeoTypeSubdistrict, GeoTypeVillage,
	}

	for _, validType := range validTypes {
		if g.Type == validType {
			return true
		}
	}
	return false
}

// CanHaveParentType checks if the current type can have the specified parent type
func (g *Geodirectory) CanHaveParentType(parentType GeoType) bool {
	validParents := map[GeoType][]GeoType{
		GeoTypeSubcontinent: {GeoTypeContinent},
		GeoTypeCountry:      {GeoTypeContinent, GeoTypeSubcontinent},
		GeoTypeState:        {GeoTypeCountry},
		GeoTypeProvince:     {GeoTypeCountry},
		GeoTypeRegency:      {GeoTypeState, GeoTypeProvince},
		GeoTypeCity:         {GeoTypeState, GeoTypeProvince, GeoTypeRegency},
		GeoTypeDistrict:     {GeoTypeCity, GeoTypeRegency},
		GeoTypeSubdistrict:  {GeoTypeDistrict},
		GeoTypeVillage:      {GeoTypeDistrict, GeoTypeSubdistrict},
	}

	allowedParents, exists := validParents[g.Type]
	if !exists {
		return false // Root types like CONTINENT don't have parents
	}

	for _, allowedParent := range allowedParents {
		if parentType == allowedParent {
			return true
		}
	}
	return false
}
