# Lipgloss Visual Enhancement Design

**Date:** 2026-02-10
**Status:** Approved
**Author:** Design collaboration with user

## Overview

Comprehensive visual enhancement of the cc-hud-go statusline using lipgloss capabilities, including Catppuccin theme support, gradient progress bars, improved spacing, and contextual table rendering.

## Goals

1. **Theme System** - Support Catppuccin colorscheme by default with user customization
2. **Visual Clarity** - Improved spacing and layout for better scannability
3. **Enhanced Rendering** - Gradient progress bars and contextual table views
4. **User Control** - Configurable thresholds and color overrides

## Design Decisions

### 1. Color System Architecture

**Catppuccin Theme Integration**

Create a new `theme` package that manages color palettes:

- **Default Theme**: Catppuccin Macchiato
- **Theme Registry**: Support all 4 Catppuccin flavors (Mocha, Macchiato, Frappé, Latte)
- **Color Mapping**: Direct mapping of Catppuccin colors to semantic meanings

**Semantic Color Mappings (Catppuccin Macchiato):**

| Semantic | Purpose | Catppuccin Color | Hex |
|----------|---------|------------------|-----|
| ColorSuccess | Healthy/good status | Green | `#a6da95` |
| ColorWarning | Caution state | Yellow | `#eed49f` |
| ColorDanger | Critical state | Red | `#ed8796` |
| ColorInput | Incoming data | Blue | `#8aadf4` |
| ColorOutput | Outgoing data | Teal | `#8bd5ca` |
| ColorCacheRead | Cache read ops | Mauve | `#c6a0f6` |
| ColorCacheWrite | Cache write ops | Pink | `#f5bde6` |
| ColorPrimary | Model/agent | Mauve | `#c6a0f6` |
| ColorHighlight | Git/duration | Sky | `#91d7e3` |
| ColorAccent | Cost/emphasis | Peach | `#f5a97f` |
| ColorMuted | Separators/static | Overlay0 | (varies by flavor) |
| ColorBright | Main text | Text | (varies by flavor) |
| ColorInfo | Information | Teal | `#8bd5ca` |

**Config Structure:**
```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#a6da95",
    "warning": "#eed49f",
    "danger": "#ed8796"
  }
}
```

Theme loads selected Catppuccin flavor, then applies semantic color overrides from config.

### 2. Gradient Progress Bar

**Replace dot-based progress** (`●●●●●○○○○○`) with segmented gradient bars.

**Specifications:**
- **Width**: 10 characters (configurable)
- **Character set**: `█▓▒░` (solid → dark → medium → light)
- **Color gradient**: Green → Yellow → Peach → Red based on percentage
  - 0-70%: Green shades
  - 70-90%: Green → Yellow → Peach transition
  - 90-100%: Peach → Red transition

**Algorithm:**
```go
func RenderGradientBar(percentage float64, width int) string {
    filled := int(percentage / 100 * float64(width))

    segments := []string{}
    for i := 0; i < width; i++ {
        // Calculate color for this segment based on overall percentage
        segmentPercent := (float64(i) / float64(width)) * 100
        color := interpolateColor(segmentPercent, percentage)

        // Determine character based on fill position
        if i < filled {
            char := getGradientChar(i, filled) // █▓▒ based on position
            segments = append(segments, style.WithColor(color).Render(char))
        } else {
            segments = append(segments, style.Muted.Render("░"))
        }
    }

    return strings.Join(segments, "")
}
```

**Usage:**
- **Context segment**: `█████▓▓▒▒░░ 54%`
- **Rate limit segment**: `████████▓▒ 87%`
- **Tasks**: Continue using table format (not progress bar)

### 3. Layout & Spacing

**Enhanced Separator Layout**

**Current:**
```
🟢 ●●●●●○○○○○ 54% │ 🌿 main ⚠5 │ 🔧 7
```

**New:**
```
🟢 █████▓▓▒▒░░ 54%  │  🌿 main ⚠5  │  🔧 7 tools
```

**Implementation:**
- Two spaces before and after each separator (`│`)
- Apply in `output/renderer.go` when joining segments
- Respect existing `LineLayout` config (compact/expanded)

**Helper function:**
```go
func JoinSegments(segments []string) string {
    return strings.Join(segments, "  │  ")
}
```

**Multi-line layout:**
When `LineLayout: "expanded"`, segments grouped logically:
```
Line 1: Model + Context
Line 2: Git + Cost
Line 3: Tools + Agent + Tasks
```

### 4. Table Rendering System

**Threshold-Based Table Views**

Segments automatically switch from compact inline to table view when data exceeds thresholds.

**Config:**
```json
{
  "display": {
    "toolsTableThreshold": 5,
    "tasksTableThreshold": 3,
    "contextTableThreshold": 999
  }
}
```

**Tools Segment:**
- **Compact** (≤5 tools): `🔧 7 (App:5 MCP:2)`
- **Table** (>5 tools):
```
┌──────────┬───────┐
│ Category │ Count │
├──────────┼───────┤
│ App      │   12  │
│ MCP      │    3  │
│ Skills   │    2  │
└──────────┴───────┘
```

**Tasks Segment:**
Shows **last 3 completed** + all active + all pending:
```
┌─────────────────────┬──────────┐
│ Task                │ Status   │
├─────────────────────┼──────────┤
│ Implement auth      │ ✓ Done   │
│ Write tests         │ ✓ Done   │
│ Add error handling  │ ✓ Done   │
│ Update docs         │ ⏵ Active │
│ Deploy staging      │ ○ Pending│
└─────────────────────┴──────────┘
```

**Context Segment:**
High threshold (999) keeps inline by default, can show table on demand.

**Decision Logic:**
```go
func (s *ToolsSegment) Render(state *state.State, cfg *config.Config) (string, error) {
    toolCount := len(state.Tools.ByCategory)

    if toolCount > cfg.Tables.ToolsThreshold {
        return s.RenderTable(state, cfg)
    }

    return s.RenderInline(state, cfg)
}
```

### 5. Component Architecture

**New Package: `theme/`**
```
theme/
├── theme.go          // Theme interface, registry, loader
├── catppuccin.go     // Catppuccin flavor definitions
└── theme_test.go
```

**Key functions:**
- `GetTheme(name string) Theme` - Returns theme by name
- `LoadFromConfig(cfg *config.Config) Theme` - Loads theme + overrides
- `Theme` interface with `GetColor(semantic string) lipgloss.Color`

**Modified: `style/` package**
- Remove hardcoded color constants
- Import colors from active theme: `style.Init(theme)`
- Add `RenderGradientBar(percentage, width) string`
- Add `RenderTable(headers, rows []string) string`

**Modified: `config/` package**
```go
type Config struct {
    Theme   string
    Colors  map[string]string  // semantic overrides
    Display DisplayConfig
    Tables  TableConfig        // NEW
    // ... existing fields
}

type TableConfig struct {
    ToolsThreshold   int `json:"toolsTableThreshold"`
    TasksThreshold   int `json:"tasksTableThreshold"`
    ContextThreshold int `json:"contextTableThreshold"`
}
```

**Modified: `segment/` implementations**
- `context.go`: Use `RenderGradientBar()` instead of dots
- `ratelimit.go`: Use `RenderGradientBar()` for 7-day usage
- `tasks.go`: Implement table view with "last 3 completed" filter
- `tools.go`: Implement table view for category breakdown

### 6. Data Flow

**Application Startup → Rendering:**

1. **Initialization** (`main.go`)
   ```go
   cfg := config.LoadFromFile(configPath)
   themeInstance := theme.LoadFromConfig(cfg)
   style.Init(themeInstance)
   ```

2. **State Parsing** (existing flow)
   - Parse stdin JSON → `state.State`
   - Parse transcript JSONL → tool usage
   - Fetch git info → git stats

3. **Segment Rendering** (`segment/segment.go`)
   ```go
   for _, seg := range segments {
       if !seg.Enabled(cfg) { continue }

       output := seg.Render(state, cfg)
   }
   ```

4. **Output Assembly** (`output/renderer.go`)
   ```go
   joined := JoinSegments(outputs) // with 2-space separators
   json.Marshal(StatuslineOutput{Text: joined})
   ```

**Key Decision Points:**
- Theme loads once at startup
- Segments decide compact vs table based on thresholds
- Color interpolation in `RenderGradientBar()`
- Tables use lipgloss.Table for consistent formatting

### 7. Error Handling & Edge Cases

**Graceful Degradation:**

**Missing/invalid theme:**
```go
theme, err := theme.GetTheme(cfg.Theme)
if err != nil {
    theme = theme.GetTheme("macchiato") // fallback
    fmt.Fprintf(os.Stderr, "warning: invalid theme '%s', using macchiato\n", cfg.Theme)
}
```

**Invalid color overrides:**
- Parse hex colors, fall back to theme default if invalid
- Log warning but continue rendering

**Terminal color support:**
- Handled by lipgloss renderer with TrueColor profile
- Automatic degradation to 256 colors

**Empty/missing data:**
- Gradient bar with 0 tokens: `░░░░░░░░░░ 0%`
- No tasks: Hide tasks segment
- No tools: Show `🔧 0` inline

**Threshold edge cases:**
- At threshold (e.g., 5 tools, threshold=5): Use inline
- Long task names: Truncate with `...` in table
- More than 3 completed: Keep last 3, sorted by completion time

**Config validation:**
```go
func (c *TableConfig) Validate() {
    if c.ToolsThreshold < 1 { c.ToolsThreshold = 5 }
    if c.TasksThreshold < 1 { c.TasksThreshold = 3 }
}
```

## Testing Strategy

**Unit Tests:**

**Theme package:**
- `TestCatppuccinFlavors()` - All 4 flavors load correctly
- `TestThemeColorMapping()` - Semantic colors map correctly
- `TestConfigOverrides()` - User overrides apply
- `TestInvalidTheme()` - Falls back to default

**Style package:**
- `TestGradientBarRendering()` - Correct characters at percentages
- `TestColorInterpolation()` - Smooth color transitions
- `TestGradientBarEdgeCases()` - 0%, 100%, invalid values
- `TestTableRendering()` - Headers, rows, borders

**Segment tests:**
- `TestContextGradientBar()` - Uses gradient instead of dots
- `TestTasksTableFiltering()` - Last 3 completed filter works
- `TestToolsTableThreshold()` - Switches at threshold
- `TestRateLimitGradientBar()` - 7-day usage gradient

**Integration tests:**
- `TestFullStatuslineWithTheme()` - End-to-end with Catppuccin
- `TestThemeSwitching()` - Theme changes apply
- `TestConfigColorOverrides()` - Custom colors throughout

**Visual regression:**
- Reference screenshots in `testdata/screenshots/`
- Compare rendered output for consistency

## Implementation Phases

### Phase 1: Theme System
- Create `theme/` package
- Implement Catppuccin flavors
- Add config support for theme selection
- Update `style/` to use theme colors

### Phase 2: Gradient Progress Bars
- Implement `RenderGradientBar()` helper
- Update `context.go` segment
- Update `ratelimit.go` segment
- Add color interpolation logic

### Phase 3: Layout & Spacing
- Update separator spacing in renderer
- Test with compact/expanded layouts
- Verify visual consistency

### Phase 4: Table Rendering
- Add `TableConfig` to config
- Implement table views for tools/tasks segments
- Add threshold detection logic
- Implement "last 3 completed" filter for tasks

### Phase 5: Testing & Polish
- Write comprehensive tests
- Create visual regression suite
- Update documentation
- Add example config files

## Configuration Example

```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#a6da95",
    "warning": "#eed49f"
  },
  "preset": "full",
  "lineLayout": "expanded",
  "display": {
    "model": true,
    "context": true,
    "git": true,
    "tools": true,
    "tasks": true,
    "rateLimits": true
  },
  "tables": {
    "toolsTableThreshold": 5,
    "tasksTableThreshold": 3,
    "contextTableThreshold": 999
  }
}
```

## Visual Examples

**Before:**
```
🟢 ●●●●●○○○○○ 54% │ 🌿 main ⚠5 │ 🔧 7
```

**After (inline):**
```
🟢 █████▓▓▒▒░░ 54%  │  🌿 main ⚠5  │  🔧 7 tools
```

**After (with tables):**
```
🟢 █████▓▓▒▒░░ 54%  │  🌿 main ⚠5

┌──────────┬───────┐
│ Category │ Count │
├──────────┼───────┤
│ App      │   12  │
│ MCP      │    3  │
│ Skills   │    2  │
└──────────┴───────┘

┌─────────────────────┬──────────┐
│ Task                │ Status   │
├─────────────────────┼──────────┤
│ Implement auth      │ ✓ Done   │
│ Write tests         │ ⏵ Active │
│ Deploy staging      │ ○ Pending│
└─────────────────────┴──────────┘
```

## Benefits

1. **Visual Consistency** - Catppuccin provides cohesive, proven color palette
2. **User Choice** - Support for 4 flavors + custom overrides
3. **Better Scannability** - Improved spacing and gradient bars
4. **Adaptive Display** - Tables appear when data is complex
5. **Semantic Colors** - Meaningful color associations remain intact
6. **Smooth Gradients** - Visual feedback on usage levels
7. **Task Focus** - Shows recent wins + current work, not endless history

## Non-Goals

- Tree rendering (tables sufficient for statusline)
- Dynamic theme switching during runtime
- Animation or transitions
- Custom theme creation UI (config-only)
