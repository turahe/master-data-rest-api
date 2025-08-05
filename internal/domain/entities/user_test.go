package entities

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	user := NewUser(email, firstName, lastName)

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.FirstName != firstName {
		t.Errorf("Expected first name %s, got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Errorf("Expected last name %s, got %s", lastName, user.LastName)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	if user.ID.String() == "" {
		t.Error("Expected user to have an ID")
	}
}

func TestUser_FullName(t *testing.T) {
	user := NewUser("test@example.com", "John", "Doe")
	expected := "John Doe"

	if user.FullName() != expected {
		t.Errorf("Expected full name %s, got %s", expected, user.FullName())
	}
}

func TestUser_Update(t *testing.T) {
	user := NewUser("test@example.com", "John", "Doe")
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Millisecond)

	newFirstName := "Jane"
	newLastName := "Smith"

	user.Update(newFirstName, newLastName)

	if user.FirstName != newFirstName {
		t.Errorf("Expected first name %s, got %s", newFirstName, user.FirstName)
	}

	if user.LastName != newLastName {
		t.Errorf("Expected last name %s, got %s", newLastName, user.LastName)
	}

	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}
}

func TestUser_Deactivate(t *testing.T) {
	user := NewUser("test@example.com", "John", "Doe")
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Millisecond)

	user.Deactivate()

	if user.IsActive {
		t.Error("Expected user to be inactive")
	}

	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}
}

func TestUser_Activate(t *testing.T) {
	user := NewUser("test@example.com", "John", "Doe")
	user.Deactivate() // First deactivate
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Millisecond)

	user.Activate()

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}
} 