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

// All returns all available segments
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
