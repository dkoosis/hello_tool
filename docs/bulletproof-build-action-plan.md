---

## Phase 1: Core Project Setup & Local Developer Experience (DX)

**Goal:** Establish a foundational Go project with a standard structure and local development tools.

1.  **Initialize Git Repository:**
    * If not already done, create a new directory for `hello-test-base`.
    * Navigate into the directory: `cd hello-test-base`
    * Initialize Git: `git init`
    * Create an initial `.gitignore` file (e.g., from [gitignore.io](https://www.toptal.com/developers/gitignore) for Go, add `hello-test-base` binary, `.env`).

2.  **Initialize Go Module:**
    * `go mod init github.com/your-username/hello-test-base` (replace with your actual module path).

3.  **Implement Standard Project Directory Structure:**
    * Create the following directories:
        * `cmd/hello-test-base/` (for `main.go`)
        * `internal/` (for private application logic, e.g., `internal/server/`, `internal/handlers/`)
        * `build/package/` (for `Dockerfile`)
        * `build/cloudbuild/` (for `cloudbuild.yaml`)
        * `scripts/` (for helper scripts, if any)
        * `configs/` (for example config files, if any)

4.  **Create Initial `Makefile`:**
    * Create a Makefile in the project root. This version incorporates build information, tool installation, improved linting/testing, and better output formatting inspired by your CowGnition project.

    * Action:
    Update GCP_PROJECT_ID, GCP_REGION, ARTIFACT_REGISTRY_REPO, and MODULE_PATH variables.
    Ensure you have a .golangci.yml configuration if you want golangci-lint fmt to use specific formatters beyond gofmt/goimports.
    The lint target uses golangci-lint run. The fmt target here uses golangci-lint fmt.

```code 

# --- Configuration ---
RESET      := \033[0m
BOLD       := \033[1m
GREEN      := \033[0;32m
YELLOW     := \033[0;33m
RED        := \033[0;31m
BLUE       := \033[0;34m
NC         := $(RESET) # No Color Alias

ICON_START := $(BLUE)▶$(NC)
ICON_OK    := $(GREEN)✓$(NC)
ICON_WARN  := $(YELLOW)⚠$(NC)
ICON_FAIL  := $(RED)✗$(NC)
ICON_INFO  := $(BLUE)ℹ$(NC)

LABEL_FMT  := "  %-20s" # Indent 2, Pad label to 20 chars, left-aligned

# --- Variables ---
SERVICE_NAME := hello-test-base
# IMPORTANT: Update this to your actual Go module path
MODULE_PATH  := github.com/your-username/hello-test-base
BINARY_NAME  := $(SERVICE_NAME)
CMD_PATH     := ./cmd/$(SERVICE_NAME)

# GCP Configuration (Update these)
GCP_PROJECT_ID         := your-gcp-project-id
GCP_REGION             := your-gcp-region # e.g., us-central1
ARTIFACT_REGISTRY_REPO := your-artifact-repo-name # e.g., hello-services

# Build information
VERSION      := $(shell git describe --tags --always --dirty --match=v* 2>/dev/null || echo "dev")
COMMIT_HASH  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
# Ensure the LDFLAGS correctly reference a package path where these vars are defined.
# Example: $(MODULE_PATH)/internal/buildinfo.Version
LDFLAGS      := -ldflags="-w -s -X $(MODULE_PATH)/internal/buildinfo.Version=$(VERSION) -X $(MODULE_PATH)/internal/buildinfo.CommitHash=$(COMMIT_HASH) -X $(MODULE_PATH)/internal/buildinfo.BuildDate=$(BUILD_DATE)"

# Tool Versions
GOLANGCILINT_VERSION := latest # Or a specific version e.g., v1.57.2
GOTESTSUM_VERSION    := latest # Or a specific version

# Phony targets
.PHONY: help all build-linux clean deps install-tools check-gomod fmt lint test test-debug \
    run-local docker-build-local docker-run-local

# --- Core Targets ---
all: check-gomod deps fmt lint test build-linux
    @printf "$(GREEN)$(BOLD)✨ All checks passed and build completed successfully! ✨$(NC)\n"

help:
    @printf "$(BLUE)$(BOLD)Available targets for $(SERVICE_NAME):$(NC)\n"
    @printf "$(LABEL_FMT) %s\n" "all" "Run all checks, format, test, and build (default)"
    @printf "$(LABEL_FMT) %s\n" "build-linux" "Build Go binary for Linux AMD64"
    @printf "$(LABEL_FMT) %s\n" "clean" "Clean build artifacts and Go caches"
    @printf "\n$(YELLOW)Dependency & Setup:$(NC)\n"
    @printf "$(LABEL_FMT) %s\n" "deps" "Tidy and download Go module dependencies"
    @printf "$(LABEL_FMT) %s\n" "install-tools" "Install/update required Go tools (golangci-lint, gotestsum)"
    @printf "$(LABEL_FMT) %s\n" "check-gomod" "Check if go.mod module path is correct"
    @printf "\n$(YELLOW)Code Quality & Formatting:$(NC)\n"
    @printf "$(LABEL_FMT) %s\n" "fmt" "Format Go code using golangci-lint fmt"
    @printf "$(LABEL_FMT) %s\n" "lint" "Run comprehensive linters with golangci-lint"
    @printf "\n$(YELLOW)Testing:$(NC)\n"
    @printf "$(LABEL_FMT) %s\n" "test" "Run tests using gotestsum (standard output)"
    @printf "$(LABEL_FMT) %s\n" "test-debug" "Run tests with verbose output (go test -v)"
    @printf "\n$(YELLOW)Local Development & Docker:$(NC)\n"
    @printf "$(LABEL_FMT) %s\n" "run-local" "Build and run the service locally"
    @printf "$(LABEL_FMT) %s\n" "docker-build-local" "Build Docker image locally (tagged as $(SERVICE_NAME):local)"
    @printf "$(LABEL_FMT) %s\n" "docker-run-local" "Run locally built Docker image"
    @printf "\n"

build-linux:
    @printf "$(ICON_START) $(BOLD)Building $(BINARY_NAME) for Linux...$(NC)\n"
    @CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME) $(CMD_PATH)/main.go && \
        printf "  $(ICON_OK) $(GREEN)Build successful: build/$(BINARY_NAME)$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)Build failed$(NC)\n" && exit 1)
    @printf "\n"

clean:
    @printf "$(ICON_START) $(BOLD)Cleaning build artifacts...$(NC)\n"
    @rm -f build/$(BINARY_NAME) coverage.out
    @go clean -cache -testcache
    @printf "  $(ICON_OK) $(GREEN)Cleaned$(NC)\n"
    @printf "\n"

# --- Dependency Management ---
deps:
    @printf "$(ICON_START) $(BOLD)Synchronizing dependencies...$(NC)\n"
    @printf "  $(ICON_INFO) Running go mod tidy...\n"
    @go mod tidy -v && printf "  $(ICON_OK) $(GREEN)Dependencies tidied successfully$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)Failed to tidy dependencies$(NC)\n" && exit 1)
    @printf "  $(ICON_INFO) Running go mod download...\n"
    @go mod download && printf "  $(ICON_OK) $(GREEN)Dependencies downloaded successfully$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)Failed to download dependencies$(NC)\n" && exit 1)
    @printf "\n"

check-gomod:
    @printf "$(ICON_START) $(BOLD)Checking go.mod module path...$(NC)\n"
    @if [ ! -f "go.mod" ]; then \
        printf "  $(ICON_FAIL) $(RED)go.mod file is missing. Run: go mod init $(MODULE_PATH)$(NC)\n"; \
        exit 1; \
    fi
    @if ! grep -q "^module $(MODULE_PATH)$$" go.mod; then \
        printf "  $(ICON_FAIL) $(RED)go.mod has incorrect module path.$(NC)\n"; \
        printf "    $(ICON_INFO) $(YELLOW)Expected: module $(MODULE_PATH)$(NC)\n"; \
        printf "    $(ICON_INFO) $(YELLOW)Found:    $$(grep '^module' go.mod)$(NC)\n"; \
        exit 1; \
    fi
    @printf "  $(ICON_OK) $(GREEN)go.mod has correct module path$(NC)\n"
    @printf "\n"

# --- Quality & Testing ---
fmt: install-tools
    @printf "$(ICON_START) $(BOLD)Formatting code using golangci-lint fmt...$(NC)\n"
    @golangci-lint fmt ./... && \
        printf "  $(ICON_OK) $(GREEN)Code formatted$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)Formatting failed (see errors above)$(NC)\n" && exit 1)
    @go mod tidy -v # Ensure go.mod is tidy after formatting
    @printf "\n"

lint: install-tools
    @printf "$(ICON_START) $(BOLD)Running golangci-lint...$(NC)\n"
    @golangci-lint run ./... && \
        printf "  $(ICON_OK) $(GREEN)golangci-lint passed$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)golangci-lint failed (see errors above)$(NC)\n" && exit 1)
    @printf "\n"

test: install-tools
    @printf "$(ICON_START) $(BOLD)Running tests with gotestsum...$(NC)\n"
    @gotestsum --format testdox -- -race -coverprofile=coverage.out -covermode=atomic ./... && \
        printf "  $(ICON_OK) $(GREEN)Tests passed$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
    @printf "\n"

test-debug: install-tools
    @printf "$(ICON_START) $(BOLD)$(YELLOW)Running tests (verbose debug mode)...$(NC)\n"
    @# LOG_LEVEL=debug (if your app uses this for more verbose test logging)
    @go test -v -race -count=1 -coverprofile=coverage.out ./... && \
        printf "  $(ICON_OK) $(GREEN)Tests finished (check output for failures)$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
    @printf "\n"

# --- Tooling & Setup ---
install-tools:
    @printf "$(ICON_START) $(BOLD)Checking/installing required Go tools...$(NC)\n"
    @printf "  $(LABEL_FMT) %s\n" "golangci-lint:" ""
    @if ! command -v golangci-lint >/dev/null 2>&1 || ([ "$(GOLANGCILINT_VERSION)" != "latest" ] && ! golangci-lint --version | grep -qF "$(GOLANGCILINT_VERSION)"); then \
        printf "    $(ICON_INFO) $(YELLOW)Installing/Updating golangci-lint@$(GOLANGCILINT_VERSION)...$(NC)\n" ;\
        GOBIN=$(shell go env GOPATH)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION) && \
        printf "    $(ICON_OK) $(GREEN)golangci-lint installed/updated: $$(golangci-lint --version)$(NC)\n" || \
        (printf "    $(ICON_FAIL) $(RED)golangci-lint installation/update failed$(NC)\n" && exit 1); \
    else \
        printf "    $(ICON_OK) $(GREEN)golangci-lint already installed: $$(golangci-lint --version)$(NC)\n"; \
    fi
    @printf "  $(LABEL_FMT) %s\n" "gotestsum:" ""
    @if ! command -v gotestsum >/dev/null 2>&1 || ([ "$(GOTESTSUM_VERSION)" != "latest" ] && ! gotestsum --version | grep -qF "$(GOTESTSUM_VERSION)"); then \
        printf "    $(ICON_INFO) $(YELLOW)Installing/Updating gotestsum@$(GOTESTSUM_VERSION)...$(NC)\n" ;\
        GOBIN=$(shell go env GOPATH)/bin go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION) && \
        printf "    $(ICON_OK) $(GREEN)gotestsum installed/updated: $$(gotestsum --version)$(NC)\n" || \
        (printf "    $(ICON_FAIL) $(RED)gotestsum installation/update failed$(NC)\n" && exit 1); \
    else \
        printf "    $(ICON_OK) $(GREEN)gotestsum already installed: $$(gotestsum --version)$(NC)\n"; \
    fi
    @printf "  $(ICON_OK) $(GREEN)Go tools check/installation complete$(NC)\n"
    @printf "\n"

# --- Local Development & Docker ---
run-local: build-linux
    @printf "$(ICON_START) $(BOLD)Running $(SERVICE_NAME) locally...$(NC)\n"
    @# You might want to pass environment variables here if needed for local run
    @# Example: API_KEY_ENV="local_dev_key" ./build/$(BINARY_NAME)
    @./build/$(BINARY_NAME)

docker-build-local:
    @printf "$(ICON_START) $(BOLD)Building Docker image $(SERVICE_NAME):local...$(NC)\n"
    @docker build -t $(SERVICE_NAME):local -f build/package/Dockerfile . && \
        printf "  $(ICON_OK) $(GREEN)Docker image built: $(SERVICE_NAME):local$(NC)\n" || \
        (printf "  $(ICON_FAIL) $(RED)$(BOLD)Docker build failed$(NC)\n" && exit 1)
    @printf "\n"

docker-run-local:
    @printf "$(ICON_START) $(BOLD)Running Docker image $(SERVICE_NAME):local...$(NC)\n"
    @# Ensure SERVICE_NAME is passed to the container if your app uses it
    @docker run --rm -p 8080:8080 -e PORT=8080 -e SERVICE_NAME=$(SERVICE_NAME) $(SERVICE_NAME):local
    @printf "\n"

```

5.  **Develop a Hello World:**

Create an internal/buildinfo/buildinfo.go file:
Go

package buildinfo

// These variables are populated by LDFLAGS during the build process
var (
	Version    string = "dev"
	CommitHash string = "unknown"
	BuildDate  string = "unknown"
)
Modify cmd/hello-test-base/main.go to use and log this build information.
Go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// IMPORTANT: Update this import path to match your MODULE_PATH from the Makefile
	// For example: "github.com/your-username/hello-test-base/internal/buildinfo"
	"YOUR_MODULE_PATH_HERE/internal/buildinfo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "hello-test-base" // Default if not set
	}

	log.Printf("Starting %s | Version: %s | Commit: %s | Built: %s",
		serviceName, buildinfo.Version, buildinfo.CommitHash, buildinfo.BuildDate)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		helloHandler(w, r, serviceName)
	})

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request, serviceNameFromMain string) {
	log.Printf("Received request for %s from %s\n", r.URL.Path, r.RemoteAddr)

	greeting := os.Getenv("APP_GREETING")
	if greeting == "" {
		greeting = fmt.Sprintf("Hello from %s!", serviceNameFromMain)
	}
	fmt.Fprintln(w, greeting) // Use Fprintln for a newline

	// Example of using a secret passed as env var (from Phase 5)
	apiKey := os.Getenv("API_KEY_ENV")
	if apiKey != "" {
		// In a real app, NEVER log the actual secret value like this.
		// This is just for demonstration.
		fmt.Fprintf(w, "Secret API Key (demo): %s (first 5 chars)\n", apiKey[:min(5, len(apiKey))])
	}

	fmt.Fprintf(w, "\nBuild Info:\n  Version: %s\n  Commit: %s\n  Date: %s\n",
		buildinfo.Version, buildinfo.CommitHash, buildinfo.BuildDate)
}

// Helper function (add this outside main/helloHandler)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
Action:
Crucially, replace YOUR_MODULE_PATH_HERE in the import statement within main.go with your actual module path (e.g., github.com/your-username/hello-test-base). This must match the MODULE_PATH variable in your Makefile.
Test locally: make run-local. You should see the build information in the logs and when you access the service.

---

## Phase 2: Containerization

**Goal:** Create an optimized and secure Docker image for the application.

1. **Create `Dockerfile`:**
    * Create `build/package/Dockerfile`.

        ```dockerfile
        # ---- Builder Stage ----
        FROM golang:1.22-alpine AS builder

        WORKDIR /app

        # Install build tools if necessary (e.g., make, git)
        # RUN apk add --no-cache make git

        # Copy Go modules and download dependencies
        # This layer is cached unless go.mod/go.sum changes
        COPY go.mod go.sum ./
        RUN go mod download -d all

        # Copy the rest of the application source code
        COPY . .

        # Get SERVICE_NAME from build arg, defaulting if not provided
        ARG SERVICE_NAME=hello-test-base
        ENV SERVICE_NAME=<span class="math-inline">\{SERVICE\_NAME\}

\# Build the application
\# Ensure your Makefile's build\-linux target outputs to /app/app\_binary or adjust path
RUN CGO\_ENABLED\=0 GOOS\=linux GOARCH\=amd64 go build \-ldflags\="\-w \-s" \-o /app/app\_binary \./cmd/</span>{SERVICE_NAME}/main.go

        # ---- Final Stage ----
        FROM gcr.io/distroless/static-nonroot
        # Or use gcr.io/distroless/base-nonroot if CGO_ENABLED=1 or you need shell access for debugging (not recommended for prod)

        WORKDIR /
        COPY --from=builder /app/app_binary /app_binary

        # Copy any other necessary assets (e.g., config files, templates) if needed
        # COPY --from=builder /app/configs /configs

        # Expose the port the application listens on
        EXPOSE 8080

        # Set environment variables (can be overridden at runtime)
        ENV PORT=8080
        # SERVICE_NAME will be passed by Cloud Run or Docker run command

        # Run as non-root user (distroless/*-nonroot images do this by default)
        # USER nonroot:nonroot # Not needed for distroless nonroot images

        ENTRYPOINT ["/app_binary"]
        ```
    * **Note:** Adjust `SERVICE_NAME` in Dockerfile `RUN` and `CMD_PATH` in Makefile if your main package path is different.

2.  **Create `.dockerignore`:**
    * Create `.dockerignore` in the project root.
        ```
        # Git
        .git
        .gitignore

        # Go build artifacts
        *.out
        build/*
        !build/package
        !build/cloudbuild

        # Local env
        .env*

        # IDE / OS specific
        .vscode/
        .idea/
        *.DS_Store
        ```

3.  **Build and Test Docker Image Locally:**
    *`make docker-build-local`
    * `make docker-run-local`
    * Open your browser to `<http://localhost>

---

## Phase 3: Basic CI/CD with Cloud Build & Artifact Registry

**Goal:** Automate building and deploying the application to Cloud Run using Cloud Build.

1. **Google Cloud Project Setup:**
    * Ensure you have a GCP Project selected: `gcloud config set project YOUR_PROJECT_ID`
    * Enable necessary APIs:
        * `gcloud services enable cloudbuild.googleapis.com`
        * `gcloud services enable artifactregistry.googleapis.com`
        * `gcloud services enable run.googleapis.com`
        * `gcloud services enable iam.googleapis.com`
        * `gcloud services enable secretmanager.googleapis.com` (for later)

2. **Create Artifact Registry Repository:**
    * `gcloud artifacts repositories create $(ARTIFACT_REGISTRY_REPO) --repository-format=docker --location=$(GCP_REGION) --description="Docker repository for hello services"`
    * Configure Docker to authenticate with Artifact Registry (if you want to push manually, Cloud Build does this automatically): `gcloud auth configure-docker $(GCP_REGION)-docker.pkg.dev`

3. **Grant Permissions to Cloud Build Service Account:**
    * Cloud Build's default SA (`[PROJECT_NUMBER]@cloudbuild.gserviceaccount.com`) often has broad permissions. For production, consider creating a dedicated SA with least privilege. For this plan, we'll describe roles for the default SA for simplicity.
    * **Permissions needed by Cloud Build SA:**
        * Artifact Registry Writer: `roles/artifactregistry.writer` (to push images)
        * Cloud Run Developer: `roles/run.developer` (to deploy and manage Cloud Run services)
        * Service Account User (to act as the Cloud Run runtime service): `roles/iam.serviceAccountUser` (on the Cloud Run runtime SA, which you'll create later)
        * Secret Manager Secret Accessor (for later, if Cloud Build needs to access secrets during build time, or if it provisions secrets): `roles/secretmanager.secretAccessor`
    * **Action:** Grant these roles to your Cloud Build service account. Example for default SA:

        ```bash
        PROJECT_NUMBER=$(gcloud projects describe $(GCP_PROJECT_ID) --format='get(projectNumber)')
        CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

        gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
            --member="serviceAccount:${CLOUD_BUILD_SA}" \
            --role="roles/artifactregistry.writer"

        gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
            --member="serviceAccount:${CLOUD_BUILD_SA}" \
            --role="roles/run.developer"

        # Note: The 'roles/iam.serviceAccountUser' permission on the Cloud Run runtime SA
        # will be granted when you create that specific SA for your Cloud Run service.
        ```

4. **Create `cloudbuild.yaml`:**
    * Create `build/cloudbuild/cloudbuild.yaml`.

        ```yaml
        steps:
          # Lint (Optional but recommended)
          # - name: 'golangci/golangci-lint:v1.55.2' # Use a specific version
          #   id: 'Lint'
          #   args: ['golangci-lint', 'run', './...']

          # Run tests
          - name: 'golang:1.22-alpine' # Match your Go version
            id: 'Run Unit Tests'
            entrypoint: /bin/sh
            args:
              - -c
              - |
                echo "Running tests..."
                go test -v -race -coverprofile=coverage.out ./...
                echo "Tests completed."

          # Build Docker image
          - name: 'gcr.io/cloud-builders/docker'
            id: 'Build Docker Image'
            args:
              - 'build'
              - '-t'
              - '${_GCP_REGION}-docker.pkg.dev/$PROJECT_ID/${_ARTIFACT_REGISTRY_REPO}/${_SERVICE_NAME}:$SHORT_SHA'
              - '-f'
              - 'build/package/Dockerfile'
              - '.'
              - '--build-arg' # Pass SERVICE_NAME to Docker build
              - 'SERVICE_NAME=${_SERVICE_NAME}'


          # Push Docker image to Artifact Registry
          - name: 'gcr.io/cloud-builders/docker'
            id: 'Push Docker Image'
            args:
              - 'push'
              - '${_GCP_REGION}-docker.pkg.dev/$PROJECT_ID/${_ARTIFACT_REGISTRY_REPO}/${_SERVICE_NAME}:$SHORT_SHA'

          # Deploy to Cloud Run
          - name: 'gcr.io/[google.com/cloudsdktool/cloud-sdk](https://google.com/cloudsdktool/cloud-sdk)'
            id: 'Deploy to Cloud Run'
            entrypoint: gcloud
            args:
              - 'run'
              - 'deploy'
              - '${_SERVICE_NAME}'
              - '--image=${_GCP_REGION}-docker.pkg.dev/$PROJECT_ID/${_ARTIFACT_REGISTRY_REPO}/${_SERVICE_NAME}:$SHORT_SHA'
              - '--region=${_GCP_REGION}'
              - '--platform=managed' # or 'gke' if deploying to GKE cluster
              - '--quiet'
              - '--allow-unauthenticated' # For public access, remove for private
              - '--port=8080' # Port your application listens on inside the container
              - '--set-env-vars=SERVICE_NAME=${_SERVICE_NAME}' # Example environment variable
              # - '--service-account=your-cloud-run-runtime-sa@${PROJECT_ID}.iam.gserviceaccount.com' # Specify runtime SA later

        # Store images in Artifact Registry
        images:
          - '${_GCP_REGION}-docker.pkg.dev/$PROJECT_ID/${_ARTIFACT_REGISTRY_REPO}/${_SERVICE_NAME}:$SHORT_SHA'

        # Define substitutions (can be overridden by triggers)
        substitutions:
          _SERVICE_NAME: 'hello-test-base'
          _GCP_REGION: 'your-gcp-region' # e.g., us-central1
          _ARTIFACT_REGISTRY_REPO: 'your-artifact-repo-name' # e.g., hello-services

        options:
          logging: CLOUD_LOGGING_ONLY
          # machineType: 'E2_HIGHCPU_8' # Optional: Specify machine type for build
        ```

    * **Action:** Update `_GCP_REGION` and `_ARTIFACT_REGISTRY_REPO` in the `substitutions` section of your `cloudbuild.yaml` and in your `Makefile` variables.

5. **Test Cloud Build Locally (Optional but Recommended):**
    * `gcloud builds submit --config build/cloudbuild/cloudbuild.yaml . --substitutions=_SERVICE_NAME=hello-test-base,_GCP_REGION=your-gcp-region,_ARTIFACT_REGISTRY_REPO=your-artifact-repo-name`
    * Check the build logs in the GCP console.

6. **Configure Cloud Build Trigger:**
    * Go to Cloud Build > Triggers in the GCP Console.
    * Connect your GitHub repository.
    * Create a new trigger:
        * Name: e.g., `deploy-hello-test-base-main`
        * Event: Push to a branch
        * Repository: Your `hello-test-base` repo
        * Branch: `^main$` (or your primary branch)
        * Configuration: Cloud Build configuration file (path to `build/cloudbuild/cloudbuild.yaml`)
        * Substitution variables (these will override those in the `cloudbuild.yaml` if set here):
            * `_SERVICE_NAME`: `hello-test-base`
            * `_GCP_REGION`: `your-gcp-region`
            * `_ARTIFACT_REGISTRY_REPO`: `your-artifact-repo-name`

7. **Commit Cloud Build Setup:**
    * `git add .`
    * `git commit -m "feat: Add Cloud Build configuration for CI/CD"`
    * `git push origin main` (This should trigger your first build if the trigger is active)
    * Verify the deployment in Cloud Run console and access the service URL.

---
---

## Phase 4: GitHub Actions for CI Orchestration (Optional - WIF)

**Goal:** Trigger Cloud Build using GitHub Actions with Workload Identity Federation for enhanced security, avoiding the need to manage service account keys in GitHub Secrets.

1. **Set up Workload Identity Federation (WIF) in GCP:**
    * Follow the official Google Cloud guide: [Configuring Workload Identity Federation](https://cloud.google.com/iam/docs/configuring-workload-identity-federation)
    * **Create a WIF Pool:** e.g., `github-wif-pool`
    * **Create a WIF Provider within the Pool:** Configure it for GitHub, specifying your GitHub organization or user.
    * **Create a dedicated GCP Service Account (SA) for GitHub Actions to impersonate:**
        * Name: e.g., `github-actions-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com`
        * Grant this SA the `roles/cloudbuild.builds.editor` role on your project (to allow it to trigger Cloud Builds).

            ```bash
            gcloud projects add-iam-policy-binding $(GCP_PROJECT_ID) \
                --member="serviceAccount:github-actions-runner@$(GCP_PROJECT_ID).iam.gserviceaccount.com" \
                --role="roles/cloudbuild.builds.editor"
            ```

    * **Grant the GitHub WIF identity permission to impersonate this SA:**
        * You'll need your `PROJECT_NUMBER`, WIF Pool ID, and the attribute mapping for your GitHub repository. The `attribute.repository` mapping allows specific repositories to impersonate the SA.

            ```bash
            PROJECT_NUMBER=$(gcloud projects describe $(GCP_PROJECT_ID) --format='get(projectNumber)')
            GCP_SA_EMAIL="github-actions-runner@$(GCP_PROJECT_ID).iam.gserviceaccount.com"
            WIF_POOL_ID="your-wif-pool-id" # Replace with your WIF Pool ID
            # Attribute mapping for a specific repository:
            REPO_ATTRIBUTE_MAPPING="attribute.repository/your-github-org-or-user/your-repo-name"
            # Or for all repositories in an organization:
            # REPO_ATTRIBUTE_MAPPING="attribute.repository_owner/your-github-org-or-user"


            gcloud iam service-accounts add-iam-policy-binding ${GCP_SA_EMAIL} \
              --project=$(GCP_PROJECT_ID) \
              --role="roles/iam.workloadIdentityUser" \
              --member="principalSet://[iam.googleapis.com/projects/$](https://iam.googleapis.com/projects/$){PROJECT_NUMBER}/locations/global/workloadIdentityPools/${WIF_POOL_ID}/${REPO_ATTRIBUTE_MAPPING}"
            ```

        * **Action:** Replace placeholders like `your-wif-pool-id`, `your-github-org-or-user`, and `your-repo-name`.

2. **Create GitHub Actions Workflow:**
    * Create `.github/workflows/google-cloud-build.yml`.

        ```yaml
        name: Trigger Google Cloud Build

        on:
          push:
            branches:
              - main # Or your trigger branch

        jobs:
          build-and-deploy:
            name: Build and Deploy
            runs-on: ubuntu-latest
            permissions:
              contents: 'read'
              id-token: 'write' # Required for OIDC token to exchange for GCP token

            steps:
              - name: Checkout code
                uses: actions/checkout@v4

              - name: Authenticate to Google Cloud
                uses: google-github-actions/auth@v2
                with:
                  workload_identity_provider: 'projects/YOUR_PROJECT_NUMBER/locations/global/workloadIdentityPools/YOUR_WIF_POOL_ID/providers/YOUR_WIF_PROVIDER_ID' # Update these details
                  service_account: 'github-actions-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com' # Update with your SA email

              - name: Set up Cloud SDK
                uses: google-github-actions/setup-gcloud@v2
                with:
                  project_id: YOUR_PROJECT_ID # Explicitly set project ID

              - name: Submit build to Google Cloud Build
                run: |-
                  gcloud builds submit --config build/cloudbuild/cloudbuild.yaml . \
                    --substitutions=_SERVICE_NAME=hello-test-base,_GCP_REGION=${{ vars.GCP_REGION }},_ARTIFACT_REGISTRY_REPO=${{ vars.ARTIFACT_REGISTRY_REPO }},COMMIT_SHA=${{ github.sha }},SHORT_SHA=${{ github.sha_short }}
                env:
                  # PROJECT_ID is set by setup-gcloud step
                  # GCP_REGION and ARTIFACT_REGISTRY_REPO can be set as GitHub repository variables
                  # Or hardcode them in the substitutions if they don't change often
                  GCP_REGION_VAR: ${{ vars.GCP_REGION }} # Example using GitHub repo variable
                  ARTIFACT_REGISTRY_REPO_VAR: ${{ vars.ARTIFACT_REGISTRY_REPO }} # Example

        ```

    * **Action:**
        * Update `YOUR_PROJECT_NUMBER`, `YOUR_WIF_POOL_ID`, `YOUR_WIF_PROVIDER_ID`, `YOUR_PROJECT_ID`, and the `service_account` email in the workflow file.
        * Consider using GitHub Repository Variables (Settings > Secrets and variables > Actions > Variables) for `GCP_REGION` and `ARTIFACT_REGISTRY_REPO` to avoid hardcoding them in the workflow. The example uses `vars.GCP_REGION` and `vars.ARTIFACT_REGISTRY_REPO`.
        * The `SHORT_SHA` substitution is added to `cloudbuild.yaml` image tagging and `gcloud builds submit`. Ensure your `cloudbuild.yaml` uses `$SHORT_SHA` for image tagging.
        * If you use this GitHub Actions workflow as your primary trigger, you might disable or make the Cloud Build native trigger manual to avoid duplicate builds.

3. **Commit GitHub Actions Workflow:**
    * `git add .github/workflows/google-cloud-build.yml`
    * `git commit -m "ci: Add GitHub Actions workflow to trigger Cloud Build via WIF"`
    * `git push`
    * Monitor the GitHub Actions run and the subsequent Cloud Build.

---
---

## Phase 5: Configuration & Secrets (Basic)

**Goal:** Demonstrate basic configuration and secret management for the Cloud Run service. This phase focuses on injecting configuration and secrets into your Cloud Run service.

1. **Environment Variables via `cloudbuild.yaml` (for non-sensitive config):**
    * You've already added `SERVICE_NAME` as an environment variable in the `deploy` step of your `build/cloudbuild/cloudbuild.yaml`:

        ```yaml
        # In deploy step args:
        # ...
              - '--set-env-vars=SERVICE_NAME=${_SERVICE_NAME}'
        # ...
        ```

    * **Action:** Add another example non-sensitive configuration. Modify the `--set-env-vars` line:

        ```yaml
        # In deploy step args:
        # ...
              - '--set-env-vars=SERVICE_NAME=${_SERVICE_NAME},APP_GREETING="Hello from your Go Service!"'
        # ...
        ```

    * **Action:** Modify your `helloHandler` in `cmd/hello-test-base/main.go` to use this new environment variable.

        ```go
        // In func helloHandler:
        greeting := os.Getenv("APP_GREETING")
        if greeting == "" {
            greeting = fmt.Sprintf("Hello from %s!", os.Getenv("SERVICE_NAME"))
        }
        fmt.Fprintln(w, greeting) // Use Fprintln for a newline
        ```

2. **Basic Secret Manager Setup (for sensitive data):**
    * **Create a secret in Google Cloud Secret Manager:**

        ```bash
        # Ensure your GCP_PROJECT_ID variable is set in your shell
        # export GCP_PROJECT_ID="your-gcp-project-id"
        echo "ThisIsMySuperSecretAPIKey123!" | gcloud secrets create MY_API_KEY \
            --data-file=- \
            --replication-policy=automatic \
            --project=$(GCP_PROJECT_ID)
        ```

    * **Create a dedicated runtime Service Account (SA) for your Cloud Run service:**
        * This SA will be used by your Cloud Run service at runtime. It's best practice for it to have minimal necessary permissions.

            ```bash
            gcloud iam service-accounts create hello-test-base-run-sa \
                --display-name="Cloud Run Runtime SA for hello-test-base" \
                --project=$(GCP_PROJECT_ID)
            ```

        * Store the email of this new SA:

            ```bash
            CLOUD_RUN_SA_EMAIL="hello-test-base-run-sa@$(GCP_PROJECT_ID).iam.gserviceaccount.com"
            ```

    * **Grant the Cloud Run Runtime SA access to the secret:**

        ```bash
        gcloud secrets add-iam-policy-binding MY_API_KEY \
            --member="serviceAccount:${CLOUD_RUN_SA_EMAIL}" \
            --role="roles/secretmanager.secretAccessor" \
            --project=$(GCP_PROJECT_ID)
        ```

    * **Grant the Cloud Build SA permission to act as (impersonate) the Cloud Run Runtime SA:**
        * This is necessary for Cloud Build to deploy the service *as* the runtime SA.

            ```bash
            PROJECT_NUMBER=$(gcloud projects describe $(GCP_PROJECT_ID) --format='get(projectNumber)')
            CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

            gcloud iam service-accounts add-iam-policy-binding ${CLOUD_RUN_SA_EMAIL} \
                --member="serviceAccount:${CLOUD_BUILD_SA}" \
                --role="roles/iam.serviceAccountUser" \
                --project=$(GCP_PROJECT_ID)
            ```

    * **Update `build/cloudbuild/cloudbuild.yaml` to use this runtime SA and mount the secret:**
        * Modify the `Deploy to Cloud Run` step:

            ```yaml
            # In deploy step args:
            # ...
              - '--service-account=${CLOUD_RUN_SA_EMAIL}' # Use the variable or hardcode
              # Option 1: Mount secret as an environment variable
              - '--set-secrets=API_KEY_ENV=MY_API_KEY:latest'
              # Option 2: Mount secret as a file (e.g., at /etc/secrets/my-api-key)
              # - '--set-secrets=/etc/secrets/my-api-key=MY_API_KEY:latest'
            # ...
            ```

        * **Action:** If you use `${CLOUD_RUN_SA_EMAIL}` as above, you'll need to pass it as a substitution to `gcloud builds submit` or define it in the `substitutions` block of `cloudbuild.yaml`. For simplicity, you can hardcode the SA email directly in the `cloudbuild.yaml` for now if you prefer.
            Example with substitution:

            ```yaml
            # In cloudbuild.yaml substitutions:
            # _CLOUD_RUN_SA_EMAIL: 'hello-test-base-run-sa@your-gcp-project-id.iam.gserviceaccount.com'
            # And in deploy args:
            # - '--service-account=${_CLOUD_RUN_SA_EMAIL}'
            ```

    * **Modify `cmd/hello-test-base/main.go` to read and display the secret (for demonstration only; be careful with logging real secrets):**

        ```go
        // In func helloHandler:
        apiKey := os.Getenv("API_KEY_ENV")
        if apiKey != "" {
            // In a real app, NEVER log the actual secret value like this.
            // This is just for demonstration.
            fmt.Fprintf(w, "Secret API Key (demo): %s (first 5 chars)\n", apiKey[:min(5, len(apiKey))])
        }

        // Helper function (add this outside main/helloHandler)
        func min(a, b int) int {
            if a < b {
                return a
            }
            return b
        }
        ```

3. **Commit Configuration and Secret Changes:**
    * `git add cmd/hello-test-base/main.go build/cloudbuild/cloudbuild.yaml`
    * `git commit -m "feat: Add env var config and basic secret handling"`
    * `git push`
    * Trigger the build (either via GitHub Actions or manually if you're using the Cloud Build trigger directly).
    * Verify the deployment. Check the Cloud Run service logs and access the URL to see if the greeting and the (partial) API key are displayed.

---
---

## Phase 6: Testing Enhancements

**Goal:** Improve testing practices within the CI pipeline, focusing on automated checks that provide confidence in deployments.

1. **Test Coverage Report (Optional but Recommended):**
    * **Goal:** To understand what percentage of your code is covered by tests.
    * **Action:** Modify the "Run Unit Tests" step in `build/cloudbuild/cloudbuild.yaml` to generate an HTML coverage report.

        ```yaml
        # In 'Run Unit Tests' step (inside args -> -c -> script block):
        # ...
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
          go tool cover -html=coverage.out -o coverage.html
          echo "Test coverage report generated: coverage.html"
        # ...
        ```

    * **Further Action (Optional):** You can choose to upload this `coverage.html` file as an artifact to a Google Cloud Storage bucket for later review. This would involve adding another step to your `cloudbuild.yaml` using `gsutil cp`.

        ```yaml
        # Example step to upload coverage report (add after test step)
        # - name: 'gcr.io/cloud-builders/gsutil'
        #   id: 'Upload Coverage Report'
        #   args:
        #     - 'cp'
        #     - 'coverage.html'
        #     - 'gs://your-build-artifacts-bucket/${BUILD_ID}/coverage.html' # Replace with your bucket name
        ```

        * If you add this, ensure your Cloud Build service account has `roles/storage.objectCreator` on the target GCS bucket.

2. **Simple Post-Deployment Smoke Test:**
    * **Goal:** To perform a quick, automated check that the newly deployed service is responsive and serving basic requests.
    * **Action:** Add a new step in `build/cloudbuild/cloudbuild.yaml` *after* the "Deploy to Cloud Run" step.

        ```yaml
        # After 'Deploy to Cloud Run' step
        - name: 'curlimages/curl:latest' # A small image with curl
          id: 'Smoke Test Deployed Service'
          entrypoint: /bin/sh # Or 'ash' if /bin/sh is not present in the image
          args:
            - -c # Use -c to execute a shell script
            - |
              set -e # Exit immediately if a command exits with a non-zero status.
              echo "Attempting to get service URL..."
              # Use --project to be explicit if Cloud Build SA has access to multiple projects
              SERVICE_URL=$$(gcloud run services describe ${_SERVICE_NAME} --platform=managed --region=${_GCP_REGION} --project=$PROJECT_ID --format='value(status.url)')

              if [ -z "$$SERVICE_URL" ]; then
                echo "ERROR: Failed to retrieve Service URL for ${_SERVICE_NAME}."
                exit 1
              fi
              echo "Service URL: $$SERVICE_URL"

              echo "Performing smoke test..."
              # Retry a few times in case the service is still starting up
              RETRY_COUNT=0
              MAX_RETRIES=5
              RETRY_DELAY=10 # seconds

              until curl --fail --silent --show-error --connect-timeout 5 --max-time 10 "$$SERVICE_URL"; do
                RETRY_COUNT=$$((RETRY_COUNT + 1))
                if [ $$RETRY_COUNT -ge $$MAX_RETRIES ]; then
                  echo "ERROR: Smoke test failed after $$MAX_RETRIES retries."
                  exit 1
                fi
                echo "Smoke test attempt $$RETRY_COUNT failed. Retrying in $$RETRY_DELAY seconds..."
                sleep $$RETRY_DELAY
              done
              echo "Smoke test passed for $$SERVICE_URL!"
        ```

        * **Note:**
            * `$$SERVICE_URL` and other `$$` variables are used because Cloud Build performs its own variable substitution first. The double dollar signs ensure that the shell script sees a single dollar sign and performs its own variable expansion.
            * The `set -e` command ensures the script exits if any command fails.
            * The retry loop makes the smoke test more resilient to initial cold starts of the Cloud Run service. Adjust `MAX_RETRIES` and `RETRY_DELAY` as needed.
            * Ensure `$PROJECT_ID` is available in this step. Cloud Build usually makes it available as a default substitution.

3. **Commit Testing Enhancements:**
    * `git add build/cloudbuild/cloudbuild.yaml`
    * `git commit -m "feat: Add test coverage generation and post-deployment smoke test"`
    * `git push`
    * Monitor the build. Check the logs for the "Run Unit Tests" step to confirm coverage generation (if you don't upload it, it's just generated within the build environment). Crucially, check the "Smoke Test Deployed Service" step logs to ensure it passes or provides useful error messages if it fails.

---
---

## Phase 7: Documentation

**Goal:** Ensure the project is well-documented for current and future developers (including your future self!). Good documentation is key to maintainability and ease of use.

1. **Create/Update `README.md`:**
    * Your project should have a `README.md` file in the root directory.
    * **Action:** Add or expand the following sections in your `README.md`:
        * **Project Overview:**
            * A brief description of `hello-test-base`.
            * What problem does it solve or what does it do?
        * **Prerequisites:**
            * List necessary tools and their versions (e.g., Go 1.22+, Docker, gcloud CLI, Make).
            * Link to installation guides if helpful.
        * **Getting Started / Local Development Setup:**
            * How to clone the repository.
            * Any one-time setup steps (e.g., `cp .env.example .env` if you add an example env file).
            * Key `Makefile` targets for local development:
                * `make fmt`: Format code.
                * `make lint`: Run linters.
                * `make test`: Run unit tests.
                * `make build-linux`: Build the binary for Linux.
                * `make run-local`: Build and run the service locally.
                * `make docker-build-local`: Build the Docker image locally.
                * `make docker-run-local`: Run the locally built Docker image.
                * Explain how to access the service locally (e.g., `http://localhost:8080`).
        * **Configuration:**
            * Explain how the application is configured (e.g., environment variables).
            * List important environment variables (e.g., `PORT`, `SERVICE_NAME`, `APP_GREETING`, `API_KEY_ENV`).
            * Mention how secrets are handled (e.g., via Google Secret Manager, mounted as environment variables in Cloud Run).
        * **Building:**
            * Briefly mention that `make build-linux` builds the binary and `make docker-build-local` builds the image.
        * **Testing:**
            * How to run tests (`make test`).
            * Mention test coverage generation (e.g., `coverage.out`, `coverage.html`).
        * **Deployment:**
            * Explain that CI/CD is handled by Cloud Build, triggered by pushes to the `main` branch (or via GitHub Actions).
            * Mention key Cloud Build substitutions (`_SERVICE_NAME`, `_GCP_REGION`, `_ARTIFACT_REGISTRY_REPO`).
            * (Optional) If you have a manual way to trigger deployments (e.g., a `make deploy-dev` target that calls `gcloud builds submit`), document it.
        * **Troubleshooting (Optional):**
            * Common issues and their solutions.
            * How to check logs (local Docker, Cloud Run).

2. **Review and Refine `Makefile` Help Text:**
    * Ensure the `make help` target in your `Makefile` accurately reflects all available and useful targets and provides a clear, concise description for each.

3. **Add Comments to Code and Configuration Files:**
    * Go through your Go code (`main.go`, test files), `Dockerfile`, and `cloudbuild.yaml`.
    * **Action:** Add comments to explain non-obvious logic, configuration choices, or important steps. Good comments make the codebase much easier to understand.

4. **Commit Documentation:**
    * `git add README.md Makefile` (if you updated help text)
    * `git commit -m "docs: Add comprehensive README and review Make targets"`
    * If you added significant comments to other files, add them to the commit as well.
    * `git push`

---
