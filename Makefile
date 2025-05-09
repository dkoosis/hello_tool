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
SERVICE_NAME           := hello-tool-base # Unified name for Cloud Run service and primary image artifact
REGION                 := us-central1 # Target region for deployment
ARTIFACT_REGISTRY_REPO := my-go-apps  # Your Artifact Registry repository name

# Build-time variables for version injection (optional, but good practice)
# To use these, declare corresponding string variables in your main.go's main package.
# For example:
# package main
# var (
#   version    string
#   commitHash string
#   buildDate  string
# )
# func main() {
#   if version != "" {
#     log.Printf("App Version: %s, Commit: %s, Build Date: %s", version, commitHash, buildDate)
#   }
#   // ... rest of your main function
# }
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
# Linker flags to inject version information into the binary
LDFLAGS     := -ldflags "-s -w -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}"

# --- Core Targets ---

# Default target: format, lint, and build the application
all: fmt lint build

# Build the application binary for Linux (typical for containerized deployments)
build:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Building $(BINARY_NAME) for linux/amd64...$(NC)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE) && \
	    printf "  $(ICON_OK) $(GREEN)Build successful: $(PWD)/$(BINARY_NAME)$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Build failed$(NC)\n" && exit 1)

# Clean build artifacts and Go build caches
clean:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Cleaning build artifacts...$(NC)\n"
	@rm -f $(BINARY_NAME)
	@go clean -cache -testcache
	@printf "  $(ICON_OK) $(GREEN)Cleaned$(NC)\n"

# --- Dependency Management ---

# Synchronize Go module dependencies
deps:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Synchronizing dependencies...$(NC)\n"
	@go mod tidy
	@go mod download
	@printf "  $(ICON_OK) $(GREEN)Dependencies synchronized$(NC)\n"

# --- Quality Assurance ---

# Format Go source code using gofmt
fmt:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Formatting code with gofmt...$(NC)\n"
	@go fmt ./...
	@printf "  $(ICON_OK) $(GREEN)Code formatted$(NC)\n"

# Run basic Go linter (go vet)
lint:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running basic linter (go vet)...$(NC)\n"
	@go vet ./... && \
	    printf "  $(ICON_OK) $(GREEN)go vet passed$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)go vet found issues$(NC)\n" && exit 1)

# --- Deployment ---

# Deploy the application to Google Cloud Run via Cloud Build, using cloudbuild.yaml
deploy: build # Ensures local build works, though Cloud Build performs its own build.
	@printf "$(ICON_START) $(BOLD)$(BLUE)Deploying $(SERVICE_NAME) to Google Cloud via Cloud Build...$(NC)\n"
	@if [ -z "$(PROJECT_ID)" ]; then \
		printf "  $(ICON_FAIL) $(RED)Error: Google Cloud Project ID not found. Set it via 'gcloud config set project YOUR_PROJECT_ID' or ensure gcloud is configured.$(NC)\n"; \
		exit 1; \
	fi
	@printf "  $(ICON_INFO) Project: $(PROJECT_ID)\n"
	@printf "  $(ICON_INFO) Service: $(SERVICE_NAME)\n"
	@printf "  $(ICON_INFO) Region: $(REGION)\n"
	@printf "  $(ICON_INFO) Artifact Registry Repo: $(ARTIFACT_REGISTRY_REPO)\n"
	@gcloud builds submit . \
		--config=cloudbuild.yaml \
		--project=$(PROJECT_ID) \
		--substitutions="_SERVICE_NAME=$(SERVICE_NAME),_REGION=$(REGION),_ARTIFACT_REGISTRY_REPO=$(ARTIFACT_REGISTRY_REPO)" && \
	    printf "  $(ICON_OK) $(GREEN)Cloud Build deployment triggered successfully.$(NC)\n" || \
	    (printf "  $(ICON_FAIL) $(RED)Cloud Build deployment trigger failed.$(NC)\n" && exit 1)
	@printf "  $(ICON_INFO) Monitor build logs at: https://console.cloud.google.com/cloud-build/builds?project=$(PROJECT_ID)\n"
	@# Attempt to construct the service URL - actual URL confirmed in Cloud Run console or via 'gcloud run services describe'
	@# This requires the project number, which can be fetched.
	@PROJECT_NUMBER=$(shell gcloud projects describe $(PROJECT_ID) --format="value(projectNumber)" 2>/dev/null); \
	if [ -n "$$PROJECT_NUMBER" ]; then \
		SERVICE_URL_HOSTNAME=$(SERVICE_NAME)-$$PROJECT_NUMBER-$(REGION).a.run.app; \
		printf "  $(ICON_INFO) Expected service URL (approximate): https://%s\n" "$$SERVICE_URL_HOSTNAME"; \
	else \
		printf "  $(ICON_WARN) $(YELLOW)Could not determine project number to form approximate service URL.$(NC)\n"; \
	fi


# --- Help ---

# Display available Make targets and their descriptions
help:
	@printf "$(BLUE)$(BOLD)hello-tool-base Make Targets:$(NC)\n"
	@printf "  %-20s %s\n" "all" "Format, lint, and build the application (default)"
	@printf "  %-20s %s\n" "build" "Build the application binary ($(BINARY_NAME)) for Linux"
	@printf "  %-20s %s\n" "clean" "Clean build artifacts and caches"
	@printf "  %-20s %s\n" "deps" "Tidy and download Go module dependencies"
	@printf "  %-20s %s\n" "fmt" "Format Go source code using gofmt"
	@printf "  %-20s %s\n" "lint" "Run 'go vet' checks"
	@printf "  %-20s %s\n" "deploy" "Deploy the application to Google Cloud Run via Cloud Build"
	@printf "  %-20s %s\n" "help" "Display this help message"

