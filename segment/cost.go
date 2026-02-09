package segment

import (
	"fmt"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
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

	// Show cost with money icon - Orange/Accent (emphasis on cost)
	if st.Cost.TotalUSD > 0 {
		costStr := fmt.Sprintf("💰$%.4f", st.Cost.TotalUSD)
		costStyle := style.GetRenderer().NewStyle().Foreground(style.ColorAccent).Bold(true)
		parts = append(parts, costStyle.Render(costStr))
	}

	// Show duration with clock icon - Cyan (time tracking)
	if st.Cost.DurationMs > 0 {
		durationSec := st.Cost.DurationMs / 1000
		mins := durationSec / 60
		secs := durationSec % 60

		durationStr := ""
		if mins > 0 {
			durationStr = fmt.Sprintf("⏱ %dm%ds", mins, secs)
		} else {
			durationStr = fmt.Sprintf("⏱ %ds", secs)
		}
		durationStyle := style.GetRenderer().NewStyle().Foreground(style.ColorHighlight)
		parts = append(parts, durationStyle.Render(durationStr))
	}

	// Show lines changed with code icon - Split colors for +/-
	if st.Cost.LinesAdded > 0 || st.Cost.LinesRemoved > 0 {
		addStyle := style.GetRenderer().NewStyle().Foreground(style.ColorSuccess)
		removeStyle := style.GetRenderer().NewStyle().Foreground(style.ColorDanger)
		linesStr := fmt.Sprintf("📝 %s%s%s",
			addStyle.Render(fmt.Sprintf("+%d", st.Cost.LinesAdded)),
			style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
			removeStyle.Render(fmt.Sprintf("-%d", st.Cost.LinesRemoved)),
		)
		parts = append(parts, linesStr)
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
