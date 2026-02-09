package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestRender(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"
	s.Context.UsedTokens = 5000
	s.Context.TotalTokens = 10000

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Should have segments array
	if _, ok := result["segments"]; !ok {
		t.Error("output missing 'segments' field")
	}
}

func TestRenderWithDisabledSegments(t *testing.T) {
	cfg := config.Minimal()
	s := state.New()
	s.Model.Name = "claude-sonnet-4.5"

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should still be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Check that model segment is included (enabled in minimal)
	if !strings.Contains(output, "model") {
		t.Error("expected model segment in output")
	}
}

func TestRenderEmptyState(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	output, err := Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should be valid JSON even with empty state
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}
