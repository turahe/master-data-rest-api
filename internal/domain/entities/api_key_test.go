package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIKey(t *testing.T) {
	// Given
	name := "Test API Key"
	key := "test-key-123"

	// When
	apiKey := NewAPIKey(name, key)

	// Then
	assert.NotNil(t, apiKey)
	assert.Equal(t, name, apiKey.Name)
	assert.Equal(t, key, apiKey.Key)
	assert.True(t, apiKey.IsActive)
	assert.NotEqual(t, uuid.Nil, apiKey.ID)
	assert.Nil(t, apiKey.Description)
	assert.Nil(t, apiKey.ExpiresAt)
	assert.Nil(t, apiKey.LastUsedAt)
	assert.Nil(t, apiKey.DeletedAt)
}

func TestAPIKey_SetDescription(t *testing.T) {
	// Given
	apiKey := NewAPIKey("Test", "key")
	description := "Test API key for unit testing"

	// When
	apiKey.SetDescription(description)

	// Then
	require.NotNil(t, apiKey.Description)
	assert.Equal(t, description, *apiKey.Description)
}

func TestAPIKey_SetExpiration(t *testing.T) {
	// Given
	apiKey := NewAPIKey("Test", "key")
	expiration := time.Now().Add(24 * time.Hour)

	// When
	apiKey.SetExpiration(expiration)

	// Then
	require.NotNil(t, apiKey.ExpiresAt)
	assert.Equal(t, expiration, *apiKey.ExpiresAt)
}

func TestAPIKey_IsExpired(t *testing.T) {
	t.Run("not expired with future expiration", func(t *testing.T) {
		// Given
		apiKey := NewAPIKey("Test", "key")
		futureTime := time.Now().Add(24 * time.Hour)
		apiKey.SetExpiration(futureTime)

		// When
		expired := apiKey.IsExpired()

		// Then
		assert.False(t, expired)
	})

	t.Run("expired with past expiration", func(t *testing.T) {
		// Given
		apiKey := NewAPIKey("Test", "key")
		pastTime := time.Now().Add(-24 * time.Hour)
		apiKey.SetExpiration(pastTime)

		// When
		expired := apiKey.IsExpired()

		// Then
		assert.True(t, expired)
	})

	t.Run("not expired when no expiration set", func(t *testing.T) {
		// Given
		apiKey := NewAPIKey("Test", "key")

		// When
		expired := apiKey.IsExpired()

		// Then
		assert.False(t, expired)
	})
}

func TestAPIKey_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		setupKey func() *APIKey
		expected bool
	}{
		{
			name: "valid active key without expiration",
			setupKey: func() *APIKey {
				return NewAPIKey("Test", "key")
			},
			expected: true,
		},
		{
			name: "valid active key with future expiration",
			setupKey: func() *APIKey {
				key := NewAPIKey("Test", "key")
				key.SetExpiration(time.Now().Add(24 * time.Hour))
				return key
			},
			expected: true,
		},
		{
			name: "invalid inactive key",
			setupKey: func() *APIKey {
				key := NewAPIKey("Test", "key")
				key.Deactivate()
				return key
			},
			expected: false,
		},
		{
			name: "invalid expired key",
			setupKey: func() *APIKey {
				key := NewAPIKey("Test", "key")
				key.SetExpiration(time.Now().Add(-24 * time.Hour))
				return key
			},
			expected: false,
		},
		{
			name: "invalid inactive and expired key",
			setupKey: func() *APIKey {
				key := NewAPIKey("Test", "key")
				key.Deactivate()
				key.SetExpiration(time.Now().Add(-24 * time.Hour))
				return key
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			apiKey := tt.setupKey()

			// When
			isValid := apiKey.IsValid()

			// Then
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestAPIKey_UpdateLastUsed(t *testing.T) {
	// Given
	apiKey := NewAPIKey("Test", "key")
	beforeUpdate := time.Now()

	// When
	apiKey.UpdateLastUsed()

	// Then
	require.NotNil(t, apiKey.LastUsedAt)
	assert.True(t, apiKey.LastUsedAt.After(beforeUpdate) || apiKey.LastUsedAt.Equal(beforeUpdate))
}

func TestAPIKey_Deactivate(t *testing.T) {
	// Given
	apiKey := NewAPIKey("Test", "key")
	assert.True(t, apiKey.IsActive) // Verify it starts active

	// When
	apiKey.Deactivate()

	// Then
	assert.False(t, apiKey.IsActive)
}

func TestAPIKey_Activate(t *testing.T) {
	// Given
	apiKey := NewAPIKey("Test", "key")
	apiKey.IsActive = false

	// When
	apiKey.Activate()

	// Then
	assert.True(t, apiKey.IsActive)
}

func TestAPIKey_TableName(t *testing.T) {
	// Given
	apiKey := APIKey{}

	// When
	tableName := apiKey.TableName()

	// Then
	assert.Equal(t, "tm_api_keys", tableName)
}
