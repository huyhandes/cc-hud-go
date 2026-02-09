package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var (
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type ContextSegment struct{}

func (c *ContextSegment) ID() string {
	return "context"
}

func (c *ContextSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Context
}

func (c *ContextSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Context.TotalTokens == 0 {
		return "", nil
	}

	percentage := s.Context.Percentage

	// Choose color based on thresholds
	var style lipgloss.Style
	if percentage < 70 {
		style = greenStyle
	} else if percentage < 90 {
		style = yellowStyle
	} else {
		style = redStyle
	}

	// Build progress bar
	barWidth := 10
	filled := int(percentage / 10)
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	var display string
	if cfg.ContextValue == "tokens" {
		display = fmt.Sprintf("%d/%d", s.Context.UsedTokens, s.Context.TotalTokens)
	} else {
		display = fmt.Sprintf("%.0f%%", percentage)
	}

	return style.Render(fmt.Sprintf("[%s] %s", bar, display)), nil
}
