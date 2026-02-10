# Release Notes - v0.2.0

**Release Date:** 2026-02-10

## 🎨 Visual Enhancements

### Catppuccin Theme System
- **4 Beautiful Themes:** Macchiato (default), Mocha, Frappe, and Latte
- **Semantic Color System:** 13 semantic colors for consistent, meaningful UI
- **Custom Color Overrides:** Override any color via config
- **Gradient Progress Bars:** Smooth color transitions for context and rate limit indicators
- **Lipgloss Integration:** Beautiful terminal UI with proper color rendering

### Enhanced Layout
- **Multi-line Layout:** Organized 4-line display grouping related metrics
- **Better Spacing:** Improved visual separation between segments
- **Inline Lipgloss Boxes:** Styled boxes for tools and tasks (replaces plain tables)
- **Color-coded Elements:** All UI elements use semantic colors for better readability

## 🐛 Bug Fixes

### Task Tracking
- **Fixed 1-based Task ID Indexing:** Claude Code uses 1-based task IDs ("1", "2", "3"), but parser was treating them as 0-based array indices, causing the first task to never be marked complete
- **Backward Compatibility:** Supports both 0-based (legacy) and 1-based (current) task IDs

### Skill Tracking
- **Skill Name Extraction:** Skills now display with full names (e.g., "superpowers:using-git-worktrees") instead of generic "Skill" counter
- **Proper Categorization:** Skills tracked separately in their own category

## 🧹 Code Quality

### Linting & Formatting
- **Zero Lint Issues:** Fixed all golangci-lint warnings
- **Code Formatting:** Applied gofmt across entire codebase
- **Deprecated API Removal:** Removed deprecated lipgloss `Copy()` calls
- **Unused Code Cleanup:** Removed unused variables and functions

### Testing
- **70+ Tests:** Comprehensive test coverage (62-100% per package)
- **New Test Coverage:** Added tests for task indexing and skill tracking fixes
- **All Tests Passing:** Full test suite verification

## 📚 Documentation

### New Documentation
- **Theme Documentation:** Complete guide to themes and customization in README
- **Example Configs:** 5 pre-configured examples for all themes + custom colors
- **Bug Fix Documentation:** Detailed analysis of fixes in docs/BUG_FIXES.md
- **Color Scheme Reference:** Semantic color system documentation

### Updated Examples
- **Config Examples:** ~/.claude/cc-hud-go/config.json examples for each theme
- **Visual Previews:** Screenshots showing different themes and layouts

## 🔧 Configuration Changes

### New Config Options
- **Theme Selection:** `theme` field (macchiato/mocha/frappe/latte)
- **Color Overrides:** `colors` map for custom color values
- **Table Thresholds:** High default thresholds (999) to prefer lipgloss inline views

### Preset Updates
- **Full Preset:** All features with beautiful themes (default)
- **Essential Preset:** Core metrics with compact layout
- **Minimal Preset:** Minimal information display

## 📦 Installation

```bash
# Build from source
git clone https://github.com/huybui/cc-hud-go.git
cd cc-hud-go
make build
make install

# Or download release binary
curl -L https://github.com/huybui/cc-hud-go/releases/download/v0.2.0/cc-hud-go -o ~/.local/bin/cc-hud-go
chmod +x ~/.local/bin/cc-hud-go
```

## 🔄 Upgrading from v0.1.0

### Breaking Changes
None - fully backward compatible

### Recommended Actions
1. Update config to use new theme system (optional)
2. Review example configs in `examples/` directory
3. Rebuild and reinstall binary: `make build && make install`

### Migration Guide
Your existing config will continue to work. To take advantage of new themes:

```json
{
  "preset": "full",
  "theme": "macchiato",
  "colors": {
    "success": "#a6e3a1",
    "warning": "#f9e2af",
    "danger": "#f38ba8"
  }
}
```

## 📊 Statistics

- **Commits:** 20+ commits since v0.1.0
- **Files Changed:** 40+ files
- **Tests Added:** 10+ new tests
- **Lines Added:** 2000+ lines
- **Lines Removed:** 500+ lines (cleanup)

## 🙏 Credits

- Inspired by [Oh My Posh Claude Segment](https://ohmyposh.dev/docs/segments/cli/claude)
- Powered by [Charm's Lipgloss](https://github.com/charmbracelet/lipgloss)
- Themed with [Catppuccin](https://github.com/catppuccin/catppuccin)

## 🐛 Known Issues

None at this time. Report issues at: https://github.com/huybui/cc-hud-go/issues

## 🚀 What's Next

Future roadmap for v0.3.0:
- Additional theme families (Nord, Dracula, etc.)
- Interactive configuration wizard
- Plugin system for custom segments
- Performance optimizations
- More granular display controls

---

**Full Changelog:** https://github.com/huybui/cc-hud-go/compare/v0.1.0...v0.2.0
