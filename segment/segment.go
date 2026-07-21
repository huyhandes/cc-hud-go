package segment

import "github.com/huyhandes/cc-hud-go/state"

// Segment represents a displayable statusline segment.
// On/off is controlled by self-gating: Render returns "" when there is
// nothing to show.
type Segment interface {
	ID() string
	Render(s *state.State) (string, error)
}

// All returns all available segments in display order
func All() []Segment {
	return []Segment{
		&ModelSegment{},
		&CavemanSegment{},
		&PonytailSegment{},
		&GitSegment{},
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
