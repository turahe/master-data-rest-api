package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
	"github.com/turahe/master-data-rest-api/internal/ports/secondary"
)

// UserRepository implements the UserRepository interface for MySQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new MySQL user repository
func NewUserRepository(db *sql.DB) secondary.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, first_name, last_name, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// GetByID retrieves a user by ID from the database
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT id, email, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`
	
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email from the database
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE email = ?
	`
	
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &user, nil
}

// GetAll retrieves all users with pagination from the database
func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, email, first_name, last_name, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()
	
	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	
	return users, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = ?, first_name = ?, last_name = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// Delete deletes a user by ID from the database
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// Exists checks if a user exists by email in the database
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	return count > 0, nil
}

// Count returns the total number of users in the database
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`
	
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	return count, nil
} 