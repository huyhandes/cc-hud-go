package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/huybui/cc-hud-go/state"
)

type ToolCategory int

const (
	CategoryApp ToolCategory = iota
	CategoryInternal
	CategoryCustom
	CategoryMCP
	CategorySkill
)

var appTools = map[string]bool{
	"read":      true,
	"write":     true,
	"edit":      true,
	"bash":      true,
	"glob":      true,
	"grep":      true,
	"task":      true,
	"webfetch":  true,
	"websearch": true,
}

// StdinData represents the JSON structure from Claude Code statusline API
type StdinData struct {
	SessionID      string `json:"session_id"`
	CWD            string `json:"cwd"`
	TranscriptPath string `json:"transcript_path"`
	Version        string `json:"version"`
	Model          struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`
	ContextWindow struct {
		TotalInputTokens    int     `json:"total_input_tokens"`
		TotalOutputTokens   int     `json:"total_output_tokens"`
		ContextWindowSize   int     `json:"context_window_size"`
		UsedPercentage      float64 `json:"used_percentage"`
		RemainingPercentage float64 `json:"remaining_percentage"`
		CurrentUsage        *struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"current_usage"`
	} `json:"context_window"`
	Cost *struct {
		TotalCostUSD       float64 `json:"total_cost_usd"`
		TotalDurationMs    int64   `json:"total_duration_ms"`
		TotalAPIDurationMs int64   `json:"total_api_duration_ms"`
		TotalLinesAdded    int     `json:"total_lines_added"`
		TotalLinesRemoved  int     `json:"total_lines_removed"`
	} `json:"cost,omitempty"`
	Exceeds200KTokens bool `json:"exceeds_200k_tokens"`
	OutputStyle       *struct {
		Name string `json:"name"`
	} `json:"output_style,omitempty"`
	Vim *struct {
		Mode string `json:"mode"`
	} `json:"vim,omitempty"`
	Agent *struct {
		Name string `json:"name"`
	} `json:"agent,omitempty"`
}

// ParseStdin parses stdin JSON from Claude Code and updates state
func ParseStdin(data []byte, s *state.State) error {
	var stdin StdinData
	if err := json.Unmarshal(data, &stdin); err != nil {
		return err
	}

	// Update session info
	s.Session.ID = stdin.SessionID
	s.Session.TranscriptPath = stdin.TranscriptPath

	// Update model info
	s.Model.Name = stdin.Model.DisplayName
	if s.Model.Name == "" {
		s.Model.Name = stdin.Model.ID
	}
	// Infer plan type from model ID (Pro/Max/Team indicators not in API)
	s.Model.PlanType = ""

	// Update context - use current usage if available, otherwise use totals
	s.Context.TotalInputTokens = stdin.ContextWindow.TotalInputTokens
	s.Context.TotalOutputTokens = stdin.ContextWindow.TotalOutputTokens
	s.Context.TotalTokens = stdin.ContextWindow.ContextWindowSize

	if stdin.ContextWindow.CurrentUsage != nil {
		// Calculate used tokens from current usage (input only, as per docs)
		s.Context.CurrentInputTokens = stdin.ContextWindow.CurrentUsage.InputTokens
		s.Context.CacheCreateTokens = stdin.ContextWindow.CurrentUsage.CacheCreationInputTokens
		s.Context.CacheReadTokens = stdin.ContextWindow.CurrentUsage.CacheReadInputTokens

		usedTokens := s.Context.CurrentInputTokens + s.Context.CacheCreateTokens + s.Context.CacheReadTokens
		s.Context.UsedTokens = usedTokens
	} else {
		// Fallback to total tokens
		s.Context.UsedTokens = stdin.ContextWindow.TotalInputTokens
	}

	// Update agent info if present
	if stdin.Agent != nil {
		s.Agents.ActiveAgent = stdin.Agent.Name
	}

	// Update cost info if present
	if stdin.Cost != nil {
		s.Cost.TotalUSD = stdin.Cost.TotalCostUSD
		s.Cost.DurationMs = stdin.Cost.TotalDurationMs
		s.Cost.APIDurationMs = stdin.Cost.TotalAPIDurationMs
		s.Cost.LinesAdded = stdin.Cost.TotalLinesAdded
		s.Cost.LinesRemoved = stdin.Cost.TotalLinesRemoved
	}

	// Update rate limits - not provided in API, keep existing values
	// Rate limits data is not in the Claude Code API

	return nil
}

// CategorizeTool determines the category of a tool by name
func CategorizeTool(name string) ToolCategory {
	lower := strings.ToLower(name)

	// Check for MCP pattern
	if strings.HasPrefix(lower, "mcp__") {
		return CategoryMCP
	}

	// Check for Skill
	if lower == "skill" {
		return CategorySkill
	}

	// Check for internal (Bash is special) - must check before appTools
	if lower == "bash" {
		return CategoryInternal
	}

	// Check for app tools
	if appTools[lower] {
		return CategoryApp
	}

	// Everything else is custom
	return CategoryCustom
}

type TranscriptLine struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"` // For backward compatibility with simple format
	Message *MessageWrapper `json:"message"`
}

type MessageWrapper struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// TaskItem represents a tracked task from TodoWrite/TaskCreate/TaskUpdate
type TaskItem struct {
	ID      string // Task ID from TaskCreate/TaskUpdate
	Content string // Subject or description
	Status  string // pending, in_progress, completed, deleted
}

// TaskTracker maintains task state during transcript parsing
type TaskTracker struct {
	Tasks     []TaskItem
	TaskIDMap map[string]int // Maps task ID to index in Tasks slice
}

// normalizeTaskStatus normalizes various status strings to standard values
func normalizeTaskStatus(status string) string {
	switch strings.ToLower(status) {
	case "completed", "complete", "done":
		return "completed"
	case "in_progress", "running":
		return "in_progress"
	case "pending", "not_started", "":
		return "pending"
	case "deleted":
		return "deleted"
	default:
		return "pending"
	}
}

// processTaskTool processes task-related tool calls (TodoWrite, TaskCreate, TaskUpdate)
func processTaskTool(block ContentBlock, tracker *TaskTracker) {
	switch block.Name {
	case "TodoWrite":
		// Bulk replace all tasks
		todos, ok := block.Input["todos"].([]interface{})
		if !ok {
			return
		}

		// Clear existing tasks
		tracker.Tasks = nil
		tracker.TaskIDMap = make(map[string]int)

		// Add all tasks from TodoWrite
		for _, t := range todos {
			todo, ok := t.(map[string]interface{})
			if !ok {
				continue
			}

			content := ""
			if c, ok := todo["content"].(string); ok {
				content = c
			}

			status := "pending"
			if s, ok := todo["status"].(string); ok {
				status = normalizeTaskStatus(s)
			}

			if status != "deleted" {
				tracker.Tasks = append(tracker.Tasks, TaskItem{
					Content: content,
					Status:  status,
				})
			}
		}

	case "TaskCreate":
		// Create a new task
		subject := ""
		if s, ok := block.Input["subject"].(string); ok {
			subject = s
		}
		description := ""
		if d, ok := block.Input["description"].(string); ok {
			description = d
		}

		content := subject
		if content == "" {
			content = description
		}
		if content == "" {
			content = "Untitled task"
		}

		status := "pending"
		if s, ok := block.Input["status"].(string); ok {
			status = normalizeTaskStatus(s)
		}

		// Add task
		index := len(tracker.Tasks)
		tracker.Tasks = append(tracker.Tasks, TaskItem{
			Content: content,
			Status:  status,
		})

		// Map task ID to index for later updates
		if taskID, ok := block.Input["taskId"].(string); ok {
			tracker.TaskIDMap[taskID] = index
		} else if taskID, ok := block.Input["taskId"].(float64); ok {
			tracker.TaskIDMap[fmt.Sprintf("%.0f", taskID)] = index
		}
		// Also map by block ID as fallback
		if block.ID != "" {
			tracker.TaskIDMap[block.ID] = index
		}

	case "TaskUpdate":
		// Update an existing task
		taskID := resolveTaskID(block.Input["taskId"])
		if taskID == "" {
			return
		}

		// Find task by ID or index
		index := -1
		if idx, ok := tracker.TaskIDMap[taskID]; ok {
			index = idx
		} else if idxNum, err := strconv.Atoi(taskID); err == nil && idxNum >= 0 && idxNum < len(tracker.Tasks) {
			index = idxNum
		}

		if index < 0 || index >= len(tracker.Tasks) {
			return
		}

		// Update status if provided
		if s, ok := block.Input["status"].(string); ok {
			newStatus := normalizeTaskStatus(s)

			// Handle deletion by removing from list
			if newStatus == "deleted" {
				// Remove task by index
				tracker.Tasks = append(tracker.Tasks[:index], tracker.Tasks[index+1:]...)
				// Rebuild ID map
				newMap := make(map[string]int)
				for id, idx := range tracker.TaskIDMap {
					if idx < index {
						newMap[id] = idx
					} else if idx > index {
						newMap[id] = idx - 1
					}
				}
				tracker.TaskIDMap = newMap
				return
			}

			tracker.Tasks[index].Status = newStatus
		}

		// Update content if provided
		if s, ok := block.Input["subject"].(string); ok && s != "" {
			tracker.Tasks[index].Content = s
		} else if d, ok := block.Input["description"].(string); ok && d != "" {
			tracker.Tasks[index].Content = d
		}
	}
}

// resolveTaskID extracts task ID from various input types
func resolveTaskID(taskID interface{}) string {
	switch v := taskID.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

// updateStateFromTasks updates state task counts from tracker
func updateStateFromTasks(tracker *TaskTracker, s *state.State) {
	s.Tasks.Pending = 0
	s.Tasks.InProgress = 0
	s.Tasks.Completed = 0
	s.Tasks.Details = make([]state.Task, 0, len(tracker.Tasks))

	for _, task := range tracker.Tasks {
		// Add to details list
		s.Tasks.Details = append(s.Tasks.Details, state.Task{
			Subject: task.Content,
			Status:  task.Status,
		})

		// Update counts
		switch task.Status {
		case "pending":
			s.Tasks.Pending++
		case "in_progress":
			s.Tasks.InProgress++
		case "completed":
			s.Tasks.Completed++
		}
	}
}

// ParseTranscriptLine parses a single JSONL line and updates state
func ParseTranscriptLine(data []byte, s *state.State) error {
	return ParseTranscriptLineWithTracker(data, s, nil)
}

// ParseTranscriptLineWithTracker parses a single JSONL line with task tracking
func ParseTranscriptLineWithTracker(data []byte, s *state.State, tracker *TaskTracker) error {
	var line TranscriptLine
	if err := json.Unmarshal(data, &line); err != nil {
		return err
	}

	// Handle nested message structure (assistant messages with tool calls)
	if line.Message != nil && len(line.Message.Content) > 0 {
		for _, block := range line.Message.Content {
			if block.Type == "tool_use" {
				// Process task-related tools if tracker is provided
				if tracker != nil && (block.Name == "TodoWrite" || block.Name == "TaskCreate" || block.Name == "TaskUpdate") {
					processTaskTool(block, tracker)
				}

				// Process regular tool tracking
				category := CategorizeTool(block.Name)
				switch category {
				case CategoryApp:
					s.Tools.AppTools[block.Name]++
				case CategoryInternal:
					s.Tools.InternalTools[block.Name]++
				case CategoryCustom:
					s.Tools.CustomTools[block.Name]++
				case CategoryMCP:
					parts := strings.Split(block.Name, "__")
					if len(parts) >= 3 {
						server := state.MCPServer{
							Name: parts[1],
							Type: "mcp",
						}
						if s.Tools.MCPTools[server] == nil {
							s.Tools.MCPTools[server] = make(map[string]int)
						}
						toolName := strings.Join(parts[2:], "__")
						s.Tools.MCPTools[server][toolName]++
					}
				case CategorySkill:
					s.Tools.AppTools["Skill"]++
				}
			}
		}
		return nil
	}

	// Handle simple format (backward compatibility for old tool_use events)
	if line.Type != "tool_use" {
		return nil
	}

	category := CategorizeTool(line.Name)

	switch category {
	case CategoryApp:
		s.Tools.AppTools[line.Name]++

	case CategoryInternal:
		s.Tools.InternalTools[line.Name]++

	case CategoryCustom:
		s.Tools.CustomTools[line.Name]++

	case CategoryMCP:
		// Parse MCP tool name: mcp__<server>__<tool>
		parts := strings.Split(line.Name, "__")
		if len(parts) >= 3 {
			server := state.MCPServer{
				Name: parts[1],
				Type: "mcp",
			}

			if s.Tools.MCPTools[server] == nil {
				s.Tools.MCPTools[server] = make(map[string]int)
			}

			toolName := strings.Join(parts[2:], "__")
			s.Tools.MCPTools[server][toolName]++
		}

	case CategorySkill:
		// Skills need additional parsing from the tool parameters
		// For now, just count as app tool
		s.Tools.AppTools["Skill"]++
	}

	return nil
}

// ParseTranscript reads and parses the entire transcript file
func ParseTranscript(path string, s *state.State) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Initialize task tracker
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer size for large transcript lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // Max 1MB per line

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		// Ignore errors from individual lines, just continue
		_ = ParseTranscriptLineWithTracker(line, s, tracker)
	}

	// Update state with final task counts
	updateStateFromTasks(tracker, s)

	return scanner.Err()
}
