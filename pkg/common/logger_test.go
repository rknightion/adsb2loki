package common

import (
	"testing"
	"time"
)

func TestLogEntry(t *testing.T) {
	now := time.Now()
	entry := LogEntry{
		Timestamp: now,
		Labels: map[string]string{
			"service": "test",
			"env":     "dev",
		},
		Line: `{"test": "data"}`,
	}

	// Test timestamp
	if entry.Timestamp != now {
		t.Errorf("Expected timestamp %v, got %v", now, entry.Timestamp)
	}

	// Test labels
	if entry.Labels["service"] != "test" {
		t.Errorf("Expected service label 'test', got %s", entry.Labels["service"])
	}
	if entry.Labels["env"] != "dev" {
		t.Errorf("Expected env label 'dev', got %s", entry.Labels["env"])
	}

	// Test line
	expectedLine := `{"test": "data"}`
	if entry.Line != expectedLine {
		t.Errorf("Expected line '%s', got %s", expectedLine, entry.Line)
	}
}
