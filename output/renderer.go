package output

import (
	"strings"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/segment"
	"github.com/huybui/cc-hud-go/state"
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

	// Render each segment type on its own line for better readability
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

		// Add the segment output
		lines = append(lines, text)

		// Add blank line after major sections for visual grouping
		id := seg.ID()
		if id == "context" || id == "cost" || id == "tasks" {
			lines = append(lines, "")
		}
	}

	// Remove trailing blank lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
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
