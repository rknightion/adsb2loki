# Variables
BINARY_NAME=adsb2loki
DOCKER_IMAGE=ghcr.io/rknightion/adsb2loki
VERSION?=latest

# Build the Go binary
.PHONY: build
build:
	go build -o $(BINARY_NAME) -v

# Run the application
.PHONY: run
run:
	go run main.go

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	go tool cover -html=coverage.txt -o coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet
.PHONY: vet
vet:
	go vet ./...

# Run all linters
.PHONY: lint
lint: fmt vet
	@echo "Linting complete"

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Download dependencies
.PHONY: deps
deps:
	go mod download

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run -e MODE=$${MODE:-loki} \
	           -e LOKI_URL=$${LOKI_URL} \
	           -e AIRCRAFT_JSON_URL=$${AIRCRAFT_JSON_URL} \
	           -e OTEL_EXPORTER_OTLP_ENDPOINT=$${OTEL_EXPORTER_OTLP_ENDPOINT} \
	           $(DOCKER_IMAGE):$(VERSION)

# Push Docker image
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)

# Run with docker-compose
.PHONY: compose-up
compose-up:
	docker-compose up -d

# Stop docker-compose
.PHONY: compose-down
compose-down:
	docker-compose down

# View docker-compose logs
.PHONY: compose-logs
compose-logs:
	docker-compose logs -f

# All-in-one command to ensure code quality
.PHONY: check
check: tidy lint test
	@echo "All checks passed!"

# Help command
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the Go binary"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run all linters"
	@echo "  tidy           - Tidy dependencies"
	@echo "  deps           - Download dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-push    - Push Docker image"
	@echo "  compose-up     - Start services with docker-compose"
	@echo "  compose-down   - Stop services with docker-compose"
	@echo "  compose-logs   - View docker-compose logs"
	@echo "  check          - Run all quality checks"
	@echo "  help           - Show this help message"

.DEFAULT_GOAL := help 