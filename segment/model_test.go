package segment

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestModelSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"
	s.Model.PlanType = "Pro"

	seg := &ModelSegment{}

	if seg.ID() != "model" {
		t.Errorf("expected ID 'model', got '%s'", seg.ID())
	}

	if !seg.Enabled(cfg) {
		t.Error("expected segment to be enabled by default")
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "claude-sonnet-4.5") {
		t.Errorf("expected output to contain model name, got '%s'", output)
	}

	if !strings.Contains(output, "Pro") {
		t.Errorf("expected output to contain plan type, got '%s'", output)
	}
}

func TestModelSegmentDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.Display.Model = false

	seg := &ModelSegment{}

	if seg.Enabled(cfg) {
		t.Error("expected segment to be disabled")
	}
}
