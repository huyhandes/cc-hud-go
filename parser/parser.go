package parser

import (
	"encoding/json"

	"github.com/huybui/cc-hud-go/state"
)

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
