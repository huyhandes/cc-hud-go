package segment

import (
	"fmt"
	"strings"

	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
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

	// Branch name with icon - Cyan (highlight color)
	if cfg.Git.ShowBranch {
		branchIcon := "🌿"
		branchStyle := style.GetRenderer().NewStyle().Foreground(style.ColorHighlight).Bold(true)
		parts = append(parts, branchStyle.Render(fmt.Sprintf("%s %s", branchIcon, s.Git.Branch)))
	}

	// Dirty indicator with warning icon - Orange (warning)
	if cfg.Git.ShowDirty && s.Git.DirtyFiles > 0 {
		dirtyStyle := style.GetRenderer().NewStyle().Foreground(style.ColorWarning)
		parts = append(parts, dirtyStyle.Render(fmt.Sprintf("⚠%d", s.Git.DirtyFiles)))
	}

	// Ahead/behind with colored arrows
	if cfg.Git.ShowAheadBehind {
		if s.Git.Ahead > 0 {
			// Ahead - Emerald/Green (good, pushing forward)
			aheadStyle := style.GetRenderer().NewStyle().Foreground(style.ColorSuccess)
			parts = append(parts, aheadStyle.Render(fmt.Sprintf("↑%d", s.Git.Ahead)))
		}
		if s.Git.Behind > 0 {
			// Behind - Red (needs attention)
			behindStyle := style.GetRenderer().NewStyle().Foreground(style.ColorDanger)
			parts = append(parts, behindStyle.Render(fmt.Sprintf("↓%d", s.Git.Behind)))
		}
	}

	// File stats with diverse colors
	if cfg.Git.ShowFileStats {
		if s.Git.Added > 0 {
			// Added - Green (new/positive)
			addedStyle := style.GetRenderer().NewStyle().Foreground(style.ColorSuccess)
			parts = append(parts, addedStyle.Render(fmt.Sprintf("+%d", s.Git.Added)))
		}
		if s.Git.Modified > 0 {
			// Modified - Teal (changed/neutral)
			modStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
			parts = append(parts, modStyle.Render(fmt.Sprintf("~%d", s.Git.Modified)))
		}
		if s.Git.Deleted > 0 {
			// Deleted - Red (removed/negative)
			delStyle := style.GetRenderer().NewStyle().Foreground(style.ColorDanger)
			parts = append(parts, delStyle.Render(fmt.Sprintf("-%d", s.Git.Deleted)))
		}
	}

	return strings.Join(parts, " "), nil
}
