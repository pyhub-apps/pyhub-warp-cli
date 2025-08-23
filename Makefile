# Get version from headver
VERSION := $(shell ./scripts/headver.sh)
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags
LDFLAGS := -ldflags "-s -w \
	-X main.version=$(VERSION) \
	-X main.gitCommit=$(GIT_COMMIT) \
	-X main.buildDate=$(BUILD_DATE)"

.PHONY: all build test clean install release-snapshot release version

all: test build

# Build the binary
build:
	@echo "Building sejong $(VERSION)..."
	go build $(LDFLAGS) -o sejong ./cmd/sejong

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Install the binary to $GOPATH/bin
install: build
	@echo "Installing sejong to GOPATH/bin..."
	cp sejong $(shell go env GOPATH)/bin/sejong

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f sejong
	rm -f coverage.txt coverage.html
	rm -rf dist/

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...

# Show current version
version:
	@echo "Current version: $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"

# Build for all platforms (local testing)
release-snapshot:
	@echo "Building release snapshot..."
	goreleaser release --snapshot --clean

# Create a release (requires git tag)
release:
	@echo "Creating release..."
	goreleaser release --clean

# Development build with race detector
dev:
	@echo "Building development version with race detector..."
	go build -race $(LDFLAGS) -o sejong ./cmd/sejong

# Quick test of the binary
run: build
	./sejong version

# Bump head version
bump-head:
	./scripts/headver.sh --bump-head

# Set specific head version
set-head:
	@read -p "Enter new head version: " version; \
	./scripts/headver.sh --set-head $$version