package entities

import (
	"time"

	"github.com/google/uuid"
)

// Bank represents a bank entity
type Bank struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Alias     string    `json:"alias" db:"alias"`
	Company   string    `json:"company" db:"company"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the Bank entity
func (b *Bank) TableName() string {
	return "tm_banks"
}

// GenerateID generates a new UUID for the bank if not set
func (b *Bank) GenerateID() {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
}

// NewBank creates a new Bank instance
func NewBank(name, alias, company, code string) *Bank {
	return &Bank{
		ID:        uuid.New(),
		Name:      name,
		Alias:     alias,
		Company:   company,
		Code:      code,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetName sets the name of the bank
func (b *Bank) SetName(name string) {
	b.Name = name
	b.UpdatedAt = time.Now()
}

// SetAlias sets the alias of the bank
func (b *Bank) SetAlias(alias string) {
	b.Alias = alias
	b.UpdatedAt = time.Now()
}

// SetCompany sets the company of the bank
func (b *Bank) SetCompany(company string) {
	b.Company = company
	b.UpdatedAt = time.Now()
}

// SetCode sets the code of the bank
func (b *Bank) SetCode(code string) {
	b.Code = code
	b.UpdatedAt = time.Now()
}

// IsValid validates the bank entity
func (b *Bank) IsValid() bool {
	return b.Name != "" && b.Code != ""
}
