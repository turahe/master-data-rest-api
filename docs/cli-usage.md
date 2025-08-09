# Master Data REST API - CLI Usage Guide

The Master Data REST API now uses Cobra CLI framework for a modern, intuitive command-line experience.

## Installation & Building

### Windows
```bash
# Build Windows executable
go build -o bin/master-data-api.exe main.go

# Or use the PowerShell build script
PowerShell -ExecutionPolicy Bypass -File scripts/build.ps1 -Clean

# Or use Makefile
make build-windows
```

### Linux/macOS
```bash
# Build for current platform
go build -o bin/master-data-api main.go

# Or use Makefile
make build        # Current platform
make build-linux  # Linux
make build-darwin # macOS
make build-all    # All platforms
```

## Available Commands

### 🚀 Server Management

#### Start the API Server
```bash
# Start with default settings
master-data-api serve

# Start with custom host and port
master-data-api serve --host localhost --port 9090

# Start without running migrations
master-data-api serve --migrate=false

# Start with debug logging
master-data-api serve --log-level debug --log-format json
```

### 🗄️ Database Management

#### Run Migrations
```bash
# Run all pending migrations
master-data-api migrate up

# Rollback the last migration
master-data-api migrate down

# Rollback specific number of migrations
master-data-api migrate down --step 2

# Check migration status
master-data-api migrate status

# Use custom migrations directory
master-data-api migrate up --migrations-dir ./custom-migrations
```

#### Database Seeding
```bash
# Seed database with sample data (uses configs/data by default)
master-data-api seed

# Seed specific data type
master-data-api seed --name languages
master-data-api seed --name banks
master-data-api seed --name currencies
master-data-api seed --name geodirectories

# TRUNCATE existing data and seed fresh (fast bulk deletion)
master-data-api seed --clear

# Seed specific data with TRUNCATE
master-data-api seed --name languages --clear

# Use custom data directory
master-data-api seed --data-dir ./custom-data

# Only seed without clearing
master-data-api seed --seed-only

# Combine flags to seed specific type with custom directory
master-data-api seed --name languages --data-dir ./my-data
```

**Available Flags:**
- `--name, -n`: Seed specific data type (languages, banks, currencies, geodirectories)
- `--data-dir, -d`: Directory containing seed data files (default: `configs/data`)
- `--clear, -c`: **TRUNCATE** existing data before seeding for efficient bulk deletion (default: `false`)
- `--seed-only`: Only seed data, don't clear existing data (default: `false`)

**Performance Features:**
- **🚀 TRUNCATE Operations**: Uses `TRUNCATE TABLE` instead of `DELETE` for fast bulk data clearing
- **📊 Progress Tracking**: Real-time progress logging during data seeding
- **🔄 Upsert Logic**: Prevents duplicate records by checking existing data before insertion
- **📁 Custom Data Sources**: Support for loading data from custom directories

**Available Seed Data:**
- **Languages** (185 records) - ISO language codes with names
- **Banks** (142 records) - Indonesian bank master data
- **Currencies** (168 records) - World currencies with symbols
- **Countries** (247 records) - World countries  
- **Geodirectories** - Indonesian administrative hierarchy:
  - Provinces (34 records)
  - Cities/Regencies (514 records) - Auto-classified by KOTA/KAB prefix
  - Districts (7,000+ records)
  - Villages (83,000+ records)

### 🔑 API Key Management

#### Create API Keys
```bash
# Create basic API key
master-data-api create-api-key

# Create with custom name and description
master-data-api create-api-key --name "Production Key" --description "API key for production environment"

# Create with expiration date
master-data-api create-api-key --expires "2024-12-31T23:59:59Z"

# Create with all options
master-data-api create-api-key \
  --name "Development Key" \
  --description "Key for development team" \
  --expires "2024-06-30T23:59:59Z"
```

### 🔍 Search Index Management

#### Initialize Search Indexes
```bash
# Initialize all search indexes for Meilisearch
master-data-api search init

# Check initialization status
master-data-api search health
```

#### Reindex Data
```bash
# Reindex all data to search engine
master-data-api search reindex

# View indexing statistics
master-data-api search stats
```

**Search Commands:**
- `search init` - Initialize search indexes for all entities
- `search reindex` - Reindex all data for search functionality
- `search health` - Check search service connectivity and status
- `search stats` - Display search index statistics and document counts

**Requirements:**
- Meilisearch server must be running and accessible
- Configure `MEILISEARCH_HOST` and `MEILISEARCH_API_KEY` environment variables
- Ensure database contains data before reindexing

### 📊 Utility Commands

#### Version Information
```bash
# Show detailed version information
master-data-api version
```

#### Help System
```bash
# General help
master-data-api --help

# Command-specific help
master-data-api serve --help
master-data-api migrate --help
master-data-api create-api-key --help
master-data-api seed --help
master-data-api search --help
```

## Global Flags

These flags work with all commands:

```bash
--config string       # Config file (default is .env)
--log-level string    # Log level: trace, debug, info, warn, error (default "info")
--log-format string   # Log format: text, json (default "text")
--log-output string   # Log output: stdout, stderr, or file path (default "stdout")
```

### Examples with Global Flags
```bash
# Use custom config file
master-data-api serve --config ./production.env

# Enable debug logging with JSON format
master-data-api serve --log-level debug --log-format json

# Log to file
master-data-api serve --log-output ./logs/api.log

# Combine multiple flags
master-data-api migrate up --log-level trace --config ./staging.env
```

## Environment Configuration

The CLI respects the following environment variables:

### Application Settings
```bash
APP_NAME=Master Data REST API
APP_ENV=development
APP_VERSION=1.0.0
```

### Server Configuration
```bash
APP_HOST=localhost
APP_PORT=8080
```

### Database Configuration
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=master_data
DB_USER=appuser
DB_PASSWORD=apppassword
DB_SSL_MODE=disable

# Database Logging
DB_LOG_LEVEL=info
DB_LOG_QUERIES=true
DB_LOG_SLOW_QUERY=100ms
```

### Logging Configuration
```bash
LOG_LEVEL=info
LOG_FORMAT=text
LOG_OUTPUT=stdout
```

## Development Workflow

### 1. Setup Development Environment
```bash
# Copy environment file
cp env.example .env

# Run migrations
master-data-api migrate up

# Create initial API key
master-data-api create-api-key --name "Development Key"

# Start server in development mode
master-data-api serve --log-level debug
```

### 2. Database Management
```bash
# Check current migration status
master-data-api migrate status

# Run new migrations
master-data-api migrate up

# Rollback if needed
master-data-api migrate down --step 1
```

### 3. Production Deployment
```bash
# Build for production
make build-all

# Run with production config
master-data-api serve --config ./production.env --log-format json
```

## Makefile Integration

The project includes an updated Makefile with Cobra integration:

```bash
# Build commands
make build          # Current platform
make build-windows  # Windows (.exe)
make build-linux    # Linux
make build-darwin   # macOS
make build-all      # All platforms

# Development commands
make run            # Start server
make run-dev        # Start with debug logging
make test           # Run tests
make lint           # Run linter
make format         # Format code

# Database commands
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make migrate-status # Check status

# Utility commands
make create-api-key # Create API key
make seed           # Seed database
make seed-clear     # Seed with clearing
```

## PowerShell Build Script (Windows)

For Windows users, use the PowerShell build script:

```powershell
# Build all platforms
.\scripts\build.ps1

# Build with clean
.\scripts\build.ps1 -Clean

# Build specific version
.\scripts\build.ps1 -Version "2.0.0"

# Build to custom directory
.\scripts\build.ps1 -BuildDir "release"
```

## Troubleshooting

### Common Issues

1. **Command not found**
   ```bash
   # Make sure you've built the executable
   go build -o bin/master-data-api.exe main.go
   ```

2. **Database connection errors**
   ```bash
   # Check your .env file configuration
   master-data-api migrate status --log-level debug
   ```

3. **Permission errors (Windows)**
   ```bash
   # Run PowerShell as Administrator or use:
   PowerShell -ExecutionPolicy Bypass -File scripts/build.ps1
   ```

4. **Migration failures**
   ```bash
   # Check migration status first
   master-data-api migrate status
   
   # Force to specific version if needed
   # (This requires manual intervention)
   ```

### Debug Mode

Enable debug mode for troubleshooting:

```bash
# Enable debug logging
master-data-api serve --log-level debug --log-format json

# Enable trace logging (most verbose)
master-data-api migrate up --log-level trace
```

## Migration from Old CLI

If you're migrating from the old CLI structure:

| Old Command | New Command |
|-------------|-------------|
| `go run cmd/api/main.go` | `master-data-api serve` |
| `go run cmd/create-api-key/main.go` | `master-data-api create-api-key` |
| `go run cmd/seeder/main.go` | `master-data-api seed` |
| Manual migration commands | `master-data-api migrate up/down/status` |

The new CLI provides better error handling, consistent flags, comprehensive help, and improved user experience.
