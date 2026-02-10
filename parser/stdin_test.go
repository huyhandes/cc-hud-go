package parser

import (
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
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

func TestParseStdinWithRateLimits(t *testing.T) {
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
			"total_input_tokens": 10000,
			"total_output_tokens": 2000,
			"context_window_size": 200000,
			"used_percentage": 6.0,
			"remaining_percentage": 94.0
		},
		"rate_limits": {
			"hourly_used": 10,
			"hourly_total": 50,
			"seven_day_used": 450,
			"seven_day_total": 1000
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)

	if err != nil {
		t.Fatalf("ParseStdin failed: %v", err)
	}

	if s.RateLimits.HourlyUsed != 10 {
		t.Errorf("expected HourlyUsed 10, got %d", s.RateLimits.HourlyUsed)
	}

	if s.RateLimits.HourlyTotal != 50 {
		t.Errorf("expected HourlyTotal 50, got %d", s.RateLimits.HourlyTotal)
	}

	if s.RateLimits.SevenDayUsed != 450 {
		t.Errorf("expected SevenDayUsed 450, got %d", s.RateLimits.SevenDayUsed)
	}

	if s.RateLimits.SevenDayTotal != 1000 {
		t.Errorf("expected SevenDayTotal 1000, got %d", s.RateLimits.SevenDayTotal)
	}
}

func TestParseStdinInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid JSON", `{"model": invalid}`},
		{"empty input", ``},
		{"partial JSON", `{"model": "test"`},
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

func TestParseStdinModelNameFallback(t *testing.T) {
	// display_name empty -> should use id
	input := `{
		"session_id": "test",
		"model": {"id": "claude-opus-4-6", "display_name": ""},
		"workspace": {},
		"context_window": {"context_window_size": 200000}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Model.Name != "claude-opus-4-6" {
		t.Errorf("expected model ID fallback, got %q", s.Model.Name)
	}
}

func TestParseStdinWithCost(t *testing.T) {
	input := `{
		"session_id": "test",
		"model": {"id": "test", "display_name": "Test"},
		"workspace": {},
		"context_window": {"context_window_size": 200000},
		"cost": {
			"total_cost_usd": 0.567,
			"total_duration_ms": 120000,
			"total_api_duration_ms": 90000,
			"total_lines_added": 100,
			"total_lines_removed": 25
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Cost.TotalUSD != 0.567 {
		t.Errorf("expected cost 0.567, got %f", s.Cost.TotalUSD)
	}
	if s.Cost.DurationMs != 120000 {
		t.Errorf("expected duration 120000, got %d", s.Cost.DurationMs)
	}
	if s.Cost.APIDurationMs != 90000 {
		t.Errorf("expected api duration 90000, got %d", s.Cost.APIDurationMs)
	}
	if s.Cost.LinesAdded != 100 {
		t.Errorf("expected 100 lines added, got %d", s.Cost.LinesAdded)
	}
	if s.Cost.LinesRemoved != 25 {
		t.Errorf("expected 25 lines removed, got %d", s.Cost.LinesRemoved)
	}
}

func TestParseStdinNilOptionalFields(t *testing.T) {
	// No cost, no agent, no rate_limits, no current_usage
	input := `{
		"session_id": "bare",
		"model": {"id": "test", "display_name": "Minimal"},
		"workspace": {},
		"context_window": {
			"total_input_tokens": 5000,
			"context_window_size": 200000
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Cost.TotalUSD != 0 {
		t.Error("expected zero cost when cost field is nil")
	}
	if s.Agents.ActiveAgent != "" {
		t.Error("expected no agent when agent field is nil")
	}
	if s.RateLimits.HourlyTotal != 0 {
		t.Error("expected zero rate limits when field is nil")
	}
	// Without current_usage, UsedTokens should fallback to TotalInputTokens
	if s.Context.UsedTokens != 5000 {
		t.Errorf("expected UsedTokens=5000 fallback, got %d", s.Context.UsedTokens)
	}
}
