package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rknightion/adsb2loki/pkg/common"
)

// Client represents a Loki client
type Client struct {
	url    string
	client *http.Client
}

// NewClient creates a new Loki client
func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PushLogs sends log entries to Loki
func (c *Client) PushLogs(ctx context.Context, entries []common.LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// Create the request payload
	streams := make([]map[string]interface{}, 0)
	for _, entry := range entries {
		// Create value array - timestamp, line, and optionally structured metadata
		value := []interface{}{
			fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
			entry.Line,
		}

		// Add structured metadata if present
		if len(entry.StructuredMetadata) > 0 {
			value = append(value, entry.StructuredMetadata)
		}

		stream := map[string]interface{}{
			"stream": entry.Labels,
			"values": [][]interface{}{value},
		}
		streams = append(streams, stream)
	}

	payload := map[string]interface{}{
		"streams": streams,
	}

	// Marshal the payload
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/loki/api/v1/push", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
