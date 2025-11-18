.PHONY: build clean install test help deps release

# Build configuration
BINARY_NAME=android-offline
VERSION?=1.0.0
BUILD_DIR=bin
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

build: deps ## Build for current platform
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/android-offline

build-all: deps ## Build for all platforms
	@echo "Building for all platforms..."
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			output="$(BUILD_DIR)/$(BINARY_NAME)-$$platform-$$arch"; \
			if [ "$$platform" = "windows" ]; then output="$$output.exe"; fi; \
			echo "Building $$platform/$$arch..."; \
			GOOS=$$platform GOARCH=$$arch $(GOBUILD) $(LDFLAGS) -o $$output ./cmd/android-offline || exit 1; \
		done; \
	done
	@echo "Build complete! Binaries in $(BUILD_DIR)/"

release: build-all ## Create release packages
	@echo "Creating release packages..."
	@mkdir -p releases
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			name="$(BINARY_NAME)-$(VERSION)-$$platform-$$arch"; \
			binary="$(BUILD_DIR)/$(BINARY_NAME)-$$platform-$$arch"; \
			if [ "$$platform" = "windows" ]; then binary="$$binary.exe"; fi; \
			if [ -f "$$binary" ]; then \
				tar -czf "releases/$$name.tar.gz" -C $(BUILD_DIR) $$(basename $$binary) README.md LICENSE; \
				echo "Created releases/$$name.tar.gz"; \
			fi; \
		done; \
	done

install: build ## Install binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete! Run '$(BINARY_NAME) --help' to get started"

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) releases
	@echo "Clean complete"

run: build ## Build and run
	$(BUILD_DIR)/$(BINARY_NAME) --help

fmt: ## Format Go code
	$(GOCMD) fmt ./...

lint: ## Run linter
	golangci-lint run ./...

.DEFAULT_GOAL := help
