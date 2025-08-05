package gorm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/ports/secondary"
	"gorm.io/gorm"
)

// UserRepository implements the UserRepository interface using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository
func NewUserRepository(db *gorm.DB) secondary.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a user by ID from the database
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email from the database
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", result.Error)
	}
	return &user, nil
}

// GetAll retrieves all users with pagination from the database
func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []entities.User
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all users: %w", result.Error)
	}

	// Convert slice of entities to slice of pointers
	userPtrs := make([]*entities.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	return userPtrs, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID from the database
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Exists checks if a user exists by email in the database
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check user existence: %w", result.Error)
	}

	return count > 0, nil
}

// Count returns the total number of users in the database
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count users: %w", result.Error)
	}

	return count, nil
}
