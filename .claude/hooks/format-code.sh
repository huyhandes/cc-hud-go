#!/bin/bash
# Format code hook - runs after Write/Edit operations
# Formats Go files with gofmt/goimports, Python files with ruff
# Returns structured JSON feedback

input=$(cat)

# Extract the file path from the JSON input
file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty')

if [ -z "$file_path" ]; then
    echo '{"status":"skipped","reason":"No file path provided"}'
    exit 0
fi

# Check if file exists
if [ ! -f "$file_path" ]; then
    echo '{"status":"skipped","reason":"File does not exist"}'
    exit 0
fi

result=""
formatted=false

# Format Go files
if [[ "$file_path" == *.go ]]; then
    # Prefer goimports if available (handles imports + formatting)
    if command -v goimports &>/dev/null; then
        before=$(cat "$file_path")
        goimports -w "$file_path" 2>/dev/null
        after=$(cat "$file_path")
        if [ "$before" != "$after" ]; then
            formatted=true
            result="goimports"
        fi
    elif command -v gofmt &>/dev/null; then
        before=$(cat "$file_path")
        gofmt -w "$file_path" 2>/dev/null
        after=$(cat "$file_path")
        if [ "$before" != "$after" ]; then
            formatted=true
            result="gofmt"
        fi
    fi
fi

# Format Python files with ruff
if [[ "$file_path" == *.py ]]; then
    RUFF="${RUFF:-/home/coder/.local/share/uv/tools/ruff/bin/ruff}"
    if [ -x "$RUFF" ] || command -v ruff &>/dev/null; then
        RUFF="${RUFF:-ruff}"
        before=$(cat "$file_path")
        "$RUFF" format "$file_path" 2>/dev/null
        "$RUFF" check "$file_path" --fix 2>/dev/null
        after=$(cat "$file_path")
        if [ "$before" != "$after" ]; then
            formatted=true
            result="ruff"
        fi
    fi
fi

# Format TypeScript/JavaScript with prettier
if [[ "$file_path" == *.ts || "$file_path" == *.tsx || "$file_path" == *.js || "$file_path" == *.jsx ]]; then
    if command -v prettier &>/dev/null; then
        before=$(cat "$file_path")
        prettier --write "$file_path" 2>/dev/null
        after=$(cat "$file_path")
        if [ "$before" != "$after" ]; then
            formatted=true
            result="prettier"
        fi
    fi
fi

if [ "$formatted" = true ]; then
    cat << EOF
{
  "status": "formatted",
  "file": "$file_path",
  "formatter": "$result",
  "message": "File was auto-formatted with $result"
}
EOF
else
    cat << EOF
{
  "status": "unchanged",
  "file": "$file_path",
  "message": "No formatting changes needed"
}
EOF
fi

exit 0
