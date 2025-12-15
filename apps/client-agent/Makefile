.PHONY: all build clean install deps build-all build-linux build-macos build-windows test

# Variables
BINARY_NAME=phd-client-agent
VERSION=1.0.0
BUILD_DIR=build
GO=go
MAIN_PATH=./cmd/agent

# Default target
all: clean deps build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Build for current platform
build:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all: build-linux build-macos build-windows
	@echo "All builds complete!"

# Build for Linux (amd64)
build-linux:
	@echo "Building for Linux amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

# Build for Linux (arm64)
build-linux-arm64:
	@echo "Building for Linux arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "Linux ARM64 build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64"

# Build for macOS (amd64 - Intel)
build-macos:
	@echo "Building for macOS amd64 (Intel)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "macOS Intel build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64"

# Build for macOS (arm64 - Apple Silicon)
build-macos-arm64:
	@echo "Building for macOS arm64 (Apple Silicon)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "macOS Apple Silicon build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64"

# Build for Windows (amd64)
build-windows:
	@echo "Building for Windows amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Build with version info
build-release:
	@echo "Building release version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)/release
	@$(MAKE) build-linux && mv $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BUILD_DIR)/release/
	@$(MAKE) build-macos && mv $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/release/
	@$(MAKE) build-macos-arm64 && mv $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(BUILD_DIR)/release/
	@$(MAKE) build-windows && mv $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BUILD_DIR)/release/
	@echo "Release builds complete in $(BUILD_DIR)/release/"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race ./...

# Run
run:
	@echo "Running client agent..."
	$(GO) run $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# Display help
help:
	@echo "PHD Client Agent - Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              - Clean, install deps, and build for current platform"
	@echo "  deps             - Install Go dependencies"
	@echo "  build            - Build for current platform"
	@echo "  build-all        - Build for all platforms (Linux, macOS, Windows)"
	@echo "  build-linux      - Build for Linux amd64"
	@echo "  build-macos      - Build for macOS amd64 (Intel)"
	@echo "  build-macos-arm64 - Build for macOS arm64 (Apple Silicon)"
	@echo "  build-windows    - Build for Windows amd64"
	@echo "  build-release    - Build release versions for all platforms"
	@echo "  test             - Run tests"
	@echo "  run              - Run the application"
	@echo "  clean            - Remove build artifacts"
	@echo "  install          - Install binary to /usr/local/bin"
	@echo "  help             - Display this help message"
