package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// RenderTable renders a table with headers and rows using lipgloss table
func RenderTable(headers []string, rows [][]string) string {
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(renderer.NewStyle().Foreground(ColorMuted)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return renderer.NewStyle().Foreground(ColorBright).Bold(true).Padding(0, 1)
			}
			return renderer.NewStyle().Padding(0, 1)
		})

	return t.Render()
}
