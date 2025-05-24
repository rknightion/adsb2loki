package common

import (
	"context"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp          time.Time
	Labels             map[string]string
	Line               string
	StructuredMetadata map[string]string // Optional structured metadata
}

// Logger interface that can be implemented by different backends
type Logger interface {
	PushLogs(ctx context.Context, entries []LogEntry) error
}
