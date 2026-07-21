# 1. Default-only mode (no user config)

## Status

Accepted

## Context

cc-hud-go shipped a JSON config file with per-segment toggles, color overrides, and theme variants reachable only via config. Every option defaulted to "on." Nobody used it.

The `Config` struct threaded through every `Segment.Render` signature; `Segment.Enabled(cfg)` duplicated the empty-string self-gating segments already did. Maintaining it cost real complexity: the `LoadFromFile` path, the example configs, and 157 lines of `CONFIG.md` documenting options that equal the default.

## Decision

Delete the `config` package entirely. Hardcoded defaults live where they are used. No config file. Self-gating (empty-string return) is the only on/off mechanism.

## Consequences

Users cannot customize. That is the point — the tool ships opinionated and correct. Future knobs become CLI flags (see ADR 0002), not config keys. Net complexity drop is large; net flexibility drop is zero in practice.

## Alternatives Considered

- **Keep config but trim to one or two keys.** Rejected: the threading cost stays, and the surviving keys would just re-earn their deletion later.
- **Env vars.** Rejected: a statusline tool runs once per refresh; flags are the natural surface.
