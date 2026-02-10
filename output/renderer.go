package output

import (
	"fmt"
	"strings"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/segment"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
)

// Render generates plain text output for the statusline
// Returns plain text that Claude Code will display
func Render(s *state.State, cfg *config.Config) (string, error) {
	// Update derived fields before rendering
	s.UpdateDerived()

	// Check if multi-line layout is requested
	if cfg.LineLayout == "multiline" || cfg.LineLayout == "expanded" {
		return renderMultiLine(s, cfg)
	}

	// Single line layout (original)
	return renderSingleLine(s, cfg)
}

func renderSingleLine(s *state.State, cfg *config.Config) (string, error) {
	var parts []string

	// Render all segments
	for _, seg := range segment.All() {
		if !seg.Enabled(cfg) {
			continue
		}

		text, err := seg.Render(s, cfg)
		if err != nil {
			return "", err
		}

		if text == "" {
			continue
		}

		parts = append(parts, text)
	}

	return joinSegments(parts), nil
}

func renderMultiLine(s *state.State, cfg *config.Config) (string, error) {
	// Custom 4-line layout as requested
	var lines []string

	// Line 1: Model | Context | Token
	line1 := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if id == "model" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				line1 = append(line1, text)
			}
		}
	}
	// Add context bar inline
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		line1 = append(line1, renderContextBar(s))
	}
	if len(line1) > 0 {
		lines = append(lines, joinSegments(line1))
	}

	// Line 2: Input/Output tokens | Cost (without file changes)
	line2 := []string{}
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		line2 = append(line2, renderTokenDetails(s))
	}
	for _, seg := range segment.All() {
		if seg.ID() == "cost" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				line2 = append(line2, text)
			}
		}
	}
	if len(line2) > 0 {
		lines = append(lines, joinSegments(line2))
	}

	// Line 3: Git | File changes
	line3 := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if id == "git" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				line3 = append(line3, text)
			}
		}
	}
	// Add file changes to git line
	if s.Cost.LinesAdded > 0 || s.Cost.LinesRemoved > 0 {
		line3 = append(line3, renderFileChanges(s))
	}
	if len(line3) > 0 {
		lines = append(lines, joinSegments(line3))
	}

	// Line 4+: Each tool/task segment on its own line
	for _, seg := range segment.All() {
		id := seg.ID()
		if (id == "tools" || id == "tasks" || id == "agent" || id == "ratelimit") && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				lines = append(lines, text)
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}

// renderContextBar renders just the progress bar and percentage
func renderContextBar(s *state.State) string {
	percentage := s.Context.Percentage
	bar := renderGradientBar(percentage, 10)

	color := getColorForPercentage(percentage)
	percentageText := fmt.Sprintf("%.0f%%", percentage)

	return fmt.Sprintf("%s %s", bar, colorize(percentageText, color))
}

// renderTokenDetails renders token breakdown with colors
func renderTokenDetails(s *state.State) string {
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	details := []string{}

	// Input tokens - Blue
	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	details = append(details, fmt.Sprintf("📥 %s", inStyle.Render(formatTokens(s.Context.TotalInputTokens))))

	// Output tokens - Emerald/Green
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)
	details = append(details, fmt.Sprintf("📤 %s", outStyle.Render(formatTokens(s.Context.TotalOutputTokens))))

	// Cache stats if available
	if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
		cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
		cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)
		details = append(details, fmt.Sprintf("💾 %s%s%s",
			cacheReadStyle.Render("R:"+formatTokens(s.Context.CacheReadTokens)),
			style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
			cacheWriteStyle.Render("W:"+formatTokens(s.Context.CacheCreateTokens))))
	}

	// Total context size - Muted gray
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorMuted)
	details = append(details, fmt.Sprintf("⚡ %s", totalStyle.Render(formatTokens(s.Context.TotalTokens))))

	return strings.Join(details, " ")
}

// renderFileChanges renders file changes (lines added/removed)
func renderFileChanges(s *state.State) string {
	if s.Cost.LinesAdded == 0 && s.Cost.LinesRemoved == 0 {
		return ""
	}

	addStyle := style.GetRenderer().NewStyle().Foreground(style.ColorSuccess)
	removeStyle := style.GetRenderer().NewStyle().Foreground(style.ColorDanger)

	return fmt.Sprintf("📝 %s%s%s",
		addStyle.Render(fmt.Sprintf("+%d", s.Cost.LinesAdded)),
		style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
		removeStyle.Render(fmt.Sprintf("-%d", s.Cost.LinesRemoved)),
	)
}

// Helper functions for gradient bar rendering
func renderGradientBar(percentage float64, width int) string {
	if width <= 0 {
		width = 10
	}
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	filled := int(percentage / 100 * float64(width))
	if filled > width {
		filled = width
	}

	var bar strings.Builder
	for i := 0; i < width; i++ {
		if i < filled {
			char := getGradientChar(i, filled)
			bar.WriteString(char)
		} else {
			bar.WriteString("░")
		}
	}
	return bar.String()
}

func getGradientChar(position, filled int) string {
	if filled == 0 {
		return "░"
	}
	progress := float64(position) / float64(filled)
	if progress < 0.3 {
		return "█"
	} else if progress < 0.6 {
		return "▓"
	} else {
		return "▒"
	}
}

func getColorForPercentage(percentage float64) string {
	if percentage >= 90 {
		return "danger"
	} else if percentage >= 70 {
		return "warning"
	}
	return "success"
}

func colorize(text, colorName string) string {
	// Simple ANSI color codes (will be replaced by actual theme colors in real output)
	return text
}

// joinSegments joins segment outputs with two-space separators
func joinSegments(segments []string) string {
	// Filter out empty segments
	nonEmpty := make([]string, 0, len(segments))
	for _, seg := range segments {
		if strings.TrimSpace(seg) != "" {
			nonEmpty = append(nonEmpty, seg)
		}
	}

	return strings.Join(nonEmpty, "  │  ")
}
