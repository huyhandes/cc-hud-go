package state

import "time"

// State holds all current session data
type State struct {
	Model      ModelInfo
	Context    ContextInfo
	RateLimits RateLimitInfo
	Git        GitInfo
	Tools      ToolsState
	Agents     AgentInfo
	Tasks      TaskInfo
	Session    SessionInfo
	Cost       CostInfo
}

type ModelInfo struct {
	Name     string
	PlanType string
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
	HourlyUsed       int
	HourlyTotal      int
	SevenDayUsed     int
	SevenDayTotal    int
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

type ToolsState struct {
	AppTools      map[string]int
	InternalTools map[string]int
	CustomTools   map[string]int
	MCPTools      map[MCPServer]map[string]int
	Skills        map[string]SkillUsage
}

type MCPServer struct {
	Name string
	Type string
}

type SkillUsage struct {
	Count int
}

type AgentInfo struct {
	ActiveAgent string
	TaskDesc    string
}

type Task struct {
	Subject string
	Status  string
}

type TaskInfo struct {
	Pending    int
	InProgress int
	Completed  int
	Details    []Task
}

type SessionInfo struct {
	ID             string
	TranscriptPath string
	StartTime      time.Time
	Duration       time.Duration
}

type CostInfo struct {
	TotalUSD      float64
	DurationMs    int64
	APIDurationMs int64
	LinesAdded    int
	LinesRemoved  int
}

// New creates a new State with initialized maps
func New() *State {
	return &State{
		Tools: ToolsState{
			AppTools:      make(map[string]int),
			InternalTools: make(map[string]int),
			CustomTools:   make(map[string]int),
			MCPTools:      make(map[MCPServer]map[string]int),
			Skills:        make(map[string]SkillUsage),
		},
		Session: SessionInfo{
			StartTime: time.Now(),
		},
	}
}

// UpdateDerived updates calculated fields like duration and percentage
func (s *State) UpdateDerived() {
	// Update session duration
	s.Session.Duration = time.Since(s.Session.StartTime)

	// Update context percentage
	if s.Context.TotalTokens > 0 {
		s.Context.Percentage = float64(s.Context.UsedTokens) / float64(s.Context.TotalTokens) * 100.0
	}
}
