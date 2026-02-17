# Color Scheme Guide

## Design Philosophy

Colors are organized into **semantic groups** where each color has a specific meaning, making the statusline intuitive and easy to scan.

## Color Groups

### 🎯 Status Colors (Usage Indicators)
| Color | Hex | Usage | Meaning |
|-------|-----|-------|---------|
| 🟢 Green | `#10B981` | 0-70% usage | Healthy, plenty of capacity |
| 🟡 Yellow | `#F59E0B` | 70-90% usage | Caution, approaching limit |
| 🔴 Red | `#EF4444` | 90-100% usage | Critical, near capacity |

### 📊 Data Flow Colors
| Color | Hex | Element | Meaning |
|-------|-----|---------|---------|
| 🔵 Blue | `#3B82F6` | Input tokens (📥) | Incoming data/requests |
| 🟢 Emerald | `#10B981` | Output tokens (📤) | Outgoing data/responses |

### 💾 Storage Layer Colors
| Color | Hex | Element | Meaning |
|-------|-----|---------|---------|
| 🟣 Purple | `#8B5CF6` | Cache read (R:) | Reading from storage |
| 🩷 Pink | `#EC4899` | Cache write (W:) | Writing to storage |

### 🎨 Primary UI Colors
| Color | Hex | Element | Meaning |
|-------|-----|---------|---------|
| 🟣 Purple | `#7C3AED` | Model, Agent | AI/identity |
| 🔵 Cyan | `#06B6D4` | Git branch, Duration | Highlights, time |
| 🟠 Orange | `#F59E0B` | Cost, Warnings | Emphasis, attention |
| 🔷 Teal | `#14B8A6` | Tools, Modified files | Information, changes |

### 📁 Git Status Colors
| Color | Hex | Element | Meaning |
|-------|-----|---------|---------|
| 🔵 Cyan Bold | `#06B6D4` | Branch name | Current location |
| 🟠 Orange | `#F59E0B` | Dirty files (⚠) | Uncommitted changes |
| 🟢 Green | `#10B981` | Ahead (↑), Added (+) | Progress, additions |
| 🔴 Red | `#EF4444` | Behind (↓), Deleted (-) | Needs sync, removals |
| 🔷 Teal | `#14B8A6` | Modified (~) | Changed files |

### 💰 Cost Metrics Colors
| Color | Hex | Element | Meaning |
|-------|-----|---------|---------|
| 🟠 Orange Bold | `#F59E0B` | Cost ($) | Financial emphasis |
| 🔵 Cyan | `#06B6D4` | Duration (⏱) | Time tracking |
| 🟢 Green | `#10B981` | Lines added (+) | Productivity gain |
| 🔴 Red | `#EF4444` | Lines removed (-) | Code reduction |

### ⚙️ Utility Colors
| Color | Hex | Usage | Meaning |
|-------|-----|-------|---------|
| ⚫ Gray | `#6B7280` | Separators, context size | Muted, static info |

## Example Output

```
🤖 Sonnet 4.5 │ 🟢 ●●●●●○○○○○ 54% 📥 108k 📤 20k 💾 R:5k/W:5k ⚡ 200k
🌿 main ⚠5 ↑14 ~5 │ 💰$13.7793 ⏱ 51m44s 📝 +758/-366
🔧 7 (App:5 MCP:1) │ 👤 code-reviewer
```

### Color Breakdown

**Line 1 - Context Information:**
- 🤖 Model: Purple (identity)
- 🟢 Status: Green (healthy usage)
- 54%: Green (matches status)
- 📥 108k: Blue (input data)
- 📤 20k: Emerald (output data)
- 💾 R:5k: Purple (cache read)
- 💾 W:5k: Pink (cache write)
- ⚡ 200k: Gray (static constant)

**Line 2 - Development Status:**
- 🌿 main: Cyan bold (git branch)
- ⚠5: Orange (dirty files warning)
- ↑14: Green (ahead commits)
- ~5: Teal (modified files)
- 💰$13.7793: Orange bold (cost)
- ⏱ 51m44s: Cyan (duration)
- +758: Green (lines added)
- -366: Red (lines removed)

**Line 3 - Activity:**
- 🔧 7: Teal (tools usage)
- 👤 code-reviewer: Purple italic (agent)

## Design Principles

1. **Semantic Consistency**: Related concepts use similar colors
2. **Visual Hierarchy**: Important info (cost, warnings) uses bold/bright colors
3. **At-a-glance Scanning**: Each metric has distinct color for quick identification
4. **Color Psychology**:
   - Green = positive/healthy
   - Red = warning/needs attention
   - Blue = incoming/passive
   - Purple = AI/processing
   - Orange = cost/emphasis

## Accessibility

- High contrast ratios for terminal visibility
- Color meanings reinforced with icons
- Status indicators use both color AND symbols (🟢🟡🔴)

## Themes

cc-hud-go features beautiful [Catppuccin](https://github.com/catppuccin/catppuccin) color palettes with gradient progress bars and smart adaptive layouts.

### Available Themes

**🌙 Macchiato** (default) - Dark theme with purple accents
```json
{ "theme": "macchiato" }
```

**🌑 Mocha** - Darkest variant with rich, deep colors
```json
{ "theme": "mocha" }
```

**🌆 Frappe** - Medium-dark with warmer tones
```json
{ "theme": "frappe" }
```

**☀️ Latte** - Light theme for bright environments
```json
{ "theme": "latte" }
```

### Custom Colors

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

**Available semantic colors:**
- `success` - Completed states, positive indicators (green)
- `warning` - Warnings, medium thresholds (yellow/orange)
- `danger` - Errors, high thresholds (red)
- `input` - Input tokens (blue)
- `output` - Output tokens (emerald)
- `cacheRead` - Cache read tokens (purple)
- `cacheWrite` - Cache write tokens (pink)
- `primary` - Main brand color (purple)
- `highlight` - Highlights and accents (cyan)
- `accent` - Secondary accents (orange)
- `muted` - Borders, subtle elements (gray)
- `bright` - Bright text (white/cream)
- `info` - Informational elements (teal)

### Example Configs

Pre-configured examples are available in the [`examples/`](examples/) directory:

```bash
# Copy Macchiato theme (default)
cp examples/config-macchiato.json ~/.claude/cc-hud-go/config.json

# Copy Mocha theme (darkest)
cp examples/config-mocha.json ~/.claude/cc-hud-go/config.json

# Copy Frappe theme (medium-dark)
cp examples/config-frappe.json ~/.claude/cc-hud-go/config.json

# Copy Latte theme (light)
cp examples/config-latte.json ~/.claude/cc-hud-go/config.json

# Copy custom colors example
cp examples/config-custom-colors.json ~/.claude/cc-hud-go/config.json
```

See [`examples/README.md`](examples/README.md) for detailed theme documentation and customization guide.

## Visual Features

**Gradient Progress Bars** - Smooth color transitions using Unicode block characters:
```
█▓▒░░░░░░░ 35% 📥 89k 📤 12k   ← Green (healthy)
█▓▓▓▒▒▒░░░ 75% 📥 150k 📤 25k  ← Yellow (warning)
█▓▓▓▓▓▒▒▒░ 95% 📥 190k 📤 38k  ← Red (danger)
```

**Multi-line Layout** - Clean 4-line display grouping related metrics:
```
Line 1: 🤖 Sonnet 4.5 │ █▓▒▒░░░░░░ 59%
Line 2: 📥 89k 📤 12k 💾 R:45k/W:23k ⚡ 200k │ 💰$0.0234  │  ⏱ 2m34s
Line 3: 🌿 main (dirty:2) │ 📝 +45/-12
Line 4: ╭─ 📦 App 23  🔌 MCP 2  ⚡ Skills 1 ─╮
```

## Smart Adaptive Layouts

Automatic switching between inline lipgloss boxes and table views based on configurable thresholds:

- Below threshold: Compact inline display with styled boxes (╭╮╰╯)
- Above threshold: Detailed table view with box-drawing characters (┌┬┐)
- Configurable thresholds per segment type (default: 999 for lipgloss boxes)
