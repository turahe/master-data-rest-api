package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/turahe/master-data-rest-api/internal/domain/entities"
)

// MockBankRepository is a mock implementation of BankRepository
type MockBankRepository struct {
	mock.Mock
}

func (m *MockBankRepository) Create(ctx context.Context, bank *entities.Bank) error {
	args := m.Called(ctx, bank)
	return args.Error(0)
}

func (m *MockBankRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Bank, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) GetByCode(ctx context.Context, code string) (*entities.Bank, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) GetByName(ctx context.Context, name string) (*entities.Bank, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) GetByAlias(ctx context.Context, alias string) (*entities.Bank, error) {
	args := m.Called(ctx, alias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) GetByCompany(ctx context.Context, company string, limit, offset int) ([]*entities.Bank, error) {
	args := m.Called(ctx, company, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Bank, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) Update(ctx context.Context, bank *entities.Bank) error {
	args := m.Called(ctx, bank)
	return args.Error(0)
}

func (m *MockBankRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBankRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Bank, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Bank), args.Error(1)
}

func (m *MockBankRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockBankRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockBankRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestNewBankService(t *testing.T) {
	// Given
	mockRepo := &MockBankRepository{}

	// When
	service := NewBankService(mockRepo)

	// Then
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.bankRepo)
}

func TestBankService_CreateBank(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		name := "Test Bank"
		alias := "TB"
		company := "Test Bank Corp"
		code := "001"

		mockRepo.On("ExistsByCode", ctx, code).Return(false, nil)
		mockRepo.On("ExistsByName", ctx, name).Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Bank")).Return(nil)

		// When
		bank, err := service.CreateBank(ctx, name, alias, company, code)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, bank)
		assert.Equal(t, name, bank.Name)
		assert.Equal(t, alias, bank.Alias)
		assert.Equal(t, company, bank.Company)
		assert.Equal(t, code, bank.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		// When
		bank, err := service.CreateBank(ctx, "", "alias", "company", "001")

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		assert.Contains(t, err.Error(), "bank name is required")
	})

	t.Run("empty code", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		// When
		bank, err := service.CreateBank(ctx, "Test Bank", "alias", "company", "")

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		assert.Contains(t, err.Error(), "bank code is required")
	})

	t.Run("duplicate code", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		code := "001"
		mockRepo.On("ExistsByCode", ctx, code).Return(true, nil)

		// When
		bank, err := service.CreateBank(ctx, "Test Bank", "alias", "company", code)

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		name := "Test Bank"
		code := "001"
		mockRepo.On("ExistsByCode", ctx, code).Return(false, nil)
		mockRepo.On("ExistsByName", ctx, name).Return(true, nil)

		// When
		bank, err := service.CreateBank(ctx, name, "alias", "company", code)

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_GetBankByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bankID := uuid.New()
		expectedBank := &entities.Bank{
			ID:   bankID,
			Name: "Test Bank",
			Code: "001",
		}

		mockRepo.On("GetByID", ctx, bankID).Return(expectedBank, nil)

		// When
		bank, err := service.GetBankByID(ctx, bankID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBank, bank)
		mockRepo.AssertExpectations(t)
	})

	t.Run("bank not found", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bankID := uuid.New()

		mockRepo.On("GetByID", ctx, bankID).Return(nil, assert.AnError)

		// When
		bank, err := service.GetBankByID(ctx, bankID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_GetBankByCode(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		code := "001"
		expectedBank := &entities.Bank{
			ID:   uuid.New(),
			Name: "Test Bank",
			Code: code,
		}

		mockRepo.On("GetByCode", ctx, code).Return(expectedBank, nil)

		// When
		bank, err := service.GetBankByCode(ctx, code)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBank, bank)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty code", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByCode", ctx, "").Return(nil, assert.AnError)

		// When
		bank, err := service.GetBankByCode(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Nil(t, bank)
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_GetAllBanks(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		limit, offset := 10, 0
		expectedBanks := []*entities.Bank{
			{ID: uuid.New(), Name: "Bank 1", Code: "001"},
			{ID: uuid.New(), Name: "Bank 2", Code: "002"},
		}

		mockRepo.On("GetAll", ctx, limit, offset).Return(expectedBanks, nil)

		// When
		banks, err := service.GetAllBanks(ctx, limit, offset)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBanks, banks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("with zero limit", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		limit, offset := 0, 0
		expectedBanks := []*entities.Bank{}

		mockRepo.On("GetAll", ctx, limit, offset).Return(expectedBanks, nil)

		// When
		banks, err := service.GetAllBanks(ctx, limit, offset)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBanks, banks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("with negative offset", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		limit, offset := 10, -1
		expectedBanks := []*entities.Bank{}

		mockRepo.On("GetAll", ctx, limit, offset).Return(expectedBanks, nil)

		// When
		banks, err := service.GetAllBanks(ctx, limit, offset)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBanks, banks)
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_UpdateBank(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bank := &entities.Bank{
			ID:      uuid.New(),
			Name:    "Updated Bank",
			Alias:   "UB",
			Company: "Updated Company",
			Code:    "002",
		}

		// Mock the existence checks - return error to indicate bank doesn't exist
		mockRepo.On("GetByCode", ctx, bank.Code).Return(nil, assert.AnError)
		mockRepo.On("GetByName", ctx, bank.Name).Return(nil, assert.AnError)
		mockRepo.On("Update", ctx, bank).Return(nil)

		// When
		err := service.UpdateBank(ctx, bank)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bank := &entities.Bank{
			ID:   uuid.New(),
			Name: "Test Bank",
			Code: "001",
		}

		// Mock the existence checks - return error to indicate bank doesn't exist
		mockRepo.On("GetByCode", ctx, bank.Code).Return(nil, assert.AnError)
		mockRepo.On("GetByName", ctx, bank.Name).Return(nil, assert.AnError)
		mockRepo.On("Update", ctx, bank).Return(assert.AnError)

		// When
		err := service.UpdateBank(ctx, bank)

		// Then
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_DeleteBank(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bankID := uuid.New()

		mockRepo.On("Delete", ctx, bankID).Return(nil)

		// When
		err := service.DeleteBank(ctx, bankID)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deletion error", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		bankID := uuid.New()

		mockRepo.On("Delete", ctx, bankID).Return(assert.AnError)

		// When
		err := service.DeleteBank(ctx, bankID)

		// Then
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBankService_SearchBanks(t *testing.T) {
	t.Run("successful search", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		query := "test"
		limit, offset := 10, 0
		expectedBanks := []*entities.Bank{
			{ID: uuid.New(), Name: "Test Bank", Code: "001"},
		}

		mockRepo.On("Search", ctx, query, limit, offset).Return(expectedBanks, nil)

		// When
		banks, err := service.SearchBanks(ctx, query, limit, offset)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBanks, banks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty query", func(t *testing.T) {
		// Given
		mockRepo := &MockBankRepository{}
		service := NewBankService(mockRepo)
		ctx := context.Background()
		query := ""
		limit, offset := 10, 0
		expectedBanks := []*entities.Bank{}

		mockRepo.On("Search", ctx, query, limit, offset).Return(expectedBanks, nil)

		// When
		banks, err := service.SearchBanks(ctx, query, limit, offset)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedBanks, banks)
		mockRepo.AssertExpectations(t)
	})
}
