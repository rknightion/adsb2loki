package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/rknightion/adsb2loki/pkg/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Client represents an OpenTelemetry client for logging and metrics
type Client struct {
	logger          log.Logger
	loggerProvider  *sdklog.LoggerProvider
	meterProvider   *sdkmetric.MeterProvider
	aircraftCounter metric.Int64Counter
	fetchDuration   metric.Float64Histogram
	pushErrors      metric.Int64Counter
}

// NewClient creates a new OpenTelemetry client
func NewClient(ctx context.Context, serviceName string) (*Client, error) {
	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			attribute.String("service.name", serviceName),
			attribute.String("service.version", "1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create log exporter - uses OTEL_EXPORTER_OTLP_LOGS_ENDPOINT or OTEL_EXPORTER_OTLP_ENDPOINT
	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	// Create log provider
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)

	// Create metric exporter - uses OTEL_EXPORTER_OTLP_METRICS_ENDPOINT or OTEL_EXPORTER_OTLP_ENDPOINT
	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)

	// Set global providers
	otel.SetMeterProvider(meterProvider)

	// Create logger
	logger := loggerProvider.Logger("adsb2loki")

	// Create metrics
	meter := meterProvider.Meter("adsb2loki")

	aircraftCounter, err := meter.Int64Counter(
		"adsb.aircraft.count",
		metric.WithDescription("Number of aircraft processed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create aircraft counter: %w", err)
	}

	fetchDuration, err := meter.Float64Histogram(
		"adsb.fetch.duration",
		metric.WithDescription("Duration of aircraft data fetch operations"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch duration histogram: %w", err)
	}

	pushErrors, err := meter.Int64Counter(
		"adsb.push.errors",
		metric.WithDescription("Number of errors pushing data"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create push errors counter: %w", err)
	}

	return &Client{
		logger:          logger,
		loggerProvider:  loggerProvider,
		meterProvider:   meterProvider,
		aircraftCounter: aircraftCounter,
		fetchDuration:   fetchDuration,
		pushErrors:      pushErrors,
	}, nil
}

// PushLogs pushes log entries via OpenTelemetry
func (c *Client) PushLogs(ctx context.Context, entries []common.LogEntry) error {
	// Record aircraft count metric
	c.aircraftCounter.Add(ctx, int64(len(entries)))

	// Emit logs
	for _, entry := range entries {
		// Create a log record
		record := log.Record{}
		record.SetTimestamp(entry.Timestamp)
		record.SetBody(log.StringValue(entry.Line))
		record.SetSeverity(log.SeverityInfo)

		// Add attributes from labels
		attrs := make([]log.KeyValue, 0, len(entry.Labels))
		for k, v := range entry.Labels {
			attrs = append(attrs, log.String(k, v))
		}
		record.AddAttributes(attrs...)

		// Emit the log
		c.logger.Emit(ctx, record)
	}

	return nil
}

// RecordFetchDuration records the duration of a fetch operation
func (c *Client) RecordFetchDuration(ctx context.Context, duration time.Duration) {
	c.fetchDuration.Record(ctx, duration.Seconds())
}

// RecordPushError increments the push error counter
func (c *Client) RecordPushError(ctx context.Context) {
	c.pushErrors.Add(ctx, 1)
}

// Shutdown gracefully shuts down the OpenTelemetry providers
func (c *Client) Shutdown(ctx context.Context) error {
	// Shutdown logger provider
	if err := c.loggerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown logger provider: %w", err)
	}

	// Shutdown meter provider
	if err := c.meterProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown meter provider: %w", err)
	}

	return nil
}
