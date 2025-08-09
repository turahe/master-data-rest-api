.PHONY: help build build-windows build-linux build-darwin build-all run test test-coverage test-short test-clean clean migrate-up migrate-down lint format docker-build docker-run swagger swagger-init swagger-gen swagger-serve swagger-validate swagger-clean swagger-watch install-tools

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
	@echo "Documentation Commands:"
	@echo "  swagger       - Generate and serve Swagger documentation"
	@echo "  swagger-init  - Initialize Swagger documentation"
	@echo "  swagger-gen   - Generate Swagger documentation files"
	@echo "  swagger-serve - Serve Swagger documentation locally"
	@echo "  swagger-validate - Validate Swagger documentation"
	@echo "  swagger-clean - Clean generated Swagger files"
	@echo "  swagger-watch - Watch for changes and auto-regenerate docs"
	@echo "  install-tools - Install required development tools"
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

# Swagger Documentation Commands
swagger: swagger-gen swagger-serve
	@echo "✅ Swagger documentation generated and served!"

swagger-init:
	@echo "Initializing Swagger documentation..."
	@if [ ! -f docs/docs.go ]; then \
		mkdir -p docs; \
		swag init -g cmd/serve.go -o docs --parseDependency --parseInternal; \
	else \
		echo "Swagger already initialized. Use 'make swagger-gen' to regenerate."; \
	fi

swagger-gen:
	@echo "Generating Swagger documentation..."
	@mkdir -p docs
	swag init -g cmd/serve.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger documentation generated in docs/ directory"
	@echo "📁 Generated files:"
	@echo "   - docs/docs.go"
	@echo "   - docs/swagger.json"
	@echo "   - docs/swagger.yaml"

swagger-serve:
	@echo "Starting API server with Swagger UI..."
	@echo "🌐 Swagger UI will be available at: http://localhost:8080/swagger/index.html"
	@echo "📖 API Documentation: http://localhost:8080/swagger/doc.json"
	@echo "🚀 Starting server..."
	go run main.go serve

swagger-validate:
	@echo "Validating Swagger documentation..."
	@if [ ! -f docs/swagger.json ]; then \
		echo "❌ Swagger documentation not found. Run 'make swagger-gen' first."; \
		exit 1; \
	fi
	@echo "✅ Checking JSON syntax..."
	@cat docs/swagger.json | python -m json.tool > /dev/null && echo "✅ swagger.json is valid JSON"
	@echo "✅ Checking YAML syntax..."
	@python -c "import yaml; yaml.safe_load(open('docs/swagger.yaml'))" && echo "✅ swagger.yaml is valid YAML"
	@echo "📊 API Statistics:"
	@echo "   - Endpoints: $$(grep -c '"' docs/swagger.json | head -1)"
	@echo "   - Tags: $$(grep -o '"tags":\[.*\]' docs/swagger.json | wc -l)"

swagger-clean:
	@echo "Cleaning generated Swagger files..."
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
	@echo "✅ Swagger documentation files cleaned"

swagger-watch:
	@echo "Starting Swagger documentation in watch mode..."
	@echo "🔄 Watching for changes and regenerating documentation..."
	@while true; do \
		find . -name "*.go" -path "./internal/adapters/primary/http/*" -newer docs/swagger.json 2>/dev/null | grep -q . && make swagger-gen; \
		sleep 2; \
	done

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	@echo "✅ Development tools installed successfully!" 