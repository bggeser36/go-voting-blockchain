# Contributing to Go Voting Blockchain

Thank you for your interest in contributing to the Go Voting Blockchain project! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our commitment to providing a welcoming and inclusive environment for all contributors, regardless of background, experience level, gender identity, race, ethnicity, religion, or other characteristics.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your feature or bugfix
4. Make your changes
5. Test your changes thoroughly
6. Submit a pull request

## How to Contribute

### Reporting Bugs

- Use the GitHub issue tracker
- Include detailed reproduction steps
- Provide system information (OS, Go version, etc.)
- Include relevant logs or error messages

### Suggesting Enhancements

- Use the GitHub issue tracker with the "enhancement" label
- Clearly describe the proposed feature
- Explain why it would be valuable
- Consider implementation complexity

### Code Contributions

- Follow the existing code style
- Write comprehensive tests
- Update documentation as needed
- Ensure all tests pass

## Development Setup

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (for local development)
- Redis (optional, for caching)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-voting-blockchain.git
   cd go-voting-blockchain
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the development environment:
   ```bash
   docker-compose up -d
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

5. Run tests:
   ```bash
   go test ./...
   ```

## Pull Request Process

1. **Fork and Branch**: Create a feature branch from `main`
2. **Commit Messages**: Use clear, descriptive commit messages
3. **Test Coverage**: Ensure adequate test coverage for new code
4. **Documentation**: Update README.md and code comments as needed
5. **Pull Request**: Create a PR with a clear description of changes

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write comprehensive comments for public APIs
- Keep functions small and focused
- Handle errors explicitly

### Example:
```go
// CreatePoll creates a new voting poll in the blockchain
func (bc *Blockchain) CreatePoll(poll *models.Poll) error {
    bc.mu.Lock()
    defer bc.mu.Unlock()
    
    // Validation logic here
    if poll.Title == "" {
        return fmt.Errorf("poll title cannot be empty")
    }
    
    // Implementation here
    return nil
}
```

## Testing

### Unit Tests
- Write unit tests for all new functions
- Aim for >80% code coverage
- Use table-driven tests where appropriate

### Integration Tests
- Test API endpoints with real database
- Test blockchain operations end-to-end
- Use Docker for consistent test environments

### Example Test:
```go
func TestCreatePoll(t *testing.T) {
    tests := []struct {
        name    string
        poll    *models.Poll
        wantErr bool
    }{
        {
            name: "valid poll",
            poll: &models.Poll{
                Title: "Test Poll",
                Options: []string{"A", "B"},
            },
            wantErr: false,
        },
        {
            name: "empty title",
            poll: &models.Poll{
                Title: "",
                Options: []string{"A", "B"},
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            bc := blockchain.NewBlockchain(3)
            err := bc.CreatePoll(tt.poll)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreatePoll() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

### Code Documentation
- Document all public APIs
- Use Go doc conventions
- Include examples for complex functions

### README Updates
- Keep README.md current
- Include setup instructions
- Document API endpoints
- Provide usage examples

## Security Considerations

- Never commit secrets or API keys
- Use environment variables for configuration
- Validate all input data
- Follow secure coding practices
- Report security vulnerabilities privately

## Getting Help

- Check existing issues and discussions
- Join our community discussions
- Contact maintainers for guidance

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- GitHub contributor graphs

Thank you for contributing to Go Voting Blockchain! 🚀
