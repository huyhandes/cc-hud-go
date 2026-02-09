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
