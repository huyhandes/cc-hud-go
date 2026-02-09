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

	// Build enhanced progress bar
	barWidth := 10
	filled := int(percentage / 10)
	if filled > barWidth {
		filled = barWidth
	}

	filledBar := strings.Repeat("●", filled)
	emptyBar := strings.Repeat("○", barWidth-filled)
	bar := filledBar + emptyBar

	// Format tokens in thousands (k)
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	// Main display with bar and percentage
	mainDisplay := fmt.Sprintf("%s %s %s",
		icon,
		barStyle.Render(bar),
		style.ContextStyle.Render(fmt.Sprintf("%.0f%%", percentage)),
	)

	// Detailed token breakdown
	details := []string{}

	// Input/Output tokens
	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorSuccess)
	details = append(details,
		fmt.Sprintf("📥 %s", inStyle.Render(formatTokens(s.Context.TotalInputTokens))),
	)
	details = append(details,
		fmt.Sprintf("📤 %s", outStyle.Render(formatTokens(s.Context.TotalOutputTokens))),
	)

	// Cache stats if available
	if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
		cacheStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCyan)
		details = append(details,
			fmt.Sprintf("💾 %s", cacheStyle.Render(
				fmt.Sprintf("R:%s W:%s",
					formatTokens(s.Context.CacheReadTokens),
					formatTokens(s.Context.CacheCreateTokens),
				),
			)),
		)
	}

	// Total context size
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorMuted)
	details = append(details,
		fmt.Sprintf("⚡ %s", totalStyle.Render(formatTokens(s.Context.TotalTokens))),
	)

	return fmt.Sprintf("%s %s", mainDisplay, strings.Join(details, " ")), nil
}
