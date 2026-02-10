package output

import (
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

	separator := style.Separator()
	return strings.Join(parts, " "+separator+" "), nil
}

func renderMultiLine(s *state.State, cfg *config.Config) (string, error) {
	separator := style.Separator()

	// Line 1: Model and Context (most important)
	line1Parts := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if (id == "model" || id == "context") && seg.Enabled(cfg) {
			text, err := seg.Render(s, cfg)
			if err != nil {
				return "", err
			}
			if text != "" {
				line1Parts = append(line1Parts, text)
			}
		}
	}

	// Line 2: Git, Cost, and operational metrics
	line2Parts := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if (id == "git" || id == "cost") && seg.Enabled(cfg) {
			text, err := seg.Render(s, cfg)
			if err != nil {
				return "", err
			}
			if text != "" {
				line2Parts = append(line2Parts, text)
			}
		}
	}

	// Line 3: Tools, Agent, and other info (excluding tasks)
	line3Parts := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if (id == "tools" || id == "agent" || id == "ratelimit") && seg.Enabled(cfg) {
			text, err := seg.Render(s, cfg)
			if err != nil {
				return "", err
			}
			if text != "" {
				line3Parts = append(line3Parts, text)
			}
		}
	}

	// Line 4: Tasks (dedicated line for dashboard display)
	line4Parts := []string{}
	for _, seg := range segment.All() {
		id := seg.ID()
		if id == "tasks" && seg.Enabled(cfg) {
			text, err := seg.Render(s, cfg)
			if err != nil {
				return "", err
			}
			if text != "" {
				line4Parts = append(line4Parts, text)
			}
		}
	}

	// Build output
	var lines []string
	if len(line1Parts) > 0 {
		lines = append(lines, strings.Join(line1Parts, " "+separator+" "))
	}
	if len(line2Parts) > 0 {
		lines = append(lines, strings.Join(line2Parts, " "+separator+" "))
	}
	if len(line3Parts) > 0 {
		lines = append(lines, strings.Join(line3Parts, " "+separator+" "))
	}
	if len(line4Parts) > 0 {
		lines = append(lines, strings.Join(line4Parts, " "+separator+" "))
	}

	return strings.Join(lines, "\n"), nil
}
