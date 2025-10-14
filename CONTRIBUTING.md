# Contributing to htmlemail

Thank you for your interest in contributing to htmlemail! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)
- [Documentation](#documentation)

## Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors. By participating, you agree to:

- Be respectful and inclusive
- Focus on constructive feedback
- Accept responsibility for mistakes
- Show empathy towards other contributors
- Help create a positive community

## Getting Started

### Prerequisites

- Go 1.19 or later
- Git
- Make (optional, for using the Makefile)

### Quick Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/go-htmlemail.git
   cd go-htmlemail
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Installing Dependencies

```bash
go mod download
```

### Building the Project

```bash
# Build the main package
go build ./...

# Or use the Makefile if available
make build
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Or use the Makefile
make test
```

### Running Examples

```bash
# Run the example application
cd example
go run main.go
```

## Project Structure

```
go-htmlemail/
├── htmlemail.go          # Main package code
├── htmlemail_test.go      # Unit tests
├── example/               # Example usage
│   └── main.go
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies
├── Makefile              # Build automation
├── README.md             # Project documentation
├── LICENSE               # License information
└── CONTRIBUTING.md       # This file
```

## Coding Standards

### Go Code Style

This project follows standard Go conventions:

- Use `gofmt` for code formatting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Write clear, concise comments
- Keep functions small and focused

### Code Formatting

```bash
# Format all Go files
gofmt -w .

# Check for formatting issues
gofmt -d .
```

### Linting

```bash
# Install golangci-lint if not already installed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linting
golangci-lint run
```

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/modifications
- `chore`: Maintenance tasks

Examples:
```
feat: add support for custom template functions
fix: resolve memory leak in template caching
docs: update API reference for new methods
```

## Testing

### Writing Tests

- Write unit tests for all public functions
- Use table-driven tests for multiple test cases
- Test edge cases and error conditions
- Aim for high test coverage (>80%)

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "TEST",
            wantErr:  false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

### Running Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run benchmarks with memory allocation info
go test -bench=. -benchmem ./...
```

## Submitting Changes

### Pull Request Process

1. Ensure your code follows the coding standards
2. Write or update tests for your changes
3. Update documentation if needed
4. Ensure all tests pass
5. Create a pull request with a clear description

### Pull Request Template

When creating a PR, include:

- **Description**: What changes were made and why
- **Type of change**: Bug fix, feature, documentation, etc.
- **Testing**: How the changes were tested
- **Breaking changes**: Any breaking changes and migration guide
- **Screenshots**: If UI changes are involved

### Review Process

- All PRs require review before merging
- Address review comments promptly
- Keep PRs focused on a single feature or fix
- Squash commits when merging

## Reporting Issues

### Bug Reports

When reporting bugs, include:

- **Go version**: `go version`
- **OS**: Your operating system and version
- **Steps to reproduce**: Clear steps to reproduce the issue
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Code sample**: Minimal code to reproduce the issue

### Feature Requests

For feature requests, include:

- **Use case**: Why do you need this feature?
- **Proposed solution**: How should it work?
- **Alternatives**: Other solutions you've considered

## Documentation

### Updating README

- Keep examples up to date
- Document new features clearly
- Update API reference when adding new functions

### Code Documentation

- Add comments to exported functions and types
- Use `go doc` compatible comment format
- Include examples in function comments when helpful

### Example

```go
// CalculateTotal calculates the total price including tax
// Example:
//
//	total := CalculateTotal(100.0, 0.08)
//	fmt.Println(total) // Output: 108.0
func CalculateTotal(price, taxRate float64) float64 {
    return price * (1 + taxRate)
}
```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing](https://golang.org/pkg/testing/)

Thank you for contributing to htmlemail! 🎉