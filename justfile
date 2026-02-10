# Get version from git tags or use dev
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
binary_name := "cc-hud-go"

# Show available recipes
default:
    @just --list

# Build the binary with version info
build:
    @echo "Building {{binary_name}} {{version}}..."
    go build -ldflags "-X github.com/huyhandes/cc-hud-go/version.Version={{version}}" -o {{binary_name}} .
    @echo "Build complete: ./{{binary_name}}"

# Build and install to ~/.local/bin
install: build
    @echo "Installing {{binary_name}} {{version}}..."
    @mkdir -p ~/.local/bin
    @cp {{binary_name}} ~/.local/bin/{{binary_name}}
    @echo "Installed to ~/.local/bin/{{binary_name}}"

# Build and run the application
run: build
    ./{{binary_name}}

# Run all tests
test:
    go test -v ./...

# Run tests with coverage report
test-coverage:
    go test -cover ./...

# Remove built binaries
clean:
    @echo "Cleaning..."
    rm -f {{binary_name}}
    rm -rf dist/
    @echo "Clean complete"

# Show version information
show-version:
    @go run -ldflags "-X github.com/huyhandes/cc-hud-go/version.Version={{version}}" . --version

# Format code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run golangci-lint (if installed)
lint:
    #!/usr/bin/env bash
    if command -v golangci-lint >/dev/null 2>&1; then
        golangci-lint run
    else
        echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi

# Run fmt, vet, and test
check: fmt vet test
    @echo "All checks passed!"

# Build for all platforms (mimics release)
build-all:
    #!/usr/bin/env bash
    echo "Building for all platforms..."
    mkdir -p dist
    for platform in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/arm64"; do
        GOOS=${platform%/*}
        GOARCH=${platform#*/}
        output="dist/{{binary_name}}-$GOOS-$GOARCH"
        [ "$GOOS" = "windows" ] && output="$output.exe"
        echo "Building for $GOOS/$GOARCH..."
        env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X github.com/huyhandes/cc-hud-go/version.Version={{version}}" -o "$output" .
    done
    echo "All builds complete in dist/"
