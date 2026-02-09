package segment

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestContextSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &ContextSegment{}

	if seg.ID() != "context" {
		t.Errorf("expected ID 'context', got '%s'", seg.ID())
	}

	// Test green threshold (<70%)
	s.Context.UsedTokens = 50000
	s.Context.TotalTokens = 200000
	s.UpdateDerived()

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "25") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	// Test yellow threshold (70-90%)
	s.Context.UsedTokens = 160000
	s.UpdateDerived()

	output, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "80") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}

	// Test red threshold (>90%)
	s.Context.UsedTokens = 190000
	s.UpdateDerived()

	output, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "95") {
		t.Errorf("expected percentage in output, got '%s'", output)
	}
}
