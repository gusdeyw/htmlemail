# Makefile for htmlemail Go package

.PHONY: all build test test-verbose test-race test-cover lint fmt vet clean deps deps-update example help

# Default target
all: fmt vet lint test build

# Build the package
build:
	@echo "Building package..."
	go build -v ./...

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean ./...
	rm -f coverage.out coverage.html

# Install/update dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Run example
example:
	@echo "Running example..."
	cd example && go run main.go

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Run security scanner
security:
	@echo "Running security scanner..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: make install-tools"; \
	fi

# Check if go.mod is tidy
mod-tidy-check:
	@echo "Checking if go.mod is tidy..."
	go mod tidy
	git diff --exit-code go.mod go.sum

# CI pipeline simulation
ci: fmt vet lint test-race security
	@echo "CI pipeline completed successfully!"

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run fmt, vet, lint, test, and build"
	@echo "  build        - Build the package"
	@echo "  test         - Run tests"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  test-race    - Run tests with race detection"
	@echo "  test-cover   - Run tests with coverage report"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  example      - Run example program"
	@echo "  install-tools- Install development tools"
	@echo "  security     - Run security scanner (requires gosec)"
	@echo "  mod-tidy-check - Check if go.mod is tidy"
	@echo "  ci           - Run CI pipeline simulation"
	@echo "  help         - Show this help message"