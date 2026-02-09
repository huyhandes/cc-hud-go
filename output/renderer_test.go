package output

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestRender(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "Sonnet 4.5"
	s.Context.UsedTokens = 5000
	s.Context.TotalTokens = 10000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should be plain text, not empty
	if output == "" {
		t.Error("expected non-empty output")
	}

	// Should contain model name
	if !strings.Contains(output, "Sonnet") {
		t.Errorf("expected output to contain model name, got: %s", output)
	}

	// Should contain separator
	if !strings.Contains(output, "|") {
		t.Errorf("expected output to contain separator '|', got: %s", output)
	}
}

func TestRenderWithDisabledSegments(t *testing.T) {
	cfg := config.Minimal()
	s := state.New()
	s.Model.Name = "Sonnet 4.5"

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should be plain text
	if output == "" {
		t.Error("expected non-empty output")
	}

	// With minimal config, should still have some content
	if !strings.Contains(output, "Sonnet") {
		t.Errorf("expected model in output, got: %s", output)
	}
}

func TestRenderEmptyState(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should return plain text (could be empty or have default segments)
	// Output is a string, no JSON parsing needed
	if output == "" {
		// Empty output is acceptable for empty state
		t.Logf("Empty output for empty state (acceptable)")
	} else {
		// If there's output, it should be plain text
		t.Logf("Output: %s", output)
	}
}
