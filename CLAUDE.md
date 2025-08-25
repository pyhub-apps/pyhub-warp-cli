# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sejong-CLI is a command-line utility tool for searching South Korean law information using the National Law Information Center Open API. The project aims to provide developers, researchers, and legal professionals with quick access to legal information directly from their terminal.

## Development Setup

### Initialize Project (First Time)
```bash
# Initialize Go module
go mod init github.com/allieus/sejong-cli

# Add dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/olekukonko/tablewriter@latest
```

### Build and Run Commands
```bash
# Build the project
go build -o sejong cmd/sejong/main.go

# Run without building
go run cmd/sejong/main.go

# Build for multiple platforms
GOOS=windows GOARCH=amd64 go build -o sejong.exe cmd/sejong/main.go
GOOS=darwin GOARCH=arm64 go build -o sejong cmd/sejong/main.go
GOOS=linux GOARCH=amd64 go build -o sejong-linux cmd/sejong/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestName ./...

# Format code
go fmt ./...

# Vet code for suspicious constructs
go vet ./...

# Generate release binaries
goreleaser release --snapshot --clean  # local testing
goreleaser release  # actual release (requires git tag)
```

## Architecture Guidelines

### Project Structure
The application should follow a clean architecture pattern with the following structure:
- `cmd/sejong/` - Entry point and command definitions
- `internal/api/` - API client for National Law Information Center
- `internal/config/` - Configuration management (API keys, settings)
- `internal/output/` - Output formatting (table, JSON)

### Key Design Decisions

1. **Command Structure**: Use Cobra CLI framework for command handling
   - Main command: `sejong`
   - Subcommands: `law`, `config`
   - Support for flags like `--format json`

2. **Configuration Storage**: Store API keys and settings in user's home directory
   - Location: `~/.sejong/config.yaml` or similar
   - Use viper for configuration management

3. **Error Handling**: Provide user-friendly error messages
   - When API key is missing, show clear setup instructions
   - Include the API registration URL and exact setup commands
   - Gracefully handle network errors and API failures

4. **Output Formats**: Support both human-readable table and machine-readable JSON
   - Default: Formatted table output for terminal viewing
   - JSON: For scripting and automation (`--format json`)

### API Integration

The National Law Information Center API requires:
- Authentication key (stored in config)
- Base URL: `https://www.law.go.kr/DRF/lawSearch.do` (primary endpoint)
- Response parsing: Handle XML or JSON responses based on API version
- Request parameters:
  - `OC`: API key (required)
  - `target`: "law" for law search
  - `query`: search keyword
  - `type`: response format (XML/JSON)

### Development Principles

1. **Single Binary Distribution**: Compile to a single executable without runtime dependencies
2. **Cross-Platform Support**: Ensure compatibility with Windows (amd64/arm64) and macOS (arm64)
3. **Offline-First Config**: Configuration should work without network access once API key is set
4. **Progressive Enhancement**: Start with core search functionality, expand features incrementally

## Critical Implementation Notes

### First-Time User Experience
When a user runs `sejong law` without an API key configured:
1. Check for API key in configuration
2. If missing, display a helpful message that includes:
   - Why the key is needed
   - Direct link to get the key: `https://open.law.go.kr/LSO/openApi/cuAskList.do`
   - Exact command to set the key: `sejong config set law.key <YOUR_KEY>`

### Command Examples
```bash
# Search for laws
sejong law "개인정보 보호법"

# Set API key
sejong config set law.key YOUR_API_KEY_HERE

# Get current API key
sejong config get law.key

# Export results as JSON
sejong law "도로교통법" --format json > laws.json
```

## Testing Strategy

1. **Unit Tests**: Test individual components (API client, config manager, formatters)
2. **Integration Tests**: Test command execution with mocked API responses
3. **Manual Testing**: Verify actual API integration with a test key

## Dependencies to Consider

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/olekukonko/tablewriter` - Table formatting
- Standard library for JSON encoding and HTTP requests

## Implementation Status

**Current Status**: Planning phase - no Go code implementation yet

The project structure needs to be created:
```
pyhub-sejong-cli/
├── cmd/
│   └── sejong/
│       └── main.go         # CLI entry point
├── internal/
│   ├── api/
│   │   └── client.go       # API client for law.go.kr
│   ├── config/
│   │   └── config.go       # Configuration management
│   └── output/
│       └── formatter.go    # Output formatting (table/JSON)
├── go.mod
├── go.sum
└── .goreleaser.yml         # Release configuration
```