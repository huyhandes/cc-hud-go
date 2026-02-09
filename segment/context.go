package segment

import (
	"fmt"
	"strings"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
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

	// Choose style and icon based on thresholds
	var barStyle = style.ProgressGood
	var icon = "🟢"

	if percentage >= 90 {
		barStyle = style.ProgressDanger
		icon = "🔴"
	} else if percentage >= 70 {
		barStyle = style.ProgressWarning
		icon = "🟡"
	}

	// Build enhanced progress bar with gradient effect
	barWidth := 10
	filled := int(percentage / 10)
	if filled > barWidth {
		filled = barWidth
	}

	// Use different characters for better visual effect
	filledBar := strings.Repeat("●", filled)
	emptyBar := strings.Repeat("○", barWidth-filled)
	bar := filledBar + emptyBar

	var display string
	if cfg.ContextValue == "tokens" {
		display = fmt.Sprintf("%dk/%dk", s.Context.UsedTokens/1000, s.Context.TotalTokens/1000)
	} else {
		display = fmt.Sprintf("%.0f%%", percentage)
	}

	// Format with icon and styled bar
	return fmt.Sprintf("%s %s %s",
		icon,
		barStyle.Render(bar),
		style.ContextStyle.Render(display),
	), nil
}
