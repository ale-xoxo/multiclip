# MultiClip - macOS Menu Bar Clipboard Manager
# Makefile for building and managing the application

# Variables
APP_NAME = multiclip
SRC_DIR = multiclip/src
BUILD_DIR = multiclip/build
ASSETS_DIR = multiclip/assets
MAIN_FILE = $(SRC_DIR)/main.go
MODULE_DIR = multiclip
BINARY = $(BUILD_DIR)/$(APP_NAME)
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -s -w -X main.Version=$(VERSION)

# Go build settings
export CGO_ENABLED = 1
GOOS = darwin
GOARCH = amd64

# Default target
.PHONY: all
all: clean build-app

# Build the application
.PHONY: build build-app
build-app: $(BUILD_DIR)
	@echo "Building $(APP_NAME) for macOS..."
	cd $(MODULE_DIR) && GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="$(LDFLAGS)" \
		-o ../$(BINARY) \
		src/main.go
	@echo "✅ Build complete: $(BINARY)"

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run the application
.PHONY: run
run: build-app
	@echo "Starting $(APP_NAME)..."
	@./$(BINARY)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing Go dependencies..."
	cd $(MODULE_DIR) && go mod tidy
	cd $(MODULE_DIR) && go mod download
	@echo "✅ Dependencies installed"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
.PHONY: lint
lint:
	@echo "Linting Go code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, using go vet instead"; \
		go vet ./...; \
	fi
	@echo "✅ Linting complete"

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✅ Tests complete"

# Create macOS app bundle
.PHONY: bundle
bundle: build-app
	@echo "Creating macOS app bundle..."
	@mkdir -p $(BUILD_DIR)/MultiClip.app/Contents/{MacOS,Resources}
	@cp $(BINARY) $(BUILD_DIR)/MultiClip.app/Contents/MacOS/
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '<plist version="1.0">' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '<dict>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <key>CFBundleExecutable</key>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <string>$(APP_NAME)</string>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <key>CFBundleIdentifier</key>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <string>com.example.multiclip</string>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <key>CFBundleName</key>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <string>MultiClip</string>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <key>CFBundleVersion</key>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <string>$(VERSION)</string>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <key>LSUIElement</key>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '  <string>1</string>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '</dict>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo '</plist>' >> $(BUILD_DIR)/MultiClip.app/Contents/Info.plist
	@echo "✅ App bundle created: $(BUILD_DIR)/MultiClip.app"

# Build for release (optimized)
.PHONY: release
release: clean fmt lint test
	@echo "Building release version..."
	@$(MAKE) build-app
	@$(MAKE) bundle
	@echo "✅ Release build complete"

# Alias for build-app
.PHONY: build
build: build-app

# Install to /Applications (requires sudo)
.PHONY: install
install: bundle
	@echo "Installing MultiClip to /Applications..."
	@sudo cp -R $(BUILD_DIR)/MultiClip.app /Applications/
	@echo "✅ MultiClip installed to /Applications"

# Uninstall from /Applications (requires sudo)
.PHONY: uninstall
uninstall:
	@echo "Uninstalling MultiClip from /Applications..."
	@sudo rm -rf /Applications/MultiClip.app
	@echo "✅ MultiClip uninstalled"

# Development mode - watch for changes and rebuild
.PHONY: dev
dev:
	@echo "Starting development mode (requires fswatch)..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o $(SRC_DIR) | while read f; do \
			echo "File changed, rebuilding..."; \
			$(MAKE) build-app; \
		done; \
	else \
		echo "❌ fswatch not installed. Install with: brew install fswatch"; \
		exit 1; \
	fi

# Show build information
.PHONY: info
info:
	@echo "📋 MultiClip Build Information"
	@echo "==============================="
	@echo "App Name:     $(APP_NAME)"
	@echo "Version:      $(VERSION)"
	@echo "Source Dir:   $(SRC_DIR)"
	@echo "Build Dir:    $(BUILD_DIR)"
	@echo "Binary:       $(BINARY)"
	@echo "GOOS:         $(GOOS)"
	@echo "GOARCH:       $(GOARCH)"
	@echo "CGO_ENABLED:  $(CGO_ENABLED)"
	@echo "LDFLAGS:      $(LDFLAGS)"

# Help target
.PHONY: help
help:
	@echo "📋 MultiClip Makefile Commands"
	@echo "=============================="
	@echo "build       - Build the application binary"
	@echo "clean       - Remove build artifacts"
	@echo "run         - Build and run the application"
	@echo "deps        - Install Go dependencies"
	@echo "fmt         - Format Go code"
	@echo "lint        - Lint Go code"
	@echo "test        - Run tests"
	@echo "bundle      - Create macOS app bundle"
	@echo "release     - Build optimized release version"
	@echo "install     - Install to /Applications (requires sudo)"
	@echo "uninstall   - Remove from /Applications (requires sudo)"
	@echo "dev         - Watch for changes and rebuild (requires fswatch)"
	@echo "info        - Show build information"
	@echo "help        - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build the app"
	@echo "  make run            # Build and run"
	@echo "  make release        # Full release build"
	@echo "  make install        # Install to Applications"