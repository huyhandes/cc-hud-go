# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-02-10

### Added
- **Catppuccin Theme System**: 4 beautiful themes (Macchiato, Mocha, Frappe, Latte)
- **Semantic Color Palette**: 13 semantic colors for consistent UI
- **Gradient Progress Bars**: Smooth color transitions for context and rate limit indicators
- **Theme Configuration**: `theme` and `colors` fields in config
- **Skill Name Tracking**: Skills now show with full names (e.g., "superpowers:using-git-worktrees")
- **Example Configs**: 5 pre-configured examples for all themes
- **Comprehensive Documentation**: Theme guide, color scheme reference, bug fix analysis
- **Test Coverage**: 10+ new tests for task indexing and skill tracking

### Fixed
- **Task ID Indexing Bug**: Fixed off-by-one error where 1-based task IDs were treated as 0-based array indices
- **Skill Tracking**: Parser now extracts skill names from input parameters instead of using generic "Skill" counter
- **Gradient Colors**: Restored proper gradient bar coloring with theme integration
- **Layout Issues**: Improved multi-line layout with better spacing and organization

### Changed
- **Multi-line Layout**: Reorganized into clean 4-line display grouping related metrics
- **Table Thresholds**: Increased defaults to 999 to prefer lipgloss inline views over plain tables
- **Context Display**: Removed traffic light icon for cleaner look
- **Code Quality**: Zero golangci-lint issues, full gofmt formatting
- **Deprecated API**: Removed lipgloss `Copy()` calls (use direct style methods instead)

### Removed
- Unused `currentTheme` variable in style package
- Traffic light icon (🟢🟡🔴) from context segment

## [0.1.0] - 2025-01-XX

### Added
- Initial release of cc-hud-go
- Basic statusline integration with Claude Code
- Segment system architecture
  - Model and plan type display
  - Context window tracking with token usage
  - Git branch, status, and file statistics
  - Cost tracking and session duration
  - Tool usage categorization (App/MCP/Custom/Skills)
  - Task progress tracking
  - Active agent display
  - API rate limit monitoring
- Configuration system with three presets (Full, Essential, Minimal)
- Dual parsing system (stdin JSON + transcript JSONL)
- State management with derived field calculation
- Git integration via command execution
- Comprehensive test suite (60+ tests)
- Make-based build system
- Version detection from git tags

### Features
- **Segments**: Modular display components
- **State Tracking**: Centralized session state management
- **Configuration**: Preset modes with granular control
- **Parser System**: Dual approach for stdin and transcript
- **Testing**: TDD approach with high coverage

---

## Release Links

- [0.2.0] - https://github.com/huyhandes/cc-hud-go/releases/tag/v0.2.0
- [0.1.0] - https://github.com/huyhandes/cc-hud-go/releases/tag/v0.1.0

## Upgrade Notes

### Upgrading to 0.2.0 from 0.1.0

**Breaking Changes**: None - fully backward compatible

**Recommended Actions**:
1. Update config to use new theme system (optional)
2. Review example configs in `examples/` directory
3. Rebuild and reinstall: `make build && make install`

**New Features to Try**:
```json
{
  "preset": "full",
  "theme": "macchiato",
  "colors": {
    "success": "#a6e3a1",
    "warning": "#f9e2af"
  }
}
```

## Contribution Guidelines

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Support

- **Documentation**: https://github.com/huyhandes/cc-hud-go
- **Issues**: https://github.com/huyhandes/cc-hud-go/issues
- **Discussions**: https://github.com/huyhandes/cc-hud-go/discussions
