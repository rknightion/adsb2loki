package loki

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rknightion/adsb2loki/pkg/common"
)

func TestNewClient(t *testing.T) {
	url := "http://localhost:3100"
	client := NewClient(url)

	if client.url != url {
		t.Errorf("Expected URL %s, got %s", url, client.url)
	}

	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if client.client.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", client.client.Timeout)
	}
}

func TestPushLogs(t *testing.T) {
	// Create a test server
	var receivedPayload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/loki/api/v1/push" {
			t.Errorf("Expected path /loki/api/v1/push, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read and verify the payload
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		if err := json.Unmarshal(body, &receivedPayload); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)

	// Create test entries
	now := time.Now()
	entries := []common.LogEntry{
		{
			Timestamp: now,
			Labels: map[string]string{
				"service": "test",
				"env":     "dev",
			},
			Line: `{"aircraft": "test123"}`,
		},
	}

	// Push logs
	ctx := context.Background()
	err := client.PushLogs(ctx, entries)
	if err != nil {
		t.Fatalf("Failed to push logs: %v", err)
	}

	// Verify the payload
	streams, ok := receivedPayload["streams"].([]interface{})
	if !ok || len(streams) != 1 {
		t.Fatalf("Expected 1 stream, got %v", receivedPayload["streams"])
	}

	stream := streams[0].(map[string]interface{})
	labels := stream["stream"].(map[string]interface{})
	if labels["service"] != "test" {
		t.Errorf("Expected service label 'test', got %v", labels["service"])
	}

	values := stream["values"].([]interface{})
	if len(values) != 1 {
		t.Fatalf("Expected 1 value, got %d", len(values))
	}
}

func TestPushLogsEmptyEntries(t *testing.T) {
	client := NewClient("http://localhost:3100")

	// Should return nil for empty entries
	err := client.PushLogs(context.Background(), []common.LogEntry{})
	if err != nil {
		t.Errorf("Expected nil error for empty entries, got %v", err)
	}
}

func TestPushLogsServerError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	entries := []common.LogEntry{
		{
			Timestamp: time.Now(),
			Labels:    map[string]string{"test": "test"},
			Line:      "test",
		},
	}

	// Should not return error even on server error (current implementation doesn't check response)
	err := client.PushLogs(context.Background(), entries)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
