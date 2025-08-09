# 🤝 Contributing to Master Data REST API

Thank you for your interest in contributing to the Master Data REST API! We welcome contributions from the community and are pleased to have them.

## 🎯 Quick Start for Contributors

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/master-data-rest-api.git`
3. **Create** a branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes: `make test`
6. **Commit** your changes: `git commit -m 'Add amazing feature'`
7. **Push** to your branch: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [Contribution Guidelines](#contribution-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Documentation](#documentation)
- [Testing](#testing)
- [Code Style](#code-style)

## 📜 Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- **Be respectful** and inclusive
- **Be collaborative** and help others learn
- **Be constructive** in discussions and feedback
- **Be patient** with beginners and questions
- **Focus on what is best** for the community and project

## 💻 Development Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 13 or higher
- Docker (optional, for testing)
- Git

### Setup Steps

1. **Fork and clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/master-data-rest-api.git
   cd master-data-rest-api
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment**
   ```bash
   cp env.example .env
   # Edit .env with your local configuration
   ```

4. **Run database migrations**
   ```bash
   go run main.go migrate up
   ```

5. **Start development server**
   ```bash
   go run main.go serve --log-level debug
   ```

### Development Tools

Install helpful development tools:

```bash
# Install development tools
make install-tools

# Or manually:
go install github.com/cosmtrek/air@latest              # Hot reloading
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
go install github.com/swaggo/swag/cmd/swag@latest     # Swagger docs
```

## 🎯 Contribution Guidelines

### Types of Contributions

We welcome several types of contributions:

#### 🐛 Bug Fixes
- Fix issues in existing functionality
- Improve error handling
- Resolve performance problems

#### ✨ New Features
- Add new API endpoints
- Implement new CLI commands
- Extend existing functionality

#### 📚 Documentation
- Improve existing documentation
- Add new guides or examples
- Fix typos or unclear explanations

#### 🧪 Testing
- Add unit tests
- Improve test coverage
- Add integration tests

#### 🔧 Infrastructure
- Improve build processes
- Enhance Docker setup
- Optimize database queries

### Before You Start

1. **Check existing issues** - Look for existing issues or discussions
2. **Create an issue** - For new features or major changes, create an issue first
3. **Discuss the approach** - Get feedback on your proposed solution
4. **Start small** - Begin with smaller contributions to get familiar with the codebase

## 🔄 Pull Request Process

### Before Submitting

1. **Update documentation** - Ensure all relevant docs are updated
2. **Add tests** - Include tests for new functionality
3. **Run tests** - Ensure all tests pass: `make test`
4. **Lint code** - Run linter: `make lint`
5. **Update CHANGELOG** - Add entry for your changes (if applicable)

### PR Title and Description

Use clear, descriptive titles and provide detailed descriptions:

```
feat: add currency activation/deactivation endpoints

- Add POST /currencies/{id}/activate endpoint
- Add POST /currencies/{id}/deactivate endpoint  
- Update currency service with status management
- Include tests for new functionality
- Update API documentation

Fixes #123
```

### PR Template

When creating a PR, include:

- **What** - What changes does this PR make?
- **Why** - Why are these changes needed?
- **How** - How were the changes implemented?
- **Testing** - How was this tested?
- **Screenshots** - If applicable, include screenshots
- **Breaking Changes** - List any breaking changes

## 🐛 Issue Reporting

### Bug Reports

When reporting bugs, include:

1. **Description** - Clear description of the issue
2. **Steps to Reproduce** - Detailed steps to reproduce the bug
3. **Expected Behavior** - What you expected to happen
4. **Actual Behavior** - What actually happened
5. **Environment** - OS, Go version, PostgreSQL version
6. **Logs** - Relevant log output or error messages

### Feature Requests

When requesting features, include:

1. **Problem Statement** - What problem does this solve?
2. **Proposed Solution** - How should this be implemented?
3. **Use Cases** - When would this be used?
4. **Alternatives** - What alternatives did you consider?

## 📚 Documentation

### Documentation Standards

- Use clear, concise language
- Include code examples where applicable
- Update both inline code comments and markdown docs
- Follow existing documentation style and structure

### Documentation Types

1. **Code Comments** - Document public APIs and complex logic
2. **API Documentation** - Update Swagger annotations
3. **CLI Documentation** - Update help text and examples
4. **User Guides** - Create or update user-facing documentation

### Updating Documentation

When making changes that affect:

- **API endpoints** - Update Swagger comments and regenerate docs
- **CLI commands** - Update help text and CLI usage guide
- **Configuration** - Update environment variable documentation
- **Installation** - Update installation and setup guides

## 🧪 Testing

### Test Requirements

- **Unit Tests** - All new functionality must have unit tests
- **Integration Tests** - Add integration tests for complex features
- **Test Coverage** - Maintain or improve test coverage (target: 80%+)

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/domain/services/...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

Follow existing test patterns:

```go
func TestServiceMethod(t *testing.T) {
    // Given
    service := setupTestService()
    
    // When
    result, err := service.Method(testData)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}
```

## 🎨 Code Style

### Go Style Guidelines

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Keep functions small and focused
- Document public APIs

### Naming Conventions

- **Files**: `snake_case.go`
- **Packages**: `lowercase`, short and meaningful
- **Types**: `PascalCase`
- **Functions/Methods**: `PascalCase` (public), `camelCase` (private)
- **Variables**: `camelCase`
- **Constants**: `PascalCase` or `UPPER_CASE`

### Code Organization

Follow the hexagonal architecture pattern:

```
internal/
├── domain/           # Business logic
│   ├── entities/     # Domain entities
│   ├── repositories/ # Repository interfaces
│   └── services/     # Business services
└── adapters/         # Implementation
    ├── primary/      # Incoming adapters (HTTP, CLI)
    └── secondary/    # Outgoing adapters (Database, APIs)
```

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

Examples:
```
feat(api): add currency activation endpoints
fix(cli): resolve migration rollback issue
docs(readme): update installation instructions
```

## 🔍 Code Review Process

### For Contributors

- **Be responsive** to feedback
- **Address all comments** before requesting re-review
- **Test thoroughly** after making changes
- **Keep PRs focused** - one feature or fix per PR

### Review Criteria

PRs will be reviewed for:

1. **Functionality** - Does it work as intended?
2. **Code Quality** - Is it well-written and maintainable?
3. **Tests** - Are there adequate tests?
4. **Documentation** - Is it properly documented?
5. **Performance** - Does it impact performance?
6. **Security** - Are there any security concerns?

## 🏷️ Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for bug fixes (backwards compatible)

### Release Notes

Each release includes:

- **New Features** - What's new in this version
- **Bug Fixes** - What issues were resolved
- **Breaking Changes** - Any backwards incompatible changes
- **Upgrade Guide** - How to upgrade from previous version

## 🙋 Getting Help

### Community Support

- **GitHub Issues** - For bug reports and feature requests
- **Discussions** - For questions and general discussion
- **Discord/Slack** - For real-time chat (if available)

### Maintainer Contact

For sensitive issues or questions:
- Email: [maintainer@example.com]
- Direct message on GitHub

## 🎉 Recognition

### Contributors

All contributors are recognized in:

- **README.md** - Contributors section
- **CHANGELOG.md** - Release notes
- **GitHub** - Automatic contribution tracking

### Types of Recognition

- **Code Contributions** - Commits, PRs
- **Documentation** - Docs improvements
- **Issue Triage** - Helping with issues
- **Community Support** - Helping other users

## 📋 Checklist for Contributors

Before submitting a PR:

- [ ] Code follows project style guidelines
- [ ] Self-review of code completed
- [ ] Tests added for new functionality
- [ ] All tests pass locally
- [ ] Documentation updated (if applicable)
- [ ] Commit messages follow conventional format
- [ ] PR description is clear and complete
- [ ] No breaking changes (or clearly documented)

---

<div align="center">
  <strong>🚀 Thank you for contributing to Master Data REST API!</strong>
  <br>
  <sub>Every contribution, no matter how small, makes a difference</sub>
</div>
