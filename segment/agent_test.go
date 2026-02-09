package segment

import (
	"strings"
	"testing"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

func TestAgentSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Agents.ActiveAgent = "test-agent"
	s.Agents.TaskDesc = "exploring codebase"

	seg := &AgentSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "Agent:") {
		t.Errorf("expected 'Agent:' prefix in output, got '%s'", output)
	}

	if !strings.Contains(output, "test-agent") {
		t.Errorf("expected agent name in output, got '%s'", output)
	}
}

func TestAgentSegmentEmpty(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &AgentSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no agent, got '%s'", output)
	}
}
