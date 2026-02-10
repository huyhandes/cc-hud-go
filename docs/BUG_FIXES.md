# Bug Fixes - Session Task and Skill Tracking

## Issues Reported

After auto-compaction at 95% context usage, three issues were observed:

### 1. Context Auto-Compact Threshold (85% vs 95%) ✅ Expected Behavior

**User observation:** Context bar shows 85% when compaction happened, but docs say it triggers at 95%.

**Explanation:** This is **correct behavior**. The docs are accurate:
- Compaction **triggers** when context reaches 95%
- After compaction frees up space, usage drops to ~85%
- The statusline shows the **current** state (85% after compaction), not the trigger point

**No fix needed** - working as designed.

### 2. TODO Task Showing 1 Pending ✅ Fixed

**User observation:** TaskList shows "No tasks found" but statusline shows 1 pending task.

**Root cause:** Off-by-one error in task ID indexing. Claude Code uses **1-based task IDs** ("1", "2", "3"), but parser was treating them as **0-based array indices**.

**Impact:**
- 18 tasks created (indices 0-17)
- TaskUpdate with taskId="1" updated index **1** (task #2) instead of index **0** (task #1)
- TaskUpdate with taskId="18" tried index **18** (out of bounds!) and failed
- Result: Task at index 0 stayed pending

**Fix in parser.go:305-318:**
```go
// Find task by ID or index
index := -1
if idx, ok := tracker.TaskIDMap[taskID]; ok {
    index = idx
} else if idxNum, err := strconv.Atoi(taskID); err == nil {
    // Support both 0-based (backward compat) and 1-based (Claude Code standard) indexing
    // Claude Code uses 1-based task IDs ("1", "2", "3")
    // For backward compatibility, also support taskId=0 as direct index 0
    if idxNum == 0 && idxNum < len(tracker.Tasks) {
        // taskId=0: treat as 0-based index 0
        index = 0
    } else if idxNum >= 1 && idxNum <= len(tracker.Tasks) {
        // taskId >= 1: treat as 1-based task ID, convert to 0-based array index
        index = idxNum - 1
    }
}
```

**Testing:**
- Added `TestParseTaskUpdate1BasedTaskID` to verify fix
- Updated `TestParseTranscriptWithMixedOperations` to use 1-based task IDs
- All 70+ tests passing

### 3. Skills Not Showing ✅ Fixed

**User observation:** Skills section not showing activated skills (only showed generic "Skill" counter).

**Root cause:** Parser detected `Skill` tool calls but didn't extract the actual skill name from `input.skill` parameter.

**Before:**
```go
case CategorySkill:
    s.Tools.AppTools["Skill"]++  // Generic counter
```

**After in parser.go:430-438:**
```go
case CategorySkill:
    // Extract skill name from input parameters
    if skillName, ok := block.Input["skill"].(string); ok && skillName != "" {
        usage := s.Tools.Skills[skillName]
        usage.Count++
        s.Tools.Skills[skillName] = usage
    } else {
        // Fallback for skills without name
        s.Tools.AppTools["Skill"]++
    }
```

**Result:** Skills now show with their full names:
- `superpowers:using-git-worktrees: 1`
- `superpowers:brainstorming: 2`
- etc.

**Testing:**
- Added `TestParseTranscriptLineSkillTracking`
- Added `TestParseTranscriptLineSkillTrackingMultiple`
- Added `TestParseTranscriptLineSkillTrackingFallback`

## Files Changed

- `parser/parser.go` - Task indexing fix and skill tracking
- `parser/tasks_test.go` - Added test for 1-based task ID fix
- `parser/transcript_test.go` - Added tests for skill tracking

## Build Info

```bash
make build
cp cc-hud-go ~/.local/bin/
```

Version: v0.1.0-34-gbc71761-dirty

## Verification

After installing the fixed binary, the statusline should now:
1. Show correct context percentage (85% after compaction is expected)
2. Show 0 pending tasks (all 18 tasks now correctly marked completed)
3. Display skill names instead of generic "Skill" counter

## Summary

- ✅ Issue #1: Expected behavior, no fix needed
- ✅ Issue #2: Fixed 1-based task ID indexing bug
- ✅ Issue #3: Fixed skill name extraction

All tests passing (70+ tests, 62-100% coverage).
