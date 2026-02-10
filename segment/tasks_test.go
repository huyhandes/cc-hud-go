package segment

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
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

	// Check for dashboard header
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
