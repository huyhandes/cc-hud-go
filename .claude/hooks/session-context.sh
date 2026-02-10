#!/bin/bash
# Session context injection - provides Claude with project state on start

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$PROJECT_DIR" || exit 0

# Collect git info
branch=$(git branch --show-current 2>/dev/null || echo "unknown")
recent_commits=$(git log --oneline -5 --format='%h %s' 2>/dev/null | head -5 || echo "")
uncommitted=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
last_commit_time=$(git log -1 --format='%cr' 2>/dev/null || echo "unknown")

# Check for failing tests (quick scan of test cache)
# Only run if there's a test cache, otherwise skip
failing_packages=""
if [ -f "go.mod" ]; then
    # Quick check: look for recent test failures in go test cache
    # This is fast because it uses cached results
    failing_packages=$(go test ./... -json -count=1 -timeout 30s 2>/dev/null | \
        jq -r 'select(.Action=="fail" and .Test==null) | .Package' 2>/dev/null | \
        sort -u | head -5 | tr '\n' ', ' || echo "")
fi

# Check for TODO/FIXME in recently changed files
todos=""
if [ "$uncommitted" -gt 0 ]; then
    todos=$(git diff --name-only 2>/dev/null | head -10 | xargs grep -l -E 'TODO|FIXME' 2>/dev/null | tr '\n' ', ' || echo "")
fi

# Get open tasks from any task tracking
open_issues=""
if command -v gh &>/dev/null && [ -d ".git" ]; then
    open_issues=$(gh issue list --limit 3 --json number,title 2>/dev/null | jq -r '.[] | "#\(.number): \(.title)"' 2>/dev/null | tr '\n' '; ' || echo "")
fi

cat << EOF
{
  "project_context": {
    "branch": "$branch",
    "last_commit": "$last_commit_time",
    "uncommitted_files": $uncommitted,
    "recent_commits": [
$(echo "$recent_commits" | while read -r line; do echo "      \"$line\","; done | sed '$ s/,$//')
    ],
    "failing_packages": "${failing_packages%,}",
    "files_with_todos": "${todos%,}",
    "open_issues": "$open_issues"
  }
}
EOF
