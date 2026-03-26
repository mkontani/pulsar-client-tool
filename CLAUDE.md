# CLAUDE.md - AI Assistant Guide for pulsar-client-tool

## Project Overview

`pulsar-client-tool` is a Go-based CLI tool for Apache Pulsar, designed as a simpler alternative to the Java-based `pulsar-client`. It provides `produce` and `consume` commands with intuitive flags, environment variable support, and text/JSON output.

**Language:** Go 1.24+
**Module:** `github.com/mkontani/pulsar-client-tool`
**Repository:** `mkontani/pulsar-client-tool` (private)
**Default Branch:** `main`

## Project Structure

```
pulsar-client-tool/
├── CLAUDE.md                           # This file
├── README.md                           # User documentation
├── Makefile                            # Build targets
├── go.mod / go.sum                     # Go module deps
├── main.go                             # Entrypoint → cmd.Execute()
├── cmd/
│   ├── root.go                         # Root command, global flags, signal handling
│   ├── produce.go                      # produce subcommand
│   ├── consume.go                      # consume subcommand
│   └── version.go                      # version subcommand
└── internal/
    ├── config/
    │   ├── config.go                   # Config struct, Validate(), ClientOptions()
    │   └── config_test.go
    ├── client/
    │   └── client.go                   # PulsarClient interface + factory
    ├── producer/
    │   └── producer.go                 # Send logic (string, file, stdin)
    ├── consumer/
    │   ├── consumer.go                 # Receive loop with graceful shutdown
    │   └── consumer_test.go
    └── output/
        ├── output.go                   # Text/JSON message formatter
        └── output_test.go
```

## Architecture

- **cobra** CLI framework with subcommands (`produce`, `consume`, `version`)
- **Thin cmd layer** — cobra files parse flags and delegate to `internal/` packages
- **`internal/client.PulsarClient`** interface wraps `pulsar.Client` for testability
- **`io.Reader`/`io.Writer`** used throughout for testable I/O
- **`context.Context`** propagated from root command (with `signal.NotifyContext`) for Ctrl+C handling
- **No global state** — config built in `PersistentPreRunE` and passed via context

### Key Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/apache/pulsar-client-go` — Pulsar client library

## Development Workflow

### Common Commands

```bash
make build        # Build binary → ./pulsar-client-tool
make test         # go test ./...
make vet          # go vet ./...
make fmt          # gofmt -w .
make lint         # golangci-lint run
make clean        # Remove binary
go mod tidy       # After changing dependencies
```

### Quick Verification

```bash
go build ./... && go vet ./... && go test ./...
```

## CLI Commands

### produce

Sends messages to a Pulsar topic. Message source priority:
1. `--message "text"` — literal string
2. `--file path` — one message per line
3. stdin (if piped)

Key flags: `-t/--topic` (required), `-m/--message`, `-f/--file`, `-k/--key`, `-p/--property` (repeatable), `-n/--num-messages`, `--rate`

### consume

Subscribes and receives messages. Runs until Ctrl+C or `--num-messages` reached.

Key flags: `-t/--topic` (required), `-s/--subscription` (required), `-S/--subscription-type`, `-n/--num-messages`

### Global flags

- `--service-url` / `PULSAR_SERVICE_URL` (default: `pulsar://localhost:6650`)
- `--auth-token` / `PULSAR_AUTH_TOKEN`
- `--tls-cert`, `--tls-key`, `--tls-ca`
- `--timeout` (default: `30s`)
- `-o/--output` (`text` or `json`)

## Code Conventions

### Go Standards

- `gofmt` for formatting, `go vet` for correctness
- Error messages: lowercase, no trailing punctuation
- Wrap errors with context: `fmt.Errorf("operation: %w", err)`
- `context.Context` as first parameter for I/O functions
- Prefer returning errors over panicking

### Naming

- Package names: short, lowercase, no underscores
- `PulsarClient` interface in `internal/client/` for testability
- Test files: `*_test.go` alongside source

### Testing

- Table-driven tests preferred
- Standard library `testing` package (no external assertion library required)
- Unit tests run without a live Pulsar instance
- Integration tests (future): use `//go:build integration` tag

## Configuration Priority

1. CLI flags (highest)
2. Environment variables (`PULSAR_SERVICE_URL`, `PULSAR_AUTH_TOKEN`)
3. Sensible defaults (lowest)

## AI Assistant Guidelines

### When Making Changes

- Read existing code before modifying — follow established patterns
- Run `go build ./...` after changes to verify compilation
- Run `go test ./...` after changes to verify tests pass
- Run `go mod tidy` if dependencies changed
- Keep changes focused — one concern per commit

### Commit Messages

- Conventional style: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- Under 72 characters, imperative mood

### What to Avoid

- Do not commit secrets, credentials, or connection strings
- Do not add unnecessary dependencies — prefer the standard library
- Do not modify `go.sum` manually — use `go mod tidy`
- Do not ignore linter warnings without justification
