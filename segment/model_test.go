package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
)

func TestModelSegment(t *testing.T) {
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"

	seg := &ModelSegment{}

	if seg.ID() != "model" {
		t.Errorf("expected ID 'model', got '%s'", seg.ID())
	}

	output, err := seg.Render(s)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "claude-sonnet-4.5") {
		t.Errorf("expected output to contain model name, got '%s'", output)
	}

	t.Logf("Model output: %s", output)
}

func TestModelSegmentEmpty(t *testing.T) {
	s := state.New()
	seg := &ModelSegment{}

	output, err := seg.Render(s)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output with no model name, got '%s'", output)
	}
}
