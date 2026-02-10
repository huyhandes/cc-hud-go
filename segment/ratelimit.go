package segment

import (
	"fmt"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
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

	sevenDayPercentage := float64(s.RateLimits.SevenDayUsed) / float64(s.RateLimits.SevenDayTotal) * 100.0

	// Build gradient progress bar
	bar := style.RenderGradientBar(sevenDayPercentage, 10)

	// Update the output format to include the gradient bar
	output := fmt.Sprintf("📊 %s %.0f%%",
		bar,
		sevenDayPercentage,
	)

	return output, nil
}
