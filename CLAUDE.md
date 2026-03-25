# CLAUDE.md - AI Assistant Guide for pulsar-client-tool

## Project Overview

`pulsar-client-tool` is a CLI tool for interacting with Apache Pulsar messaging clusters. It provides command-line utilities for producing, consuming, and managing Pulsar topics and messages.

**Language:** Go
**Repository:** `mkontani/pulsar-client-tool` (private)
**Default Branch:** `main`

## Project Status

This project is in early development. The repository is being initialized and the codebase is being built out.

## Expected Project Structure

```
pulsar-client-tool/
├── CLAUDE.md           # This file - AI assistant guide
├── README.md           # Project documentation
├── go.mod              # Go module definition
├── go.sum              # Go dependency checksums
├── main.go             # Application entrypoint
├── cmd/                # CLI command definitions (cobra/urfave)
├── internal/           # Internal packages (not exported)
│   ├── client/         # Pulsar client wrapper
│   ├── config/         # Configuration handling
│   └── util/           # Shared utilities
├── pkg/                # Public/reusable packages
├── Makefile            # Build and development tasks
├── Dockerfile          # Container build
├── .github/            # GitHub Actions CI/CD
│   └── workflows/
└── testdata/           # Test fixtures
```

## Development Workflow

### Prerequisites

- Go 1.21+ (use the version specified in `go.mod`)
- Apache Pulsar instance for integration testing (or use the standalone Docker image)

### Common Commands

```bash
# Build
go build -o pulsar-client-tool .

# Run tests
go test ./...

# Run tests with race detector
go test -race ./...

# Lint (if golangci-lint is configured)
golangci-lint run

# Format code
gofmt -w .
go vet ./...

# Tidy dependencies
go mod tidy
```

### Build with Makefile (when available)

```bash
make build        # Build binary
make test         # Run tests
make lint         # Run linter
make clean        # Clean build artifacts
```

## Code Conventions

### Go Standards

- Follow standard Go project layout conventions
- Use `gofmt` for formatting (no exceptions)
- Run `go vet` to catch common issues
- Exported functions and types must have doc comments
- Error messages should be lowercase, without trailing punctuation
- Prefer returning errors over panicking
- Use `context.Context` as the first parameter for functions that do I/O

### Naming

- Package names: short, lowercase, no underscores (e.g., `client`, `config`)
- Interfaces: use `-er` suffix where appropriate (e.g., `Producer`, `Consumer`)
- Test files: `*_test.go` alongside the code they test
- Constants: `CamelCase` for exported, `camelCase` for unexported

### Error Handling

- Wrap errors with context using `fmt.Errorf("operation: %w", err)`
- Check errors immediately after the call that produces them
- Do not ignore errors silently — handle or explicitly document why it's safe

### Testing

- Table-driven tests are preferred for multiple test cases
- Use `testify` or standard library assertions
- Name test functions descriptively: `TestProducer_SendMessage_WithTimeout`
- Place test helpers in `_test.go` files or a `testutil` package

## Pulsar-Specific Conventions

- Use the official Apache Pulsar Go client: `github.com/apache/pulsar-client-go`
- Connection configuration (service URL, auth) should be configurable via flags, env vars, and config file
- Support common Pulsar operations: produce, consume, list topics, manage subscriptions
- Handle Pulsar client lifecycle properly — always close clients and producers/consumers

## Configuration Priority

Configuration should follow this precedence (highest to lowest):
1. CLI flags
2. Environment variables (prefixed, e.g., `PULSAR_SERVICE_URL`)
3. Configuration file (e.g., `~/.pulsar-client-tool.yaml`)
4. Sensible defaults

## AI Assistant Guidelines

### When Making Changes

- Read existing code before modifying — understand the patterns in use
- Run `go build ./...` after changes to verify compilation
- Run `go test ./...` after changes to verify tests pass
- Run `go mod tidy` if dependencies were added or removed
- Keep changes focused — one concern per commit

### Commit Messages

- Use conventional style: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- Keep the subject line under 72 characters
- Use imperative mood: "Add feature" not "Added feature"

### What to Avoid

- Do not commit secrets, credentials, or connection strings
- Do not add unnecessary dependencies — prefer the standard library
- Do not modify `go.sum` manually — use `go mod tidy`
- Do not ignore linter warnings without justification
