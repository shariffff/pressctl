.PHONY: build install clean test

# Binary name
BINARY_NAME=press

# Build information
VERSION?=dev
COMMIT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
LDFLAGS=-ldflags "-X github.com/pressctl/cli/cmd.Version=$(VERSION) -X github.com/pressctl/cli/cmd.CommitSHA=$(COMMIT_SHA) -X github.com/pressctl/cli/cmd.BuildDate=$(BUILD_DATE)"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete. Run 'press --help' to get started."

# Install to current user's bin directory (no sudo required)
install-user: build
	@echo "Installing $(BINARY_NAME) to ~/bin..."
	@mkdir -p ~/bin
	@cp $(BINARY_NAME) ~/bin/
	@echo "Installation complete."
	@echo "Make sure ~/bin is in your PATH."
	@echo "Add this to your ~/.bashrc or ~/.zshrc if not already present:"
	@echo '  export PATH="$$HOME/bin:$$PATH"'

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

# Run the binary
run: build
	@./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed. Run: brew install golangci-lint"

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary (requires Go 1.24+)"
	@echo "  install            - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user       - Install to ~/bin (no sudo)"
	@echo "  test               - Run tests"
	@echo "  clean              - Remove build artifacts"
	@echo "  run                - Build and run the binary"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
