package parser

import (
	"testing"

	"github.com/huybui/cc-hud-go/state"
)

func TestCategorizeTool(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     ToolCategory
	}{
		{"App tool", "Read", CategoryApp},
		{"App tool lowercase", "read", CategoryApp},
		{"Bash", "Bash", CategoryInternal},
		{"MCP tool", "mcp__claude_ai_Atlassian__getConfluencePage", CategoryMCP},
		{"Skill", "Skill", CategorySkill},
		{"Custom", "MyCustomTool", CategoryCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeTool(tt.toolName)
			if got != tt.want {
				t.Errorf("CategorizeTool(%s) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

func TestParseTranscriptLine(t *testing.T) {
	line := `{"type":"tool_use","name":"Read","id":"tool_123"}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	if s.Tools.AppTools["Read"] != 1 {
		t.Errorf("expected Read count 1, got %d", s.Tools.AppTools["Read"])
	}
}

func TestParseTranscriptLineMCP(t *testing.T) {
	line := `{"type":"tool_use","name":"mcp__claude_ai_Atlassian__getConfluencePage"}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	// Check MCP tools map
	found := false
	for server, tools := range s.Tools.MCPTools {
		if server.Name == "claude_ai_Atlassian" {
			if tools["getConfluencePage"] == 1 {
				found = true
			}
		}
	}

	if !found {
		t.Error("expected MCP tool to be tracked")
	}
}
