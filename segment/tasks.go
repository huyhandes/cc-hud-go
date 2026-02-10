package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
)

type TasksSegment struct{}

func (t *TasksSegment) ID() string {
	return "tasks"
}

func (t *TasksSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tasks
}

func (t *TasksSegment) getTotalCount(s *state.State) int {
	return s.Tasks.Completed + s.Tasks.InProgress + s.Tasks.Pending
}

func (t *TasksSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	total := t.getTotalCount(s)

	if total == 0 {
		return "", nil
	}

	// Check if we should use table view
	if total > cfg.Tables.TasksThreshold && len(s.Tasks.Details) > 0 {
		return t.renderTable(s, cfg)
	}

	// Inline dashboard view
	return t.renderInline(s, cfg)
}

func (t *TasksSegment) renderInline(s *state.State, cfg *config.Config) (string, error) {

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

func (t *TasksSegment) truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func (t *TasksSegment) renderTable(s *state.State, cfg *config.Config) (string, error) {
	headers := []string{"Task", "Status"}
	rows := [][]string{}

	// Add all pending tasks
	for _, task := range s.Tasks.Details {
		if task.Status == "pending" {
			subject := t.truncate(task.Subject, 50)
			rows = append(rows, []string{subject, "⏳ Pending"})
		}
	}

	// Add all in-progress tasks
	for _, task := range s.Tasks.Details {
		if task.Status == "in_progress" {
			subject := t.truncate(task.Subject, 50)
			rows = append(rows, []string{subject, "🔄 Active"})
		}
	}

	// Add last 3 completed tasks
	completedTasks := []state.Task{}
	for _, task := range s.Tasks.Details {
		if task.Status == "completed" {
			completedTasks = append(completedTasks, task)
		}
	}

	// Get last 3 completed (most recent)
	startIdx := 0
	if len(completedTasks) > 3 {
		startIdx = len(completedTasks) - 3
	}
	for i := startIdx; i < len(completedTasks); i++ {
		subject := t.truncate(completedTasks[i].Subject, 50)
		rows = append(rows, []string{subject, "✅ Done"})
	}

	return style.RenderTable(headers, rows), nil
}
