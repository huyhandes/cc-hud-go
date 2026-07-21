package segment

import (
	"fmt"
	"time"

	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type RateLimitSegment struct{}

func (r *RateLimitSegment) ID() string {
	return "ratelimit"
}

func (r *RateLimitSegment) Render(s *state.State) string {
	// 7-day limit only. 5-hour limit is rendered by FiveHourSegment.
	// OAuth API is the single source of truth; render empty on failure.
	if s.RateLimits.SevenDayPercent <= 0 {
		return ""
	}

	bar7d := style.RenderGradientBar(s.RateLimits.SevenDayPercent, 10)
	percentStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(s.RateLimits.SevenDayPercent))
	return fmt.Sprintf("📊 %s %s", bar7d, percentStyle.Render(fmt.Sprintf("%.0f%%", s.RateLimits.SevenDayPercent)))
}

// FiveHourSegment displays 5-hour rate limit with elapsed time
type FiveHourSegment struct{}

func (f *FiveHourSegment) ID() string {
	return "fivehour"
}

func (f *FiveHourSegment) Render(s *state.State) string {
	// Only render if OAuth data available
	if s.RateLimits.FiveHourPercent <= 0 {
		return ""
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

	return fmt.Sprintf("⏱️ %s %s%s", bar5h, percentStyle.Render(fmt.Sprintf("%.0f%%", s.RateLimits.FiveHourPercent)), timeInfo)
}
