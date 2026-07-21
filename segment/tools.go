package segment

import (
	"fmt"

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
	rows := t.buildRows(s, cfg)
	if len(rows) == 0 {
		return "", nil
	}
	return style.RenderTable([]string{"Tool", "Count", "Last Used"}, rows), nil
}

func (t *ToolsSegment) buildRows(s *state.State, cfg *config.Config) [][]string {
	var rows [][]string

	if total := sumMap(s.Tools.AppTools); total > 0 {
		rows = append(rows, []string{
			"📦 App",
			fmt.Sprintf("%d", total),
			truncate(s.Tools.LastUsed["app"], 24),
		})
	}

	if total := sumMap(s.Tools.InternalTools); total > 0 {
		rows = append(rows, []string{
			"🖥 Shell",
			fmt.Sprintf("%d", total),
			truncate(s.Tools.LastUsed["internal"], 24),
		})
	}

	if total := sumMCP(s.Tools.MCPTools); total > 0 && cfg.Tools.ShowMCP {
		rows = append(rows, []string{
			"🔌 MCP",
			fmt.Sprintf("%d", total),
			truncate(s.Tools.LastUsed["mcp"], 24),
		})
	}

	if total := sumSkills(s.Tools.Skills); total > 0 && cfg.Tools.ShowSkills {
		rows = append(rows, []string{
			"⚡ Skills",
			fmt.Sprintf("%d", total),
			truncate(s.Tools.LastUsed["skills"], 24),
		})
	}

	if total := sumMap(s.Tools.CustomTools); total > 0 {
		rows = append(rows, []string{
			"🎨 Custom",
			fmt.Sprintf("%d", total),
			truncate(s.Tools.LastUsed["custom"], 24),
		})
	}

	return rows
}

func sumMap(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}

func sumMCP(m map[state.MCPServer]map[string]int) int {
	total := 0
	for _, tools := range m {
		for _, v := range tools {
			total += v
		}
	}
	return total
}

func sumSkills(m map[string]state.SkillUsage) int {
	total := 0
	for _, u := range m {
		total += u.Count
	}
	return total
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}
