package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
)

var tasksStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green

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

	// Basic format: "Tasks: 2/5"
	return tasksStyle.Render(fmt.Sprintf("Tasks: %d/%d", s.Tasks.Completed, total)), nil
}
