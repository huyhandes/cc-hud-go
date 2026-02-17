# Build Guide

## Prerequisites

- Go 1.24.2+
- [Just](https://github.com/casey/just) (task runner)
- Git

## Quick Install (Pre-built Binaries)

Download from the [latest release](https://github.com/huyhandes/cc-hud-go/releases/latest):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-arm64.tar.gz | tar xz
sudo mv cc-hud-go-darwin-arm64 /usr/local/bin/cc-hud-go

# macOS (Intel)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-darwin-amd64.tar.gz | tar xz
sudo mv cc-hud-go-darwin-amd64 /usr/local/bin/cc-hud-go

# Linux (amd64)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-linux-amd64.tar.gz | tar xz
sudo mv cc-hud-go-linux-amd64 /usr/local/bin/cc-hud-go

# Linux (arm64)
curl -L https://github.com/huyhandes/cc-hud-go/releases/latest/download/cc-hud-go-linux-arm64.tar.gz | tar xz
sudo mv cc-hud-go-linux-arm64 /usr/local/bin/cc-hud-go
```

Available: Linux (`amd64`, `arm64`), macOS (`amd64`, `arm64`), Windows (`amd64`, `arm64`)

## Build from Source

```bash
git clone git@github.com:huyhandes/cc-hud-go.git
cd cc-hud-go
just build          # builds with version info from git tags
just install        # installs to ~/.local/bin
```

Or with `go install`:
```bash
go install github.com/huyhandes/cc-hud-go@latest
```

## Just Commands

| Command | Description |
|---------|-------------|
| `just build` | Build with version from git tags |
| `just install` | Build and install to `~/.local/bin` |
| `just test` | Run all tests (verbose) |
| `just test-coverage` | Run tests with coverage report |
| `just check` | Format, vet, and test |
| `just fmt` | Format code |
| `just vet` | Run `go vet` |
| `just lint` | Run `golangci-lint run` |
| `just clean` | Remove build artifacts |
| `just build-all` | Build for all platforms |

## Running Tests

```bash
just test                                                         # all tests
just test-coverage                                                # with coverage report
# For targeted runs, use go test directly (just test has no arg passthrough):
go test -run TestModelSegment ./segment                           # specific test
go test -v ./segment                                              # verbose, one package
```

## Creating a Release

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions will automatically:
1. Build binaries for all platforms (Linux/macOS/Windows x amd64/arm64)
2. Create compressed archives (`.tar.gz` for Unix, `.zip` for Windows)
3. Generate SHA256 checksums
4. Publish the GitHub release with auto-generated notes

You can also trigger a release manually from the Actions tab.

## Version Information

| Build type | Version shown |
|-----------|---------------|
| Release build | Tagged version (e.g. `v0.2.0`) |
| Dev build | `git describe` output (e.g. `v0.2.0-3-gabc1234-dirty`) |
| Without git | Falls back to `dev` |
