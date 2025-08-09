# Master Data REST API

[![Docker Pulls](https://img.shields.io/docker/pulls/turahe/master-data-rest-api)](https://hub.docker.com/r/turahe/master-data-rest-api)
[![Docker Image Size](https://img.shields.io/docker/image-size/turahe/master-data-rest-api/latest)](https://hub.docker.com/r/turahe/master-data-rest-api)
[![Docker Image Version](https://img.shields.io/docker/v/turahe/master-data-rest-api?sort=semver)](https://hub.docker.com/r/turahe/master-data-rest-api)

A high-performance REST API built with Go, Fiber, and PostgreSQL for managing master data including geodirectories, banks, currencies, and languages. Features hexagonal architecture, comprehensive testing, and production-ready containerization.

## 🚀 Quick Start

### Pull and Run

```bash
# Pull the latest image
docker pull turahe/master-data-rest-api:latest

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

# Start all services
docker-compose up -d
```

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

## 🌐 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### API Documentation
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI JSON**: `http://localhost:8080/swagger/doc.json`

### Core Endpoints
- **Geodirectories**: `/api/v1/geodirectories`
- **Banks**: `/api/v1/banks`
- **Currencies**: `/api/v1/currencies`
- **Languages**: `/api/v1/languages`

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

### 3. Run Database Migrations
```bash
docker run --rm \
  --network master-data-network \
  -e DB_HOST=postgres \
  -e DB_USER=appuser \
  -e DB_PASSWORD=apppassword \
  -e DB_NAME=master_data \
  turahe/master-data-rest-api:latest migrate up
```

### 4. Start API Application
```bash
docker run -d \
  --name master-data-api \
  --network master-data-network \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=appuser \
  -e DB_PASSWORD=apppassword \
  -e DB_NAME=master_data \
  turahe/master-data-rest-api:latest
```

### 5. Create API Key
```bash
docker exec -it master-data-api ./main create-api-key
```

### 6. Test API
```bash
# Health check
curl http://localhost:8080/health

# Get banks (with API key)
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/banks
```

## 🛠️ Advanced Usage

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

- **Framework**: Go Fiber (high-performance HTTP framework)
- **Database**: PostgreSQL with pgx driver
- **Architecture**: Hexagonal (Ports & Adapters)
- **Authentication**: API Key based
- **Logging**: Structured logging with Logrus
- **Testing**: Comprehensive unit tests with testify

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
