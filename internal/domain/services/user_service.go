package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/repositories"
)

// UserService contains business logic for user operations
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *UserService) CreateUser(ctx context.Context, email, firstName, lastName string) (*entities.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create new user
	user := entities.NewUser(email, firstName, lastName)

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	users, err := s.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName string) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	user.Update(user.Email, firstName, lastName)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeactivateUser deactivates a user
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user for deactivation: %w", err)
	}

	user.Deactivate()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// ActivateUser activates a user
func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user for activation: %w", err)
	}

	user.Activate()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
