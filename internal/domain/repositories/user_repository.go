package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// GetAll retrieves all users with optional pagination
	GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Exists checks if a user exists by email
	Exists(ctx context.Context, email string) (bool, error)
	
	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
} 