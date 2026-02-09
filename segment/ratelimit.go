package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

type RateLimitSegment struct{}

func (r *RateLimitSegment) ID() string {
	return "ratelimit"
}

func (r *RateLimitSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.RateLimits
}

func (r *RateLimitSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.RateLimits.SevenDayTotal == 0 {
		return "", nil
	}

	percentage := float64(s.RateLimits.SevenDayUsed) / float64(s.RateLimits.SevenDayTotal) * 100.0

	// Color based on threshold
	var style lipgloss.Style
	if percentage >= float64(cfg.SevenDayThreshold) {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
	} else if percentage >= float64(cfg.SevenDayThreshold)-10 {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	}

	return style.Render(fmt.Sprintf("Rate: %.0f%%", percentage)), nil
}
