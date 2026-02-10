package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/huyhandes/cc-hud-go/state"
)

// TaskItem represents a tracked task from TodoWrite/TaskCreate/TaskUpdate
type TaskItem struct {
	ID      string
	Content string
	Status  string
}

// TaskTracker maintains task state during transcript parsing
type TaskTracker struct {
	Tasks     []TaskItem
	TaskIDMap map[string]int
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
		todos, ok := block.Input["todos"].([]interface{})
		if !ok {
			return
		}

		tracker.Tasks = nil
		tracker.TaskIDMap = make(map[string]int)

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

		index := len(tracker.Tasks)
		tracker.Tasks = append(tracker.Tasks, TaskItem{
			Content: content,
			Status:  status,
		})

		if taskID, ok := block.Input["taskId"].(string); ok {
			tracker.TaskIDMap[taskID] = index
		} else if taskID, ok := block.Input["taskId"].(float64); ok {
			tracker.TaskIDMap[fmt.Sprintf("%.0f", taskID)] = index
		}
		if block.ID != "" {
			tracker.TaskIDMap[block.ID] = index
		}

	case "TaskUpdate":
		taskID := resolveTaskID(block.Input["taskId"])
		if taskID == "" {
			return
		}

		index := -1
		if idx, ok := tracker.TaskIDMap[taskID]; ok {
			index = idx
		} else if idxNum, err := strconv.Atoi(taskID); err == nil {
			if idxNum == 0 && idxNum < len(tracker.Tasks) {
				index = 0
			} else if idxNum >= 1 && idxNum <= len(tracker.Tasks) {
				index = idxNum - 1
			}
		}

		if index < 0 || index >= len(tracker.Tasks) {
			return
		}

		if s, ok := block.Input["status"].(string); ok {
			newStatus := normalizeTaskStatus(s)

			if newStatus == "deleted" {
				tracker.Tasks = append(tracker.Tasks[:index], tracker.Tasks[index+1:]...)
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
		s.Tasks.Details = append(s.Tasks.Details, state.Task{
			Subject: task.Content,
			Status:  task.Status,
		})

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
