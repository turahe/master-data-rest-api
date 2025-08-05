package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/domain/services"
	"github.com/turahe/master-data-rest-api/internal/ports/primary"
)

// UserApplicationService implements the UserHandler interface
type UserApplicationService struct {
	userService *services.UserService
}

// NewUserApplicationService creates a new user application service
func NewUserApplicationService(userService *services.UserService) primary.UserHandler {
	return &UserApplicationService{
		userService: userService,
	}
}

// CreateUser creates a new user
func (s *UserApplicationService) CreateUser(ctx context.Context, request primary.CreateUserRequest) (*entities.User, error) {
	return s.userService.CreateUser(ctx, request.Email, request.FirstName, request.LastName)
}

// GetUserByID retrieves a user by ID
func (s *UserApplicationService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userService.GetUserByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserApplicationService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.userService.GetUserByEmail(ctx, email)
}

// GetAllUsers retrieves all users with pagination
func (s *UserApplicationService) GetAllUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	return s.userService.GetAllUsers(ctx, limit, offset)
}

// UpdateUser updates user information
func (s *UserApplicationService) UpdateUser(ctx context.Context, id uuid.UUID, request primary.UpdateUserRequest) (*entities.User, error) {
	return s.userService.UpdateUser(ctx, id, request.FirstName, request.LastName)
}

// DeactivateUser deactivates a user
func (s *UserApplicationService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	return s.userService.DeactivateUser(ctx, id)
}

// ActivateUser activates a user
func (s *UserApplicationService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	return s.userService.ActivateUser(ctx, id)
}

// DeleteUser deletes a user
func (s *UserApplicationService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userService.DeleteUser(ctx, id)
} 