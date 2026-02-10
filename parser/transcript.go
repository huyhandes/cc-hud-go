package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/huyhandes/cc-hud-go/state"
)

// TranscriptLine represents a single line from the transcript JSONL
type TranscriptLine struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Message *MessageWrapper `json:"message"`
}

// MessageWrapper wraps the message content array
type MessageWrapper struct {
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents a single content block (tool_use, text, etc.)
type ContentBlock struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
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

	if line.Message != nil && len(line.Message.Content) > 0 {
		for _, block := range line.Message.Content {
			if block.Type == "tool_use" {
				if tracker != nil && (block.Name == "TodoWrite" || block.Name == "TaskCreate" || block.Name == "TaskUpdate") {
					processTaskTool(block, tracker)
				}

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
					if skillName, ok := block.Input["skill"].(string); ok && skillName != "" {
						usage := s.Tools.Skills[skillName]
						usage.Count++
						s.Tools.Skills[skillName] = usage
					} else {
						s.Tools.AppTools["Skill"]++
					}
				}
			}
		}
		return nil
	}

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

	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		_ = ParseTranscriptLineWithTracker(line, s, tracker)
	}

	updateStateFromTasks(tracker, s)

	return scanner.Err()
}
