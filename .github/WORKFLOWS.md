# GitHub Workflows

This document describes the GitHub Actions workflows configured for this project.

## Release Workflow (`release.yml`)

**Trigger:** Push of version tags (e.g., `v1.0.0`, `v2.1.3`)

**Purpose:** Automated multi-platform builds and GitHub releases.

### What it does:

1. **Multi-platform Build**
   - Builds binaries for:
     - Linux: amd64, arm64
     - macOS (Darwin): amd64, arm64
     - Windows: amd64, arm64
   - Uses `-ldflags` to strip debug info and embed version
   - Creates compressed archives:
     - `.tar.gz` for Unix systems
     - `.zip` for Windows

2. **Checksums**
   - Generates SHA256 checksums for all artifacts
   - Outputs to `checksums.txt`

3. **GitHub Release**
   - Creates a new release with auto-generated release notes
   - Uploads all binaries and checksums as release assets
   - Not marked as draft or prerelease by default

### Usage:

```bash
# Create and push a version tag
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

You can also manually trigger this workflow from the Actions tab.

### Version Embedding:

The workflow embeds the version tag into the binary using:
```bash
go build -ldflags="-s -w -X main.version=v1.0.0"
```

Make sure your `main.go` has a version variable to capture this:
```go
var version = "dev"
```

---

## CI Workflow (`ci.yml`)

**Trigger:**
- Push to `main` branch
- Pull requests targeting `main` branch

**Purpose:** Continuous integration testing and code quality checks.

### Jobs:

#### 1. **Test Job**
- **Matrix Testing** across:
  - OS: Ubuntu, macOS, Windows
  - Go version: 1.24.2
- Runs all tests with race detection
- Generates coverage reports
- Uploads coverage to Codecov (optional)
- Uses Go module caching for faster builds

#### 2. **Lint Job**
- Runs `golangci-lint` with latest version
- Checks code style and common issues
- Linux-only (faster, consistent results)

#### 3. **Build Job**
- Verifies the code builds successfully
- Quick smoke test before merge

### Local Testing:

Before pushing, you can run the same checks locally:

```bash
# Run tests
go test -v -race -coverprofile=coverage.txt ./...

# Run linter
golangci-lint run

# Build
go build -v ./...
```

---

## Required Secrets

### Optional Secrets:

- `CODECOV_TOKEN` - For uploading coverage reports to Codecov
  - The workflow will continue even if this fails
  - Only needed if you want coverage tracking

### Automatic Secrets:

- `GITHUB_TOKEN` - Automatically provided by GitHub Actions
  - Used for creating releases
  - No configuration needed

---

## Workflow Status Badges

Add these to your README to show build status:

```markdown
[![CI](https://github.com/huyhandes/cc-hud-go/actions/workflows/ci.yml/badge.svg)](https://github.com/huyhandes/cc-hud-go/actions/workflows/ci.yml)
[![Release](https://github.com/huyhandes/cc-hud-go/actions/workflows/release.yml/badge.svg)](https://github.com/huyhandes/cc-hud-go/actions/workflows/release.yml)
```

---

## Troubleshooting

### Release workflow not triggering

- Ensure you're pushing tags, not just creating them locally:
  ```bash
  git push origin v1.0.0
  ```
- Check that the tag matches the pattern `v*` (must start with 'v')

### Build failures

- Check Go version compatibility (should be 1.24.2+)
- Ensure all dependencies are in `go.mod`
- Test locally with: `go build -o test-binary .`

### Permission errors in release

- The workflow needs `contents: write` permission
- This is configured in the workflow file
- If using organization settings, ensure Actions have write access

---

## Customization

### Adding new platforms

Edit `.github/workflows/release.yml` and add to the `platforms` array:

```bash
platforms=(
  "linux/amd64"
  "linux/arm64"
  "your/platform"  # Add here
)
```

### Changing Go version

Update the `go-version` in both workflow files:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.25.0'  # Update here
```

### Adding pre-release tags

To create pre-release versions (alpha, beta, rc):

```bash
git tag -a v1.0.0-beta.1 -m "Beta release"
git push origin v1.0.0-beta.1
```

Modify the workflow to mark these as pre-release:

```yaml
- name: Create Release
  uses: softprops/action-gh-release@v2
  with:
    prerelease: ${{ contains(github.ref, '-') }}
```
