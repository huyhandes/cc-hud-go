package segment

import (
	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/format"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
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

	if st.Cost.TotalUSD > 0 {
		costStyle := style.GetRenderer().NewStyle().Foreground(style.ColorAccent).Bold(true)
		parts = append(parts, costStyle.Render("💰"+format.Cost(st.Cost.TotalUSD)))
	}

	if st.Cost.DurationMs > 0 {
		durationStyle := style.GetRenderer().NewStyle().Foreground(style.ColorHighlight)
		parts = append(parts, durationStyle.Render("⏱ "+format.Duration(st.Cost.DurationMs)))
	}

	// File changes moved to git line in multi-line layout
	// (removed from here to avoid duplication)

	if len(parts) == 0 {
		return "", nil
	}

	// Join with separator for better visual separation
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "  │  "
		}
		result += part
	}

	return result, nil
}
