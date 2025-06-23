package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestCLICommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "role show command",
			args:     []string{"gcp-iam", "role", "show", "testRole"},
			expected: "Role 'testRole' not found",
		},
		{
			name:     "role search command",
			args:     []string{"gcp-iam", "role", "search", "test"},
			expected: "Found",
		},
		{
			name:     "role compare command",
			args:     []string{"gcp-iam", "role", "compare", "role1"},
			expected: "gcp-iam role compare role1",
		},
		{
			name:     "permission show command",
			args:     []string{"gcp-iam", "permission", "show", "testPerm"},
			expected: "gcp-iam permission show testPerm",
		},
		{
			name:     "permission search command",
			args:     []string{"gcp-iam", "permission", "search", "test"},
			expected: "gcp-iam permission search test",
		},
		{
			name:     "update command",
			args:     []string{"gcp-iam", "update"},
			expected: "Fetching permissions for all roles",
		},
		{
			name:     "info command",
			args:     []string{"gcp-iam", "info"},
			expected: "GCP IAM Configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			r, w, _ := os.Pipe()
			os.Stdout = w

			err := cmd.Run(context.Background(), tt.args)

			w.Close()
			buf.ReadFrom(r)
			output := strings.TrimSpace(buf.String())

			if err != nil {
				t.Errorf("Command failed: %v", err)
				return
			}

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain '%s', got '%s'", tt.expected, output)
			}
		})
	}
}
