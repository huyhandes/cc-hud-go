package parser

import (
	"testing"
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
