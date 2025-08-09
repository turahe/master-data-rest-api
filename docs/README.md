# 📚 Master Data REST API Documentation

Welcome to the comprehensive documentation for the Master Data REST API. This documentation is organized by topic to help you quickly find the information you need.

## 📖 Quick Navigation

### 🚀 Getting Started
- [Main README](../README.md) - Overview, quick start, and basic setup
- [Installation Guide](installation.md) - Detailed installation instructions
- [Configuration Guide](#configuration-guide) - Environment setup and configuration

### 🏗️ Architecture & Design
- [Hexagonal Architecture](hexagonal-architecture.md) - Detailed architecture explanation
- [Project Structure](#project-structure) - Code organization and patterns
- [Database Design](#database-design) - Schema and relationships

### 🎯 CLI & Tools
- [CLI Usage Guide](cli-usage.md) - Comprehensive command-line interface documentation
- [Migration Guide](#migration-guide) - Database migration management
- [Build Scripts](#build-scripts) - Cross-platform building

### 🔧 Technical Guides
- [Database Logging](database-logging.md) - Query logging and monitoring
- [Meilisearch Integration](meilisearch.md) - Search engine setup and usage
- [Testing Guide](testing.md) - Unit testing with comprehensive coverage
- [API Authentication](#api-authentication) - Security and access control
- [Performance Tuning](#performance-tuning) - Optimization guidelines

### 🌐 API Documentation
- [REST API Reference](#rest-api-reference) - Complete endpoint documentation
- [Swagger UI](http://localhost:8080/swagger/) - Interactive API explorer
- [API Examples](api-examples.md) - Common usage patterns and code samples

### 💻 Development
- [Development Setup](#development-setup) - Local development environment
- [Testing Guide](#testing-guide) - Unit and integration testing
- [Contributing](#contributing) - How to contribute to the project

### 🚀 Deployment
- [Deployment Guide](#deployment-guide) - Production deployment strategies
- [Docker Setup](#docker-setup) - Containerization and orchestration
- [Monitoring](#monitoring) - Logging and observability

## 📋 Documentation Index

### Core Documentation Files

| File | Description | Status |
|------|-------------|--------|
| [`README.md`](../README.md) | Main project documentation | ✅ Complete |
| [`installation.md`](installation.md) | Detailed installation guide | ✅ Complete |
| [`cli-usage.md`](cli-usage.md) | CLI command reference | ✅ Complete |
| [`database-logging.md`](database-logging.md) | Database logging features | ✅ Complete |
| [`testing.md`](testing.md) | Testing strategy and guidelines | ✅ Complete |
| [`api-examples.md`](api-examples.md) | API usage examples and patterns | ✅ Complete |
| [`meilisearch.md`](meilisearch.md) | Search engine integration guide | ✅ Complete |
| [`hexagonal-architecture.md`](hexagonal-architecture.md) | Architecture deep dive | ✅ Complete |

### API Documentation

| Type | Location | Description |
|------|----------|-------------|
| Interactive Swagger | `http://localhost:8080/swagger/` | Live API documentation |
| OpenAPI JSON | [`swagger.json`](swagger.json) | Machine-readable API spec |
| OpenAPI YAML | [`swagger.yaml`](swagger.yaml) | Human-readable API spec |

### Generated Documentation

| File | Description | Auto-Generated |
|------|-------------|----------------|
| [`docs.go`](docs.go) | Swagger Go bindings | ✅ Yes |
| Coverage reports | Test coverage analysis | ✅ Yes (when running tests) |

## 🎯 Quick Reference

### Common Commands
```bash
# Start the server
./master-data-api serve

# Run migrations
./master-data-api migrate up

# Create API key
./master-data-api create-api-key --name "My Key"

# Seed data with TRUNCATE optimization
./master-data-api seed --clear

# Initialize search indexes
./master-data-api search init

# Get help for any command
./master-data-api [command] --help
```

### Important URLs (Development)
- **API Base**: `http://localhost:8080/api/v1`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/health`

### Key Configuration Files
- [`.env.example`](../env.example) - Environment variables template
- [`Makefile`](../Makefile) - Build and development commands
- [`docker-compose.yml`](../docker-compose.yml) - Docker setup

## 📊 Documentation Status

### ✅ Complete Documentation
- [x] Main README with quick start
- [x] Detailed installation guide
- [x] CLI usage guide with examples
- [x] Database logging configuration
- [x] API examples and usage patterns
- [x] Architecture overview
- [x] Auto-generated API documentation
- [x] Meilisearch integration guide
- [x] Testing guide with unit tests
- [x] Docker Hub documentation

### 🚧 Planned Documentation
- [ ] Performance tuning guide
- [ ] Security best practices
- [ ] Troubleshooting guide
- [ ] API rate limiting guide
- [ ] Backup and recovery procedures
- [ ] Docker deployment guide

## 🛠️ Documentation Maintenance

### For Developers
- Update CLI docs when adding new commands
- Regenerate Swagger docs after API changes: `swag init`
- Keep architecture docs in sync with code changes
- Add examples for new features

### For Contributors
- Follow markdown standards
- Include code examples where applicable
- Update the documentation index when adding new docs
- Test all code examples before committing

## 📱 Documentation Tools

### Generation
- **Swagger/OpenAPI**: Auto-generated from code annotations
- **CLI Help**: Auto-generated from Cobra commands
- **Markdown**: Manual documentation using standard markdown

### Validation
- **Markdown Lint**: For consistent formatting
- **Link Checker**: Validates internal and external links
- **Code Examples**: All examples are tested in CI

## 🔗 External Resources

### Go & Framework Documentation
- [Go Documentation](https://golang.org/doc/)
- [Fiber Framework](https://docs.gofiber.io/)
- [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [Cobra CLI](https://cobra.dev/)

### Architecture & Patterns
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [PostgreSQL Nested Sets](https://www.postgresql.org/docs/current/ltree.html)

### Database & Migrations
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://www.postgresql.org/docs/current/ddl-alter.html)

## 📝 Document Templates

### API Endpoint Documentation Template
```markdown
### POST /api/v1/resource

Description of what this endpoint does.

**Request:**
```json
{
  "field": "value"
}
```

**Response:**
```json
{
  "success": true,
  "data": {...}
}
```

**Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/resource" \
     -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"field": "value"}'
```
```

### CLI Command Documentation Template
```markdown
### command-name

Description of what this command does.

**Usage:**
```bash
./master-data-api command-name [flags]
```

**Flags:**
- `--flag-name` - Description of flag
- `--another-flag` - Description of another flag

**Examples:**
```bash
# Basic usage
./master-data-api command-name

# With flags
./master-data-api command-name --flag-name value
```
```

---

<div align="center">
  <strong>📚 Keep documentation up-to-date and helpful!</strong>
  <br>
  <sub>Good documentation is as important as good code</sub>
</div>
