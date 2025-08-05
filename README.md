# Master Data REST API

A Go-based REST API built with hexagonal architecture (ports and adapters pattern) for managing master data.

## Architecture Overview

This project follows the hexagonal architecture pattern, which separates the business logic from external concerns:

```
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── domain/            # Business logic and entities
│   ├── ports/             # Interfaces (ports) for adapters
│   └── adapters/          # Implementation of ports
│       ├── primary/       # Driving adapters (HTTP, CLI)
│       └── secondary/     # Driven adapters (Database, External APIs)
├── pkg/                   # Public packages that can be imported
├── migrations/            # Database migration files
├── configs/              # Configuration files
└── docs/                 # Documentation
```

## Project Structure

### Domain Layer (Core Business Logic)
- `internal/domain/entities/` - Business entities
- `internal/domain/repositories/` - Repository interfaces
- `internal/domain/services/` - Business services
- `internal/domain/valueobjects/` - Value objects

### Ports Layer (Interfaces)
- `internal/ports/primary/` - Primary (driving) port interfaces
- `internal/ports/secondary/` - Secondary (driven) port interfaces

### Adapters Layer (Implementations)
- `internal/adapters/primary/` - HTTP handlers, CLI commands
- `internal/adapters/secondary/` - Database implementations, external API clients

## Getting Started

### Prerequisites
- Go 1.21 or higher
- MySQL 8.0 or higher
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd master-data-rest-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations:
```bash
make migrate-up
```

5. Start the application:
```bash
make run
```

## Available Commands

- `make run` - Start the application
- `make test` - Run tests
- `make build` - Build the application
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make lint` - Run linter
- `make format` - Format code

## API Documentation

Once the application is running, you can access:
- API endpoints: `http://localhost:8080/api/v1`
- Health check: `http://localhost:8080/health`
- Swagger documentation: `http://localhost:8080/swagger/index.html`

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/domain/...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 