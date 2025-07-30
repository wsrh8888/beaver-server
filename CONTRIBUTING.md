# 🤝 Contributing to Beaver IM

Thank you for your interest in contributing to Beaver IM! This document provides guidelines and information for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)
- [Documentation](#documentation)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

### Our Standards

- **Be respectful** - Treat everyone with respect
- **Be collaborative** - Work together to achieve common goals
- **Be constructive** - Provide helpful feedback and suggestions
- **Be professional** - Maintain professional behavior in all interactions

## 🎯 How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates.

**Bug Report Template:**

```markdown
## Bug Description
Brief description of the bug

## Steps to Reproduce
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior
What you expected to happen

## Actual Behavior
What actually happened

## Environment
- OS: [e.g. Windows 10, macOS 12.0]
- Go Version: [e.g. 1.21.0]
- Database: [e.g. MySQL 8.0]
- Redis: [e.g. 6.0]

## Additional Information
Any additional context, logs, or screenshots
```

### 💡 Suggesting Enhancements

**Feature Request Template:**

```markdown
## Problem Statement
Clear description of the problem this feature would solve

## Proposed Solution
Description of the proposed solution

## Alternative Solutions
Any alternative solutions you've considered

## Additional Context
Any other context, screenshots, or examples
```

### 🔧 Code Contributions

We welcome code contributions! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** your changes
5. **Submit** a pull request

## 🛠️ Development Setup

### Prerequisites

- Go >= 1.21
- MySQL >= 8.0
- Redis >= 6.0
- ETCD >= 3.5
- Docker >= 20.0

### Local Development

1. **Clone your fork**
```bash
git clone https://github.com/YOUR_USERNAME/beaver-server.git
cd beaver-server
```

2. **Add upstream remote**
```bash
git remote add upstream https://github.com/wsrh8888/beaver-server.git
```

3. **Install dependencies**
```bash
go mod tidy
```

4. **Start development environment**
```bash
# Start infrastructure services
docker-compose -f build/docker-compose.yaml up -d

# Initialize database
go run main.go -db
```

5. **Run tests**
```bash
go test ./...
```

## 📝 Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for code formatting
- Run `golint` to check code quality
- Follow naming conventions:
  - `camelCase` for variables and functions
  - `PascalCase` for exported names
  - `snake_case` for file names

### Project Structure

```
app/
├── service_name/
│   ├── service_api/          # HTTP API layer
│   │   ├── internal/
│   │   │   ├── handler/      # HTTP handlers
│   │   │   ├── logic/        # Business logic
│   │   │   ├── svc/          # Service context
│   │   │   └── types/        # Request/response types
│   │   └── etc/              # Configuration files
│   └── service_rpc/          # gRPC service layer
│       ├── internal/
│       │   ├── logic/        # Business logic
│       │   ├── server/       # gRPC server
│       │   └── svc/          # Service context
│       └── etc/              # Configuration files
```

### Error Handling

- Always check errors and handle them appropriately
- Use custom error types for business logic errors
- Provide meaningful error messages
- Log errors with appropriate levels

### Logging

```go
// Use structured logging
logx.Info("user registered", logx.Field("user_id", userID))
logx.Error("failed to send message", logx.Field("error", err))
```

### Testing

- Write unit tests for all business logic
- Aim for >80% code coverage
- Use table-driven tests for multiple scenarios
- Mock external dependencies

```go
func TestUserService_Register(t *testing.T) {
    tests := []struct {
        name    string
        input   RegisterRequest
        want    *RegisterResponse
        wantErr bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 📝 Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

### Examples

```bash
feat(auth): add email verification functionality
fix(chat): resolve message delivery issue
docs: update API documentation
refactor(user): simplify user registration logic
test(api): add integration tests for user endpoints
```

## 🔄 Pull Request Process

### Before Submitting

1. **Update your branch**
```bash
git fetch upstream
git rebase upstream/main
```

2. **Run tests**
```bash
go test ./...
go vet ./...
golint ./...
```

3. **Check formatting**
```bash
gofmt -s -w .
```

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

### Review Process

1. **Automated Checks**
   - CI/CD pipeline runs tests
   - Code coverage is checked
   - Linting and formatting are verified

2. **Code Review**
   - At least one maintainer must approve
   - Address all review comments
   - Update documentation if needed

3. **Merge**
   - Squash commits if requested
   - Use conventional commit message
   - Delete feature branch after merge

## 📚 Documentation

### Code Documentation

- Document all exported functions and types
- Use clear and concise comments
- Include examples for complex functions

```go
// UserService handles user-related operations
type UserService struct {
    // ...
}

// Register creates a new user account
// Returns the created user ID and any validation errors
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
    // Implementation
}
```

### API Documentation

- Update API documentation for new endpoints
- Include request/response examples
- Document error codes and messages

### README Updates

- Update README.md for new features
- Add installation instructions for new dependencies
- Update configuration examples

## 🏷️ Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH`
- `MAJOR`: Breaking changes
- `MINOR`: New features (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Changelog is updated
- [ ] Version is bumped
- [ ] Release notes are written
- [ ] Docker images are built and pushed

## 🆘 Getting Help

- **Issues**: [GitHub Issues](https://github.com/wsrh8888/beaver-server/issues)
- **Discussions**: [GitHub Discussions](https://github.com/wsrh8888/beaver-server/discussions)
- **Email**: [751135385@qq.com](mailto:751135385@qq.com)
- **QQ Group**: [1013328597](https://qm.qq.com/q/82rbf7QBzO)

## 🙏 Recognition

Contributors will be recognized in:

- Project README.md
- Release notes
- Contributor hall of fame
- GitHub contributors page

---

Thank you for contributing to Beaver IM! 🦫 