# Makefile for redup
# Based on scopy project structure

.PHONY: help build test clean install uninstall run release snapshot

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the project"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build files"
	@echo "  install   - Install locally"
	@echo "  uninstall - Uninstall"
	@echo "  run       - Run with arguments (use ARGS='arg1 arg2')"
	@echo "  release   - Create a new release"
	@echo "  snapshot  - Create a snapshot release"
	@echo "  help      - Show this help"

# Build the project
build:
	go build -o redup

# Run tests
test:
	go test ./pkg/...

# Clean build files
clean:
	rm -f redup
	rm -rf dist/

# Install locally
install:
	go install

# Uninstall
uninstall:
	go clean -i

# Run the application with arguments
run:
	@if [ -z "$(ARGS)" ]; then \
		echo "Usage: make run ARGS='arg1 arg2'"; \
		echo "Example: make run ARGS='--dry-run .'"; \
		exit 1; \
	fi
	./redup $(ARGS)

# Create a new release
release:
	@echo "Creating a new release..."
	@bin/release.sh

# Create a snapshot release
snapshot:
	@echo "Creating a snapshot release..."
	@bin/release.sh --snapshot

# Update dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Check for security vulnerabilities
security:
	gosec ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@echo "Documentation is in README.md and docs/README.md"

# Development setup
dev-setup: deps
	@echo "Development environment setup complete"
	@echo "Run 'make build' to build the project"
	@echo "Run 'make test' to run tests"
	@echo "Run 'make run ARGS=\"--help\"' to see usage"
