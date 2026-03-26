# pulsar-client-tool

A user-friendly CLI for Apache Pulsar, written in Go.

Provides a simpler alternative to the Java-based `pulsar-client` tool with intuitive commands for producing and consuming messages, delayed/scheduled delivery, and connectivity checks.

## Installation

```bash
go install github.com/mkontani/pulsar-client-tool@latest
```

Or build from source:

```bash
git clone https://github.com/mkontani/pulsar-client-tool.git
cd pulsar-client-tool
make build
```

## Usage

### Produce messages

```bash
# Send a single message
pulsar-client-tool produce -t my-topic -m "hello world"

# Send from stdin
echo "hello" | pulsar-client-tool produce -t my-topic

# Send from file (one message per line)
pulsar-client-tool produce -t my-topic -f messages.txt

# Send with key and properties
pulsar-client-tool produce -t my-topic -m "hello" -k my-key -p env=prod -p version=1

# Send 100 messages at 10 msg/sec
pulsar-client-tool produce -t my-topic -m "load test" -n 100 --rate 10

# Delayed delivery (deliver after 30 seconds)
pulsar-client-tool produce -t my-topic -m "delayed" --deliver-after 30s

# Scheduled delivery (deliver at a specific time)
pulsar-client-tool produce -t my-topic -m "scheduled" --deliver-at "2026-04-01T12:00:00Z"
```

### Consume messages

```bash
# Consume indefinitely (Ctrl+C to stop)
pulsar-client-tool consume -t my-topic -s my-sub

# Consume 10 messages then exit
pulsar-client-tool consume -t my-topic -s my-sub -n 10

# JSON output (pipe to jq)
pulsar-client-tool consume -t my-topic -s my-sub -o json | jq .

# Shared subscription
pulsar-client-tool consume -t my-topic -s my-sub -S Shared

# Consume from the beginning of the topic
pulsar-client-tool consume -t my-topic -s my-sub --initial-position earliest
```

### Check connectivity

```bash
# Verify the Pulsar cluster is reachable
pulsar-client-tool ping

# Check a specific service URL
pulsar-client-tool ping --service-url pulsar://broker:6650
```

### Global flags

| Flag | Env Variable | Default | Description |
|------|-------------|---------|-------------|
| `--service-url` | `PULSAR_SERVICE_URL` | `pulsar://localhost:6650` | Pulsar service URL |
| `--auth-token` | `PULSAR_AUTH_TOKEN` | | Authentication token |
| `--tls-cert` | | | TLS client certificate |
| `--tls-key` | | | TLS client key |
| `--tls-ca` | | | TLS CA certificate |
| `--timeout` | | `30s` | Connection/operation timeout |
| `-o, --output` | | `text` | Output format (`text` or `json`) |

## Configuration

Environment variables are supported as defaults for global flags:

```bash
export PULSAR_SERVICE_URL=pulsar+ssl://pulsar.example.com:6651
export PULSAR_AUTH_TOKEN=eyJhbGci...
pulsar-client-tool consume -t my-topic -s my-sub
```

## Development

```bash
make build    # Build binary
make test     # Run tests
make vet      # Run go vet
make fmt      # Format code
make clean    # Remove binary
```

## License

MIT
