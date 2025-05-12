# Root Makefile for hello-tool-base
# A modular Makefile system for robust local development and CI/CD integration.
# This version organizes makefiles under the build/ directory to align with project structure.

# --- Force bash shell for consistent behavior ---
SHELL := /bin/bash
.SHELLFLAGS := -e -o pipefail -c

# --- Include configuration files ---
include build/makefiles/config.mk      # Variables and configuration
include build/makefiles/colors.mk      # Color and formatting definitions
include build/makefiles/tools.mk       # Tool installation and validation
include build/makefiles/build.mk       # Build targets
include build/makefiles/checks.mk      # Validation checks
include build/makefiles/test.mk        # Testing targets
include build/makefiles/deploy.mk      # Deployment targets

# --- Default target ---
.PHONY: all
all: tree check-gomod deps fmt golangci-lint lint-yaml check-line-length test build
	@printf "$(GREEN)$(BOLD)✨ All general checks passed and build completed successfully! ✨$(NC)\n"

# --- Help ---
.PHONY: help
help:
	@printf "$(BLUE)$(BOLD)$(SERVICE_NAME) Makefile Targets:$(NC)\n"
	@printf "$(YELLOW)Core Targets:$(NC)\n"
	@printf "  %-25s %s\n" "all" "Run general checks, format, tests, and build (default, no vulnerability scan)"
	@printf "  %-25s %s\n" "build" "Build the application binary ($(BINARY_NAME)) for Linux"
	@printf "  %-25s %s\n" "clean" "Clean build artifacts and caches"
	@printf "  %-25s %s\n" "deps" "Tidy and download Go module dependencies"
	@printf "  %-25s %s\n" "tree" "Generate project directory tree view"
	@printf "\n$(YELLOW)Tool Management:$(NC)\n"
	@printf "  %-25s %s\n" "install-tools" "Install/update required Go tools"
	@printf "  %-25s %s\n" "check-env" "Check environment and required tools"
	@printf "\n$(YELLOW)Code Quality:$(NC)\n"
	@printf "  %-25s %s\n" "fmt" "Format code using golangci-lint fmt"
	@printf "  %-25s %s\n" "lint" "Run basic 'go vet' checks"
	@printf "  %-25s %s\n" "golangci-lint" "Run comprehensive linters with golangci-lint"
	@printf "  %-25s %s\n" "lint-yaml" "Lint YAML files using yamllint"
	@printf "  %-25s %s\n" "check-gomod" "Check if go.mod module path is correct"
	@printf "  %-25s %s\n" "check-line-length" "Check Go file line counts (W:$(WARN_LINES), F:$(FAIL_LINES))"
	@printf "\n$(YELLOW)Testing & Security:$(NC)\n"
	@printf "  %-25s %s\n" "test" "Run tests using gotestsum"
	@printf "  %-25s %s\n" "test-debug" "Run tests with verbose output"
	@printf "  %-25s %s\n" "check-vulns" "Scan for known vulnerabilities in dependencies"
	@printf "\n$(YELLOW)Deployment:$(NC)\n"
	@printf "  %-25s %s\n" "deploy" "Run all checks, build, then deploy to Google Cloud"
	@printf "  %-25s %s\n" "health-check" "Check deployed service health"
	@printf "\n$(YELLOW)Other:$(NC)\n"
	@printf "  %-25s %s\n" "help" "Display this help message"
