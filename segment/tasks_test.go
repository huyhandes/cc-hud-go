package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestTasksSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Tasks.Completed = 2
	s.Tasks.InProgress = 1
	s.Tasks.Pending = 2

	seg := &TasksSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// With 5 total tasks and default threshold of 3, should show inline dashboard
	if !strings.Contains(output, "Tasks Dashboard") {
		t.Errorf("expected 'Tasks Dashboard' in output, got '%s'", output)
	}

	// Check for task counts
	if !strings.Contains(output, "Todo") && !strings.Contains(output, "2") {
		t.Errorf("expected pending tasks in output, got '%s'", output)
	}

	if !strings.Contains(output, "Completed") && !strings.Contains(output, "2") {
		t.Errorf("expected completed count in output, got '%s'", output)
	}
}

func TestTasksSegmentTableThreshold(t *testing.T) {
	// Below threshold - should be inline dashboard
	s := state.New()
	s.Tasks.Pending = 1
	s.Tasks.InProgress = 1
	s.Tasks.Completed = 1

	cfg := config.Default()
	cfg.Tables.TasksThreshold = 5

	seg := &TasksSegment{}
	result, err := seg.Render(s, cfg)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be inline dashboard (no table borders)
	if strings.Contains(result, "┌") {
		t.Error("Expected inline dashboard below threshold, got table")
	}

	if !strings.Contains(result, "Tasks Dashboard") {
		t.Error("Expected dashboard format")
	}

	// Above threshold with task details - should be table
	s.Tasks.Pending = 2
	s.Tasks.InProgress = 1
	s.Tasks.Completed = 5
	s.Tasks.Details = []state.Task{
		{Subject: "Task 1", Status: "pending"},
		{Subject: "Task 2", Status: "pending"},
		{Subject: "Task 3", Status: "in_progress"},
		{Subject: "Task 4", Status: "completed"},
		{Subject: "Task 5", Status: "completed"},
		{Subject: "Task 6", Status: "completed"},
		{Subject: "Task 7", Status: "completed"},
		{Subject: "Task 8 with a very long name that should be truncated", Status: "completed"},
	}

	cfg.Tables.TasksThreshold = 3

	result, err = seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should be table format
	if !strings.Contains(result, "┌") {
		t.Error("Expected table format above threshold")
	}

	// Should contain task subjects
	if !strings.Contains(result, "Task 1") {
		t.Error("Expected pending task in table")
	}

	// Should show last 3 completed (tasks 6, 7, 8)
	if !strings.Contains(result, "Task 6") || !strings.Contains(result, "Task 7") {
		t.Error("Expected last 3 completed tasks in table")
	}

	// Should not show earlier completed tasks (task 4, 5)
	if strings.Contains(result, "Task 4") || strings.Contains(result, "Task 5") {
		t.Error("Should only show last 3 completed tasks")
	}
}
