package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGeodirectory(t *testing.T) {
	// Given
	name := "Jakarta"
	geoType := GeoTypeCity

	// When
	geo := NewGeodirectory(name, geoType)

	// Then
	assert.NotNil(t, geo)
	assert.Equal(t, name, geo.Name)
	assert.Equal(t, geoType, geo.Type)
	assert.NotEqual(t, uuid.Nil, geo.ID)
	assert.False(t, geo.CreatedAt.IsZero())
	assert.False(t, geo.UpdatedAt.IsZero())
}

func TestGeodirectory_GenerateID(t *testing.T) {
	t.Run("should generate ID when nil", func(t *testing.T) {
		// Given
		geo := &Geodirectory{ID: uuid.Nil}

		// When
		geo.GenerateID()

		// Then
		assert.NotEqual(t, uuid.Nil, geo.ID)
	})

	t.Run("should not overwrite existing ID", func(t *testing.T) {
		// Given
		existingID := uuid.New()
		geo := &Geodirectory{ID: existingID}

		// When
		geo.GenerateID()

		// Then
		assert.Equal(t, existingID, geo.ID)
	})
}

func TestGeodirectory_SetCode(t *testing.T) {
	// Given
	geo := NewGeodirectory("Test", GeoTypeCity)
	code := "JKT001"
	originalTime := geo.UpdatedAt

	// Small delay to ensure time difference
	time.Sleep(1 * time.Millisecond)

	// When
	geo.SetCode(code)

	// Then
	require.NotNil(t, geo.Code)
	assert.Equal(t, code, *geo.Code)
	assert.True(t, geo.UpdatedAt.After(originalTime))
}

func TestGeodirectory_SetPostalCode(t *testing.T) {
	// Given
	geo := NewGeodirectory("Test", GeoTypeCity)
	postalCode := "12345"
	originalTime := geo.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	geo.SetPostalCode(postalCode)

	// Then
	require.NotNil(t, geo.PostalCode)
	assert.Equal(t, postalCode, *geo.PostalCode)
	assert.True(t, geo.UpdatedAt.After(originalTime))
}

func TestGeodirectory_SetCoordinates(t *testing.T) {
	// Given
	geo := NewGeodirectory("Test", GeoTypeCity)
	latitude := "-6.2088"
	longitude := "106.8456"
	originalTime := geo.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	geo.SetCoordinates(latitude, longitude)

	// Then
	require.NotNil(t, geo.Latitude)
	require.NotNil(t, geo.Longitude)
	assert.Equal(t, latitude, *geo.Latitude)
	assert.Equal(t, longitude, *geo.Longitude)
	assert.True(t, geo.UpdatedAt.After(originalTime))
}

func TestGeodirectory_SetParent(t *testing.T) {
	// Given
	geo := NewGeodirectory("Test", GeoTypeCity)
	parentID := uuid.New()
	originalTime := geo.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	geo.SetParent(parentID)

	// Then
	require.NotNil(t, geo.ParentID)
	assert.Equal(t, parentID, *geo.ParentID)
	assert.True(t, geo.UpdatedAt.After(originalTime))
}

func TestGeodirectory_SetNestedSetValues(t *testing.T) {
	// Given
	geo := NewGeodirectory("Test", GeoTypeCity)
	left, right, ordering := 1, 10, 1
	originalTime := geo.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	geo.SetNestedSetValues(left, right, ordering)

	// Then
	require.NotNil(t, geo.RecordLeft)
	require.NotNil(t, geo.RecordRight)
	require.NotNil(t, geo.RecordOrdering)
	assert.Equal(t, left, *geo.RecordLeft)
	assert.Equal(t, right, *geo.RecordRight)
	assert.Equal(t, ordering, *geo.RecordOrdering)
	assert.True(t, geo.UpdatedAt.After(originalTime))
}

func TestGeodirectory_IsLeaf(t *testing.T) {
	tests := []struct {
		name        string
		recordLeft  *int
		recordRight *int
		expected    bool
	}{
		{
			name:        "leaf node",
			recordLeft:  intPtr(1),
			recordRight: intPtr(2),
			expected:    true,
		},
		{
			name:        "non-leaf node",
			recordLeft:  intPtr(1),
			recordRight: intPtr(10),
			expected:    false,
		},
		{
			name:        "nil left",
			recordLeft:  nil,
			recordRight: intPtr(2),
			expected:    false,
		},
		{
			name:        "nil right",
			recordLeft:  intPtr(1),
			recordRight: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			geo := &Geodirectory{
				RecordLeft:  tt.recordLeft,
				RecordRight: tt.recordRight,
			}

			// When
			result := geo.IsLeaf()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeodirectory_HasChildren(t *testing.T) {
	// Given
	leafNode := &Geodirectory{
		RecordLeft:  intPtr(1),
		RecordRight: intPtr(2),
	}
	parentNode := &Geodirectory{
		RecordLeft:  intPtr(1),
		RecordRight: intPtr(10),
	}

	// When & Then
	assert.False(t, leafNode.HasChildren())
	assert.True(t, parentNode.HasChildren())
}

func TestGeodirectory_IsRoot(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uuid.UUID
		expected bool
	}{
		{
			name:     "root node",
			parentID: nil,
			expected: true,
		},
		{
			name:     "child node",
			parentID: uuidPtr(uuid.New()),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			geo := &Geodirectory{ParentID: tt.parentID}

			// When
			result := geo.IsRoot()

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeodirectory_GetHierarchyLevel(t *testing.T) {
	tests := []struct {
		geoType  GeoType
		expected int
	}{
		{GeoTypeContinent, 1},
		{GeoTypeSubcontinent, 2},
		{GeoTypeCountry, 3},
		{GeoTypeState, 4},
		{GeoTypeProvince, 4},
		{GeoTypeRegency, 5},
		{GeoTypeCity, 6},
		{GeoTypeDistrict, 7},
		{GeoTypeSubdistrict, 8},
		{GeoTypeVillage, 9},
		{GeoType("INVALID"), 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.geoType), func(t *testing.T) {
			// Given
			geo := &Geodirectory{Type: tt.geoType}

			// When
			level := geo.GetHierarchyLevel()

			// Then
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestGeodirectory_GetFullPath(t *testing.T) {
	t.Run("root node", func(t *testing.T) {
		// Given
		geo := &Geodirectory{Name: "Asia"}

		// When
		path := geo.GetFullPath()

		// Then
		assert.Equal(t, "Asia", path)
	})

	t.Run("nested node", func(t *testing.T) {
		// Given
		continent := &Geodirectory{Name: "Asia"}
		country := &Geodirectory{Name: "Indonesia", Parent: continent}
		province := &Geodirectory{Name: "West Java", Parent: country}
		city := &Geodirectory{Name: "Jakarta", Parent: province}

		// When
		path := city.GetFullPath()

		// Then
		assert.Equal(t, "Asia > Indonesia > West Java > Jakarta", path)
	})
}

func TestGeodirectory_ValidateType(t *testing.T) {
	validTypes := []GeoType{
		GeoTypeContinent, GeoTypeSubcontinent, GeoTypeCountry,
		GeoTypeState, GeoTypeProvince, GeoTypeRegency,
		GeoTypeCity, GeoTypeDistrict, GeoTypeSubdistrict, GeoTypeVillage,
	}

	for _, geoType := range validTypes {
		t.Run("valid_"+string(geoType), func(t *testing.T) {
			// Given
			geo := &Geodirectory{Type: geoType}

			// When
			isValid := geo.ValidateType()

			// Then
			assert.True(t, isValid)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		// Given
		geo := &Geodirectory{Type: GeoType("INVALID")}

		// When
		isValid := geo.ValidateType()

		// Then
		assert.False(t, isValid)
	})
}

func TestGeodirectory_CanHaveParentType(t *testing.T) {
	tests := []struct {
		childType  GeoType
		parentType GeoType
		expected   bool
	}{
		// Valid parent-child relationships
		{GeoTypeSubcontinent, GeoTypeContinent, true},
		{GeoTypeCountry, GeoTypeContinent, true},
		{GeoTypeCountry, GeoTypeSubcontinent, true},
		{GeoTypeState, GeoTypeCountry, true},
		{GeoTypeProvince, GeoTypeCountry, true},
		{GeoTypeRegency, GeoTypeState, true},
		{GeoTypeRegency, GeoTypeProvince, true},
		{GeoTypeCity, GeoTypeState, true},
		{GeoTypeCity, GeoTypeProvince, true},
		{GeoTypeCity, GeoTypeRegency, true},
		{GeoTypeDistrict, GeoTypeCity, true},
		{GeoTypeDistrict, GeoTypeRegency, true},
		{GeoTypeSubdistrict, GeoTypeDistrict, true},
		{GeoTypeVillage, GeoTypeDistrict, true},
		{GeoTypeVillage, GeoTypeSubdistrict, true},

		// Invalid parent-child relationships
		{GeoTypeContinent, GeoTypeCountry, false}, // Continent can't have a parent
		{GeoTypeCountry, GeoTypeCity, false},      // Country can't be child of city
		{GeoTypeCity, GeoTypeVillage, false},      // City can't be child of village
		{GeoTypeVillage, GeoTypeContinent, false}, // Village can't be child of continent
	}

	for _, tt := range tests {
		t.Run(string(tt.childType)+"_"+string(tt.parentType), func(t *testing.T) {
			// Given
			geo := &Geodirectory{Type: tt.childType}

			// When
			canHave := geo.CanHaveParentType(tt.parentType)

			// Then
			assert.Equal(t, tt.expected, canHave)
		})
	}
}

func TestGeodirectory_TableName(t *testing.T) {
	// Given
	geo := &Geodirectory{}

	// When
	tableName := geo.TableName()

	// Then
	assert.Equal(t, "tm_geodirectories", tableName)
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
