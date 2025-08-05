.PHONY: help build run test clean migrate-up migrate-down lint format docker-build docker-run seed seed-clear seed-both

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  seed         - Seed data from JSON files"
	@echo "  seed-clear   - Clear all seeded data"
	@echo "  seed-both    - Clear and seed data"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app cmd/api/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/api/main.go

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

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?multiStatements=true" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?multiStatements=true" down

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