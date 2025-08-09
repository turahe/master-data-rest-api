package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBank(t *testing.T) {
	// Given
	name := "Bank Central Asia"
	alias := "BCA"
	company := "PT Bank Central Asia Tbk"
	code := "014"

	// When
	bank := NewBank(name, alias, company, code)

	// Then
	assert.NotNil(t, bank)
	assert.Equal(t, name, bank.Name)
	assert.Equal(t, alias, bank.Alias)
	assert.Equal(t, company, bank.Company)
	assert.Equal(t, code, bank.Code)
	assert.NotEqual(t, uuid.Nil, bank.ID)
	assert.False(t, bank.CreatedAt.IsZero())
	assert.False(t, bank.UpdatedAt.IsZero())
}

func TestBank_GenerateID(t *testing.T) {
	t.Run("should generate ID when nil", func(t *testing.T) {
		// Given
		bank := &Bank{ID: uuid.Nil}

		// When
		bank.GenerateID()

		// Then
		assert.NotEqual(t, uuid.Nil, bank.ID)
	})

	t.Run("should not overwrite existing ID", func(t *testing.T) {
		// Given
		existingID := uuid.New()
		bank := &Bank{ID: existingID}

		// When
		bank.GenerateID()

		// Then
		assert.Equal(t, existingID, bank.ID)
	})
}

func TestBank_SetName(t *testing.T) {
	// Given
	bank := NewBank("Old Name", "ALIAS", "Company", "001")
	newName := "New Bank Name"
	originalTime := bank.UpdatedAt

	// Small delay to ensure time difference
	time.Sleep(1 * time.Millisecond)

	// When
	bank.SetName(newName)

	// Then
	assert.Equal(t, newName, bank.Name)
	assert.True(t, bank.UpdatedAt.After(originalTime))
}

func TestBank_SetAlias(t *testing.T) {
	// Given
	bank := NewBank("Name", "OLD", "Company", "001")
	newAlias := "NEW"
	originalTime := bank.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	bank.SetAlias(newAlias)

	// Then
	assert.Equal(t, newAlias, bank.Alias)
	assert.True(t, bank.UpdatedAt.After(originalTime))
}

func TestBank_SetCompany(t *testing.T) {
	// Given
	bank := NewBank("Name", "ALIAS", "Old Company", "001")
	newCompany := "New Company Ltd"
	originalTime := bank.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	bank.SetCompany(newCompany)

	// Then
	assert.Equal(t, newCompany, bank.Company)
	assert.True(t, bank.UpdatedAt.After(originalTime))
}

func TestBank_SetCode(t *testing.T) {
	// Given
	bank := NewBank("Name", "ALIAS", "Company", "001")
	newCode := "999"
	originalTime := bank.UpdatedAt

	time.Sleep(1 * time.Millisecond)

	// When
	bank.SetCode(newCode)

	// Then
	assert.Equal(t, newCode, bank.Code)
	assert.True(t, bank.UpdatedAt.After(originalTime))
}

func TestBank_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		bankName string
		code     string
		expected bool
	}{
		{
			name:     "valid bank",
			bankName: "Test Bank",
			code:     "001",
			expected: true,
		},
		{
			name:     "empty name",
			bankName: "",
			code:     "001",
			expected: false,
		},
		{
			name:     "empty code",
			bankName: "Test Bank",
			code:     "",
			expected: false,
		},
		{
			name:     "both empty",
			bankName: "",
			code:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			bank := &Bank{
				Name: tt.bankName,
				Code: tt.code,
			}

			// When
			isValid := bank.IsValid()

			// Then
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestBank_TableName(t *testing.T) {
	// Given
	bank := &Bank{}

	// When
	tableName := bank.TableName()

	// Then
	assert.Equal(t, "tm_banks", tableName)
}
