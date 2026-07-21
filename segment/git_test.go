package segment

import (
	"strings"
	"testing"

	"github.com/huyhandes/cc-hud-go/state"
)

func TestGitSegment(t *testing.T) {
	s := state.New()
	s.Git.Branch = "main"
	s.Git.DirtyFiles = 3
	s.Git.Ahead = 2
	s.Git.Behind = 1

	seg := &GitSegment{}

	if seg.ID() != "git" {
		t.Errorf("expected ID 'git', got '%s'", seg.ID())
	}

	output := seg.Render(s)
	if !strings.Contains(output, "main") {
		t.Errorf("expected branch name in output, got '%s'", output)
	}

	if !strings.Contains(output, "3") {
		t.Errorf("expected dirty files count in output, got '%s'", output)
	}
}

func TestGitSegmentNoBranch(t *testing.T) {
	s := state.New()

	seg := &GitSegment{}

	output := seg.Render(s)
	if output != "" {
		t.Errorf("expected empty output with no branch, got '%s'", output)
	}
}
