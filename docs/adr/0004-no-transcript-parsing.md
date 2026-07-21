# 4. Transcript JSONL no longer parsed

## Status

Accepted

## Context

The transcript parser scanned a potentially large JSONL file on every statusline refresh to populate tool-call counts, skill usage, and task progress. Those three displays are being removed (no longer carry their weight).

## Decision

Delete `parser/transcript.go`, `parser/tool.go`, `parser/task.go` and their tests. `parser/` retains only `stdin.go` + `stdin_test.go`. The `TranscriptPath` field is removed from `SessionInfo` and from `StdinData` — nothing consumes it.

## Consequences

Every statusline refresh skips a file scan. The `parser` package shrinks to a single concern. No feature is lost because the data it produced is no longer displayed.

## Alternatives Considered

- **Keep parsing but cache.** Rejected: the consumers are gone; caching dead data is pure cost.
- **Keep the parser as a library for future use.** Rejected: YAGNI; reintroducing it is straightforward if a future display needs it.
