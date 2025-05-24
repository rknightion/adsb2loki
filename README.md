# adsb2loki

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

### Enhanced CI Features

The GitHub Actions workflow includes these advanced features:

#### 🧪 Testing & Quality
- **Test Reporting**: JUnit XML test results displayed in GitHub UI
- **Coverage Tracking**: Code coverage reports with Codecov integration
- **PR Comments**: Automatic comments on PRs with test results and coverage
- **Benchmarks**: Performance benchmarks run on every commit
- **golangci-lint**: Comprehensive linting with multiple linters

#### 🔐 Security & Supply Chain
- **Container Signing**: Images signed with cosign for verification
- **SBOM Generation**: Software Bill of Materials for all artifacts
- **Dual Vulnerability Scanning**: Both Trivy and Grype scanners
- **Artifact Attestation**: Cryptographic proof of build provenance
- **Binary Signing**: Release binaries signed with cosign

#### 📊 Reporting & Metrics
- **Binary Size Tracking**: Size reports for all platform builds
- **Test Result Visualization**: Test results shown in PR checks
- **Coverage Badges**: Dynamic coverage badges (requires GIST_SECRET)
- **Benchmark History**: Performance tracking over time

#### 🚀 Release Automation
- **Automatic Changelog**: Generated from conventional commits
- **Multi-Platform Releases**: Binaries for all supported platforms
- **Signed Artifacts**: All release files include signatures
- **Rich Release Notes**: Includes sizes, verification instructions

### Verifying Artifacts

To verify a downloaded binary:
```bash
cosign verify-blob \
  --certificate adsb2loki-linux-amd64.pem \
  --signature adsb2loki-linux-amd64.sig \
  adsb2loki-linux-amd64
```

To verify container images:
```bash
cosign verify ghcr.io/rknightion/adsb2loki:latest
```

### Required Secrets

For full functionality, configure these secrets:
- `GIST_SECRET`: GitHub token with gist scope for coverage badges
- `CODECOV_TOKEN`: (Optional) For private repos on Codecov

### Container Registry

Images are automatically published to:
```