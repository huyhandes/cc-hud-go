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

	if !strings.Contains(output, "2/5") {
		t.Errorf("expected progress ratio in output, got '%s'", output)
	}
}
