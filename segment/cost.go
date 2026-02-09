package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var (
	costYellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	costGreenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

type CostSegment struct{}

func (s CostSegment) ID() string {
	return "cost"
}

func (s CostSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Duration // Reuse Duration config for now
}

func (s CostSegment) Render(st *state.State, cfg *config.Config) (string, error) {
	if st.Cost.TotalUSD == 0 && st.Cost.DurationMs == 0 {
		return "", nil
	}

	parts := []string{}

	// Show cost if available
	if st.Cost.TotalUSD > 0 {
		costStr := fmt.Sprintf("$%.4f", st.Cost.TotalUSD)
		parts = append(parts, costYellowStyle.Render(costStr))
	}

	// Show duration if available
	if st.Cost.DurationMs > 0 {
		durationSec := st.Cost.DurationMs / 1000
		mins := durationSec / 60
		secs := durationSec % 60

		durationStr := ""
		if mins > 0 {
			durationStr = fmt.Sprintf("%dm%ds", mins, secs)
		} else {
			durationStr = fmt.Sprintf("%ds", secs)
		}
		parts = append(parts, durationStr)
	}

	// Show lines changed if available
	if st.Cost.LinesAdded > 0 || st.Cost.LinesRemoved > 0 {
		linesStr := fmt.Sprintf("+%d/-%d", st.Cost.LinesAdded, st.Cost.LinesRemoved)
		parts = append(parts, costGreenStyle.Render(linesStr))
	}

	if len(parts) == 0 {
		return "", nil
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}

	return result, nil
}
