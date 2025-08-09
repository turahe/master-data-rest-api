package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCurrency(t *testing.T) {
	// Given
	name := "Indonesian Rupiah"
	code := "IDR"
	decimalPlaces := 0

	// When
	currency := NewCurrency(name, code, decimalPlaces)

	// Then
	assert.NotNil(t, currency)
	assert.Equal(t, name, currency.Name)
	assert.Equal(t, code, currency.Code)
	assert.Equal(t, decimalPlaces, currency.DecimalPlaces)
	assert.True(t, currency.IsActive)
	assert.NotEqual(t, uuid.Nil, currency.ID)
	assert.False(t, currency.CreatedAt.IsZero())
	assert.False(t, currency.UpdatedAt.IsZero())
}

func TestCurrency_GenerateID(t *testing.T) {
	t.Run("should generate ID when nil", func(t *testing.T) {
		// Given
		currency := &Currency{ID: uuid.Nil}

		// When
		currency.GenerateID()

		// Then
		assert.NotEqual(t, uuid.Nil, currency.ID)
	})

	t.Run("should not overwrite existing ID", func(t *testing.T) {
		// Given
		existingID := uuid.New()
		currency := &Currency{ID: existingID}

		// When
		currency.GenerateID()

		// Then
		assert.Equal(t, existingID, currency.ID)
	})
}

func TestCurrency_SetName(t *testing.T) {
	// Given
	currency := NewCurrency("Old Name", "OLD", 2)
	newName := "New Currency Name"
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.SetName(newName)

	// Then
	assert.Equal(t, newName, currency.Name)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_SetCode(t *testing.T) {
	// Given
	currency := NewCurrency("Name", "OLD", 2)
	newCode := "NEW"
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.SetCode(newCode)

	// Then
	assert.Equal(t, newCode, currency.Code)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_SetSymbol(t *testing.T) {
	// Given
	currency := NewCurrency("Name", "USD", 2)
	symbol := "$"
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.SetSymbol(symbol)

	// Then
	require.NotNil(t, currency.Symbol)
	assert.Equal(t, symbol, *currency.Symbol)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_SetDecimalPlaces(t *testing.T) {
	// Given
	currency := NewCurrency("Name", "USD", 2)
	newPlaces := 4
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.SetDecimalPlaces(newPlaces)

	// Then
	assert.Equal(t, newPlaces, currency.DecimalPlaces)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_Activate(t *testing.T) {
	// Given
	currency := NewCurrency("Name", "USD", 2)
	currency.IsActive = false
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.Activate()

	// Then
	assert.True(t, currency.IsActive)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_Deactivate(t *testing.T) {
	// Given
	currency := NewCurrency("Name", "USD", 2)
	originalTime := currency.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	currency.Deactivate()

	// Then
	assert.False(t, currency.IsActive)
	assert.True(t, currency.UpdatedAt.After(originalTime))
}

func TestCurrency_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		currencyName  string
		code          string
		decimalPlaces int
		expected      bool
	}{
		{
			name:          "valid currency",
			currencyName:  "US Dollar",
			code:          "USD",
			decimalPlaces: 2,
			expected:      true,
		},
		{
			name:          "empty name",
			currencyName:  "",
			code:          "USD",
			decimalPlaces: 2,
			expected:      false,
		},
		{
			name:          "empty code",
			currencyName:  "US Dollar",
			code:          "",
			decimalPlaces: 2,
			expected:      false,
		},
		{
			name:          "code too long",
			currencyName:  "US Dollar",
			code:          "USDX",
			decimalPlaces: 2,
			expected:      false,
		},
		{
			name:          "negative decimal places",
			currencyName:  "US Dollar",
			code:          "USD",
			decimalPlaces: -1,
			expected:      false,
		},
		{
			name:          "zero decimal places",
			currencyName:  "Indonesian Rupiah",
			code:          "IDR",
			decimalPlaces: 0,
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			currency := &Currency{
				Name:          tt.currencyName,
				Code:          tt.code,
				DecimalPlaces: tt.decimalPlaces,
			}

			// When
			isValid := currency.IsValid()

			// Then
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestCurrency_GetDisplayName(t *testing.T) {
	t.Run("with symbol", func(t *testing.T) {
		// Given
		currency := NewCurrency("US Dollar", "USD", 2)
		symbol := "$"
		currency.SetSymbol(symbol)

		// When
		displayName := currency.GetDisplayName()

		// Then
		assert.Equal(t, "US Dollar ($)", displayName)
	})

	t.Run("without symbol", func(t *testing.T) {
		// Given
		currency := NewCurrency("Indonesian Rupiah", "IDR", 0)

		// When
		displayName := currency.GetDisplayName()

		// Then
		assert.Equal(t, "Indonesian Rupiah", displayName)
	})

	t.Run("with empty symbol", func(t *testing.T) {
		// Given
		currency := NewCurrency("Euro", "EUR", 2)
		emptySymbol := ""
		currency.Symbol = &emptySymbol

		// When
		displayName := currency.GetDisplayName()

		// Then
		assert.Equal(t, "Euro", displayName)
	})
}

func TestCurrency_TableName(t *testing.T) {
	// Given
	currency := &Currency{}

	// When
	tableName := currency.TableName()

	// Then
	assert.Equal(t, "tm_currencies", tableName)
}
