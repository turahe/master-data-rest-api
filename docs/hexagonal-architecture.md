# Hexagonal Architecture Implementation

This document explains the hexagonal architecture (also known as ports and adapters pattern) implementation in this Go project.

## Overview

Hexagonal architecture separates the business logic from external concerns by defining clear boundaries through interfaces (ports) and their implementations (adapters). This creates a highly testable and maintainable codebase.

## Architecture Layers

### 1. Domain Layer (Core Business Logic)

Located in `internal/domain/`, this is the heart of the application containing:

- **Entities** (`internal/domain/entities/`): Business objects with behavior
- **Value Objects** (`internal/domain/valueobjects/`): Immutable objects representing concepts
- **Services** (`internal/domain/services/`): Business logic and orchestration
- **Repository Interfaces** (`internal/domain/repositories/`): Data access contracts

**Key Principles:**
- Contains pure business logic
- No dependencies on external frameworks
- Defines the core domain model
- Implements business rules and validation

### 2. Ports Layer (Interfaces)

Located in `internal/ports/`, defines the contracts between layers:

- **Primary Ports** (`internal/ports/primary/`): Interfaces for driving adapters (HTTP, CLI)
- **Secondary Ports** (`internal/ports/secondary/`): Interfaces for driven adapters (Database, External APIs)

**Key Principles:**
- Defines clear contracts
- Enables dependency inversion
- Facilitates testing through mocking
- Separates concerns

### 3. Adapters Layer (Implementations)

Located in `internal/adapters/`, implements the ports:

- **Primary Adapters** (`internal/adapters/primary/`):
  - HTTP handlers
  - CLI commands
  - Application services
- **Secondary Adapters** (`internal/adapters/secondary/`):
  - Database implementations
  - External API clients
  - File system operations

**Key Principles:**
- Implements port interfaces
- Handles external concerns
- Can be easily swapped
- Contains framework-specific code

## Dependency Flow

```
HTTP Handler (Primary Adapter)
    ↓
Application Service (Primary Adapter)
    ↓
Domain Service (Domain)
    ↓
Repository Interface (Port)
    ↓
Database Implementation (Secondary Adapter)
```

## Benefits

### 1. Testability
- Business logic can be tested in isolation
- Dependencies can be easily mocked
- Unit tests are fast and reliable

### 2. Maintainability
- Clear separation of concerns
- Easy to understand and modify
- Changes in one layer don't affect others

### 3. Flexibility
- Easy to swap implementations
- Framework agnostic domain logic
- Multiple interfaces for the same functionality

### 4. Scalability
- Easy to add new features
- Clear boundaries for team collaboration
- Consistent patterns across the codebase

## Example: User Management

### Domain Entity
```go
// internal/domain/entities/user.go
type User struct {
    ID        uuid.UUID
    Email     string
    FirstName string
    LastName  string
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *User) Update(firstName, lastName string) {
    u.FirstName = firstName
    u.LastName = lastName
    u.UpdatedAt = time.Now()
}
```

### Repository Interface
```go
// internal/domain/repositories/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Domain Service
```go
// internal/domain/services/user_service.go
type UserService struct {
    userRepo repositories.UserRepository
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName string) (*entities.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    user.Update(firstName, lastName)
    
    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### Primary Port
```go
// internal/ports/primary/user_handler.go
type UserHandler interface {
    UpdateUser(ctx context.Context, id uuid.UUID, request UpdateUserRequest) (*entities.User, error)
}

type UpdateUserRequest struct {
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}
```

### Primary Adapter (HTTP Handler)
```go
// internal/adapters/primary/http/user_handler.go
type UserHTTPHandler struct {
    userService primary.UserHandler
}

func (h *UserHTTPHandler) UpdateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var request primary.UpdateUserRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    user, err := h.userService.UpdateUser(c.Request.Context(), id, request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}
```

### Secondary Adapter (Database Implementation)
```go
// internal/adapters/secondary/database/mysql/user_repository.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
    query := `
        UPDATE users 
        SET first_name = ?, last_name = ?, updated_at = ?
        WHERE id = ?
    `
    
    result, err := r.db.ExecContext(ctx, query,
        user.FirstName,
        user.LastName,
        user.UpdatedAt,
        user.ID,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}
```

## Testing Strategy

### 1. Domain Layer Testing
- Test entities and value objects in isolation
- Test domain services with mocked repositories
- Focus on business logic validation

### 2. Adapter Testing
- Test HTTP handlers with mocked services
- Test database implementations with test databases
- Use integration tests for external dependencies

### 3. End-to-End Testing
- Test complete user journeys
- Use test containers for external services
- Validate API contracts

## Best Practices

### 1. Dependency Injection
- Use interfaces for all dependencies
- Inject dependencies through constructors
- Avoid global state

### 2. Error Handling
- Define domain-specific errors
- Handle errors at appropriate layers
- Provide meaningful error messages

### 3. Validation
- Validate input at adapter boundaries
- Use domain rules for business validation
- Return clear validation errors

### 4. Logging and Monitoring
- Log at adapter boundaries
- Use structured logging
- Monitor external dependencies

## Migration Strategy

When migrating from a traditional layered architecture:

1. **Extract Domain Logic**: Move business logic to domain layer
2. **Define Interfaces**: Create ports for external dependencies
3. **Implement Adapters**: Create implementations for ports
4. **Update Dependencies**: Wire everything together
5. **Add Tests**: Ensure comprehensive test coverage

## Conclusion

Hexagonal architecture provides a robust foundation for building maintainable and testable applications. By clearly separating concerns and defining explicit boundaries, it enables teams to work efficiently and deliver high-quality software. 