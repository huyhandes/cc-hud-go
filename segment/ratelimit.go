package segment

import (
	"fmt"
	"time"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type RateLimitSegment struct{}

func (r *RateLimitSegment) ID() string {
	return "ratelimit"
}

func (r *RateLimitSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.RateLimits
}

func (r *RateLimitSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	// This segment now only renders 7d limit
	// 5h limit is rendered separately by FiveHourSegment

	// Prefer OAuth API data (more accurate)
	if s.RateLimits.SevenDayPercent > 0 {
		bar7d := style.RenderGradientBar(s.RateLimits.SevenDayPercent, 10)
		percentStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(s.RateLimits.SevenDayPercent))
		return fmt.Sprintf("📊 %s %s", bar7d, percentStyle.Render(fmt.Sprintf("%.0f%%", s.RateLimits.SevenDayPercent))), nil
	}

	// Fallback to stdin data (if provided)
	if s.RateLimits.SevenDayTotal == 0 {
		return "", nil
	}

	sevenDayPercentage := float64(s.RateLimits.SevenDayUsed) / float64(s.RateLimits.SevenDayTotal) * 100.0
	bar := style.RenderGradientBar(sevenDayPercentage, 10)
	percentStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(sevenDayPercentage))
	return fmt.Sprintf("📊 %s %s", bar, percentStyle.Render(fmt.Sprintf("%.0f%%", sevenDayPercentage))), nil
}

// FiveHourSegment displays 5-hour rate limit with elapsed time
type FiveHourSegment struct{}

func (f *FiveHourSegment) ID() string {
	return "fivehour"
}

func (f *FiveHourSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.RateLimits
}

func (f *FiveHourSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	// Only render if OAuth data available
	if s.RateLimits.FiveHourPercent <= 0 {
		return "", nil
	}

	bar5h := style.RenderGradientBar(s.RateLimits.FiveHourPercent, 10)
	percentStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(s.RateLimits.FiveHourPercent))

	// Calculate time remaining in 5h window
	timeInfo := ""
	if s.RateLimits.FiveHourResetsAt != "" {
		if resetTime, err := time.Parse(time.RFC3339, s.RateLimits.FiveHourResetsAt); err == nil {
			now := time.Now()
			if resetTime.After(now) {
				remaining := resetTime.Sub(now)

				// Format remaining time
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60

				if hours > 0 {
					timeInfo = fmt.Sprintf(" (%dh%dm)", hours, minutes)
				} else {
					timeInfo = fmt.Sprintf(" (%dm)", minutes)
				}
			}
		}
	}

	return fmt.Sprintf("⏱️ %s %s%s", bar5h, percentStyle.Render(fmt.Sprintf("%.0f%%", s.RateLimits.FiveHourPercent)), timeInfo), nil
}
