# 3. OAuth always fetched (no silent stdin rate-limit fallback)

## Status

Accepted

## Context

Rate-limit data previously had two sources: a `rate_limits` block in stdin (a stale snapshot Claude Code passed in) and the OAuth API (fresh). When both were wired, a stale stdin value could disagree with fresh OAuth data. The old code gated the OAuth call behind `cfg.Display.FetchOAuth`.

## Decision

OAuth is always fetched on every refresh. On failure it silently no-ops and the rate-limit segments render empty (the established `FiveHourSegment` empty-return pattern). The stdin `rate_limits` block is no longer parsed. OAuth is the single source of truth.

## Consequences

Rate-limit data is always fresh when OAuth is reachable, always absent (never stale) when it isn't. One code path, one source of truth. Adds one HTTP call per statusline refresh — acceptable for a tool that refreshes on user interaction, not in a tight loop.

## Alternatives Considered

- **Keep stdin fallback.** Rejected: stale/wrong data is worse than no data on a statusline.
- **Cache OAuth result for N seconds.** Rejected as out of scope; the call is cheap enough.
