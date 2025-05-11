# Makefile for hello-tool-base

# Specify phony targets (targets not associated with files)
.PHONY: all build clean deps fmt lint deploy help

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
ICON_FAIL  := $(RED)✗$(NC)
ICON_INFO  := $(BLUE)ℹ$(NC)

# --- Variables ---
# Application specific variables
BINARY_NAME  := app  # The name of the compiled executable
MAIN_PACKAGE := .    # The location of the main Go package (current directory)
MODULE_PATH  := github.com/dkoosis/hello-tool-base # Expected Go module path, ensure it matches go.mod

# Google Cloud Platform variables
# Attempt to get PROJECT_ID from gcloud config; ensure gcloud CLI is configured and authenticated.
PROJECT_ID             := $(shell gcloud config get-value project 2>/dev/null)
# These are still defined for other potential uses or for easy restoration later,
# but _SERVICE_NAME, _REGION, _ARTIFACT_REGISTRY_REPO are NOT passed as substitutions in the deploy target below for this test.
SERVICE_NAME           := hello-tool-base
REGION                 := us-central1
ARTIFACT_REGISTRY_REPO := my-go-apps

# Build-time variables for version injection
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS     := -ldflags "-s -w -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}"

# --- Core Targets ---
all: fmt lint build

build:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Building $(BINARY_NAME) for linux/amd64...$(NC)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE) && \
	    printf "  $(ICON_OK) $(GREEN)Build successful: $(PWD)/$(BINARY_NAME)$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Build failed$(NC)\n" && exit 1)

clean:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Cleaning build artifacts...$(NC)\n"
	@rm -f $(BINARY_NAME)
	@go clean -cache -testcache
	@printf "  $(ICON_OK) $(GREEN)Cleaned$(NC)\n"

deps:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Synchronizing dependencies...$(NC)\n"
	@go mod tidy
	@go mod download
	@printf "  $(ICON_OK) $(GREEN)Dependencies synchronized$(NC)\n"

fmt:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Formatting code with gofmt...$(NC)\n"
	@go fmt ./...
	@printf "  $(ICON_OK) $(GREEN)Code formatted$(NC)\n"

lint:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running basic linter (go vet)...$(NC)\n"
	@go vet ./... && \
	    printf "  $(ICON_OK) $(GREEN)go vet passed$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)go vet found issues$(NC)\n" && exit 1)

# --- Deployment ---
# MODIFIED FOR TEST: Removed --substitutions flag from gcloud builds submit
deploy: build
	@printf "$(ICON_START) $(BOLD)$(BLUE)Deploying $(SERVICE_NAME) to Google Cloud via Cloud Build (using hardcoded cloudbuild.yaml)...$(NC)\n"
	@if [ -z "$(PROJECT_ID)" ]; then \
		printf "  $(ICON_FAIL) $(RED)Error: Google Cloud Project ID not found. Set it via 'gcloud config set project YOUR_PROJECT_ID' or ensure gcloud is configured.$(NC)\n"; \
		exit 1; \
	fi
	@printf "  $(ICON_INFO) Project: $(PROJECT_ID)\n"
	@printf "  $(ICON_INFO) Service: $(SERVICE_NAME) (Note: cloudbuild.yaml has hardcoded names for this test)\n"
	@printf "  $(ICON_INFO) Region: $(REGION) (Note: cloudbuild.yaml has hardcoded names for this test)\n"
	@printf "  $(ICON_INFO) Artifact Registry Repo: $(ARTIFACT_REGISTRY_REPO) (Note: cloudbuild.yaml has hardcoded names for this test)\n"
	@gcloud builds submit . \
		--config=cloudbuild.yaml \
		--project=$(PROJECT_ID) && \
	    printf "  $(ICON_OK) $(GREEN)Cloud Build triggered successfully.$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Cloud Build trigger failed.$(NC)\n" && exit 1)
	@printf "  $(ICON_INFO) Monitor build logs at: https://console.cloud.google.com/cloud-build/builds?project=$(PROJECT_ID)\n"
	@PROJECT_NUMBER=$(shell gcloud projects describe $(PROJECT_ID) --format="value(projectNumber)" 2>/dev/null); \
	if [ -n "$$PROJECT_NUMBER" ]; then \
		SERVICE_URL_HOSTNAME=$(SERVICE_NAME)-$$PROJECT_NUMBER-$(REGION).a.run.app; \
		printf "  $(ICON_INFO) Expected service URL (approximate, based on Makefile vars): https://%s\n" "$$SERVICE_URL_HOSTNAME"; \
	else \
		printf "  $(ICON_WARN) $(YELLOW)Could not determine project number to form approximate service URL.$(NC)\n"; \
	fi

help:
	@printf "$(BLUE)$(BOLD)hello-tool-base Make Targets:$(NC)\n"
	@printf "  %-20s %s\n" "all" "Format, lint, and build the application (default)"
	@printf "  %-20s %s\n" "build" "Build the application binary ($(BINARY_NAME)) for Linux"
	@printf "  %-20s %s\n" "clean" "Clean build artifacts and caches"
	@printf "  %-20s %s\n" "deps" "Tidy and download Go module dependencies"
	@printf "  %-20s %s\n" "fmt" "Format Go source code using gofmt"
	@printf "  %-20s %s\n" "lint" "Run 'go vet' checks"
	@printf "  %-20s %s\n" "deploy" "Deploy to Google Cloud (uses cloudbuild.yaml with hardcoded names for current test)"
	@printf "  %-20s %s\n" "help" "Display this help message"
