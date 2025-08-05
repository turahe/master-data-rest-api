package primary

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// UserHandler defines the interface for user operations from the primary side
type UserHandler interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, request CreateUserRequest) (*entities.User, error)
	
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	
	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// GetAllUsers retrieves all users with pagination
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
	
	// UpdateUser updates user information
	UpdateUser(ctx context.Context, id uuid.UUID, request UpdateUserRequest) (*entities.User, error)
	
	// DeactivateUser deactivates a user
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	
	// ActivateUser activates a user
	ActivateUser(ctx context.Context, id uuid.UUID) error
	
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// CreateUserRequest represents the request for creating a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
} 