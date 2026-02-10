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

	// With 26 total tools (23+2+1) and default threshold of 5, should show table
	// Should show category counts in table format
	if !strings.Contains(output, "App") {
		t.Errorf("expected 'App' category in output, got '%s'", output)
	}

	// Should show MCP and Skills categories
	if !strings.Contains(output, "MCP") {
		t.Errorf("expected 'MCP' category in output, got '%s'", output)
	}

	if !strings.Contains(output, "Skills") {
		t.Errorf("expected 'Skills' category in output, got '%s'", output)
	}

	// Should contain table borders since above threshold
	if !strings.Contains(output, "┌") {
		t.Errorf("expected table format for 26 tools (threshold=5), got '%s'", output)
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

func TestToolsSegmentTableThreshold(t *testing.T) {
	// Below threshold - should be inline
	s := state.New()
	s.Tools.AppTools["Read"] = 3
	s.Tools.AppTools["Edit"] = 1
	s.Tools.Skills["test"] = state.SkillUsage{Count: 1}

	cfg := config.Default()
	cfg.Tables.ToolsThreshold = 5

	seg := &ToolsSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be inline (no table borders)
	if strings.Contains(result, "┌") {
		t.Error("Expected inline format below threshold")
	}

	// Above threshold - should be table
	s.Tools.AppTools["Read"] = 10
	s.Tools.MCPTools[state.MCPServer{Name: "github", Type: "mcp"}] = map[string]int{
		"create_issue": 3,
	}
	s.Tools.Skills["test2"] = state.SkillUsage{Count: 2}
	s.Tools.CustomTools["custom"] = 1

	result, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be table format
	if !strings.Contains(result, "┌") {
		t.Error("Expected table format above threshold")
	}
}
