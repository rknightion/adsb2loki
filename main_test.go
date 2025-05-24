package main

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
		setEnv       bool
	}{
		{
			name:         "environment variable exists",
			envKey:       "TEST_VAR",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
			setEnv:       true,
		},
		{
			name:         "environment variable does not exist",
			envKey:       "NON_EXISTENT_VAR",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
			setEnv:       false,
		},
		{
			name:         "empty environment variable",
			envKey:       "EMPTY_VAR",
			envValue:     "",
			defaultValue: "default",
			expected:     "",
			setEnv:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			// Test the function
			result := getEnvOrDefault(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestModeSelection(t *testing.T) {
	// Test that default mode is "loki"
	os.Unsetenv("MODE")
	mode := getEnvOrDefault("MODE", "loki")
	if mode != "loki" {
		t.Errorf("Expected default mode 'loki', got %s", mode)
	}

	// Test setting mode to "otel"
	os.Setenv("MODE", "otel")
	defer os.Unsetenv("MODE")
	mode = getEnvOrDefault("MODE", "loki")
	if mode != "otel" {
		t.Errorf("Expected mode 'otel', got %s", mode)
	}

	// Test case insensitive mode
	os.Setenv("MODE", "OTEL")
	mode = getEnvOrDefault("MODE", "loki")
	if mode != "OTEL" {
		t.Errorf("Expected mode 'OTEL', got %s", mode)
	}
}
