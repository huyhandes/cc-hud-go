package parser

import "strings"

// ToolCategory represents the type of tool being used
type ToolCategory int

const (
	CategoryApp ToolCategory = iota
	CategoryInternal
	CategoryCustom
	CategoryMCP
	CategorySkill
)

var appTools = map[string]bool{
	"read":      true,
	"write":     true,
	"edit":      true,
	"bash":      true,
	"glob":      true,
	"grep":      true,
	"task":      true,
	"webfetch":  true,
	"websearch": true,
}

// CategorizeTool determines the category of a tool by name
func CategorizeTool(name string) ToolCategory {
	lower := strings.ToLower(name)

	if strings.HasPrefix(lower, "mcp__") {
		return CategoryMCP
	}
	if lower == "skill" {
		return CategorySkill
	}
	if lower == "bash" {
		return CategoryInternal
	}
	if appTools[lower] {
		return CategoryApp
	}
	return CategoryCustom
}
