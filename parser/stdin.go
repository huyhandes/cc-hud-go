package parser

import (
	"encoding/json"

	"github.com/huyhandes/cc-hud-go/state"
)

// StdinData represents the JSON structure from Claude Code statusline API
type StdinData struct {
	TranscriptPath string `json:"transcript_path"`
	Model          struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow struct {
		TotalInputTokens  int     `json:"total_input_tokens"`
		TotalOutputTokens int     `json:"total_output_tokens"`
		ContextWindowSize int     `json:"context_window_size"`
		UsedPercentage    float64 `json:"used_percentage"`
		CurrentUsage      *struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"current_usage"`
	} `json:"context_window"`
	Cost *struct {
		TotalCostUSD      float64 `json:"total_cost_usd"`
		TotalDurationMs   int64   `json:"total_duration_ms"`
		TotalLinesAdded   int     `json:"total_lines_added"`
		TotalLinesRemoved int     `json:"total_lines_removed"`
	} `json:"cost,omitempty"`
	Agent *struct {
		Name string `json:"name"`
	} `json:"agent,omitempty"`
	RateLimits *struct {
		SevenDayUsed  int `json:"seven_day_used"`
		SevenDayTotal int `json:"seven_day_total"`
	} `json:"rate_limits,omitempty"`
}

// ParseStdin parses stdin JSON from Claude Code and updates state
func ParseStdin(data []byte, s *state.State) error {
	var stdin StdinData
	if err := json.Unmarshal(data, &stdin); err != nil {
		return err
	}

	s.Session.TranscriptPath = stdin.TranscriptPath

	s.Model.Name = stdin.Model.DisplayName
	if s.Model.Name == "" {
		s.Model.Name = stdin.Model.ID
	}

	s.Context.TotalInputTokens = stdin.ContextWindow.TotalInputTokens
	s.Context.TotalOutputTokens = stdin.ContextWindow.TotalOutputTokens
	s.Context.TotalTokens = stdin.ContextWindow.ContextWindowSize

	if stdin.ContextWindow.CurrentUsage != nil {
		s.Context.CurrentInputTokens = stdin.ContextWindow.CurrentUsage.InputTokens
		s.Context.CacheCreateTokens = stdin.ContextWindow.CurrentUsage.CacheCreationInputTokens
		s.Context.CacheReadTokens = stdin.ContextWindow.CurrentUsage.CacheReadInputTokens

		s.Context.UsedTokens = s.Context.CurrentInputTokens + s.Context.CacheCreateTokens + s.Context.CacheReadTokens
	} else {
		s.Context.UsedTokens = stdin.ContextWindow.TotalInputTokens
	}

	if stdin.Agent != nil {
		s.Agents.ActiveAgent = stdin.Agent.Name
	}

	if stdin.Cost != nil {
		s.Cost.TotalUSD = stdin.Cost.TotalCostUSD
		s.Cost.DurationMs = stdin.Cost.TotalDurationMs
		s.Cost.LinesAdded = stdin.Cost.TotalLinesAdded
		s.Cost.LinesRemoved = stdin.Cost.TotalLinesRemoved
	}

	if stdin.RateLimits != nil {
		s.RateLimits.SevenDayUsed = stdin.RateLimits.SevenDayUsed
		s.RateLimits.SevenDayTotal = stdin.RateLimits.SevenDayTotal
	}

	return nil
}
