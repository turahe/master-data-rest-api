package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "valid email with plus",
			email:       "user+test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with dash",
			email:       "user-test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with numbers",
			email:       "user123@example123.com",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
			errorMsg:    "email cannot be empty",
		},
		{
			name:        "invalid email - no @",
			email:       "testexample.com",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - no domain",
			email:       "test@",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - no local part",
			email:       "@example.com",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - no TLD",
			email:       "test@example",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - TLD too short",
			email:       "test@example.c",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - spaces",
			email:       "test @example.com",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid email - double @",
			email:       "test@@example.com",
			expectError: true,
			errorMsg:    "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			email, err := NewEmail(tt.email)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, email)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, email)
				assert.Equal(t, tt.email, email.Value())
				assert.Equal(t, tt.email, email.String())
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	// Given
	emailStr := "test@example.com"
	email, err := NewEmail(emailStr)
	assert.NoError(t, err)

	// When
	value := email.Value()

	// Then
	assert.Equal(t, emailStr, value)
}

func TestEmail_String(t *testing.T) {
	// Given
	emailStr := "test@example.com"
	email, err := NewEmail(emailStr)
	assert.NoError(t, err)

	// When
	str := email.String()

	// Then
	assert.Equal(t, emailStr, str)
}
