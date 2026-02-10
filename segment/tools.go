package segment

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/huybui/cc-hud-go/config"
	"github.com/huybui/cc-hud-go/state"
	"github.com/huybui/cc-hud-go/style"
)

type ToolsSegment struct{}

func (t *ToolsSegment) ID() string {
	return "tools"
}

func (t *ToolsSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tools
}

func (t *ToolsSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	// Count totals
	appTotal := 0
	for _, count := range s.Tools.AppTools {
		appTotal += count
	}

	internalTotal := 0
	for _, count := range s.Tools.InternalTools {
		internalTotal += count
	}

	customTotal := 0
	for _, count := range s.Tools.CustomTools {
		customTotal += count
	}

	mcpTotal := 0
	for _, tools := range s.Tools.MCPTools {
		for _, count := range tools {
			mcpTotal += count
		}
	}

	skillsTotal := 0
	for _, usage := range s.Tools.Skills {
		skillsTotal += usage.Count
	}

	total := appTotal + internalTotal + customTotal + mcpTotal + skillsTotal

	if total == 0 {
		return "", nil
	}

	// Build output with icon - Teal/Info color (distinct from other segments)
	icon := "🔧"
	toolsMainStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)

	// Simple inline display if not grouped
	if !cfg.Tools.GroupByCategory {
		return toolsMainStyle.Render(fmt.Sprintf("%s %d", icon, total)), nil
	}

	// Enhanced lipgloss display when grouped by category
	borderColor := lipgloss.Color("240")
	headerColor := lipgloss.Color("14")      // Cyan
	appColor := lipgloss.Color("12")         // Blue
	mcpColor := lipgloss.Color("13")         // Magenta
	skillsColor := lipgloss.Color("11")      // Yellow
	customColor := lipgloss.Color("10")      // Green

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(headerColor).
		Width(22)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(16)

	countStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Right).
		Width(6)

	// Build header
	header := headerStyle.Render(fmt.Sprintf("%s Tool Usage (%d)", icon, total))

	// Build rows for each category
	var rows []string

	if appTotal > 0 {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  📦 App"),
			countStyle.Copy().Foreground(appColor).Render(fmt.Sprintf("%d", appTotal)),
		)
		rows = append(rows, row)
	}

	if mcpTotal > 0 && cfg.Tools.ShowMCP {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  🔌 MCP"),
			countStyle.Copy().Foreground(mcpColor).Render(fmt.Sprintf("%d", mcpTotal)),
		)
		rows = append(rows, row)
	}

	if skillsTotal > 0 && cfg.Tools.ShowSkills {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  ⚡ Skills"),
			countStyle.Copy().Foreground(skillsColor).Render(fmt.Sprintf("%d", skillsTotal)),
		)
		rows = append(rows, row)
	}

	if customTotal > 0 {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  🎨 Custom"),
			countStyle.Copy().Foreground(customColor).Render(fmt.Sprintf("%d", customTotal)),
		)
		rows = append(rows, row)
	}

	// If no categories to show, just show total
	if len(rows) == 0 {
		return toolsMainStyle.Render(fmt.Sprintf("%s %d", icon, total)), nil
	}

	// Combine header and rows
	var contentParts []string
	contentParts = append(contentParts, header)
	contentParts = append(contentParts, rows...)

	content := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	// Create bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return boxStyle.Render(content), nil
}
