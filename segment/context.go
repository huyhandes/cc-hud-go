package segment

import (
	"fmt"
	"strings"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/format"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
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

	// Build gradient progress bar
	bar := style.RenderGradientBar(percentage, 10)

	percentageStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(percentage))

	// Detailed token breakdown with semantic colors
	details := []string{}

	// Input tokens - Blue (incoming data)
	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	details = append(details,
		fmt.Sprintf("📥 %s", inStyle.Render(format.Tokens(s.Context.TotalInputTokens))),
	)

	// Output tokens - Emerald/Green (outgoing data)
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)
	details = append(details,
		fmt.Sprintf("📤 %s", outStyle.Render(format.Tokens(s.Context.TotalOutputTokens))),
	)

	// Cache stats if available - Different colors for Read vs Write
	if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
		cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
		cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)

		details = append(details,
			fmt.Sprintf("💾 %s%s%s",
				cacheReadStyle.Render("R:"+format.Tokens(s.Context.CacheReadTokens)),
				style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
				cacheWriteStyle.Render("W:"+format.Tokens(s.Context.CacheCreateTokens)),
			),
		)
	}

	// Total context size - Muted gray (static constant)
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorMuted)
	details = append(details,
		fmt.Sprintf("⚡ %s", totalStyle.Render(format.Tokens(s.Context.TotalTokens))),
	)

	// Single line format for use in custom layouts
	return fmt.Sprintf("%s %s %s",
		bar,
		percentageStyle.Render(fmt.Sprintf("%.0f%%", percentage)),
		strings.Join(details, " "),
	), nil
}
