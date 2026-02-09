package segment

import (
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

// Segment represents a displayable statusline segment
type Segment interface {
	ID() string
	Render(s *state.State, cfg *config.Config) (string, error)
	Enabled(cfg *config.Config) bool
}

type dummySegment struct{}

func (d dummySegment) ID() string { return "dummy" }
func (d dummySegment) Enabled(cfg *config.Config) bool { return true }
func (d dummySegment) Render(s *state.State, cfg *config.Config) (string, error) {
	return "dummy", nil
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
		&RateLimitSegment{},
	}
}
