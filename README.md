# adsb2loki

‼️🚨🔔 ALL credit to this goes to @burnettdev (https://github.com/burnettdev/flightaware2loki). For real I can't even write a hello world in Go I just forked this repo and had cursor pimp it to high heaven the actual underlying code is alllllll him

🚀 A Go service that fetches aircraft data from ADS-B sources (like FlightAware's SkyAware) and streams it to Grafana Loki or OpenTelemetry for real-time monitoring and analysis.

## Features

- Real-time aircraft data collection from FlightAware SkyAware
- Automatic data streaming to Grafana Loki
- Configurable via environment variables
- Graceful shutdown handling
- Efficient batch processing of aircraft data

## Prerequisites

- Go 1.21 or later
- FlightAware SkyAware instance
- Grafana Loki instance

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Required for both modes
AIRCRAFT_JSON_URL=http://your-skyaware-instance/skyaware/data/aircraft.json

# Mode selection (defaults to 'loki' for backward compatibility)
MODE=loki  # or 'otel' for OpenTelemetry

# Required for Loki mode
LOKI_URL=http://your-loki-instance:3100

# Required for OpenTelemetry mode (standard OTEL env vars)
# OTEL_EXPORTER_OTLP_ENDPOINT=http://your-otel-collector:4318
# OTEL_EXPORTER_OTLP_LOGS_ENDPOINT=http://your-otel-collector:4318/v1/logs
# OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://your-otel-collector:4318/v1/metrics
```

### Operating Modes

The service supports two modes of operation:

1. **Loki Mode** (default) - Pushes logs directly to Grafana Loki
2. **OpenTelemetry Mode** - Sends logs and metrics via OpenTelemetry Protocol (OTLP)

Set the mode using the `MODE` environment variable:
- `MODE=loki` - Use Loki HTTP API (default)
- `MODE=otel` - Use OpenTelemetry exporters

### OpenTelemetry Configuration

When running in OpenTelemetry mode, the service uses standard OTEL environment variables:

- `OTEL_EXPORTER_OTLP_ENDPOINT` - Base endpoint for all signals
- `OTEL_EXPORTER_OTLP_LOGS_ENDPOINT` - Specific endpoint for logs
- `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` - Specific endpoint for metrics
- `OTEL_EXPORTER_OTLP_HEADERS` - Headers to include in requests
- `OTEL_EXPORTER_OTLP_TIMEOUT` - Export timeout (default: 10s)

The service exports the following metrics in OpenTelemetry mode:
- `adsb.aircraft.count` - Number of aircraft processed
- `adsb.fetch.duration` - Duration of aircraft data fetch operations
- `adsb.push.errors` - Number of errors pushing data

## Installation

1. Clone the repository:
```bash
git clone https://github.com/rknightion/adsb2loki.git
cd adsb2loki
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build
```

## Usage

Run the application:
```bash
./adsb2loki
```

The service will:
- Fetch aircraft data every 5 seconds
- Push the data to Loki with appropriate labels
- Log any errors that occur during the process

## Data Structure

Each aircraft entry in Loki includes:
- Timestamp from FlightAware
- Labels for easy querying
- Full aircraft data as JSON

## Contributing

Feel free to open issues or submit pull requests!

## License

MIT License - feel free to use this project for whatever you'd like!

## Docker Support

### Building the Docker Image

Build the image locally:
```bash
docker build -t adsb2loki:latest .
```

### Running with Docker

Using docker run:
```bash
# Loki mode
docker run -e MODE=loki \
           -e LOKI_URL=http://your-loki-instance:3100 \
           -e AIRCRAFT_JSON_URL=http://your-skyaware-instance/skyaware/data/aircraft.json \
           ghcr.io/rknightion/adsb2loki:latest

# OpenTelemetry mode
docker run -e MODE=otel \
           -e OTEL_EXPORTER_OTLP_ENDPOINT=http://your-otel-collector:4318 \
           -e AIRCRAFT_JSON_URL=http://your-skyaware-instance/skyaware/data/aircraft.json \
           ghcr.io/rknightion/adsb2loki:latest
```

Using docker-compose:
```bash
docker-compose up -d
```

### Multi-Architecture Support

The CI/CD pipeline automatically builds images for both `linux/amd64` and `linux/arm64` architectures, making it compatible with both x86 and ARM systems (including Raspberry Pi).

## CI/CD with GitHub Actions

This project includes a comprehensive CI/CD pipeline that:

- **Tests**: Runs Go tests, vet, and fmt checks
- **Builds**: Creates multi-architecture Docker images
- **Publishes**: Pushes images to GitHub Container Registry (ghcr.io)
- **Security**: Scans images for vulnerabilities using Trivy
- **Tagging**: Supports semantic versioning and branch-based tags

### Container Registry

Images are automatically published to:
```
ghcr.io/rknightion/adsb2loki
```

Available tags:
- `latest` - Latest stable release from main branch
- `main` - Latest build from main branch
- `v1.2.3` - Specific version tags
- `main-sha-abc123` - Commit-specific builds

### Using Pre-built Images

Pull the latest image:
```bash
docker pull ghcr.io/rknightion/adsb2loki:latest
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Running Tests with Coverage

```bash
go test -v ./... -cover
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

## Multi-Architecture Support

The project supports multiple architectures for both binaries and Docker images:

### Binary Releases

Pre-built binaries are available for:
- **Linux**: amd64, arm64, armv7 (32-bit ARM), armv6 (Raspberry Pi)
- **Windows**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)

### Docker Images

Multi-architecture Docker images support:
- `linux/amd64` - Standard x86-64 servers
- `linux/arm64` - 64-bit ARM (AWS Graviton, Apple Silicon under emulation)
- `linux/arm/v7` - 32-bit ARM (newer Raspberry Pi models)
- `linux/arm/v6` - 32-bit ARM (older Raspberry Pi models)

The correct architecture will be automatically selected when you pull the image.

## Testing

The project includes comprehensive unit tests for all major components:

- **Common Package**: Data structure tests
- **Loki Package**: HTTP client and payload formatting tests
- **OpenTelemetry Package**: Basic client creation tests
- **FlightAware Package**: Data fetching and transformation tests
- **Models Package**: JSON marshaling/unmarshaling tests
- **Main Package**: Configuration and environment variable tests

Run tests with:
```bash
make test
```

Or with coverage:
```bash
make test-coverage
```

## Troubleshooting

### JSON Parsing Errors

If you see errors like:
```
json: cannot unmarshal string into Go struct field .aircraft.alt_baro of type int
```

This is because ADS-B data can contain mixed types for certain fields (e.g., altitude can be a number or the string "ground"). The application handles these cases automatically by accepting both string and numeric values for fields like:
- `alt_baro` / `alt_geom` - Can be numeric altitude or "ground"
- `gs` / `ias` / `tas` - Speed fields that might be missing or non-numeric
- `baro_rate` / `geom_rate` - Rate fields that might be missing or non-numeric
