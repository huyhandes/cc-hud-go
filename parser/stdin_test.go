package parser

import (
	"testing"

	"github.com/huybui/cc-hud-go/state"
)

func TestParseStdin(t *testing.T) {
	input := `{
		"model": "claude-sonnet-4.5",
		"planType": "Pro",
		"context": {
			"used": 1500,
			"total": 200000
		}
	}`

	s := state.New()
	err := ParseStdin([]byte(input), s)

	if err != nil {
		t.Fatalf("ParseStdin failed: %v", err)
	}

	if s.Model.Name != "claude-sonnet-4.5" {
		t.Errorf("expected model 'claude-sonnet-4.5', got '%s'", s.Model.Name)
	}

	if s.Model.PlanType != "Pro" {
		t.Errorf("expected plan 'Pro', got '%s'", s.Model.PlanType)
	}

	if s.Context.UsedTokens != 1500 {
		t.Errorf("expected UsedTokens 1500, got %d", s.Context.UsedTokens)
	}

	if s.Context.TotalTokens != 200000 {
		t.Errorf("expected TotalTokens 200000, got %d", s.Context.TotalTokens)
	}
}
