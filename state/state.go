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
}

type ModelInfo struct {
	Name     string
	PlanType string
}

type ContextInfo struct {
	UsedTokens  int
	TotalTokens int
	Percentage  float64
}

type RateLimitInfo struct {
	HourlyUsed    int
	HourlyTotal   int
	SevenDayUsed  int
	SevenDayTotal int
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
	Count    int
	LastUsed time.Time
	Duration time.Duration
}

type AgentInfo struct {
	ActiveAgent string
	TaskDesc    string
	ElapsedTime time.Duration
}

type TaskInfo struct {
	Pending    int
	InProgress int
	Completed  int
}

type SessionInfo struct {
	StartTime  time.Time
	Duration   time.Duration
	TokenSpeed float64
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
