package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
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

	// Define colors
	borderColor := lipgloss.Color("240")     // Gray border
	headerColor := lipgloss.Color("14")      // Cyan
	pendingColor := lipgloss.Color("11")     // Yellow
	progressColor := lipgloss.Color("12")    // Blue
	completedColor := lipgloss.Color("10")   // Green

	// Create styles for each component
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(headerColor).
		Width(20)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(18)

	pendingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(pendingColor).
		Align(lipgloss.Right).
		Width(4)

	progressStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(progressColor).
		Align(lipgloss.Right).
		Width(4)

	completedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(completedColor).
		Align(lipgloss.Right).
		Width(4)

	// Build each row
	header := headerStyle.Render("📋 Tasks Dashboard")

	pendingRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render("  ⏳ Todo"),
		pendingStyle.Render(fmt.Sprintf("%d", s.Tasks.Pending)),
	)

	progressRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render("  🔄 In Progress"),
		progressStyle.Render(fmt.Sprintf("%d", s.Tasks.InProgress)),
	)

	completedRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render("  ✅ Completed"),
		completedStyle.Render(fmt.Sprintf("%d", s.Tasks.Completed)),
	)

	// Combine all rows vertically
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		pendingRow,
		progressRow,
		completedRow,
	)

	// Create a bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return boxStyle.Render(content), nil
}
