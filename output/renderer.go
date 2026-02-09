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

	var parts []string

	// Render all segments
	for _, seg := range segment.All() {
		// Skip disabled segments
		if !seg.Enabled(cfg) {
			continue
		}

		// Render segment
		text, err := seg.Render(s, cfg)
		if err != nil {
			return "", err
		}

		// Skip empty segments
		if text == "" {
			continue
		}

		parts = append(parts, text)
	}

	// Join segments with separator
	return strings.Join(parts, " | "), nil
}
