package segment

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var (
	branchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Magenta
	dirtyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
)

type GitSegment struct{}

func (g *GitSegment) ID() string {
	return "git"
}

func (g *GitSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Git
}

func (g *GitSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	if s.Git.Branch == "" {
		return "", nil
	}

	var parts []string

	// Branch name
	if cfg.Git.ShowBranch {
		parts = append(parts, branchStyle.Render(s.Git.Branch))
	}

	// Dirty indicator
	if cfg.Git.ShowDirty && s.Git.DirtyFiles > 0 {
		parts = append(parts, dirtyStyle.Render(fmt.Sprintf("✗%d", s.Git.DirtyFiles)))
	}

	// Ahead/behind
	if cfg.Git.ShowAheadBehind {
		if s.Git.Ahead > 0 {
			parts = append(parts, fmt.Sprintf("↑%d", s.Git.Ahead))
		}
		if s.Git.Behind > 0 {
			parts = append(parts, fmt.Sprintf("↓%d", s.Git.Behind))
		}
	}

	// File stats
	if cfg.Git.ShowFileStats {
		if s.Git.Added > 0 {
			parts = append(parts, fmt.Sprintf("+%d", s.Git.Added))
		}
		if s.Git.Modified > 0 {
			parts = append(parts, fmt.Sprintf("~%d", s.Git.Modified))
		}
		if s.Git.Deleted > 0 {
			parts = append(parts, fmt.Sprintf("-%d", s.Git.Deleted))
		}
	}

	return strings.Join(parts, " "), nil
}
