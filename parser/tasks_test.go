package parser

import (
	"testing"

	"github.com/huybui/cc-hud-go/state"
)

func TestParseTodoWrite(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Simulate TodoWrite tool call with bulk task write
	line := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TodoWrite",
				"input": {
					"todos": [
						{"content": "First task", "status": "completed"},
						{"content": "Second task", "status": "in_progress"},
						{"content": "Third task", "status": "pending"}
					]
				}
			}]
		}
	}`

	_ = ParseTranscriptLineWithTracker([]byte(line), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task counts
	if s.Tasks.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", s.Tasks.Completed)
	}
	if s.Tasks.InProgress != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", s.Tasks.InProgress)
	}
	if s.Tasks.Pending != 1 {
		t.Errorf("Expected 1 pending task, got %d", s.Tasks.Pending)
	}
}

func TestParseTaskCreate(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Simulate TaskCreate tool call
	line := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"id": "task-1",
				"name": "TaskCreate",
				"input": {
					"taskId": "alpha",
					"subject": "Build authentication",
					"description": "Implement user authentication system",
					"status": "pending"
				}
			}]
		}
	}`

	_ = ParseTranscriptLineWithTracker([]byte(line), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task was created as pending
	if s.Tasks.Pending != 1 {
		t.Errorf("Expected 1 pending task, got %d", s.Tasks.Pending)
	}
	if s.Tasks.InProgress != 0 {
		t.Errorf("Expected 0 in-progress tasks, got %d", s.Tasks.InProgress)
	}
	if s.Tasks.Completed != 0 {
		t.Errorf("Expected 0 completed tasks, got %d", s.Tasks.Completed)
	}
}

func TestParseTaskUpdate(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// First, create a task
	createLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"id": "task-1",
				"name": "TaskCreate",
				"input": {
					"taskId": "alpha",
					"subject": "Test task",
					"status": "pending"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(createLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify initial state
	if s.Tasks.Pending != 1 {
		t.Fatalf("Setup failed: expected 1 pending task, got %d", s.Tasks.Pending)
	}

	// Now update it to in_progress
	updateLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"id": "task-2",
				"name": "TaskUpdate",
				"input": {
					"taskId": "alpha",
					"status": "in_progress"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(updateLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task was updated
	if s.Tasks.Pending != 0 {
		t.Errorf("Expected 0 pending tasks, got %d", s.Tasks.Pending)
	}
	if s.Tasks.InProgress != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", s.Tasks.InProgress)
	}

	// Update to completed
	completeLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"id": "task-3",
				"name": "TaskUpdate",
				"input": {
					"taskId": "alpha",
					"status": "completed"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(completeLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task was completed
	if s.Tasks.InProgress != 0 {
		t.Errorf("Expected 0 in-progress tasks, got %d", s.Tasks.InProgress)
	}
	if s.Tasks.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", s.Tasks.Completed)
	}
}

func TestParseTaskUpdateByIndex(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Create tasks using TodoWrite
	line := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TodoWrite",
				"input": {
					"todos": [
						{"content": "Task 1", "status": "pending"},
						{"content": "Task 2", "status": "pending"}
					]
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(line), s, tracker)
	updateStateFromTasks(tracker, s)

	// Update by index (0-based)
	updateLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TaskUpdate",
				"input": {
					"taskId": 0,
					"status": "completed"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(updateLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify first task was completed
	if s.Tasks.Pending != 1 {
		t.Errorf("Expected 1 pending task, got %d", s.Tasks.Pending)
	}
	if s.Tasks.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", s.Tasks.Completed)
	}
}

func TestParseTaskStatusNormalization(t *testing.T) {
	tests := []struct {
		name           string
		inputStatus    string
		expectedField  string
		expectedCount  int
	}{
		{"completed", "completed", "Completed", 1},
		{"complete", "complete", "Completed", 1},
		{"done", "done", "Completed", 1},
		{"in_progress", "in_progress", "InProgress", 1},
		{"running", "running", "InProgress", 1},
		{"pending", "pending", "Pending", 1},
		{"not_started", "not_started", "Pending", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			tracker := &TaskTracker{
				Tasks:     []TaskItem{},
				TaskIDMap: make(map[string]int),
			}

			line := `{
				"type": "assistant",
				"message": {
					"content": [{
						"type": "tool_use",
						"name": "TaskCreate",
						"input": {
							"subject": "Test",
							"status": "` + tt.inputStatus + `"
						}
					}]
				}
			}`

			_ = ParseTranscriptLineWithTracker([]byte(line), s, tracker)
			updateStateFromTasks(tracker, s)

			// Check the appropriate field
			switch tt.expectedField {
			case "Completed":
				if s.Tasks.Completed != tt.expectedCount {
					t.Errorf("Status %q: expected %d completed, got %d",
						tt.inputStatus, tt.expectedCount, s.Tasks.Completed)
				}
			case "InProgress":
				if s.Tasks.InProgress != tt.expectedCount {
					t.Errorf("Status %q: expected %d in_progress, got %d",
						tt.inputStatus, tt.expectedCount, s.Tasks.InProgress)
				}
			case "Pending":
				if s.Tasks.Pending != tt.expectedCount {
					t.Errorf("Status %q: expected %d pending, got %d",
						tt.inputStatus, tt.expectedCount, s.Tasks.Pending)
				}
			}
		})
	}
}

func TestParseTaskDeletedStatus(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Create a task
	createLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"id": "task-1",
				"name": "TaskCreate",
				"input": {
					"taskId": "alpha",
					"subject": "Test task",
					"status": "pending"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(createLine), s, tracker)
	updateStateFromTasks(tracker, s)

	if s.Tasks.Pending != 1 {
		t.Fatalf("Setup failed: expected 1 pending task, got %d", s.Tasks.Pending)
	}

	// Delete it
	updateLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TaskUpdate",
				"input": {
					"taskId": "alpha",
					"status": "deleted"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(updateLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task was removed from counts
	if s.Tasks.Pending != 0 {
		t.Errorf("Expected 0 pending tasks after deletion, got %d", s.Tasks.Pending)
	}
	if s.Tasks.InProgress != 0 {
		t.Errorf("Expected 0 in-progress tasks after deletion, got %d", s.Tasks.InProgress)
	}
	if s.Tasks.Completed != 0 {
		t.Errorf("Expected 0 completed tasks after deletion, got %d", s.Tasks.Completed)
	}
}

func TestParseTranscriptWithMixedOperations(t *testing.T) {
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Simulate a realistic sequence of operations
	lines := []string{
		// 1. Bulk write initial tasks
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "TodoWrite",
					"input": {
						"todos": [
							{"content": "Task 1", "status": "completed"},
							{"content": "Task 2", "status": "pending"}
						]
					}
				}]
			}
		}`,
		// 2. Create a new task
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"id": "task-new",
					"name": "TaskCreate",
					"input": {
						"taskId": "beta",
						"subject": "Task 3"
					}
				}]
			}
		}`,
		// 3. Update task 2 to in_progress (by index)
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "TaskUpdate",
					"input": {
						"taskId": 1,
						"status": "in_progress"
					}
				}]
			}
		}`,
	}

	for _, line := range lines {
		_ = ParseTranscriptLineWithTracker([]byte(line), s, tracker)
	}
	updateStateFromTasks(tracker, s)

	// Verify final state: 1 completed, 1 in_progress, 1 pending
	if s.Tasks.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", s.Tasks.Completed)
	}
	if s.Tasks.InProgress != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", s.Tasks.InProgress)
	}
	if s.Tasks.Pending != 1 {
		t.Errorf("Expected 1 pending task, got %d", s.Tasks.Pending)
	}
}
