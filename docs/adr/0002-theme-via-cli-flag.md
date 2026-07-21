# 2. Theme selection via `--theme` CLI flag

## Status

Accepted

## Context

The only customization users actually wanted was picking a Catppuccin variant to match their terminal. It was buried in the config file.

## Decision

Move theme selection to a `--theme` CLI flag (default `macchiato`; values `macchiato` / `mocha` / `frappe` / `latte`; unknown values fall back to macchiato via `theme.GetTheme`'s default branch). The user's `statusline.command` in their Claude Code config is itself the configuration.

## Consequences

Theme is set once in the Claude Code config, never re-edited. Typos never crash the statusline (graceful fallback). `theme.LoadThemeFromConfig` wrapper is kept (thin, honest) but the color-overrides codepath becomes unreachable in practice.

## Alternatives Considered

- **Env var `CC_HUD_THEME`.** Rejected: a flag is more discoverable and lives where the command is declared.
- **Hardcode macchiato only.** Rejected: users have real terminal-scheme variety; the four variants cost nothing at runtime.
