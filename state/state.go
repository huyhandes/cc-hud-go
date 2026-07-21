package state

import "time"

// State holds all current session data
type State struct {
	Model      ModelInfo
	Context    ContextInfo
	RateLimits RateLimitInfo
	Git        GitInfo
	Session    SessionInfo
	Cost       CostInfo
}

type ModelInfo struct {
	Name string
}

type ContextInfo struct {
	UsedTokens         int
	TotalTokens        int
	Percentage         float64
	TotalInputTokens   int
	TotalOutputTokens  int
	CacheReadTokens    int
	CacheCreateTokens  int
	CurrentInputTokens int
}

type RateLimitInfo struct {
	FiveHourPercent  float64 // From OAuth API
	SevenDayPercent  float64 // From OAuth API
	FiveHourResetsAt string  // ISO 8601 timestamp
	SevenDayResetsAt string  // ISO 8601 timestamp
}

type GitInfo struct {
	Branch     string
	DirtyFiles int
	Ahead      int
	Behind     int
	Added      int
	Modified   int
	Deleted    int
}

type SessionInfo struct {
	StartTime time.Time
}

type CostInfo struct {
	TotalUSD     float64
	DurationMs   int64
	LinesAdded   int
	LinesRemoved int
}

// New creates a new State
func New() *State {
	return &State{
		Session: SessionInfo{
			StartTime: time.Now(),
		},
	}
}

// UpdateDerived updates calculated fields like percentage
func (s *State) UpdateDerived() {
	// Update context percentage
	if s.Context.TotalTokens > 0 {
		s.Context.Percentage = float64(s.Context.UsedTokens) / float64(s.Context.TotalTokens) * 100.0
	}
}
