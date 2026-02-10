package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
)

func TestGitSegment(t *testing.T) {
	cfg := config.Default()
	s := state.New()
	s.Git.Branch = "main"
	s.Git.DirtyFiles = 3
	s.Git.Ahead = 2
	s.Git.Behind = 1

	seg := &GitSegment{}

	if seg.ID() != "git" {
		t.Errorf("expected ID 'git', got '%s'", seg.ID())
	}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !strings.Contains(output, "main") {
		t.Errorf("expected branch name in output, got '%s'", output)
	}

	if !strings.Contains(output, "3") {
		t.Errorf("expected dirty files count in output, got '%s'", output)
	}
}

func TestGitSegmentNoBranch(t *testing.T) {
	cfg := config.Default()
	s := state.New()

	seg := &GitSegment{}

	output, err := seg.Render(s, cfg)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty output with no branch, got '%s'", output)
	}
}
