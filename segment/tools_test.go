package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestToolsSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	s.Tools.AppTools["Read"] = 15
	s.Tools.AppTools["Edit"] = 8
	s.Tools.LastUsed["app"] = "Edit"

	s.Tools.MCPTools[state.MCPServer{Name: "github", Type: "mcp"}] = map[string]int{"create_issue": 2}
	s.Tools.LastUsed["mcp"] = "github/create_issue"

	s.Tools.Skills["brainstorming"] = state.SkillUsage{Count: 1}
	s.Tools.LastUsed["skills"] = "brainstorming"

	seg := &ToolsSegment{}

	if seg.ID() != "tools" {
		t.Errorf("expected ID 'tools', got '%s'", seg.ID())
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should be a table with the 3 headers
	if !strings.Contains(output, "Tool") {
		t.Errorf("expected 'Tool' header, got %q", output)
	}
	if !strings.Contains(output, "Count") {
		t.Errorf("expected 'Count' header, got %q", output)
	}
	if !strings.Contains(output, "Last Used") {
		t.Errorf("expected 'Last Used' header, got %q", output)
	}

	// Should have rows for each active category
	if !strings.Contains(output, "App") {
		t.Errorf("expected App row, got %q", output)
	}
	if !strings.Contains(output, "MCP") {
		t.Errorf("expected MCP row, got %q", output)
	}
	if !strings.Contains(output, "Skills") {
		t.Errorf("expected Skills row, got %q", output)
	}

	// Should show last used values
	if !strings.Contains(output, "Edit") {
		t.Errorf("expected last used 'Edit' in App row, got %q", output)
	}
	if !strings.Contains(output, "github/create_issue") {
		t.Errorf("expected last used MCP tool, got %q", output)
	}
	if !strings.Contains(output, "brainstorming") {
		t.Errorf("expected last used skill, got %q", output)
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
		t.Errorf("expected empty output with no tools, got %q", output)
	}
}

func TestToolsSegmentShellRow(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Tools.InternalTools["Bash"] = 5
	s.Tools.LastUsed["internal"] = "Bash"

	seg := &ToolsSegment{}
	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if !strings.Contains(output, "Shell") {
		t.Errorf("expected Shell row for internal tools, got %q", output)
	}
	if !strings.Contains(output, "5") {
		t.Errorf("expected count 5, got %q", output)
	}
}

func TestToolsSegmentTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 24, "short"},
		{"exactly24charslong!!!!!!", 24, "exactly24charslong!!!!!!"},
		{"this-is-a-very-long-skill-name-that-exceeds-limit", 24, "this-is-a-very-long-s..."},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.max)
		if got != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.expected)
		}
	}
}
