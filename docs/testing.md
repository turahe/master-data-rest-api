# Testing Guide

This document outlines the testing strategy and guidelines for the Master Data REST API.

## Testing Philosophy

We follow a comprehensive testing approach with multiple layers:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete user workflows
- **Test Coverage**: Aim for 80%+ code coverage

## Test Structure

### 1. Domain Layer Testing (100% Coverage)

#### Entities
- **Location**: `internal/domain/entities/*_test.go`
- **Coverage**: 100%
- **Tests Include**:
  - Constructor functions
  - Business logic methods
  - Validation functions
  - State changes and updates

#### Value Objects
- **Location**: `internal/domain/valueobjects/*_test.go`
- **Coverage**: 100%
- **Tests Include**:
  - Validation logic (e.g., email format)
  - Immutability guarantees
  - String representations

#### Services
- **Location**: `internal/domain/services/*_test.go`
- **Coverage**: 9.3% (only BankService tested)
- **Tests Include**:
  - Business logic validation
  - Repository interaction mocking
  - Error handling scenarios

### 2. Package Layer Testing

#### Response Package
- **Location**: `pkg/response/response_test.go`
- **Coverage**: 100%
- **Tests Include**:
  - All response types (Success, Error, Created, etc.)
  - Status code validation
  - JSON structure verification

#### Logger Package
- **Location**: `pkg/logger/logger_test.go`
- **Coverage**: 26.1%
- **Tests Include**:
  - Different log levels and formats
  - JSON vs Text formatting
  - Structured logging with fields
  - Configuration handling

## Test Commands

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests in short mode (faster)
make test-short

# Clean test cache
make test-clean
```

### Coverage Reports

Test coverage reports are generated in HTML format:

```bash
make test-coverage
# Opens coverage.html in browser
```

## Current Test Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/domain/entities` | 100% | ✅ Complete |
| `internal/domain/valueobjects` | 100% | ✅ Complete |
| `pkg/response` | 100% | ✅ Complete |
| `internal/domain/services` | 9.3% | 🔄 Partial (only BankService) |
| `pkg/logger` | 26.1% | 🔄 Partial |
| `internal/adapters/primary/http` | 0% | ❌ Missing |
| `internal/adapters/secondary/database` | 0% | ❌ Missing |

## Testing Patterns

### 1. Entity Testing

```go
func TestNewBank(t *testing.T) {
    // Given
    name := "Test Bank"
    code := "001"
    
    // When
    bank := NewBank(name, "alias", "company", code)
    
    // Then
    assert.NotNil(t, bank)
    assert.Equal(t, name, bank.Name)
    assert.Equal(t, code, bank.Code)
}
```

### 2. Service Testing with Mocks

```go
type MockBankRepository struct {
    mock.Mock
}

func (m *MockBankRepository) Create(ctx context.Context, bank *entities.Bank) error {
    args := m.Called(ctx, bank)
    return args.Error(0)
}

func TestBankService_CreateBank(t *testing.T) {
    // Given
    mockRepo := &MockBankRepository{}
    service := NewBankService(mockRepo)
    
    mockRepo.On("ExistsByCode", mock.Anything, "001").Return(false, nil)
    mockRepo.On("ExistsByName", mock.Anything, "Test Bank").Return(false, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Bank")).Return(nil)
    
    // When
    bank, err := service.CreateBank(context.Background(), "Test Bank", "TB", "Test Corp", "001")
    
    // Then
    assert.NoError(t, err)
    assert.NotNil(t, bank)
    mockRepo.AssertExpectations(t)
}
```

### 3. HTTP Response Testing

```go
func TestSuccess(t *testing.T) {
    // Given
    app := fiber.New()
    data := map[string]string{"key": "value"}
    message := "Operation successful"
    
    app.Get("/test", func(c *fiber.Ctx) error {
        return Success(c, data, message)
    })
    
    // When
    req := httptest.NewRequest("GET", "/test", nil)
    resp, err := app.Test(req)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    var response Response
    err = json.NewDecoder(resp.Body).Decode(&response)
    assert.NoError(t, err)
    assert.True(t, response.Success)
    assert.Equal(t, message, response.Message)
}
```

## Test Dependencies

The project uses the following testing libraries:

- **`github.com/stretchr/testify`**: Assertions and mocking
- **`github.com/DATA-DOG/go-sqlmock`**: Database mocking (planned)
- **`net/http/httptest`**: HTTP testing utilities

## Pending Test Implementation

### High Priority
1. **Service Layer**: Complete testing for all services (Currency, Language, Geodirectory, APIKey)
2. **Repository Layer**: Database interaction testing with sqlmock
3. **HTTP Handlers**: API endpoint testing with mocked services

### Medium Priority
1. **Middleware**: Authentication and logging middleware tests
2. **Database Layer**: Connection and migration testing
3. **Integration Tests**: Full API workflow testing

### Low Priority
1. **Performance Tests**: Load and stress testing
2. **Security Tests**: Authentication and authorization testing
3. **Contract Tests**: API contract validation

## Best Practices

### Test Organization
- One test file per source file (`*_test.go`)
- Group related tests using subtests
- Use descriptive test names that explain the scenario

### Test Structure (AAA Pattern)
```go
func TestFeature(t *testing.T) {
    // Arrange (Given)
    // Setup test data and mocks
    
    // Act (When)
    // Execute the function under test
    
    // Assert (Then)
    // Verify the results
}
```

### Mocking Guidelines
- Mock external dependencies (databases, APIs)
- Keep mocks simple and focused
- Verify mock expectations are met
- Use dependency injection for testability

### Test Data
- Use meaningful test data that represents real scenarios
- Create test helpers for common setup
- Avoid hardcoded values where possible

## Continuous Integration

Tests are automatically run on:
- Pull request creation/updates
- Commits to main branch
- Scheduled nightly builds

Coverage reports are generated and stored for tracking progress over time.

## Troubleshooting

### Common Issues

1. **Import Cycles**: Ensure test files don't create circular dependencies
2. **Mock Setup**: Verify all expected method calls are properly mocked
3. **Test Isolation**: Ensure tests don't depend on external state
4. **Time-based Tests**: Use fixed time values or time mocking for consistency

### Debugging Tests

```bash
# Run specific test
go test -run TestBankService_CreateBank ./internal/domain/services/

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

## Contributing

When adding new features:

1. Write tests first (TDD approach)
2. Ensure new code has >80% test coverage
3. Update this documentation for new testing patterns
4. Run the full test suite before committing

For test improvements or new testing utilities, please discuss with the team first to ensure consistency with our testing strategy.
