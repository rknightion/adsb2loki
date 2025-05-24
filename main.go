package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rknightion/adsb2loki/pkg/common"
	"github.com/rknightion/adsb2loki/pkg/flightaware"
	"github.com/rknightion/adsb2loki/pkg/loki"
	"github.com/rknightion/adsb2loki/pkg/otel"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger based on MODE environment variable
	mode := strings.ToLower(getEnvOrDefault("MODE", "loki"))
	var logger common.Logger
	var otelClient *otel.Client

	switch mode {
	case "otel":
		log.Println("Running in OpenTelemetry mode")
		client, err := otel.NewClient(ctx, "adsb2loki")
		if err != nil {
			log.Fatalf("Failed to create OpenTelemetry client: %v", err)
		}
		otelClient = client
		logger = client
		defer func() {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			if err := otelClient.Shutdown(shutdownCtx); err != nil {
				log.Printf("Failed to shutdown OpenTelemetry: %v", err)
			}
		}()
	case "loki":
		log.Println("Running in Loki mode")
		lokiURL := os.Getenv("LOKI_URL")
		if lokiURL == "" {
			log.Fatal("LOKI_URL environment variable is required in Loki mode")
		}
		logger = loki.NewClient(lokiURL)
	default:
		log.Fatalf("Invalid MODE '%s'. Must be 'loki' or 'otel'", mode)
	}

	// Create a ticker to fetch data periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the main loop
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			err := flightaware.FetchAndPushToLoki(ctx, logger)
			duration := time.Since(start)

			if err != nil {
				log.Printf("Error fetching and pushing data: %v", err)
				if otelClient != nil {
					otelClient.RecordPushError(ctx)
				}
			} else if otelClient != nil {
				otelClient.RecordFetchDuration(ctx, duration)
			}
		case <-sigChan:
			log.Println("Received shutdown signal, exiting...")
			return
		case <-ctx.Done():
			return
		}
	}
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
