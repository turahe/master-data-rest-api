# Master Data REST API

🐳 **Docker Hub**: `turahe/master-data-rest-api`  
🐙 **GitHub Container Registry**: `ghcr.io/turahe/master-data-rest-api`

[![Docker Pulls](https://img.shields.io/docker/pulls/turahe/master-data-rest-api)](https://hub.docker.com/r/turahe/master-data-rest-api)
[![Docker Image Size](https://img.shields.io/docker/image-size/turahe/master-data-rest-api/latest)](https://hub.docker.com/r/turahe/master-data-rest-api)
[![Docker Image Version](https://img.shields.io/docker/v/turahe/master-data-rest-api?sort=semver)](https://hub.docker.com/r/turahe/master-data-rest-api)
[![GHCR](https://img.shields.io/badge/ghcr-container-blue)](https://github.com/turahe/master-data-rest-api/pkgs/container/master-data-rest-api)

A modern, high-performance REST API built with **Go**, **Fiber**, and **PostgreSQL** for managing master data including geographical information, banks, currencies, and languages. Features **hexagonal architecture**, **Meilisearch integration**, **CLI tools**, **TRUNCATE-optimized seeding**, and production-ready containerization.

## 🚀 Quick Start

### Pull and Run

```bash
# Pull from Docker Hub
docker pull turahe/master-data-rest-api:latest

# Or pull from GitHub Container Registry
docker pull ghcr.io/turahe/master-data-rest-api:latest

# Run with PostgreSQL (requires database)
docker run -d \
  --name master-data-api \
  -p 8080:8080 \
  -e DB_HOST=your-postgres-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e DB_NAME=master_data \
  turahe/master-data-rest-api:latest
```

### With Docker Compose (Recommended)

```bash
# Create docker-compose.yml
curl -O https://raw.githubusercontent.com/turahe/master-data-rest-api/main/docker-compose.yml

# Start all services (includes PostgreSQL and Meilisearch)
docker-compose up -d

# Check service status
docker-compose ps
```

## ✨ Key Features

- 🏗️ **Hexagonal Architecture** - Clean separation of concerns with ports & adapters
- 🌍 **Geographic Hierarchy** - Nested set model for efficient hierarchical queries
- 🏦 **Master Data Management** - Banks, currencies, languages with comprehensive APIs
- 🔍 **Meilisearch Integration** - Fast, typo-tolerant search with automatic indexing
- 🚀 **TRUNCATE Optimization** - Efficient bulk data operations for seeding
- 🎯 **Modern CLI Tools** - Comprehensive command-line interface with Cobra
- 📊 **Database Logging** - Query performance monitoring and structured logging
- 🔐 **API Key Authentication** - Secure access control with Bearer token support
- 📖 **Auto-Generated Docs** - Swagger/OpenAPI 3.0 compliant documentation
- 🧪 **Comprehensive Testing** - 100% domain layer coverage with unit tests

## 📋 Available Tags

| Tag | Description | Size |
|-----|-------------|------|
| `latest` | Latest stable release | ~15MB |
| `v1.0.0` | Version 1.0.0 | ~15MB |

## 🔧 Environment Variables

### Required Variables
| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` | `postgres` |
| `DB_PORT` | PostgreSQL port | `5432` | `5432` |
| `DB_USER` | Database user | `appuser` | `postgres` |
| `DB_PASSWORD` | Database password | `apppassword` | `secretpass` |
| `DB_NAME` | Database name | `master_data` | `master_data` |

### Optional Variables
| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `APP_PORT` | Application port | `8080` | Any valid port |
| `APP_ENV` | Environment | `production` | `development`, `staging`, `production` |
| `LOG_LEVEL` | Log level | `info` | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | Log format | `json` | `json`, `text` |
| `MEILISEARCH_HOST` | Meilisearch server URL | `http://localhost:7700` | Meilisearch endpoint |
| `MEILISEARCH_API_KEY` | Meilisearch API key | `` | Search service key |

## 🌐 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### API Documentation
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI JSON**: `http://localhost:8080/swagger/doc.json`

### Core Endpoints
- **Geodirectories**: `/api/v1/geodirectories` - Hierarchical geographic data (countries → provinces → cities/regencies → districts → villages)
- **Banks**: `/api/v1/banks` - Indonesian bank master data with search capabilities
- **Currencies**: `/api/v1/currencies` - World currencies with symbols and status management
- **Languages**: `/api/v1/languages` - ISO language codes with localization support
- **API Keys**: `/api/v1/api-keys` - API key management and authentication

## 🔐 Authentication

API uses API key authentication. Create an API key:

```bash
# Using Docker exec
docker exec -it master-data-api ./main create-api-key

# Or with environment variables
docker run --rm \
  -e DB_HOST=your-postgres-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e DB_NAME=master_data \
  turahe/master-data-rest-api:latest create-api-key
```

Use the API key in requests:
```bash
# Using Authorization Bearer (recommended)
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/api/v1/banks

# Using X-API-Key header (alternative)
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/banks
```

## 📊 Complete Setup Example

### 1. Create Network
```bash
docker network create master-data-network
```

### 2. Start PostgreSQL
```bash
docker run -d \
  --name postgres \
  --network master-data-network \
  -e POSTGRES_DB=master_data \
  -e POSTGRES_USER=appuser \
  -e POSTGRES_PASSWORD=apppassword \
  -p 5432:5432 \
  postgres:17
```

### 3. Start Meilisearch (Optional but Recommended)
```bash
docker run -d \
  --name meilisearch \
  --network master-data-network \
  -e MEILI_MASTER_KEY=masterKey123 \
  -p 7700:7700 \
  getmeili/meilisearch:v1.5
```

### 4. Run Database Migrations
```bash
docker run --rm \
  --network master-data-network \
  -e DB_HOST=postgres \
  -e DB_USER=appuser \
  -e DB_PASSWORD=apppassword \
  -e DB_NAME=master_data \
  turahe/master-data-rest-api:latest migrate up
```

### 5. Start API Application
```bash
docker run -d \
  --name master-data-api \
  --network master-data-network \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=appuser \
  -e DB_PASSWORD=apppassword \
  -e DB_NAME=master_data \
  -e MEILISEARCH_HOST=http://meilisearch:7700 \
  -e MEILISEARCH_API_KEY=masterKey123 \
  turahe/master-data-rest-api:latest
```

### 6. Create API Key
```bash
docker exec -it master-data-api ./main create-api-key --name "Docker Setup"
```

### 7. Seed Sample Data (Optional)
```bash
# Seed all data types using TRUNCATE for efficient bulk operations
docker exec -it master-data-api ./main seed --clear

# Or seed specific data types
docker exec -it master-data-api ./main seed --name languages
docker exec -it master-data-api ./main seed --name banks
docker exec -it master-data-api ./main seed --name currencies
docker exec -it master-data-api ./main seed --name geodirectories
```

### 8. Initialize Search Indexes (Optional)
```bash
# Initialize Meilisearch indexes
docker exec -it master-data-api ./main search init

# Reindex all data for search
docker exec -it master-data-api ./main search reindex
```

### 9. Test API
```bash
# Health check
curl http://localhost:8080/health

# Get banks (with API key)
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/api/v1/banks

# Search geodirectories
curl -H "Authorization: Bearer your-api-key" \
     "http://localhost:8080/api/v1/geodirectories/search?q=jakarta"

# Get Swagger documentation
curl http://localhost:8080/swagger/index.html
```

## 🛠️ Advanced Usage

### CLI Tools
The container includes comprehensive CLI tools for database and application management:

```bash
# Database migrations
docker exec -it master-data-api ./main migrate up
docker exec -it master-data-api ./main migrate down --step 2
docker exec -it master-data-api ./main migrate status

# API key management
docker exec -it master-data-api ./main create-api-key --name "Production Key"
docker exec -it master-data-api ./main create-api-key --name "Temp Key" --expires "2024-12-31T23:59:59Z"

# Data seeding with TRUNCATE optimization
docker exec -it master-data-api ./main seed --clear              # TRUNCATE and seed all
docker exec -it master-data-api ./main seed --name languages     # Seed languages only
docker exec -it master-data-api ./main seed --seed-only          # Seed without clearing

# Search management
docker exec -it master-data-api ./main search init              # Initialize indexes
docker exec -it master-data-api ./main search reindex          # Reindex all data
docker exec -it master-data-api ./main search health           # Check search status
docker exec -it master-data-api ./main search stats            # View statistics

# Application info
docker exec -it master-data-api ./main version                 # Show version
docker exec -it master-data-api ./main --help                  # Show help
```

### Available Seed Data
When using the `seed` command, the following datasets are available:
- **Languages** (185 records) - ISO language codes with names
- **Banks** (142 records) - Indonesian bank master data  
- **Currencies** (168 records) - World currencies with symbols
- **Countries** (247 records) - World countries
- **Geodirectories** - Indonesian administrative hierarchy:
  - Provinces (34 records)
  - Cities/Regencies (514 records) - Auto-classified by KOTA/KAB prefix
  - Districts (7,000+ records)
  - Villages (83,000+ records)

### Custom Configuration
```bash
# Run with custom configuration
docker run -d \
  --name master-data-api \
  -p 8080:8080 \
  -e APP_PORT=3000 \
  -e LOG_LEVEL=debug \
  -e LOG_FORMAT=text \
  -e DB_HOST=your-postgres \
  -e DB_USER=your-user \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=your-database \
  -e MEILISEARCH_HOST=http://your-meilisearch:7700 \
  -e MEILISEARCH_API_KEY=your-search-key \
  turahe/master-data-rest-api:latest
```

### Volume Mounting
```bash
# Mount configuration files
docker run -d \
  --name master-data-api \
  -p 8080:8080 \
  -v /path/to/config:/app/config \
  -v /path/to/logs:/app/logs \
  turahe/master-data-rest-api:latest
```

### Health Checks
```bash
# Run with health check
docker run -d \
  --name master-data-api \
  -p 8080:8080 \
  --health-cmd="curl -f http://localhost:8080/health || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  -e DB_HOST=postgres \
  turahe/master-data-rest-api:latest
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Connection Refused
```bash
# Check if container is running
docker ps

# Check logs
docker logs master-data-api

# Verify network connectivity
docker exec -it master-data-api ping postgres
```

#### 2. Database Connection Failed
```bash
# Verify database credentials
docker exec -it master-data-api env | grep DB_

# Test database connection
docker exec -it postgres psql -U appuser -d master_data -c "SELECT 1;"
```

#### 3. Migration Issues
```bash
# Check migration status
docker exec -it master-data-api ./main migrate status

# Force migration version (if needed)
docker exec -it master-data-api ./main migrate force -v 1
```

### Logs and Debugging
```bash
# View application logs
docker logs master-data-api

# Follow logs in real-time
docker logs -f master-data-api

# Access container shell
docker exec -it master-data-api sh

# Check application version
docker exec -it master-data-api ./main version
```

## 📚 Documentation

- **GitHub Repository**: [turahe/master-data-rest-api](https://github.com/turahe/master-data-rest-api)
- **API Documentation**: Available at `/swagger/index.html` when running
- **Architecture Guide**: See repository documentation
- **Testing Guide**: Comprehensive test suite with 100% domain coverage

## 🏗️ Architecture

- **Framework**: Go Fiber v2 (high-performance HTTP framework)
- **Database**: PostgreSQL 13+ with pgx driver for optimal performance
- **Search Engine**: Meilisearch integration for fast, typo-tolerant search
- **Architecture**: Hexagonal (Ports & Adapters) for clean separation of concerns
- **Authentication**: API Key based with Bearer token support
- **Logging**: Structured logging with Logrus and database query logging
- **CLI Tools**: Cobra-powered command-line interface
- **Testing**: Comprehensive unit tests with testify (100% domain layer coverage)
- **Geographic Data**: Nested set model for hierarchical data with automatic type classification
- **Performance**: TRUNCATE-optimized seeding for efficient bulk operations

## 🔄 Updates

This image is automatically updated when new releases are published to the GitHub repository. 

- **Automated Builds**: Connected to GitHub repository
- **Multi-Stage Build**: Optimized for production use
- **Security**: Runs as non-root user
- **Size Optimized**: Alpine-based final image (~15MB)

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/turahe/master-data-rest-api/issues)
- **Documentation**: [GitHub Repository](https://github.com/turahe/master-data-rest-api)
- **Docker Hub**: [turahe/master-data-rest-api](https://hub.docker.com/r/turahe/master-data-rest-api)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/turahe/master-data-rest-api/blob/main/LICENSE) file for details.
