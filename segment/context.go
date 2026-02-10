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

	// Choose icon based on thresholds
	var icon = "🟢"
	if percentage >= 90 {
		icon = "🔴"
	} else if percentage >= 70 {
		icon = "🟡"
	}

	// Build gradient progress bar
	bar := style.RenderGradientBar(percentage, 10)

	// Format tokens in thousands (k)
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	// Main display with bar and percentage
	percentageColor := style.ColorSuccess
	if percentage >= 90 {
		percentageColor = style.ColorDanger
	} else if percentage >= 70 {
		percentageColor = style.ColorWarning
	}

	percentageStyle := style.GetRenderer().NewStyle().Foreground(percentageColor)
	mainDisplay := fmt.Sprintf("%s %s %s",
		icon,
		bar,
		percentageStyle.Render(fmt.Sprintf("%.0f%%", percentage)),
	)

	// Detailed token breakdown with semantic colors
	details := []string{}

	// Input tokens - Blue (incoming data)
	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	details = append(details,
		fmt.Sprintf("📥 %s", inStyle.Render(formatTokens(s.Context.TotalInputTokens))),
	)

	// Output tokens - Emerald/Green (outgoing data)
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)
	details = append(details,
		fmt.Sprintf("📤 %s", outStyle.Render(formatTokens(s.Context.TotalOutputTokens))),
	)

	// Cache stats if available - Different colors for Read vs Write
	if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
		cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
		cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)

		details = append(details,
			fmt.Sprintf("💾 %s%s%s",
				cacheReadStyle.Render("R:"+formatTokens(s.Context.CacheReadTokens)),
				style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
				cacheWriteStyle.Render("W:"+formatTokens(s.Context.CacheCreateTokens)),
			),
		)
	}

	// Total context size - Muted gray (static constant)
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorMuted)
	details = append(details,
		fmt.Sprintf("⚡ %s", totalStyle.Render(formatTokens(s.Context.TotalTokens))),
	)

	return fmt.Sprintf("%s %s", mainDisplay, strings.Join(details, " ")), nil
}
