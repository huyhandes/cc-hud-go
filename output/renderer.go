package output

import (
	"encoding/json"
	"strings"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/segment"
	"github.com/huybui/cc-hud-go/state"
)

// SegmentOutput represents a rendered segment
type SegmentOutput struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// StatuslineOutput represents the complete statusline output
type StatuslineOutput struct {
	Segments []SegmentOutput `json:"segments"`
	Line     string          `json:"line"`
}

// Render generates JSON output for the statusline
func Render(s *state.State, cfg *config.Config) (string, error) {
	// Update derived fields before rendering
	s.UpdateDerived()

	var segments []SegmentOutput
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

		segments = append(segments, SegmentOutput{
			ID:   seg.ID(),
			Text: text,
		})
		parts = append(parts, text)
	}

	// Build output
	output := StatuslineOutput{
		Segments: segments,
		Line:     strings.Join(parts, " | "),
	}

	// Marshal to JSON
	data, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
