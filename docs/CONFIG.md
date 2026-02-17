# Configuration Reference

Config file location: `~/.claude/cc-hud-go/config.json`

## Full Example

```json
{
  "theme": "macchiato",
  "colors": {},
  "preset": "full",
  "lineLayout": "expanded",
  "pathLevels": 2,
  "sevenDayThreshold": 80,
  "display": {
    "model": true,
    "context": true,
    "git": true,
    "tools": true,
    "agents": true,
    "tasks": true,
    "rateLimits": true,
    "duration": true,
    "fetchOAuth": true
  },
  "git": {
    "showBranch": true,
    "showDirty": true,
    "showAheadBehind": true,
    "showFileStats": true
  },
  "tools": {
    "groupByCategory": true,
    "showTopN": 5,
    "showSkills": true,
    "showMCP": true
  },
  "tables": {
    "toolsTableThreshold": 999,
    "tasksTableThreshold": 999,
    "contextTableThreshold": 999
  }
}
```

## Top-Level Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"macchiato"` | Color theme: `macchiato`, `mocha`, `frappe`, `latte` |
| `colors` | object | `{}` | Custom semantic color overrides (hex) |
| `preset` | string | `"full"` | Preset: `full`, `essential`, `minimal` |
| `lineLayout` | string | `"expanded"` | Layout: `expanded` or `compact` |
| `pathLevels` | int | `2` | Directory levels to show in path (1–3) |
| `sevenDayThreshold` | int | `80` | Warning threshold for 7-day rate limit (0–100) |

## Presets

| Preset | Description |
|--------|-------------|
| `full` | All segments enabled |
| `essential` | Core metrics only (model, context, git, cost); switches to compact layout |
| `minimal` | Minimal info (model, context); compact layout, git/tools/agents/tasks/rateLimits/duration hidden |

## Display Options

Boolean flags to enable/disable individual segments:

| Key | Default | Description |
|-----|---------|-------------|
| `model` | `true` | Model name and plan type |
| `context` | `true` | Token usage with color thresholds |
| `git` | `true` | Git branch, status, file stats |
| `tools` | `true` | Tool usage by category |
| `agents` | `true` | Active agent name and task |
| `tasks` | `true` | Task completion progress |
| `rateLimits` | `true` | 7-day API usage |
| `duration` | `true` | Session duration |
| `fetchOAuth` | `true` | Fetch OAuth token info (plan type, rate limits) |

## Git Options

| Key | Default | Description |
|-----|---------|-------------|
| `showBranch` | `true` | Current branch name |
| `showDirty` | `true` | Count of uncommitted files |
| `showAheadBehind` | `true` | Commits ahead/behind remote |
| `showFileStats` | `true` | Added/modified/deleted counts |

## Tools Options

| Key | Default | Description |
|-----|---------|-------------|
| `groupByCategory` | `true` | Group by App/MCP/Skills/Custom |
| `showTopN` | `5` | Number of top tools to show (0 = all) |
| `showSkills` | `true` | Include skill usage |
| `showMCP` | `true` | Include MCP tool usage |

## Table Options

Smart threshold: when item count exceeds the threshold, switches from styled lipgloss boxes (╭╮) to plain table view (┌┬┐).

| Key | Default | Description |
|-----|---------|-------------|
| `toolsTableThreshold` | `999` | Tool count for table view |
| `tasksTableThreshold` | `999` | Task count for table view |
| `contextTableThreshold` | `999` | Context size for table view |

Default of 999 always uses the styled lipgloss box. Lower to 3–5 for earlier table switching on large datasets.

## Custom Colors

Override any semantic color while keeping the base theme:

```json
{
  "theme": "macchiato",
  "colors": {
    "success": "#00ff00",
    "warning": "#ffaa00",
    "danger": "#ff0000",
    "primary": "#00aaff"
  }
}
```

Available semantic color keys:

| Key | Role |
|-----|------|
| `success` | Positive states, low usage (green) |
| `warning` | Caution, medium usage (yellow/orange) |
| `danger` | Critical states, high usage (red) |
| `input` | Input tokens (blue) |
| `output` | Output tokens (emerald) |
| `cacheRead` | Cache read tokens (purple) |
| `cacheWrite` | Cache write tokens (pink) |
| `primary` | Main brand color (purple) |
| `highlight` | Highlights and accents (cyan) |
| `accent` | Secondary accents (orange) |
| `muted` | Borders, subtle elements (gray) |
| `bright` | Bright text (white/cream) |
| `info` | Informational elements (teal) |

## Example Configs

Pre-configured examples in `examples/`:

```bash
cp examples/config-macchiato.json ~/.claude/cc-hud-go/config.json
cp examples/config-mocha.json ~/.claude/cc-hud-go/config.json
cp examples/config-frappe.json ~/.claude/cc-hud-go/config.json
cp examples/config-latte.json ~/.claude/cc-hud-go/config.json
cp examples/config-custom-colors.json ~/.claude/cc-hud-go/config.json
```

See `examples/README.md` for the full theme guide.
