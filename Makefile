# Makefile for hello-tool-base
# Incorporates features for robust local development and CI/CD integration.

# Specify phony targets (targets not associated with files)
.PHONY: all tree build clean deps fmt lint golangci-lint test test-debug \
        install-tools check-gomod check-line-length check deploy help

# --- Configuration ---
# Colors for output formatting
RESET   := \033[0m
BOLD    := \033[1m
GREEN   := \033[0;32m
YELLOW  := \033[0;33m
RED     := \033[0;31m
BLUE    := \033[0;34m
NC      := $(RESET) # No Color Alias

# Icons for visually distinct output
ICON_START := $(BLUE)▶$(NC)
ICON_OK    := $(GREEN)✓$(NC)
ICON_WARN  := $(YELLOW)⚠$(NC)
ICON_FAIL  := $(RED)✗$(NC)
ICON_INFO  := $(BLUE)ℹ$(NC)

# Formatting Strings for Alignment
LABEL_FMT := "   %-25s" # Indent 3, Pad label to 25 chars, left-aligned

# --- Variables ---
# Application specific variables
SERVICE_NAME := hello-tool-base
BINARY_NAME  := $(SERVICE_NAME) # Binary name will be 'hello-tool-base'

# Go module path - IMPORTANT: Ensure this matches your go.mod file
MODULE_PATH  := github.com/dkoosis/hello-tool-base

# Path to the main package, assuming main.go is in cmd/$(SERVICE_NAME)/
CMD_PATH     := ./cmd/$(SERVICE_NAME)
SCRIPT_DIR   := ./scripts

# Google Cloud Platform variables (primarily for 'deploy' target)
# Attempts to get PROJECT_ID from gcloud config.
PROJECT_ID   := $(shell gcloud config get-value project 2>/dev/null)

# Build-time variables for version injection using an internal/buildinfo package
VERSION      := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
# LDFLAGS for injecting build information into the internal/buildinfo package
LDFLAGS      := -ldflags "-s -w \
                -X $(MODULE_PATH)/internal/buildinfo.Version=$(VERSION) \
                -X $(MODULE_PATH)/internal/buildinfo.CommitHash=$(COMMIT_HASH) \
                -X $(MODULE_PATH)/internal/buildinfo.BuildDate=$(BUILD_DATE)"

# Tool Versions (for install-tools target)
GOLANGCILINT_VERSION := latest # Or a specific version e.g., v1.58.0
GOTESTSUM_VERSION    := latest # Or a specific version e.g., v1.11.0

# Line length check configuration
WARN_LINES   := 350  # Warn if lines exceed this
FAIL_LINES   := 1500 # Fail if lines exceed this
# GO_FILES for line length check - adjust if needed (e.g., to exclude more specific test files)
GO_FILES     := $(shell find . -name "*.go" -not -path "./vendor/*")

# --- Core Targets ---
# Default target: Runs all checks, formatting, tests, and then builds the application.
all: tree check-gomod deps fmt golangci-lint check-line-length test build
	@printf "$(GREEN)$(BOLD)✨ All checks passed and build completed successfully! ✨$(NC)\n"

# Generates a project directory tree view.
tree:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Generating project tree...$(NC)\n"
	@mkdir -p ./docs
	@if ! command -v tree > /dev/null; then \
	    printf "   $(ICON_WARN) $(YELLOW)'tree' command not found. Skipping tree generation.$(NC)\n"; \
	    printf "   $(ICON_INFO) $(YELLOW)To install on macOS: brew install tree$(NC)\n"; \
	    printf "   $(ICON_INFO) $(YELLOW)To install on Debian/Ubuntu: sudo apt-get install tree$(NC)\n"; \
	else \
	    tree -F -I 'vendor|.git|.idea*|*.DS_Store|$(BINARY_NAME)|coverage.out' --dirsfirst > ./docs/project_directory_tree.txt && \
	    printf "   $(ICON_OK) $(GREEN)Project tree generated at ./docs/project_directory_tree.txt$(NC)\n" || \
	    printf "   $(ICON_FAIL) $(RED)Failed to generate project tree.$(NC)\n"; \
	fi
	@printf "\n"

# Builds the application binary for Linux AMD64.
build:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Building $(BINARY_NAME) for linux/amd64...$(NC)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH) && \
	    printf "  $(ICON_OK) $(GREEN)Build successful: $(PWD)/$(BINARY_NAME)$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Build failed$(NC)\n" && exit 1)
	@printf "\n"

# Cleans build artifacts and Go caches.
clean:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Cleaning build artifacts...$(NC)\n"
	@rm -f $(BINARY_NAME) coverage.out
	@go clean -cache -testcache
	@printf "  $(ICON_OK) $(GREEN)Cleaned$(NC)\n"
	@printf "\n"

# --- Dependency Management ---
# Tidies and downloads Go module dependencies.
deps:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Synchronizing dependencies...$(NC)\n"
	@printf "  $(ICON_INFO) Running go mod tidy...\n"
	@go mod tidy -v && printf "  $(ICON_OK) $(GREEN)Dependencies tidied successfully$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Failed to tidy dependencies$(NC)\n" && exit 1)
	@printf "  $(ICON_INFO) Running go mod download...\n"
	@go mod download && printf "  $(ICON_OK) $(GREEN)Dependencies downloaded successfully$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Failed to download dependencies$(NC)\n" && exit 1)
	@printf "\n"

# Checks if the go.mod file exists and has the correct module path.
check-gomod:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Checking go.mod module path...$(NC)\n"
	@if [ ! -f "go.mod" ]; then \
	    printf "  $(ICON_FAIL) $(RED)go.mod file is missing. Run: go mod init $(MODULE_PATH)$(NC)\n"; \
	    exit 1; \
	fi
	@# Use the script if available, otherwise basic grep for a fallback check.
	@if [ -x "$(SCRIPT_DIR)/check_go_mod_path.sh" ]; then \
	    "$(SCRIPT_DIR)/check_go_mod_path.sh" $(MODULE_PATH) || exit 1; \
	elif ! grep -q "^module $(MODULE_PATH)$$" go.mod; then \
	    printf "  $(ICON_FAIL) $(RED)go.mod has incorrect module path (basic check).$(NC)\n"; \
	    printf "    $(ICON_INFO) $(YELLOW)Expected: module $(MODULE_PATH)$(NC)\n"; \
	    printf "    $(ICON_INFO) $(YELLOW)Found:    $$(grep '^module' go.mod)$(NC)\n"; \
	    exit 1; \
	else \
	    printf "  $(ICON_OK) $(GREEN)go.mod has correct module path (basic check).$(NC)\n"; \
	fi
	@printf "\n"

# --- Quality & Testing ---
# Formats Go code using golangci-lint's formatters.
fmt: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Formatting code using golangci-lint fmt...$(NC)\n"
	@golangci-lint fmt ./... && \
	    printf "  $(ICON_OK) $(GREEN)Code formatted$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)$(BOLD)Formatting failed (see errors above)$(NC)\n" && exit 1)
	@# Ensure go.mod and go.sum are tidy after formatting, as some formatters might change imports.
	@go mod tidy -v
	@printf "\n"

# Runs basic 'go vet' linter for quick checks.
lint:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running basic linter (go vet)...$(NC)\n"
	@go vet ./... && \
	    printf "  $(ICON_OK) $(GREEN)go vet passed$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)go vet found issues$(NC)\n" && exit 1)
	@printf "\n"

# Runs comprehensive linters using golangci-lint.
golangci-lint: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running golangci-lint...$(NC)\n"
	@golangci-lint run ./... && \
	    printf "  $(ICON_OK) $(GREEN)golangci-lint passed$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)$(BOLD)golangci-lint failed (see errors above)$(NC)\n" && exit 1)
	@printf "\n"

# Runs tests using gotestsum for better formatted output and race detection.
test: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running tests with gotestsum...$(NC)\n"
	@gotestsum --format testdox -- -race -coverprofile=coverage.out -covermode=atomic ./... && \
	    printf "  $(ICON_OK) $(GREEN)Tests passed$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
	@printf "\n"

# Runs tests with verbose output for debugging purposes.
test-debug: install-tools
	@printf "$(ICON_START) $(BOLD)$(YELLOW)Running tests (verbose debug mode)...$(NC)\n"
	@# LOG_LEVEL=debug can be used if your application/tests use it for more verbose logging.
	@go test -v -race -count=1 -coverprofile=coverage.out ./... && \
	    printf "  $(ICON_OK) $(GREEN)Tests finished (check output for failures)$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
	@printf "\n"

# Checks Go file line lengths using an external script.
check-line-length:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Checking file lengths (warn > $(WARN_LINES), fail > $(FAIL_LINES))...$(NC)\n"
	@if [ ! -x "$(SCRIPT_DIR)/check_file_length.sh" ]; then \
	    printf "   $(ICON_FAIL) $(RED)Error: Script '$(SCRIPT_DIR)/check_file_length.sh' not found or not executable.$(NC)\n"; \
	    exit 1; \
	fi
	@# Script handles its own output and exit status.
	@# Modified to exit 0 on warnings, but Makefile will still show a message.
	@$(SCRIPT_DIR)/check_file_length.sh $(WARN_LINES) $(FAIL_LINES) $(GO_FILES) || \
	    (printf "   $(ICON_WARN) $(YELLOW)Line length check reported issues (see script output above)$(NC)\n" && exit 0)
	@printf "\n"

# --- Tooling & Setup ---
# Installs or updates required Go development tools.
install-tools:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Checking/installing required Go tools...$(NC)\n"
	@printf $(LABEL_FMT) "golangci-lint:"
	@if ! command -v golangci-lint >/dev/null 2>&1 || ([ "$(GOLANGCILINT_VERSION)" != "latest" ] && ! golangci-lint --version | grep -qF "$(GOLANGCILINT_VERSION)"); then \
	    printf "$(ICON_INFO) $(YELLOW)Installing/Updating golangci-lint@$(GOLANGCILINT_VERSION)...$(NC)\n" ;\
	    GOBIN=$(shell go env GOPATH)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION) && \
	    printf "    $(ICON_OK) $(GREEN)golangci-lint installed/updated: $$(golangci-lint --version)$(NC)\n" || \
	    (printf "    $(ICON_FAIL) $(RED)golangci-lint installation/update failed$(NC)\n" && exit 1); \
	else \
	    printf "$(ICON_OK) $(GREEN)golangci-lint already installed: $$(golangci-lint --version)$(NC)\n"; \
	fi
	@printf $(LABEL_FMT) "gotestsum:"
	@if ! command -v gotestsum >/dev/null 2>&1 || ([ "$(GOTESTSUM_VERSION)" != "latest" ] && ! gotestsum --version | grep -qF "$(GOTESTSUM_VERSION)"); then \
	    printf "$(ICON_INFO) $(YELLOW)Installing/Updating gotestsum@$(GOTESTSUM_VERSION)...$(NC)\n" ;\
	    GOBIN=$(shell go env GOPATH)/bin go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION) && \
	    printf "    $(ICON_OK) $(GREEN)gotestsum installed/updated: $$(gotestsum --version)$(NC)\n" || \
	    (printf "    $(ICON_FAIL) $(RED)gotestsum installation/update failed$(NC)\n" && exit 1); \
	else \
	    printf "$(ICON_OK) $(GREEN)gotestsum already installed: $$(gotestsum --version)$(NC)\n"; \
	fi
	@printf "  $(ICON_OK) $(GREEN)Go tools check/installation complete$(NC)\n"
	@printf "\n"

# Checks for required local tools and environment setup.
check: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Checking for required tools and environment...$(NC)\n"
	@printf $(LABEL_FMT) "Go:"
	@if command -v go >/dev/null 2>&1; then \
	    printf "$(ICON_OK) $(GREEN)Installed$(NC) ($(shell go version))\n"; \
	else \
	    printf "$(ICON_FAIL) $(RED)Not installed$(NC)\n"; exit 1; \
	fi
	@printf $(LABEL_FMT) "tree:"
	@if command -v tree >/dev/null 2>&1; then \
	    printf "$(ICON_OK) $(GREEN)Installed$(NC)\n"; \
	else \
	    printf "$(ICON_WARN) $(YELLOW)Not installed$(NC) (needed for 'tree' target)\n"; \
	fi
	@# Add checks for other scripts if needed, e.g., check_go_bin_path.sh
	@if [ -x "$(SCRIPT_DIR)/check_go_bin_path.sh" ]; then \
		printf $(LABEL_FMT) "Go Bin Path Script:"; \
		"$(SCRIPT_DIR)/check_go_bin_path.sh" && \
		printf "$(ICON_OK) $(GREEN)check_go_bin_path.sh passed$(NC)\n" || \
		printf "$(ICON_WARN) $(YELLOW)check_go_bin_path.sh reported issues or failed$(NC)\n"; \
	else \
		printf $(LABEL_FMT) "Go Bin Path Script:"; \
		printf "$(ICON_WARN) $(YELLOW)Not found or not executable: $(SCRIPT_DIR)/check_go_bin_path.sh$(NC)\n"; \
	fi
	@printf "  $(ICON_OK) $(GREEN)Tool and environment check complete$(NC)\n"
	@printf "\n"

# --- Deployment ---
# Deploys the application to Google Cloud via Cloud Build.
# Depends on 'all' to ensure all checks, tests, and build are successful first.
deploy: all
	@printf "$(ICON_START) $(BOLD)$(BLUE)Deploying $(SERVICE_NAME) to Google Cloud...$(NC)\n"
	@if [ -z "$(PROJECT_ID)" ]; then \
	    printf "  $(ICON_FAIL) $(RED)Error: Google Cloud Project ID not found. Set via 'gcloud config set project YOUR_PROJECT_ID' or ensure gcloud is configured.$(NC)\n"; \
	    exit 1; \
	fi
	@printf "  $(ICON_INFO) Project: $(PROJECT_ID)\n"
	@# Determine the path to cloudbuild.yaml, preferring the one in build/cloudbuild/
	@CLOUDBUILD_CONFIG_PATH="./cloudbuild.yaml"; \
	if [ -f "build/cloudbuild/cloudbuild.yaml" ]; then \
	    CLOUDBUILD_CONFIG_PATH="./build/cloudbuild/cloudbuild.yaml"; \
	fi; \
	printf "  $(ICON_INFO) Using Cloud Build config: $$CLOUDBUILD_CONFIG_PATH\n";
	@# Trigger Cloud Build
	@gcloud builds submit . \
	    --config=$$CLOUDBUILD_CONFIG_PATH \
	    --project=$(PROJECT_ID) && \
	    printf "  $(ICON_OK) $(GREEN)Cloud Build triggered successfully.$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Cloud Build trigger failed.$(NC)\n" && exit 1)
	@printf "  $(ICON_INFO) Monitor build logs at: https://console.cloud.google.com/cloud-build/builds?project=$(PROJECT_ID)\n"

# --- Help ---
# Displays this help message.
help:
	@printf "$(BLUE)$(BOLD)$(SERVICE_NAME) Make Targets:$(NC)\n"
	@printf "  %-25s %s\n" "all" "Run checks, format, tests, and build (default)"
	@printf "  %-25s %s\n" "build" "Build the application binary ($(BINARY_NAME)) for Linux"
	@printf "  %-25s %s\n" "clean" "Clean build artifacts and caches"
	@printf "  %-25s %s\n" "deps" "Tidy and download Go module dependencies"
	@printf "  %-25s %s\n" "install-tools" "Install/update required Go tools (golangci-lint, gotestsum)"
	@printf "  %-25s %s\n" "check" "Check if required tools (Go, tree, etc.) are installed"
	@printf "  %-25s %s\n" "check-gomod" "Check if go.mod module path is correct"
	@printf "\n$(YELLOW)Code Quality & Testing:$(NC)\n"
	@printf "  %-25s %s\n" "fmt" "Format code using golangci-lint fmt"
	@printf "  %-25s %s\n" "lint" "Run basic 'go vet' checks (quick pre-check)"
	@printf "  %-25s %s\n" "golangci-lint" "Run comprehensive linters with golangci-lint"
	@printf "  %-25s %s\n" "test" "Run tests using gotestsum"
	@printf "  %-25s %s\n" "test-debug" "Run tests with verbose output (go test -v)"
	@printf "  %-25s %s\n" "check-line-length" "Check Go file line count (W:$(WARN_LINES), F:$(FAIL_LINES))"
	@printf "\n$(YELLOW)Other:$(NC)\n"
	@printf "  %-25s %s\n" "tree" "Generate project directory tree view in ./docs/"
	@printf "  %-25s %s\n" "deploy" "Run all checks, build, then deploy to Google Cloud"
	@printf "  %-25s %s\n" "help" "Display this help message"

