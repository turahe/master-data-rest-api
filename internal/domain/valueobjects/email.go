package valueobjects

import (
	"fmt"
	"regexp"
)

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates a new email value object
func NewEmail(email string) (*Email, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	return &Email{value: email}, nil
}

// Value returns the email string value
func (e *Email) Value() string {
	return e.value
}

// String implements the Stringer interface
func (e *Email) String() string {
	return e.value
}

// validateEmail validates the email format
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}
