.PHONY: help build build-windows build-linux build-darwin build-all run test test-coverage test-short test-clean clean migrate-up migrate-down lint format docker-build docker-run

# Default target
help:
	@echo "Master Data REST API - Available Commands"
	@echo "========================================"
	@echo "Build Commands:"
	@echo "  build         - Build for current platform"
	@echo "  build-windows - Build Windows executable (.exe)"
	@echo "  build-linux   - Build Linux executable"
	@echo "  build-darwin  - Build macOS executable"
	@echo "  build-all     - Build for all platforms"
	@echo ""
	@echo "Development Commands:"
	@echo "  run           - Run the API server"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-short    - Run tests in short mode"
	@echo "  test-clean    - Clean test cache"
	@echo "  lint          - Run linter"
	@echo "  format        - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Database Commands:"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  migrate-status- Show migration status"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo ""
	@echo "Utility Commands:"
	@echo "  create-api-key- Create a new API key"
	@echo "  seed          - Seed database with sample data"

# Build the application for current platform
build:
	@echo "Building application for current platform..."
	go build -o bin/master-data-api main.go

# Build Windows executable
build-windows:
	@echo "Building Windows executable..."
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/master-data-api.exe main.go
	@echo "✅ Windows build completed: bin/master-data-api.exe"

# Build Linux executable
build-linux:
	@echo "Building Linux executable..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/master-data-api-linux main.go
	@echo "✅ Linux build completed: bin/master-data-api-linux"

# Build macOS executable
build-darwin:
	@echo "Building macOS executable..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/master-data-api-darwin main.go
	@echo "✅ macOS build completed: bin/master-data-api-darwin"

# Build for all platforms
build-all: build-windows build-linux build-darwin
	@echo "✅ All platform builds completed!"
	@ls -la bin/

# Run the API server
run:
	@echo "Starting Master Data REST API server..."
	go run main.go serve

# Run with custom settings
run-dev:
	@echo "Starting Master Data REST API server in development mode..."
	go run main.go serve --log-level debug --log-format text

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in short mode
test-short:
	@echo "Running tests (short mode)..."
	go test ./... -short

# Clean test cache
test-clean:
	@echo "Cleaning test cache..."
	go clean -testcache

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations using Cobra commands
migrate-up:
	@echo "Running database migrations..."
	go run main.go migrate up

migrate-down:
	@echo "Rolling back database migrations..."
	go run main.go migrate down

migrate-status:
	@echo "Checking migration status..."
	go run main.go migrate status

# Create API key
create-api-key:
	@echo "Creating new API key..."
	go run main.go create-api-key

# Seed database
seed:
	@echo "Seeding database with sample data..."
	go run main.go seed

seed-clear:
	@echo "Seeding database with clearing existing data..."
	go run main.go seed --clear

# Linting and formatting
lint:
	@echo "Running linter..."
	golangci-lint run

format:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t master-data-api .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env master-data-api

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 