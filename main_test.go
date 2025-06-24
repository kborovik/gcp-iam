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
			args:     []string{"gcp-iam", "role", "show", "actions.Viewer"},
			expected: "Role: actions.Viewer",
		},
		{
			name:     "role search command",
			args:     []string{"gcp-iam", "role", "search", "actions.Viewer"},
			expected: "Found 1 roles matching 'actions.Viewer':",
		},
		{
			name:     "role compare command",
			args:     []string{"gcp-iam", "role", "compare", "role1"},
			expected: "gcp-iam role compare role1",
		},
		{
			name:     "permission show command",
			args:     []string{"gcp-iam", "permission", "show", "actions.agent.get"},
			expected: "Permission: actions.agent.get",
		},
		{
			name:     "permission search command",
			args:     []string{"gcp-iam", "permission", "search", "actions.agent.get"},
			expected: "Found 1 permissions matching 'actions.agent.get':",
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
