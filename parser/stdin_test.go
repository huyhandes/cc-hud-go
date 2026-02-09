package parser

import (
	"testing"

	"github.com/huybui/cc-hud-go/state"
)

func TestParseStdin(t *testing.T) {
	input := `{
		"session_id": "test123",
		"cwd": "/test/dir",
		"model": {
			"id": "claude-sonnet-4-5",
			"display_name": "Sonnet 4.5"
		},
		"workspace": {
			"current_dir": "/test/dir",
			"project_dir": "/test/dir"
		},
		"context_window": {
			"total_input_tokens": 15000,
			"total_output_tokens": 5000,
			"context_window_size": 200000,
			"used_percentage": 10.0,
			"remaining_percentage": 90.0,
			"current_usage": {
				"input_tokens": 8000,
				"output_tokens": 2000,
				"cache_creation_input_tokens": 5000,
				"cache_read_input_tokens": 2000
			}
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)

	if err != nil {
		t.Fatalf("ParseStdin failed: %v", err)
	}

	if s.Model.Name != "Sonnet 4.5" {
		t.Errorf("expected model 'Sonnet 4.5', got '%s'", s.Model.Name)
	}

	// UsedTokens should be sum of input_tokens + cache_creation_input_tokens + cache_read_input_tokens
	expectedUsed := 8000 + 5000 + 2000 // = 15000
	if s.Context.UsedTokens != expectedUsed {
		t.Errorf("expected UsedTokens %d, got %d", expectedUsed, s.Context.UsedTokens)
	}

	if s.Context.TotalTokens != 200000 {
		t.Errorf("expected TotalTokens 200000, got %d", s.Context.TotalTokens)
	}
}

func TestParseStdinWithAgent(t *testing.T) {
	input := `{
		"session_id": "test123",
		"cwd": "/test/dir",
		"model": {
			"id": "claude-opus-4-6",
			"display_name": "Opus"
		},
		"workspace": {
			"current_dir": "/test/dir",
			"project_dir": "/test/dir"
		},
		"context_window": {
			"total_input_tokens": 10000,
			"total_output_tokens": 2000,
			"context_window_size": 200000,
			"used_percentage": 6.0,
			"remaining_percentage": 94.0
		},
		"agent": {
			"name": "security-reviewer"
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)

	if err != nil {
		t.Fatalf("ParseStdin failed: %v", err)
	}

	if s.Model.Name != "Opus" {
		t.Errorf("expected model 'Opus', got '%s'", s.Model.Name)
	}

	if s.Agents.ActiveAgent != "security-reviewer" {
		t.Errorf("expected agent 'security-reviewer', got '%s'", s.Agents.ActiveAgent)
	}

	// Without current_usage, should use total_input_tokens
	if s.Context.UsedTokens != 10000 {
		t.Errorf("expected UsedTokens 10000, got %d", s.Context.UsedTokens)
	}
}

func TestParseStdinInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid JSON",
			input: `{"model": invalid}`,
		},
		{
			name:  "empty input",
			input: ``,
		},
		{
			name:  "partial JSON",
			input: `{"model": "test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			err := ParseStdin([]byte(tt.input), s)

			if err == nil {
				t.Error("expected error for invalid input")
			}
		})
	}
}
