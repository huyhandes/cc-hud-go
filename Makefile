.PHONY: build clean test install run help version

# Get version from git tags or use dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/huybui/cc-hud-go/version.Version=$(VERSION)"
BINARY_NAME := cc-hud-go

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary with version info
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

## install: Build and install to ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) $(VERSION)..."
	@mkdir -p ~/.local/bin
	@cp $(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@echo "Installed to ~/.local/bin/$(BINARY_NAME)"

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## test: Run all tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -cover ./...

## clean: Remove built binaries
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "Clean complete"

## version: Show version information
version:
	@go run $(LDFLAGS) . --version

## fmt: Format code
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## lint: Run golangci-lint (if installed)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## check: Run fmt, vet, and test
check: fmt vet test
	@echo "All checks passed!"

## build-all: Build for all platforms (mimics release)
build-all:
	@echo "Building for all platforms..."
	@mkdir -p dist
	@for platform in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/arm64"; do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output="dist/$(BINARY_NAME)-$$GOOS-$$GOARCH"; \
		[ "$$GOOS" = "windows" ] && output="$$output.exe"; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o "$$output" .; \
	done
	@echo "All builds complete in dist/"

.DEFAULT_GOAL := help
