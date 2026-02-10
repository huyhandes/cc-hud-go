#!/bin/bash
# Smart Go validation - only test packages with changes
# Returns structured JSON feedback for Claude

set -o pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$PROJECT_DIR" || exit 0

# Skip if not a Go project
[ ! -f "go.mod" ] && exit 0

# Get changed Go files (staged + unstaged)
CHANGED=$(git diff --name-only HEAD 2>/dev/null | grep '\.go$' || true)
STAGED=$(git diff --cached --name-only 2>/dev/null | grep '\.go$' || true)
ALL_CHANGED=$(echo -e "$CHANGED\n$STAGED" | sort -u | grep -v '^$' || true)

# If no Go files changed, run quick vet only
if [ -z "$ALL_CHANGED" ]; then
    echo '{"status":"skipped","reason":"No Go files changed"}'
    exit 0
fi

# Extract unique packages from changed files
PACKAGES=$(echo "$ALL_CHANGED" | xargs -I{} dirname {} 2>/dev/null | sort -u | sed 's|^|./|' | tr '\n' ' ')

if [ -z "$PACKAGES" ]; then
    echo '{"status":"skipped","reason":"No packages to test"}'
    exit 0
fi

errors=()
warnings=()

# Run go vet on affected packages
echo "Running go vet on: $PACKAGES" >&2
vet_output=$(go vet $PACKAGES 2>&1)
vet_exit=$?

if [ $vet_exit -ne 0 ]; then
    errors+=("go vet failed: $vet_output")
fi

# Run go test with short flag on affected packages
echo "Running go test -short on: $PACKAGES" >&2
test_output=$(go test $PACKAGES -short -timeout 60s 2>&1)
test_exit=$?

if [ $test_exit -ne 0 ]; then
    # Extract failed test names
    failed_tests=$(echo "$test_output" | grep -E "^--- FAIL:" | sed 's/--- FAIL: //' || true)
    errors+=("go test failed: $failed_tests")
fi

# Build JSON response
if [ ${#errors[@]} -gt 0 ]; then
    # Convert errors array to JSON
    error_json=$(printf '%s\n' "${errors[@]}" | jq -R . | jq -s .)
    cat << EOF
{
  "status": "failed",
  "packages_tested": "$PACKAGES",
  "errors": $error_json,
  "suggestion": "Fix the failing tests/vet issues in the affected packages before proceeding"
}
EOF
    exit 2
fi

cat << EOF
{
  "status": "passed",
  "packages_tested": "$PACKAGES",
  "message": "All affected packages pass vet and tests"
}
EOF
exit 0
