package parser

import (
	"bufio"
	"encoding/json"
	"os"
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

// StdinData represents the JSON structure from Claude Code statusline API
type StdinData struct {
	SessionID      string `json:"session_id"`
	CWD            string `json:"cwd"`
	TranscriptPath string `json:"transcript_path"`
	Version        string `json:"version"`
	Model          struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`
	ContextWindow struct {
		TotalInputTokens    int     `json:"total_input_tokens"`
		TotalOutputTokens   int     `json:"total_output_tokens"`
		ContextWindowSize   int     `json:"context_window_size"`
		UsedPercentage      float64 `json:"used_percentage"`
		RemainingPercentage float64 `json:"remaining_percentage"`
		CurrentUsage        *struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"current_usage"`
	} `json:"context_window"`
	Cost *struct {
		TotalCostUSD       float64 `json:"total_cost_usd"`
		TotalDurationMs    int64   `json:"total_duration_ms"`
		TotalAPIDurationMs int64   `json:"total_api_duration_ms"`
		TotalLinesAdded    int     `json:"total_lines_added"`
		TotalLinesRemoved  int     `json:"total_lines_removed"`
	} `json:"cost,omitempty"`
	Exceeds200KTokens bool `json:"exceeds_200k_tokens"`
	OutputStyle       *struct {
		Name string `json:"name"`
	} `json:"output_style,omitempty"`
	Vim *struct {
		Mode string `json:"mode"`
	} `json:"vim,omitempty"`
	Agent *struct {
		Name string `json:"name"`
	} `json:"agent,omitempty"`
}

// ParseStdin parses stdin JSON from Claude Code and updates state
func ParseStdin(data []byte, s *state.State) error {
	var stdin StdinData
	if err := json.Unmarshal(data, &stdin); err != nil {
		return err
	}

	// Update session info
	s.Session.ID = stdin.SessionID
	s.Session.TranscriptPath = stdin.TranscriptPath

	// Update model info
	s.Model.Name = stdin.Model.DisplayName
	if s.Model.Name == "" {
		s.Model.Name = stdin.Model.ID
	}
	// Infer plan type from model ID (Pro/Max/Team indicators not in API)
	s.Model.PlanType = ""

	// Update context - use current usage if available, otherwise use totals
	s.Context.TotalInputTokens = stdin.ContextWindow.TotalInputTokens
	s.Context.TotalOutputTokens = stdin.ContextWindow.TotalOutputTokens
	s.Context.TotalTokens = stdin.ContextWindow.ContextWindowSize

	if stdin.ContextWindow.CurrentUsage != nil {
		// Calculate used tokens from current usage (input only, as per docs)
		s.Context.CurrentInputTokens = stdin.ContextWindow.CurrentUsage.InputTokens
		s.Context.CacheCreateTokens = stdin.ContextWindow.CurrentUsage.CacheCreationInputTokens
		s.Context.CacheReadTokens = stdin.ContextWindow.CurrentUsage.CacheReadInputTokens

		usedTokens := s.Context.CurrentInputTokens + s.Context.CacheCreateTokens + s.Context.CacheReadTokens
		s.Context.UsedTokens = usedTokens
	} else {
		// Fallback to total tokens
		s.Context.UsedTokens = stdin.ContextWindow.TotalInputTokens
	}

	// Update agent info if present
	if stdin.Agent != nil {
		s.Agents.ActiveAgent = stdin.Agent.Name
	}

	// Update cost info if present
	if stdin.Cost != nil {
		s.Cost.TotalUSD = stdin.Cost.TotalCostUSD
		s.Cost.DurationMs = stdin.Cost.TotalDurationMs
		s.Cost.APIDurationMs = stdin.Cost.TotalAPIDurationMs
		s.Cost.LinesAdded = stdin.Cost.TotalLinesAdded
		s.Cost.LinesRemoved = stdin.Cost.TotalLinesRemoved
	}

	// Update rate limits - not provided in API, keep existing values
	// Rate limits data is not in the Claude Code API

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

// ParseTranscript reads and parses the entire transcript file
func ParseTranscript(path string, s *state.State) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		// Ignore errors from individual lines, just continue
		_ = ParseTranscriptLine(line, s)
	}

	return scanner.Err()
}
