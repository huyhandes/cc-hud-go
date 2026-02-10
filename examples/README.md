# Example Configurations

This directory contains example configuration files for cc-hud-go demonstrating different themes and customization options.

## Using Example Configs

Copy an example config to your Claude Code config directory:

```bash
# Copy Macchiato theme (default)
cp examples/config-macchiato.json ~/.claude/cc-hud-go/config.json

# Copy Mocha theme (darker variant)
cp examples/config-mocha.json ~/.claude/cc-hud-go/config.json

# Copy Frappe theme (medium variant)
cp examples/config-frappe.json ~/.claude/cc-hud-go/config.json

# Copy Latte theme (light variant)
cp examples/config-latte.json ~/.claude/cc-hud-go/config.json

# Copy custom colors example
cp examples/config-custom-colors.json ~/.claude/cc-hud-go/config.json
```

## Available Example Configs

### `config-macchiato.json`
Catppuccin Macchiato theme (default). A beautiful dark theme with purple accents and excellent contrast.

**Best for:** Dark terminal backgrounds, extended coding sessions

### `config-mocha.json`
Catppuccin Mocha theme. The darkest variant with deep, rich colors and high contrast.

**Best for:** Very dark terminal backgrounds, OLED displays, low-light environments

### `config-frappe.json`
Catppuccin Frappe theme. Medium-dark theme with softer, warmer tones.

**Best for:** Dark terminals with slight warmth preference

### `config-latte.json`
Catppuccin Latte theme. Light theme optimized for daylight use.

**Best for:** Light terminal backgrounds, bright environments, daytime coding

### `config-custom-colors.json`
Example showing custom color overrides. Demonstrates how to override specific theme colors while keeping the base theme intact.

**Customize these colors:**
- `success` - Green colors (completed states, positive indicators)
- `warning` - Yellow/orange colors (warnings, medium thresholds)
- `danger` - Red colors (errors, high thresholds)
- `input` - Blue for input tokens
- `output` - Emerald for output tokens
- `cacheRead` - Purple for cache reads
- `cacheWrite` - Pink for cache writes
- `primary` - Main brand color
- `highlight` - Cyan for highlights
- `accent` - Orange accents
- `muted` - Gray for borders and subtle elements
- `bright` - Bright text color
- `info` - Teal for informational elements

## Theme Features

All themes include:

✨ **Visual Enhancements**
- Gradient progress bars with smooth color transitions (█▓▒░)
- Semantic color coding for different metric types
- Contextual table rendering based on data volume
- Enhanced spacing with clean separators (│)

🎨 **Catppuccin Colors**
- Carefully curated color palettes
- Excellent contrast ratios
- Consistent semantic meaning across themes

📊 **Smart Adaptive Display**
- Automatic switching between inline and table views
- Threshold-based rendering for tools and tasks
- Filtered display (e.g., last 3 completed tasks)

## Creating Your Own Theme

To create a custom theme, start with `config-custom-colors.json` and modify the `colors` object:

```json
{
  "theme": "macchiato",  // Base theme to start from
  "colors": {
    "success": "#your-hex-color",
    "primary": "#your-hex-color"
    // ... override any semantic colors
  }
}
```

All color values should be hex codes in the format `#RRGGBB`.

## Configuration Options

All example configs use the "full" preset with:
- All display segments enabled
- Expanded multi-line layout
- Git file statistics
- Tool grouping by category
- Task progress tracking
- Rate limit monitoring

To customize further, refer to the main [Configuration Guide](../README.md#configuration).
