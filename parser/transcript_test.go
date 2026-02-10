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

func TestParseTranscriptLineSkillTracking(t *testing.T) {
	// Test that Skill tool calls extract the skill name from input parameters
	line := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "Skill",
				"input": {
					"skill": "superpowers:using-git-worktrees",
					"args": "feature-branch"
				}
			}]
		}
	}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	// Verify skill was tracked with its name
	if usage, ok := s.Tools.Skills["superpowers:using-git-worktrees"]; !ok {
		t.Errorf("Expected skill 'superpowers:using-git-worktrees' to be tracked")
	} else if usage.Count != 1 {
		t.Errorf("Expected skill count 1, got %d", usage.Count)
	}

	// Verify it wasn't counted as generic "Skill" in AppTools
	if s.Tools.AppTools["Skill"] != 0 {
		t.Errorf("Expected no generic 'Skill' count, got %d", s.Tools.AppTools["Skill"])
	}
}

func TestParseTranscriptLineSkillTrackingMultiple(t *testing.T) {
	// Test tracking multiple different skills
	lines := []string{
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "Skill",
					"input": {
						"skill": "superpowers:brainstorming"
					}
				}]
			}
		}`,
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "Skill",
					"input": {
						"skill": "superpowers:brainstorming"
					}
				}]
			}
		}`,
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "Skill",
					"input": {
						"skill": "superpowers:writing-plans"
					}
				}]
			}
		}`,
	}

	s := state.New()
	for _, line := range lines {
		_ = ParseTranscriptLine([]byte(line), s)
	}

	// Verify brainstorming was called twice
	if usage, ok := s.Tools.Skills["superpowers:brainstorming"]; !ok {
		t.Errorf("Expected skill 'superpowers:brainstorming' to be tracked")
	} else if usage.Count != 2 {
		t.Errorf("Expected brainstorming count 2, got %d", usage.Count)
	}

	// Verify writing-plans was called once
	if usage, ok := s.Tools.Skills["superpowers:writing-plans"]; !ok {
		t.Errorf("Expected skill 'superpowers:writing-plans' to be tracked")
	} else if usage.Count != 1 {
		t.Errorf("Expected writing-plans count 1, got %d", usage.Count)
	}
}

func TestParseTranscriptLineSkillTrackingFallback(t *testing.T) {
	// Test fallback for Skill calls without skill name
	line := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "Skill",
				"input": {}
			}]
		}
	}`

	s := state.New()
	err := ParseTranscriptLine([]byte(line), s)

	if err != nil {
		t.Fatalf("ParseTranscriptLine failed: %v", err)
	}

	// Verify it fell back to generic "Skill" in AppTools
	if s.Tools.AppTools["Skill"] != 1 {
		t.Errorf("Expected generic 'Skill' count 1, got %d", s.Tools.AppTools["Skill"])
	}

	// Verify no skills were tracked in Skills map
	if len(s.Tools.Skills) != 0 {
		t.Errorf("Expected no skills tracked, got %d", len(s.Tools.Skills))
	}
}
