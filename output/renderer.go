package output

import (
	"fmt"
	"strings"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/format"
	"github.com/huyhandes/cc-hud-go/segment"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
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
	var lines []string
	segs := segment.ByID()

	renderSeg := func(id string) string {
		seg, ok := segs[id]
		if !ok || !seg.Enabled(cfg) {
			return ""
		}
		text, _ := seg.Render(s, cfg)
		return text
	}

	// Line 1: Model Context Size | Context Bar | 5h Limit | 7d Limit
	line1 := []string{}

	modelAndContext := renderSeg("model")
	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		ctxSize := renderContextSize(s)
		if modelAndContext != "" {
			modelAndContext += " " + ctxSize
		} else {
			modelAndContext = ctxSize
		}
	}
	if modelAndContext != "" {
		line1 = append(line1, modelAndContext)
	}

	if cfg.Display.Context && s.Context.TotalTokens > 0 {
		line1 = append(line1, renderContextBar(s))
	}
	if text := renderSeg("fivehour"); text != "" {
		line1 = append(line1, text)
	}
	if text := renderSeg("ratelimit"); text != "" {
		line1 = append(line1, text)
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
	if text := renderSeg("git"); text != "" {
		line3 = append(line3, text)
	}
	if s.Cost.LinesAdded > 0 || s.Cost.LinesRemoved > 0 {
		line3 = append(line3, renderFileChanges(s))
	}
	if len(line3) > 0 {
		lines = append(lines, joinSegments(line3))
	}

	// Line 4+: Each tool/task segment on its own line
	for _, id := range []string{"tools", "tasks", "agent"} {
		if text := renderSeg(id); text != "" {
			lines = append(lines, text)
		}
	}

	return strings.Join(lines, "\n"), nil
}

// renderContextSize renders just the total context window size
func renderContextSize(s *state.State) string {
	totalStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
	return fmt.Sprintf("⚡ %s", totalStyle.Render(format.Tokens(s.Context.TotalTokens)))
}

// renderContextBar renders just the progress bar and percentage
func renderContextBar(s *state.State) string {
	percentage := s.Context.Percentage

	// Use style package's gradient bar which has colors
	bar := style.RenderGradientBar(percentage, 10)

	percentageStyle := style.GetRenderer().NewStyle().Foreground(style.ThresholdColor(percentage))
	percentageText := percentageStyle.Render(fmt.Sprintf("%.0f%%", percentage))

	return fmt.Sprintf("🧠 %s %s", bar, percentageText)
}

// renderIOTokens renders input/output token counts
func renderIOTokens(s *state.State) string {
	inStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInput)
	outStyle := style.GetRenderer().NewStyle().Foreground(style.ColorOutput)

	return fmt.Sprintf("📥 %s  📤 %s",
		inStyle.Render(format.Tokens(s.Context.TotalInputTokens)),
		outStyle.Render(format.Tokens(s.Context.TotalOutputTokens)))
}

// renderCacheTokens renders cache read/write token counts
func renderCacheTokens(s *state.State) string {
	cacheReadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheRead)
	cacheWriteStyle := style.GetRenderer().NewStyle().Foreground(style.ColorCacheWrite)

	return fmt.Sprintf("💾 %s%s%s",
		cacheReadStyle.Render("R:"+format.Tokens(s.Context.CacheReadTokens)),
		style.GetRenderer().NewStyle().Foreground(style.ColorMuted).Render("/"),
		cacheWriteStyle.Render("W:"+format.Tokens(s.Context.CacheCreateTokens)))
}

// renderCost renders the total cost
func renderCost(s *state.State) string {
	costStyle := style.GetRenderer().NewStyle().Foreground(style.ColorAccent).Bold(true)
	return costStyle.Render("💰" + format.Cost(s.Cost.TotalUSD))
}

// renderTime renders the session duration
func renderTime(s *state.State) string {
	durationStyle := style.GetRenderer().NewStyle().Foreground(style.ColorHighlight)
	return durationStyle.Render("⏱ " + format.Duration(s.Cost.DurationMs))
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
