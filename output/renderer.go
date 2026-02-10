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

	// Line 1: Model Context Size | Context Bar | 5h Limit | 7d Limit
	line1 := []string{}

	// Combine Model and Context Size together
	modelAndContext := ""
	for _, seg := range segment.All() {
		id := seg.ID()
		if id == "model" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				modelAndContext = text
			}
		}
	}
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		if modelAndContext != "" {
			modelAndContext += " " + renderContextSize(s)
		} else {
			modelAndContext = renderContextSize(s)
		}
	}
	if modelAndContext != "" {
		line1 = append(line1, modelAndContext)
	}

	// Add context bar
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		line1 = append(line1, renderContextBar(s))
	}
	// Add 5h rate limit inline
	for _, seg := range segment.All() {
		if seg.ID() == "fivehour" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				line1 = append(line1, text)
			}
		}
	}
	// Add 7d rate limit inline
	for _, seg := range segment.All() {
		if seg.ID() == "ratelimit" && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				line1 = append(line1, text)
			}
		}
	}
	if len(line1) > 0 {
		lines = append(lines, joinSegments(line1))
	}

	// Line 2: Input/Output | Cache Read/Write | Cost | Time
	line2 := []string{}
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		line2 = append(line2, renderIOTokens(s))
		if s.Context.CacheReadTokens > 0 || s.Context.CacheCreateTokens > 0 {
			line2 = append(line2, renderCacheTokens(s))
		}
	}
	// Add cost and time
	if s.Cost.TotalUSD > 0 {
		line2 = append(line2, renderCost(s))
	}
	if s.Cost.DurationMs > 0 {
		line2 = append(line2, renderTime(s))
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
		if (id == "tools" || id == "tasks" || id == "agent") && seg.Enabled(cfg) {
			text, _ := seg.Render(s, cfg)
			if text != "" {
				lines = append(lines, text)
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}

// renderContextSize renders just the total context window size
func renderContextSize(s *state.State) string {
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	// Use Info color (cyan) for context size - represents static information
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
	return fmt.Sprintf("⚡ %s", totalStyle.Render(formatTokens(s.Context.TotalTokens)))
}

// renderContextBar renders just the progress bar and percentage
func renderContextBar(s *state.State) string {
	percentage := s.Context.Percentage

	// Use style package's gradient bar which has colors
	bar := style.RenderGradientBar(percentage, 10)

	// Color the percentage text based on threshold
	percentageColor := style.ColorSuccess
	if percentage >= 90 {
		percentageColor = style.ColorDanger
	} else if percentage >= 70 {
		percentageColor = style.ColorWarning
	}

	percentageStyle := style.GetRenderer().NewStyle().Foreground(percentageColor)
	percentageText := percentageStyle.Render(fmt.Sprintf("%.0f%%", percentage))

	return fmt.Sprintf("🧠 %s %s", bar, percentageText)
}

// renderTokenDetails renders token breakdown with colors (legacy - kept for single line mode)
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

// renderIOTokens renders input/output token counts
func renderIOTokens(s *state.State) string {
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)

	return fmt.Sprintf("📥 %s  📤 %s",
		inStyle.Render(formatTokens(s.Context.TotalInputTokens)),
		outStyle.Render(formatTokens(s.Context.TotalOutputTokens)))
}

// renderCacheTokens renders cache read/write token counts
func renderCacheTokens(s *state.State) string {
	formatTokens := func(tokens int) string {
		if tokens >= 1000 {
			return fmt.Sprintf("%dk", tokens/1000)
		}
		return fmt.Sprintf("%d", tokens)
	}

	cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
	cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)

	return fmt.Sprintf("💾 %s%s%s",
		cacheReadStyle.Render("R:"+formatTokens(s.Context.CacheReadTokens)),
		style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
		cacheWriteStyle.Render("W:"+formatTokens(s.Context.CacheCreateTokens)))
}

// renderCost renders the total cost
func renderCost(s *state.State) string {
	costStyle := style.GetRenderer().NewStyle().Foreground(style.ColorAccent).Bold(true)
	return costStyle.Render(fmt.Sprintf("💰$%.4f", s.Cost.TotalUSD))
}

// renderTime renders the session duration
func renderTime(s *state.State) string {
	durationSec := s.Cost.DurationMs / 1000
	mins := durationSec / 60
	secs := durationSec % 60

	durationStr := ""
	if mins > 0 {
		durationStr = fmt.Sprintf("⏱ %dm%ds", mins, secs)
	} else {
		durationStr = fmt.Sprintf("⏱ %ds", secs)
	}

	durationStyle := style.GetRenderer().NewStyle().Foreground(style.ColorHighlight)
	return durationStyle.Render(durationStr)
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
