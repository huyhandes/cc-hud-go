package segment

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/huyhandes/cc-hud-go/config"
	"github.com/huyhandes/cc-hud-go/state"
	"github.com/huyhandes/cc-hud-go/style"
)

type ToolsSegment struct{}

func (t *ToolsSegment) ID() string {
	return "tools"
}

func (t *ToolsSegment) Enabled(cfg *config.Config) bool {
	return cfg.Display.Tools
}

func (t *ToolsSegment) Render(s *state.State, cfg *config.Config) (string, error) {
	toolCount := t.getTotalCount(s)

	if toolCount == 0 {
		return "", nil
	}

	// Check if we should use table view
	if toolCount > cfg.Tables.ToolsThreshold {
		return t.renderTable(s, cfg)
	}

	// Inline view
	return t.renderInline(s, cfg)
}

func (t *ToolsSegment) getTotalCount(s *state.State) int {
	total := 0

	for _, count := range s.Tools.AppTools {
		total += count
	}
	for _, count := range s.Tools.InternalTools {
		total += count
	}
	for _, count := range s.Tools.CustomTools {
		total += count
	}
	for _, tools := range s.Tools.MCPTools {
		for _, count := range tools {
			total += count
		}
	}
	for _, usage := range s.Tools.Skills {
		total += usage.Count
	}

	return total
}

func (t *ToolsSegment) renderInline(s *state.State, cfg *config.Config) (string, error) {
	toolCount := t.getTotalCount(s)

	// Simple inline display if not grouped
	if !cfg.Tools.GroupByCategory {
		icon := "🔧"
		toolsMainStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
		return toolsMainStyle.Render(fmt.Sprintf("%s %d", icon, toolCount)), nil
	}

	// Enhanced lipgloss display when grouped by category
	icon := "🔧"
	appTotal := 0
	for _, count := range s.Tools.AppTools {
		appTotal += count
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
	customTotal := 0
	for _, count := range s.Tools.CustomTools {
		customTotal += count
	}

	borderColor := lipgloss.Color("240")
	headerColor := lipgloss.Color("14") // Cyan
	appColor := lipgloss.Color("12")    // Blue
	mcpColor := lipgloss.Color("13")    // Magenta
	skillsColor := lipgloss.Color("11") // Yellow
	customColor := lipgloss.Color("10") // Green

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
	header := headerStyle.Render(fmt.Sprintf("%s Tool Usage (%d)", icon, toolCount))

	// Build rows for each category
	var rows []string

	if appTotal > 0 {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  📦 App"),
			countStyle.Foreground(appColor).Render(fmt.Sprintf("%d", appTotal)),
		)
		rows = append(rows, row)
	}

	if mcpTotal > 0 && cfg.Tools.ShowMCP {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  🔌 MCP"),
			countStyle.Foreground(mcpColor).Render(fmt.Sprintf("%d", mcpTotal)),
		)
		rows = append(rows, row)
	}

	if skillsTotal > 0 && cfg.Tools.ShowSkills {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  ⚡ Skills"),
			countStyle.Foreground(skillsColor).Render(fmt.Sprintf("%d", skillsTotal)),
		)
		rows = append(rows, row)
	}

	if customTotal > 0 {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render("  🎨 Custom"),
			countStyle.Foreground(customColor).Render(fmt.Sprintf("%d", customTotal)),
		)
		rows = append(rows, row)
	}

	// If no categories to show, just show total
	if len(rows) == 0 {
		icon := "🔧"
		toolsMainStyle := style.GetRenderer().NewStyle().Foreground(style.ColorInfo)
		return toolsMainStyle.Render(fmt.Sprintf("%s %d", icon, toolCount)), nil
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

func (t *ToolsSegment) renderTable(s *state.State, cfg *config.Config) (string, error) {
	headers := []string{"Category", "Count"}
	rows := [][]string{}

	// Collect categories with counts
	categories := make(map[string]int)

	appTotal := 0
	for _, count := range s.Tools.AppTools {
		appTotal += count
	}
	if appTotal > 0 {
		categories["App"] = appTotal
	}

	mcpTotal := 0
	for _, tools := range s.Tools.MCPTools {
		for _, count := range tools {
			mcpTotal += count
		}
	}
	if mcpTotal > 0 && cfg.Tools.ShowMCP {
		categories["MCP"] = mcpTotal
	}

	skillsTotal := 0
	for _, usage := range s.Tools.Skills {
		skillsTotal += usage.Count
	}
	if skillsTotal > 0 && cfg.Tools.ShowSkills {
		categories["Skills"] = skillsTotal
	}

	customTotal := 0
	for _, count := range s.Tools.CustomTools {
		customTotal += count
	}
	if customTotal > 0 {
		categories["Custom"] = customTotal
	}

	// Sort categories for consistent order
	categoryNames := make([]string, 0, len(categories))
	for cat := range categories {
		categoryNames = append(categoryNames, cat)
	}
	sort.Strings(categoryNames)

	for _, cat := range categoryNames {
		count := categories[cat]
		rows = append(rows, []string{cat, fmt.Sprintf("%d", count)})
	}

	return style.RenderTable(headers, rows), nil
}
