package segment

import (
	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

// Segment represents a displayable statusline segment
type Segment interface {
	ID() string
	Render(s *state.State, cfg *config.Config) (string, error)
	Enabled(cfg *config.Config) bool
}

// All returns all available segments in display order
func All() []Segment {
	return []Segment{
		&ModelSegment{},
		&ContextSegment{},
		&GitSegment{},
		&CostSegment{},
		&ToolsSegment{},
		&TasksSegment{},
		&AgentSegment{},
		&FiveHourSegment{},
		&RateLimitSegment{},
	}
}

// ByID returns a map of segment ID to Segment for O(1) lookups
func ByID() map[string]Segment {
	m := make(map[string]Segment)
	for _, seg := range All() {
		m[seg.ID()] = seg
	}
	return m
}
