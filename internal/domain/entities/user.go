package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity in the domain
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Email     string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	FirstName string    `json:"first_name" gorm:"type:varchar(255);not null"`
	LastName  string    `json:"last_name" gorm:"type:varchar(255);not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (u *User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// NewUser creates a new user instance
func NewUser(email, firstName, lastName string) *User {
	return &User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// Activate activates the user
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Update updates the user information
func (u *User) Update(email, firstName, lastName string) {
	u.Email = email
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}
