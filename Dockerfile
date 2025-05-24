# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Install certificates for HTTPS connections
RUN apk add --no-cache ca-certificates git

# Create a non-root user for the final image
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files for better caching
COPY go.mod go.sum ./

# Download dependencies - this layer will be cached unless go.mod/go.sum change
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# GOARM is set based on the TARGETVARIANT
RUN export GOARM="${TARGETVARIANT#v}" && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags="-w -s" -o adsb2loki main.go

# Final stage
FROM scratch

# Copy certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary from builder
COPY --from=builder /build/adsb2loki /adsb2loki

# Use non-root user
USER appuser

# Expose any ports if needed (adjust based on your application)
# EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/adsb2loki"] 