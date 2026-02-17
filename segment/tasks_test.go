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

	if !strings.Contains(output, "⏳") {
		t.Errorf("expected pending icon in output, got %q", output)
	}
	if !strings.Contains(output, "🔄") {
		t.Errorf("expected in-progress icon in output, got %q", output)
	}
	if !strings.Contains(output, "✅") {
		t.Errorf("expected completed icon in output, got %q", output)
	}
}

func TestTasksSegmentZeroTotal(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &TasksSegment{}
	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output for zero tasks, got %q", output)
	}
}

func TestTasksSegmentSkipsZeroCounts(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Tasks.Completed = 3

	seg := &TasksSegment{}
	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if strings.Contains(output, "⏳") {
		t.Errorf("expected no pending icon when count is 0, got %q", output)
	}
	if strings.Contains(output, "🔄") {
		t.Errorf("expected no in-progress icon when count is 0, got %q", output)
	}
	if !strings.Contains(output, "✅") {
		t.Errorf("expected completed icon in output, got %q", output)
	}
	if !strings.Contains(output, "3") {
		t.Errorf("expected count 3 in output, got %q", output)
	}
}
