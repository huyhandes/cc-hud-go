#!/bin/bash
# Secret detection hook - blocks writes containing potential secrets
# Runs as PreToolUse on Write/Edit operations

input=$(cat)

# Extract file path and content from JSON input
file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty')
content=$(echo "$input" | jq -r '.tool_input.content // .tool_input.new_string // empty')

# Skip if no content to scan
if [ -z "$content" ]; then
    echo '{"decision":"allow"}'
    exit 0
fi

# Skip scanning for certain file types
case "$file_path" in
    *.md|*.txt|*.json|*.yaml|*.yml|*.toml)
        # Config files might have placeholder patterns - be more lenient
        ;;
    *)
        ;;
esac

# Secret patterns to detect
declare -a PATTERNS=(
    # AWS
    'AKIA[0-9A-Z]{16}'
    'aws_secret_access_key\s*=\s*[A-Za-z0-9/+=]{40}'

    # Generic secrets (simplified patterns)
    # [:=]+ handles Go := and standard = assignments
    'password\s*[:=]+\s*"[^"]{8,}"'
    "password\s*[:=]+\s*'[^']{8,}'"
    'secret\s*[:=]+\s*"[^"]{8,}"'
    "secret\s*[:=]+\s*'[^']{8,}'"
    'api[_-]?key\s*[:=]+\s*"[^"]{16,}"'
    'token\s*[:=]+\s*"[^"]{20,}"'

    # Private keys
    'BEGIN (RSA|DSA|EC|OPENSSH) PRIVATE KEY'
    'BEGIN PGP PRIVATE KEY'

    # GitHub/GitLab tokens
    'gh[pousr]_[A-Za-z0-9]{36,}'
    'glpat-[A-Za-z0-9-]{20,}'

    # Slack
    'xox[baprs]-[A-Za-z0-9-]{10,}'
)

detected=()

for pattern in "${PATTERNS[@]}"; do
    if echo "$content" | grep -qiE "$pattern"; then
        # Get the matching line for context
        match=$(echo "$content" | grep -oiE ".{0,20}$pattern.{0,20}" | head -1)
        detected+=("Pattern '$pattern' matched: ...${match}...")
    fi
done

if [ ${#detected[@]} -gt 0 ]; then
    # Build JSON array of detections
    detections_json=$(printf '%s\n' "${detected[@]}" | jq -R . | jq -s .)

    cat << EOF
{
  "decision": "block",
  "reason": "Potential secrets detected in file",
  "file": "$file_path",
  "detections": $detections_json,
  "suggestion": "Remove or use environment variables for sensitive values. If this is a false positive, you can acknowledge and retry."
}
EOF
    exit 2
fi

echo '{"decision":"allow"}'
exit 0
