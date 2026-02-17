package segment

import (
	"fmt"
	"strings"

	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type TasksSegment struct{}

func (t *TasksSegment) ID() string {
	return "tasks"
}

func (t *TasksSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tasks
}

func (t *TasksSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	total := s.Tasks.Completed + s.Tasks.InProgress + s.Tasks.Pending
	if total == 0 {
		return "", nil
	}

	pendingStyle := style.GetRenderer().NewStyle().Bold(true).Foreground(style.ColorWarning)
	progressStyle := style.GetRenderer().NewStyle().Bold(true).Foreground(style.ColorInfo)
	completedStyle := style.GetRenderer().NewStyle().Bold(true).Foreground(style.ColorSuccess)

	var parts []string
	if s.Tasks.Pending > 0 {
		parts = append(parts, pendingStyle.Render(fmt.Sprintf("⏳ %d", s.Tasks.Pending)))
	}
	if s.Tasks.InProgress > 0 {
		parts = append(parts, progressStyle.Render(fmt.Sprintf("🔄 %d", s.Tasks.InProgress)))
	}
	if s.Tasks.Completed > 0 {
		parts = append(parts, completedStyle.Render(fmt.Sprintf("✅ %d", s.Tasks.Completed)))
	}

	return strings.Join(parts, "  "), nil
}
