package parser

import (
	"encoding/json"
	"strings"

	"github.com/huybui/cc-hud-go/state"
)

type ToolCategory int

const (
	CategoryApp ToolCategory = iota
	CategoryInternal
	CategoryCustom
	CategoryMCP
	CategorySkill
)

var appTools = map[string]bool{
	"read":      true,
	"write":     true,
	"edit":      true,
	"bash":      true,
	"glob":      true,
	"grep":      true,
	"task":      true,
	"webfetch":  true,
	"websearch": true,
}

// StdinData represents the JSON structure from Claude Code
type StdinData struct {
	Model    string `json:"model"`
	PlanType string `json:"planType"`
	Context  struct {
		Used  int `json:"used"`
		Total int `json:"total"`
	} `json:"context"`
	RateLimits *struct {
		HourlyUsed    int `json:"hourlyUsed"`
		HourlyTotal   int `json:"hourlyTotal"`
		SevenDayUsed  int `json:"sevenDayUsed"`
		SevenDayTotal int `json:"sevenDayTotal"`
	} `json:"rateLimits,omitempty"`
}

// ParseStdin parses stdin JSON and updates state
func ParseStdin(data []byte, s *state.State) error {
	var stdin StdinData
	if err := json.Unmarshal(data, &stdin); err != nil {
		return err
	}

	// Update model info
	s.Model.Name = stdin.Model
	s.Model.PlanType = stdin.PlanType

	// Update context
	s.Context.UsedTokens = stdin.Context.Used
	s.Context.TotalTokens = stdin.Context.Total

	// Update rate limits if present
	if stdin.RateLimits != nil {
		s.RateLimits.HourlyUsed = stdin.RateLimits.HourlyUsed
		s.RateLimits.HourlyTotal = stdin.RateLimits.HourlyTotal
		s.RateLimits.SevenDayUsed = stdin.RateLimits.SevenDayUsed
		s.RateLimits.SevenDayTotal = stdin.RateLimits.SevenDayTotal
	}

	return nil
}

// CategorizeTool determines the category of a tool by name
func CategorizeTool(name string) ToolCategory {
	lower := strings.ToLower(name)

	// Check for MCP pattern
	if strings.HasPrefix(lower, "mcp__") {
		return CategoryMCP
	}

	// Check for Skill
	if lower == "skill" {
		return CategorySkill
	}

	// Check for internal (Bash is special) - must check before appTools
	if lower == "bash" {
		return CategoryInternal
	}

	// Check for app tools
	if appTools[lower] {
		return CategoryApp
	}

	// Everything else is custom
	return CategoryCustom
}

type TranscriptLine struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ParseTranscriptLine parses a single JSONL line and updates state
func ParseTranscriptLine(data []byte, s *state.State) error {
	var line TranscriptLine
	if err := json.Unmarshal(data, &line); err != nil {
		return err
	}

	// Only process tool_use events
	if line.Type != "tool_use" {
		return nil
	}

	category := CategorizeTool(line.Name)

	switch category {
	case CategoryApp:
		s.Tools.AppTools[line.Name]++

	case CategoryInternal:
		s.Tools.InternalTools[line.Name]++

	case CategoryCustom:
		s.Tools.CustomTools[line.Name]++

	case CategoryMCP:
		// Parse MCP tool name: mcp__<server>__<tool>
		parts := strings.Split(line.Name, "__")
		if len(parts) >= 3 {
			server := state.MCPServer{
				Name: parts[1],
				Type: "mcp",
			}

			if s.Tools.MCPTools[server] == nil {
				s.Tools.MCPTools[server] = make(map[string]int)
			}

			toolName := strings.Join(parts[2:], "__")
			s.Tools.MCPTools[server][toolName]++
		}

	case CategorySkill:
		// Skills need additional parsing from the tool parameters
		// For now, just count as app tool
		s.Tools.AppTools["Skill"]++
	}

	return nil
}
