package segment

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestToolsSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	s.Tools.AppTools["Read"] = 15
	s.Tools.AppTools["Edit"] = 8
	s.Tools.MCPTools[state.MCPServer{Name: "github", Type: "mcp"}] = map[string]int{
		"create_issue": 2,
	}
	s.Tools.Skills["brainstorming"] = state.SkillUsage{Count: 1}

	seg := &ToolsSegment{}

	if seg.ID() != "tools" {
		t.Errorf("expected ID 'tools', got '%s'", seg.ID())
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Format changed to use icon, check for tool count instead
	if !strings.Contains(output, "26") {
		t.Errorf("expected tool count '26' in output, got '%s'", output)
	}

	// Should show category counts
	if !strings.Contains(output, "App:") {
		t.Errorf("expected 'App:' category in output, got '%s'", output)
	}
}

func TestToolsSegmentEmpty(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &ToolsSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no tools, got '%s'", output)
	}
}
