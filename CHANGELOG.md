# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

🎉 **Initial Release** - Master Data REST API v1.0.0

### 🚀 Added

#### Core API Features
- **RESTful API** with comprehensive master data endpoints
- **Hexagonal Architecture** implementation for maintainable and testable code
- **PostgreSQL Integration** with pgx driver for high performance
- **API Key Authentication** with Bearer token support
- **Swagger Documentation** with interactive API explorer
- **Health Check Endpoint** for monitoring and load balancers

#### Master Data Entities
- **🌍 Geodirectories** - Hierarchical geographic data with nested set model
  - Countries (247 records)
  - Provinces (34 Indonesian provinces)
  - Cities/Regencies (514 records with automatic KOTA/KAB classification)
  - Districts (7,000+ records)
  - Villages (83,000+ records)
- **🏦 Banks** - Indonesian banking institutions (142 records)
- **💰 Currencies** - World currencies with symbols (168 records)
- **🗣️ Languages** - ISO language codes and names (185 records)

#### Geographic Data Management
- **Nested Set Model** for efficient hierarchical queries
- **Automatic Type Classification** - KOTA/KAB prefix detection for cities/regencies
- **Coordinate Support** - Latitude/longitude for geographic entities
- **Postal Code Integration** for enhanced location data
- **Depth-based Hierarchies** with record_depth field

#### Search & Performance
- **🔍 Meilisearch Integration** - Fast, typo-tolerant search engine
- **Fallback Search** - Database search when Meilisearch unavailable
- **Search Index Management** - CLI tools for index initialization and maintenance
- **Optimized Queries** - Efficient database operations with proper indexing

#### Command Line Interface
- **🛠️ Cobra CLI Framework** - Modern command-line experience
- **Database Migrations** - Up/down migration support with golang-migrate
- **API Key Management** - CLI tools for creating and managing API keys
- **Data Seeding** - Comprehensive seeding with TRUNCATE optimization
- **Search Management** - Index initialization, reindexing, health checks

#### Data Management & Performance
- **🚀 TRUNCATE Operations** - Efficient bulk data clearing vs DELETE
- **📊 Progress Tracking** - Real-time progress during data operations
- **🔄 Upsert Logic** - Prevents duplicate records during seeding
- **📁 Selective Seeding** - Seed specific data types with --name flag
- **Custom Data Sources** - Support for loading from custom directories

#### Development & Testing
- **Unit Testing** - Comprehensive test coverage for entities, services, repositories
- **Integration Testing** - Database and HTTP handler testing
- **Test Utilities** - Mocking with testify and go-sqlmock
- **Development Tools** - Hot reload, linting, formatting support

#### Logging & Monitoring
- **Structured Logging** - JSON and text format support with logrus
- **Database Query Logging** - Optional SQL query logging for debugging
- **Context Propagation** - Request tracing through application layers
- **Configurable Log Levels** - Debug, info, warn, error levels

#### Deployment & Infrastructure
- **🐳 Docker Support** - Multi-stage builds with Alpine Linux
- **Docker Compose** - Complete stack with PostgreSQL and Meilisearch
- **Cross-Platform Builds** - Windows, Linux, macOS support
- **Environment Configuration** - Comprehensive .env support
- **Production Ready** - Optimized for production deployment

### 🏗️ Architecture

#### Clean Architecture Implementation
- **Domain Layer** - Entities, repositories, services (business logic)
- **Application Layer** - Use cases and application services
- **Infrastructure Layer** - Database, HTTP, external services
- **Dependency Injection** - Loose coupling with interface-based design

#### Repository Pattern
- **pgx Repositories** - High-performance PostgreSQL operations
- **Interface Abstraction** - Easy testing and implementation swapping
- **Transaction Support** - Atomic operations with proper error handling

#### HTTP Layer
- **Fiber Framework** - Fast HTTP server with middleware support
- **Structured Responses** - Consistent JSON API responses
- **Error Handling** - Comprehensive error handling with proper HTTP codes
- **Middleware Chain** - Authentication, logging, CORS support

### 📦 Installation & Setup

#### Multiple Installation Methods
- **From Source** - Go build with full development setup
- **Pre-built Binaries** - Cross-platform executable releases
- **Docker** - Container-based deployment with Docker Hub images
- **Docker Compose** - Complete stack orchestration

#### Database Setup
- **PostgreSQL 13+** - Modern PostgreSQL features and performance
- **Migration System** - Version-controlled database schema evolution
- **Seed Data** - Comprehensive sample data for immediate usage
- **Backup Scripts** - Database backup and restore utilities

#### Search Engine Setup
- **Meilisearch Integration** - Optional fast search capabilities
- **Index Management** - Automated index creation and maintenance
- **Search Configuration** - Customizable search attributes and filters

### 🔧 Configuration

#### Environment Variables
- **Database Configuration** - Connection string, SSL, pooling settings
- **Application Settings** - Host, port, environment configuration
- **Logging Configuration** - Level, format, output destinations
- **Search Configuration** - Meilisearch host and API key settings

#### Security Features
- **API Key Authentication** - Secure access control with Bearer tokens
- **Database Security** - Connection encryption and user management
- **Input Validation** - Comprehensive request validation
- **Error Sanitization** - Safe error messages without information leakage

### 📊 Performance Features

#### Database Optimizations
- **Indexed Queries** - Proper indexing for fast lookups
- **Connection Pooling** - Efficient database connection management
- **Prepared Statements** - SQL injection prevention and performance
- **Batch Operations** - Efficient bulk data operations

#### Search Optimizations
- **Fast Search** - Sub-millisecond search with Meilisearch
- **Typo Tolerance** - Intelligent search with typo correction
- **Search Filters** - Advanced filtering and faceted search
- **Search Analytics** - Built-in search statistics and monitoring

### 🧪 Testing

#### Comprehensive Test Suite
- **Unit Tests** - Individual component testing with high coverage
- **Integration Tests** - End-to-end API testing
- **Repository Tests** - Database operation testing with real connections
- **Service Tests** - Business logic testing with mocked dependencies

#### Testing Tools
- **testify** - Assertion and mocking framework
- **go-sqlmock** - Database interaction mocking
- **HTTP Testing** - API endpoint testing with test server
- **Test Coverage** - Coverage reporting and analysis

### 📚 Documentation

#### Comprehensive Documentation
- **API Documentation** - Interactive Swagger UI with live testing
- **Installation Guide** - Step-by-step setup for all platforms
- **CLI Usage Guide** - Complete command reference with examples
- **Architecture Guide** - Detailed technical architecture explanation
- **Testing Guide** - Testing strategies and best practices

#### Developer Resources
- **Code Examples** - API usage patterns and integration examples
- **Troubleshooting** - Common issues and solutions
- **Contributing Guide** - Development setup and contribution guidelines
- **Docker Documentation** - Container deployment and orchestration

### 🔗 API Endpoints

#### Core Endpoints
- `GET /health` - Health check and service status
- `GET /api/v1/geodirectories` - List geographic directories
- `GET /api/v1/geodirectories/{id}` - Get specific geodirectory
- `GET /api/v1/geodirectories/search` - Search geodirectories
- `POST /api/v1/geodirectories` - Create new geodirectory
- `PUT /api/v1/geodirectories/{id}` - Update geodirectory
- `DELETE /api/v1/geodirectories/{id}` - Delete geodirectory

#### Bank Endpoints
- `GET /api/v1/banks` - List banks with pagination
- `GET /api/v1/banks/{id}` - Get specific bank
- `GET /api/v1/banks/search` - Search banks
- `POST /api/v1/banks` - Create new bank
- `PUT /api/v1/banks/{id}` - Update bank
- `DELETE /api/v1/banks/{id}` - Delete bank

#### Currency Endpoints
- `GET /api/v1/currencies` - List currencies
- `GET /api/v1/currencies/{id}` - Get specific currency
- `GET /api/v1/currencies/search` - Search currencies
- `POST /api/v1/currencies` - Create new currency
- `PUT /api/v1/currencies/{id}` - Update currency
- `DELETE /api/v1/currencies/{id}` - Delete currency

#### Language Endpoints
- `GET /api/v1/languages` - List languages
- `GET /api/v1/languages/{id}` - Get specific language
- `GET /api/v1/languages/search` - Search languages
- `POST /api/v1/languages` - Create new language
- `PUT /api/v1/languages/{id}` - Update language
- `DELETE /api/v1/languages/{id}` - Delete language

### 🛠️ CLI Commands

#### Server Management
- `serve` - Start the HTTP server
- `version` - Display version information
- `help` - Show command help

#### Database Operations
- `migrate up` - Run pending migrations
- `migrate down` - Rollback migrations
- `migrate status` - Check migration status
- `migrate create` - Create new migration

#### Data Management
- `seed` - Seed database with sample data
- `seed --clear` - Clear existing data with TRUNCATE and reseed
- `seed --name <type>` - Seed specific data type
- `create-api-key` - Create new API key

#### Search Management
- `search init` - Initialize search indexes
- `search reindex` - Reindex all data
- `search health` - Check search service status
- `search stats` - Display search statistics

### 🐳 Docker Features

#### Container Features
- **Multi-stage Build** - Optimized image size with build and runtime stages
- **Alpine Linux** - Minimal security-focused base image
- **Non-root User** - Security best practices with dedicated application user
- **Health Checks** - Built-in container health monitoring

#### Docker Compose Stack
- **PostgreSQL** - Database service with persistent volumes
- **Meilisearch** - Search engine service
- **Application** - API service with proper networking
- **Environment Management** - Centralized configuration

### 📈 Metrics & Monitoring

#### Built-in Metrics
- **Health Endpoint** - Service health and dependency status
- **Search Statistics** - Index sizes and query performance
- **Database Metrics** - Connection pool and query statistics
- **Application Metrics** - Request counts and response times

#### Logging Features
- **Structured Logging** - Machine-readable JSON logs
- **Request Logging** - HTTP request/response logging
- **Database Logging** - Optional SQL query logging
- **Error Tracking** - Detailed error logging with stack traces

### 🔄 Migration Path

This is the initial release (v1.0.0), establishing the foundation for:
- Future API versioning
- Database schema evolution
- Feature expansion
- Performance optimizations

### 🤝 Contributing

#### Development Setup
- Go 1.21+ required
- PostgreSQL 13+ for database
- Meilisearch 1.5+ for search (optional)
- Docker for containerized development

#### Code Standards
- Hexagonal architecture principles
- Comprehensive unit testing
- Code formatting with gofmt
- Linting with golangci-lint

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast HTTP framework
- [pgx](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- [Meilisearch](https://www.meilisearch.com/) - Fast search engine
- [Cobra](https://cobra.dev/) - Modern CLI framework
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging

---

**Full Changelog**: https://github.com/turahe/master-data-rest-api/commits/v1.0.0
