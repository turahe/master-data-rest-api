package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewLanguage(t *testing.T) {
	// Given
	name := "Indonesian"
	code := "id"

	// When
	language := NewLanguage(name, code)

	// Then
	assert.NotNil(t, language)
	assert.Equal(t, name, language.Name)
	assert.Equal(t, code, language.Code)
	assert.True(t, language.IsActive)
	assert.NotEqual(t, uuid.Nil, language.ID)
	assert.False(t, language.CreatedAt.IsZero())
	assert.False(t, language.UpdatedAt.IsZero())
}

func TestLanguage_GenerateID(t *testing.T) {
	t.Run("should generate ID when nil", func(t *testing.T) {
		// Given
		language := &Language{ID: uuid.Nil}

		// When
		language.GenerateID()

		// Then
		assert.NotEqual(t, uuid.Nil, language.ID)
	})

	t.Run("should not overwrite existing ID", func(t *testing.T) {
		// Given
		existingID := uuid.New()
		language := &Language{ID: existingID}

		// When
		language.GenerateID()

		// Then
		assert.Equal(t, existingID, language.ID)
	})
}

func TestLanguage_SetName(t *testing.T) {
	// Given
	language := NewLanguage("Old Name", "old")
	newName := "New Language Name"
	originalTime := language.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	language.SetName(newName)

	// Then
	assert.Equal(t, newName, language.Name)
	assert.True(t, language.UpdatedAt.After(originalTime))
}

func TestLanguage_SetCode(t *testing.T) {
	// Given
	language := NewLanguage("Name", "old")
	newCode := "new"
	originalTime := language.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	language.SetCode(newCode)

	// Then
	assert.Equal(t, newCode, language.Code)
	assert.True(t, language.UpdatedAt.After(originalTime))
}

func TestLanguage_Activate(t *testing.T) {
	// Given
	language := NewLanguage("Name", "code")
	language.IsActive = false
	originalTime := language.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	language.Activate()

	// Then
	assert.True(t, language.IsActive)
	assert.True(t, language.UpdatedAt.After(originalTime))
}

func TestLanguage_Deactivate(t *testing.T) {
	// Given
	language := NewLanguage("Name", "code")
	originalTime := language.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	language.Deactivate()

	// Then
	assert.False(t, language.IsActive)
	assert.True(t, language.UpdatedAt.After(originalTime))
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		languageName string
		code         string
		expected     bool
	}{
		{
			name:         "valid language",
			languageName: "English",
			code:         "en",
			expected:     true,
		},
		{
			name:         "empty name",
			languageName: "",
			code:         "en",
			expected:     false,
		},
		{
			name:         "empty code",
			languageName: "English",
			code:         "",
			expected:     false,
		},
		{
			name:         "code too long",
			languageName: "English",
			code:         "english-usa",
			expected:     false,
		},
		{
			name:         "valid long code",
			languageName: "Indonesian",
			code:         "id-ID",
			expected:     true,
		},
		{
			name:         "maximum length code",
			languageName: "Test",
			code:         "1234567890",
			expected:     true,
		},
		{
			name:         "over maximum length code",
			languageName: "Test",
			code:         "12345678901",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			language := &Language{
				Name: tt.languageName,
				Code: tt.code,
			}

			// When
			isValid := language.IsValid()

			// Then
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestLanguage_TableName(t *testing.T) {
	// Given
	language := &Language{}

	// When
	tableName := language.TableName()

	// Then
	assert.Equal(t, "tm_languages", tableName)
}
