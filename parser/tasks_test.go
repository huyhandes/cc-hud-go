package parser

import (
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
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
		name          string
		inputStatus   string
		expectedField string
		expectedCount int
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
		// 3. Update task 2 to in_progress (by 1-based task ID)
		`{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "TaskUpdate",
					"input": {
						"taskId": 2,
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

func TestParseTaskUpdate1BasedTaskID(t *testing.T) {
	// Test the fix for 1-based task ID indexing bug
	// Previously: TaskUpdate with taskId="1" would update index 1 (task 2)
	// After fix: TaskUpdate with taskId="1" should update index 0 (task 1)
	s := state.New()
	tracker := &TaskTracker{
		Tasks:     []TaskItem{},
		TaskIDMap: make(map[string]int),
	}

	// Create 3 tasks (indices 0, 1, 2)
	for i := 1; i <= 3; i++ {
		createLine := `{
			"type": "assistant",
			"message": {
				"content": [{
					"type": "tool_use",
					"name": "TaskCreate",
					"input": {
						"subject": "Task ` + string(rune('0'+i)) + `"
					}
				}]
			}
		}`
		_ = ParseTranscriptLineWithTracker([]byte(createLine), s, tracker)
	}
	updateStateFromTasks(tracker, s)

	// Verify all 3 tasks are pending
	if s.Tasks.Pending != 3 {
		t.Fatalf("Setup failed: expected 3 pending tasks, got %d", s.Tasks.Pending)
	}

	// Update task with taskId="1" (should update first task at index 0)
	updateLine := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TaskUpdate",
				"input": {
					"taskId": "1",
					"status": "completed"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(updateLine), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task at index 0 (task 1) was completed
	if s.Tasks.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", s.Tasks.Completed)
	}
	if s.Tasks.Pending != 2 {
		t.Errorf("Expected 2 pending tasks, got %d", s.Tasks.Pending)
	}
	if len(s.Tasks.Details) != 3 {
		t.Fatalf("Expected 3 total tasks, got %d", len(s.Tasks.Details))
	}
	if s.Tasks.Details[0].Status != "completed" {
		t.Errorf("Expected first task (index 0) to be completed, got %s", s.Tasks.Details[0].Status)
	}
	if s.Tasks.Details[1].Status != "pending" {
		t.Errorf("Expected second task (index 1) to be pending, got %s", s.Tasks.Details[1].Status)
	}

	// Update task with taskId="3" (should update third task at index 2)
	updateLine3 := `{
		"type": "assistant",
		"message": {
			"content": [{
				"type": "tool_use",
				"name": "TaskUpdate",
				"input": {
					"taskId": "3",
					"status": "completed"
				}
			}]
		}
	}`
	_ = ParseTranscriptLineWithTracker([]byte(updateLine3), s, tracker)
	updateStateFromTasks(tracker, s)

	// Verify task at index 2 (task 3) was completed
	if s.Tasks.Completed != 2 {
		t.Errorf("Expected 2 completed tasks, got %d", s.Tasks.Completed)
	}
	if s.Tasks.Pending != 1 {
		t.Errorf("Expected 1 pending task, got %d", s.Tasks.Pending)
	}
	if s.Tasks.Details[2].Status != "completed" {
		t.Errorf("Expected third task (index 2) to be completed, got %s", s.Tasks.Details[2].Status)
	}
}
