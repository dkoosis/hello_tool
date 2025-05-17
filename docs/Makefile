# Makefile for hello-tool-base
# Uses 'fo' utility from ./scripts/fo for enhanced output formatting.

# Force bash shell for consistent behavior
SHELL := /bin/bash
.SHELLFLAGS := -e -o pipefail -c

# Specify phony targets
.PHONY: all tree build clean deps fmt lint golangci-lint lint-yaml lint-makefile test test-debug install-tools check-gomod check-line-length check deploy help check-vulns health-check

# --- fo Utility Configuration ---
# Path to your fo utility
FO := ./scripts/fo
FO_FLAGS :=
FO_PRINT_FLAGS :=
ifeq ($(CI),true)
    FO_FLAGS += --ci
    FO_PRINT_FLAGS += --ci
    # FO_FLAGS += --theme ascii_minimal # Optionally switch theme for CI
    # FO_PRINT_FLAGS += --theme ascii_minimal
endif

# --- Helper Macros using fo ---
# $(call FO_PRINT, type, icon_override, indent, message)
# if icon_override is empty, fo will use the icon defined for the type in its theme.
define FO_PRINT
    $(FO) $(FO_PRINT_FLAGS) print --type "$(1)" --icon "$(2)" --indent $(3) -- "$(4)"
endef

# $(call FO_RUN, label, command_to_run, [optional_fo_flags_for_this_run_command])
define FO_RUN
    $(FO) $(FO_FLAGS) -l "$(1)" $(3) -- $(2)
endef

# --- Variables ---
# Application specific variables
SERVICE_NAME := hello-tool-base
BINARY_NAME  := $(SERVICE_NAME)

# Go module path
MODULE_PATH  := github.com/dkoosis/hello-tool-base
CMD_PATH     := ./cmd/$(SERVICE_NAME)
SCRIPT_DIR   := ./scripts

# Build-time variables
LOCAL_VERSION     := $(shell git describe --tags --always --dirty --match=v* 2>/dev/null || echo "dev")
LOCAL_COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION           := $(or $(VERSION_ARG),$(LOCAL_VERSION))
COMMIT_HASH       := $(or $(COMMIT_ARG),$(LOCAL_COMMIT_HASH))
BUILD_DATE        := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS      := -ldflags "-s -w \
    -X $(MODULE_PATH)/internal/buildinfo.Version=$(VERSION) \
    -X $(MODULE_PATH)/internal/buildinfo.CommitHash=$(COMMIT_HASH) \
    -X $(MODULE_PATH)/internal/buildinfo.BuildDate=$(BUILD_DATE)"

# Google Cloud Platform variables
PROJECT_ID             := $(shell gcloud config get-value project 2>/dev/null)
GCP_REGION             := us-central1
ARTIFACT_REGISTRY_REPO := my-go-apps
CLOUDBUILD_CONFIG_PATH := ./build/cloudbuild/cloudbuild.yaml
GCLOUD_BUILD_SUBSTITUTIONS := _SERVICE_NAME=$(SERVICE_NAME),_REGION=$(GCP_REGION),_ARTIFACT_REGISTRY_REPO=$(ARTIFACT_REGISTRY_REPO),_MODULE_PATH=$(MODULE_PATH),_MAKEFILE_VERSION=$(LOCAL_VERSION),_MAKEFILE_COMMIT=$(LOCAL_COMMIT_HASH)

# Tool Versions (still needed for go install commands)
GOLANGCILINT_VERSION := latest
GOTESTSUM_VERSION    := latest

# Line length check configuration
WARN_LINES   := 350
FAIL_LINES   := 1500
GO_FILES     := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*")
YAML_FILES   := $(shell find . -type f \( -name "*.yaml" -o -name "*.yml" \) -not -path "./vendor/*" -not -path "./.git/*")

# --- Core Targets ---
# IMPORTANT: Recipe lines below MUST start with a TAB character.
all: lint-makefile tree check-gomod deps fmt golangci-lint lint-yaml check-line-length test build
	@$(call FO_PRINT,Success,✨,0,All general checks passed and build completed successfully!) # Keep sparkle for final summary

tree:
	@$(call FO_PRINT,H1,,0,PROJECT SETUP)
	# Removed: @$(call FO_PRINT,Info,,0,Generating project directory tree...)
	@mkdir -p ./docs
	@if ! command -v tree > /dev/null; then \
		$(call FO_PRINT,Warning,,1,'tree' command not found. Skipping tree generation.); \
		$(call FO_PRINT,Info,,2,To install on macOS: brew install tree); \
		$(call FO_PRINT,Info,,2,To install on Debian/Ubuntu: sudo apt-get install tree); \
	else \
		tree -F -I 'vendor|.git|.idea*|*.DS_Store|$(BINARY_NAME)|coverage.out' --dirsfirst > ./docs/project_directory_tree.txt && \
		$(call FO_PRINT,Success,,1,Project tree generated at ./docs/project_directory_tree.txt) || \
		$(call FO_PRINT,Error,,1,Failed to generate project tree); \
	fi

build:
	@$(call FO_PRINT,H1,,0,APPLICATION BUILD)
	@$(call FO_PRINT,H2,,0,Building $(BINARY_NAME) for linux/amd64...)
	@$(call FO_PRINT,Info,,1,Version: $(VERSION), Commit: $(COMMIT_HASH), BuildDate: $(BUILD_DATE))
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH) && \
		$(call FO_PRINT,Success,,1,Build successful: $(PWD)/$(BINARY_NAME)) || \
		($(call FO_PRINT,Error,,1,Build failed) && exit 1)

clean:
	@$(call FO_PRINT,H1,,0,CLEANUP)
	@$(call FO_PRINT,Info,,0,Cleaning build artifacts and caches...) # Using theme icon for Info
	@rm -f $(BINARY_NAME) coverage.out
	@go clean -cache -testcache
	@$(call FO_PRINT,Success,,1,Cleaned build artifacts and Go caches)

# --- Dependency Management ---
deps:
	@$(call FO_PRINT,H1,,0,DEPENDENCY MANAGEMENT)
	@$(call FO_RUN,Tidying Go module dependencies,go mod tidy -v)
	@$(call FO_RUN,Downloading Go module dependencies,go mod download)

check-gomod:
	# Removed: @$(call FO_PRINT,H2,,0,Checking go.mod module path...) # Assumes H1 already printed
	@if [ ! -f "go.mod" ]; then \
		$(call FO_PRINT,Error,,1,go.mod file is missing. Run: go mod init $(MODULE_PATH)); \
		exit 1; \
	fi
	@if [ -x "$(SCRIPT_DIR)/check_go_mod_path.sh" ]; then \
		"$(SCRIPT_DIR)/check_go_mod_path.sh" $(MODULE_PATH) && \
		$(call FO_PRINT,Success,,1,go.mod has correct module path (checked via script)) || \
		($(call FO_PRINT,Error,,1,go.mod path check script failed) && exit 1); \
	elif ! grep -q "^module $(MODULE_PATH)$$" go.mod; then \
		$(call FO_PRINT,Error,,1,go.mod has incorrect module path (basic check).); \
		$(call FO_PRINT,Info,,2,Expected: module $(MODULE_PATH)); \
		$(call FO_PRINT,Info,,2,Found:    $$(grep '^module' go.mod)); \
		exit 1; \
	else \
		$(call FO_PRINT,Success,,1,go.mod has correct module path (basic check).); \
	fi

# --- Quality & Testing ---
fmt: install-tools
	@$(call FO_PRINT,H1,,0,CODE QUALITY - FORMATTING & LINTING)
	@$(call FO_RUN,Formatting Go code (golangci-lint fmt),golangci-lint fmt ./...)
	@$(call FO_RUN,Tidying go.mod/go.sum after format,go mod tidy -v)

lint: # Basic go vet, part of the above section
	@$(call FO_RUN,Running basic Go linter (go vet),go vet ./...)

golangci-lint: install-tools # part of the above section
	@$(call FO_RUN,Running comprehensive Go linters (golangci-lint),golangci-lint run ./...)

lint-yaml: # part of the above section
	@$(call FO_PRINT,H2,,0,Linting YAML files...)
	@if ! command -v yamllint >/dev/null 2>&1; then \
		$(call FO_PRINT,Warning,,1,yamllint not found.); \
		if command -v pip >/dev/null 2>&1 || command -v pip3 >/dev/null 2>&1; then \
			PIP_CMD=$$(command -v pip || command -v pip3); \
			$(call FO_PRINT,Info,,1,Attempting to install yamllint using $$PIP_CMD...); \
			$$PIP_CMD install --user yamllint && \
			$(call FO_PRINT,Success,,1,yamllint installed successfully via $$PIP_CMD.) || \
			($(call FO_PRINT,Error,,1,Failed to install yamllint via $$PIP_CMD. Please install it manually.) && exit 1); \
		else \
			$(call FO_PRINT,Error,,1,pip/pip3 not found. Please install yamllint manually.); \
			exit 1; \
		fi; \
		if ! command -v yamllint >/dev/null 2>&1; then \
			$(call FO_PRINT,Warning,,1,yamllint installed but might not be in your PATH. Please check your environment.); \
			$(call FO_PRINT,Info,,2,Common paths include: ~/.local/bin - ensure this is in your PATH.); \
			exit 1; \
		fi; \
	else \
		$(call FO_PRINT,Success,,1,yamllint found: $$(command -v yamllint)); \
	fi
	@if [ -z "$(YAML_FILES)" ]; then \
		$(call FO_PRINT,Info,,1,No YAML files found to lint.); \
	else \
		$(call FO_PRINT,Info,,1,Checking files: $(YAML_FILES)); \
		yamllint $(YAML_FILES); \
		EXIT_CODE=$$?; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			$(call FO_PRINT,Success,,1,YAML linting passed.); \
		else \
			if [ $$EXIT_CODE -gt 1 ]; then \
				 $(call FO_PRINT,Error,,1,YAML linting failed with errors (exit code: $$EXIT_CODE). See output above.); \
				 exit $$EXIT_CODE; \
			else \
				 $(call FO_PRINT,Warning,,1,YAML linting reported warnings (exit code: $$EXIT_CODE). See output above.); \
				 $(call FO_PRINT,Success,,1,YAML linting completed with warnings.); \
			fi; \
		fi; \
	fi

test: install-tools
	@$(call FO_PRINT,H1,,0,QUALITY ASSURANCE - TESTS)
	@$(call FO_RUN,Running tests with gotestsum,gotestsum --format testdox -- -race -coverprofile=coverage.out -covermode=atomic ./...,-s)

test-debug: install-tools # part of the above section
	@$(call FO_PRINT,H2,,0,Running tests (verbose debug mode)...)
	@LOG_LEVEL=debug go test -v -race -count=1 -coverprofile=coverage.out ./... ; EXIT_CODE=$$? ; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			$(call FO_PRINT,Success,,1,Tests finished (check output for any failures if not explicitly marked)); \
		else \
			$(call FO_PRINT,Error,,1,Tests failed (see verbose output above)); \
			exit 1; \
		fi

check-line-length: # part of CODE QUALITY section
	@$(call FO_PRINT,H2,,0,Checking file lengths (warn > $(WARN_LINES), fail > $(FAIL_LINES))...)
	@if [ ! -x "$(SCRIPT_DIR)/check_file_length.sh" ]; then \
		$(call FO_PRINT,Error,,1,Error: Script '$(SCRIPT_DIR)/check_file_length.sh' not found or not executable.); \
		exit 1; \
	fi
	@$(SCRIPT_DIR)/check_file_length.sh $(WARN_LINES) $(FAIL_LINES) $(GO_FILES) || \
		($(call FO_PRINT,Warning,,1,Line length check reported issues (see script output above). This is a non-fatal warning.) && exit 0)
	@$(call FO_PRINT,Success,,1,Line length check completed.)

# --- Security Checks ---
check-vulns: install-tools
	@$(call FO_PRINT,H1,,0,SECURITY SCAN)
	@$(call FO_PRINT,H2,,0,Scanning for vulnerabilities with govulncheck...)
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		$(call FO_PRINT,Info,,1,govulncheck not found. Installing...); \
		go install golang.org/x/vuln/cmd/govulncheck@latest && \
		$(call FO_PRINT,Success,,1,govulncheck installed successfully.) || \
		($(call FO_PRINT,Error,,1,Failed to install govulncheck.) && exit 1); \
	else \
		$(call FO_PRINT,Success,,1,govulncheck is already installed.); \
	fi
	@govulncheck ./... ; EXIT_CODE=$$? ; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			$(call FO_PRINT,Success,,1,No known vulnerabilities found.); \
		else \
			$(call FO_PRINT,Error,,1,govulncheck found potential vulnerabilities. Deployment will be halted.); \
			exit 1; \
		fi

# --- Tooling & Setup ---
install-tools:
	@$(call FO_PRINT,H1,,0,DEVELOPMENT TOOLING)
	# Removed: @$(call FO_PRINT,H2,,0,Checking/installing required Go tools...)
	@$(call FO_PRINT,Info,,1,golangci-lint:)
	@if ! command -v golangci-lint >/dev/null 2>&1 || ([ "$(GOLANGCILINT_VERSION)" != "latest" ] && ! golangci-lint --version | grep -qF "$(GOLANGCILINT_VERSION)"); then \
		$(call FO_PRINT,Info,,2,Installing/Updating golangci-lint@$(GOLANGCILINT_VERSION)...) ;\
		GOBIN=$(shell go env GOPATH)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION) && \
		$(call FO_PRINT,Success,,2,golangci-lint installed/updated: $$(golangci-lint --version)) || \
		($(call FO_PRINT,Error,,2,golangci-lint installation/update failed) && exit 1); \
	else \
		$(call FO_PRINT,Success,,2,golangci-lint already installed: $$(golangci-lint --version)); \
	fi
	@$(call FO_PRINT,Info,,1,gotestsum:)
	@if ! command -v gotestsum >/dev/null 2>&1 || ([ "$(GOTESTSUM_VERSION)" != "latest" ] && ! gotestsum --version | grep -qF "$(GOTESTSUM_VERSION)"); then \
		$(call FO_PRINT,Info,,2,Installing/Updating gotestsum@$(GOTESTSUM_VERSION)...) ;\
		GOBIN=$(shell go env GOPATH)/bin go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION) && \
		$(call FO_PRINT,Success,,2,gotestsum installed/updated: $$(gotestsum --version)) || \
		($(call FO_PRINT,Error,,2,gotestsum installation/update failed) && exit 1); \
	else \
		$(call FO_PRINT,Success,,2,gotestsum already installed: $$(gotestsum --version)); \
	fi
	@$(call FO_PRINT,Success,,0,Go tools check/installation complete)

check: install-tools # part of DEVELOPMENT TOOLING section
	@$(call FO_PRINT,H2,,0,Checking for required tools and environment...)
	@$(call FO_PRINT,Info,,1,Go: $$(go version 2>/dev/null || echo "Not installed"))
	@$(call FO_PRINT,Info,,1,Tree: $$(tree --version 2>/dev/null | head -n1 || echo "Not installed"))
	@$(call FO_PRINT,Info,,1,YamlLint: $$(yamllint --version 2>/dev/null || echo "Not installed"))
	@if [ -x "$(SCRIPT_DIR)/check_go_bin_path.sh" ]; then \
		$(call FO_PRINT,Info,,1,Go Bin Path Script check:); \
		"$(SCRIPT_DIR)/check_go_bin_path.sh" && \
		$(call FO_PRINT,Success,,2,check_go_bin_path.sh passed) || \
		$(call FO_PRINT,Warning,,2,check_go_bin_path.sh reported issues or failed); \
	else \
		$(call FO_PRINT,Warning,,1,Go Bin Path Script not found or not executable: $(SCRIPT_DIR)/check_go_bin_path.sh); \
	fi
	@$(call FO_PRINT,Success,,0,Tool and environment check complete)

# --- Deployment ---
# add check-vulns
deploy: all lint-yaml
# --- Deployment ---
deploy: all check-vulns lint-yaml
	@$(call FO_PRINT,H1,,0,DEPLOYMENT TO GOOGLE CLOUD)
	@$(call FO_PRINT,H2,,0,Deploying $(SERVICE_NAME) to Google Cloud...)
	@if [ -z "$(PROJECT_ID)" ]; then \
		$(call FO_PRINT,Error,,1,Error: Google Cloud Project ID not found. Set via 'gcloud config set project YOUR_PROJECT_ID' or ensure gcloud is configured.); \
		exit 1; \
	fi
    @$(call FO_PRINT,Info,,1,Project: $(PROJECT_ID))
	@CONFIG_FILE_TO_USE=""; \
    if [ -f "$(CLOUDBUILD_CONFIG_PATH)" ]; then \
        CONFIG_FILE_TO_USE="$(CLOUDBUILD_CONFIG_PATH)"; \
    elif [ -f "./cloudbuild.yaml" ]; then \
        CONFIG_FILE_TO_USE="./cloudbuild.yaml"; \
    else \
        $(call FO_PRINT,Error,,1,Error: Cloud Build config file not found at '$(CLOUDBUILD_CONFIG_PATH)' or './cloudbuild.yaml'.); \
        exit 1; \
    fi; \
    @$(call FO_PRINT,Info,,1,Using Cloud Build config: '$$CONFIG_FILE_TO_USE')
	@$(call FO_PRINT,Info,,2,Substitutions: $(GCLOUD_BUILD_SUBSTITUTIONS))
	@$(call FO_PRINT,Info,,1,Submitting to Google Cloud Build and awaiting completion...)
	@$(call FO_PRINT,Info,,1,--- Begin gcloud output ---)
    @if gcloud builds submit . --config="$$CONFIG_FILE_TO_USE" --project=$(PROJECT_ID) --substitutions="$(GCLOUD_BUILD_SUBSTITUTIONS)"; then \
        $(call FO_PRINT,Info,,1,--- End gcloud output ---); \
        $(call FO_PRINT,Success,,1,Cloud Build completed successfully.); \
        $(call FO_PRINT,Info,,1,Monitor detailed build logs at the URL provided in the gcloud output or here: https://console.cloud.google.com/cloud-build/builds?project=$(PROJECT_ID)); \
        $(call FO_PRINT,H2,,0,Fetching deployed service URL...); \
        SERVICE_URL=$$(gcloud run services describe $(SERVICE_NAME) --platform=managed --region=$(GCP_REGION) --project=$(PROJECT_ID) --format="value(status.url)" 2>/dev/null); \
        if [ -n "$$SERVICE_URL" ]; then \
            $(call FO_PRINT,Success,,1,Service URL: $$SERVICE_URL); \
            $(MAKE) health-check HEALTH_CHECK_URL="$$SERVICE_URL/health" EXPECTED_VERSION="$(LOCAL_VERSION)" EXPECTED_COMMIT="$(LOCAL_COMMIT_HASH)"; \
        else \
            $(call FO_PRINT,Warning,,1,Could not retrieve service URL. Skipping health check. Please check the Cloud Run console.); \
        fi; \
    else \
        $(call FO_PRINT,Info,,1,--- End gcloud output ---); \
        $(call FO_PRINT,Error,,1,Cloud Build submission or execution failed. Review gcloud output above and build logs for details.); \
        exit 1; \
    fi

health-check:
ifndef HEALTH_CHECK_URL
	$(error HEALTH_CHECK_URL is not set. Usage: make health-check HEALTH_CHECK_URL=... EXPECTED_VERSION=... EXPECTED_COMMIT=...)
endif
ifndef EXPECTED_VERSION
	$(error EXPECTED_VERSION is not set)
endif
ifndef EXPECTED_COMMIT
	$(error EXPECTED_COMMIT is not set)
endif
	@$(call FO_PRINT,H1,,0,POST-DEPLOYMENT HEALTH CHECK)
	@$(call FO_PRINT,H2,,0,Performing health check on $(HEALTH_CHECK_URL)...)
	@$(call FO_PRINT,Info,,1,Expected Version: $(EXPECTED_VERSION), Commit: $(EXPECTED_COMMIT))
	@SUCCESS=false; \
	for i in $$(seq 1 5); do \
		$(call FO_PRINT,Info,,1,Health check attempt #$$i...); \
		RESPONSE=$$(curl -s -w "\n%{http_code}" "$(HEALTH_CHECK_URL)"); \
		HTTP_CODE=$$(echo "$$RESPONSE" | tail -n1); \
		BODY=$$(echo "$$RESPONSE" | sed '$$d'); \
		if [ "$$HTTP_CODE" -eq 200 ]; then \
			RESPONSE_VERSION=$$(echo "$$BODY" | jq -r .version 2>/dev/null || echo "jq_error_version"); \
			RESPONSE_COMMIT=$$(echo "$$BODY" | jq -r .commit 2>/dev/null || echo "jq_error_commit"); \
			if [ "$$RESPONSE_VERSION" = "$(EXPECTED_VERSION)" ] && [ "$$RESPONSE_COMMIT" = "$(EXPECTED_COMMIT)" ]; then \
				$(call FO_PRINT,Success,,1,Health check PASSED. Version and Commit match expected values.); \
				$(call FO_PRINT,Info,,2,Response: $$BODY); \
				SUCCESS=true; break; \
			else \
				$(call FO_PRINT,Warning,,1,Health check attempt #$$i: Service responded OK (200), but content mismatch.); \
				$(call FO_PRINT,Info,,2,Expected Version: $(EXPECTED_VERSION), Got: $$RESPONSE_VERSION); \
				$(call FO_PRINT,Info,,2,Expected Commit: $(EXPECTED_COMMIT), Got: $$RESPONSE_COMMIT); \
				$(call FO_PRINT,Info,,2,Full Response: $$BODY); \
			fi; \
		else \
			$(call FO_PRINT,Warning,,1,Health check attempt #$$i: Failed with HTTP status $$HTTP_CODE.); \
			$(call FO_PRINT,Info,,2,Response Body: $$BODY); \
		fi; \
		if [ $$i -lt 5 ]; then $(call FO_PRINT,Info,,2,Retrying in 10 seconds...); sleep 10; fi; \
	done; \
	if [ "$$SUCCESS" = "false" ]; then \
		$(call FO_PRINT,Error,,1,Post-deployment health check FAILED after 5 attempts.); \
		exit 1; \
	else \
		$(call FO_PRINT,Success,,1,Post-deployment health check SUCCEEDED.); \
	fi

debug-script-block:
	@echo "--- Debugging script block ---"
	@echo "CLOUDBUILD_CONFIG_PATH is: $(CLOUDBUILD_CONFIG_PATH)"
	@echo "FO is: $(FO)"
	@echo "FO_PRINT_FLAGS is: $(FO_PRINT_FLAGS)" # Added this for completeness
	@set -x; \
	CONFIG_FILE_TO_USE=""; \
	if [ -f "$(CLOUDBUILD_CONFIG_PATH)" ]; then \
	    CONFIG_FILE_TO_USE="$(CLOUDBUILD_CONFIG_PATH)"; \
	elif [ -f "./cloudbuild.yaml" ]; then \
	    CONFIG_FILE_TO_USE="./cloudbuild.yaml"; \
	else \
	    echo "Config file not found, about to call FO_PRINT directly"; \
	    "$(FO)" $(FO_PRINT_FLAGS) print --type "Error" --icon "" --indent 1 -- "Error: Cloud Build config file not found at '$(CLOUDBUILD_CONFIG_PATH)' or './cloudbuild.yaml'."; \
	    exit 1; \
	fi; \
	echo "Final CONFIG_FILE_TO_USE is '$$CONFIG_FILE_TO_USE'"; \
	set +x

# --- Help ---
# Keeping help target with direct printf for now.
help:
	@printf "\033[1m\033[0;34m%-20s %s\033[0m\n" "$(SERVICE_NAME) Makefile" "Development Targets (using fo)"
	@printf "\033[0;34m-------------------------------------------------------------\033[0m\n"
	@printf "  %-25s %s\n" "all" "Run general checks, format, tests, and build (default)"
	@printf "  %-25s %s\n" "build" "Build the application binary ($(BINARY_NAME)) for Linux"
	@printf "  %-25s %s\n" "clean" "Clean build artifacts and caches"
	@printf "  %-25s %s\n" "deps" "Tidy and download Go module dependencies"
	@printf "  %-25s %s\n" "install-tools" "Install/update required Go tools (golangci-lint, gotestsum)"
	@printf "  %-25s %s\n" "check" "Check if required tools (Go, tree, yamllint, etc.) are installed"
	@printf "  %-25s %s\n" "check-gomod" "Check if go.mod module path is correct"
	@printf "\n\033[1mCode Quality & Testing:\033[0m\n"
	@printf "  %-25s %s\n" "fmt" "Format code using golangci-lint fmt"
	@printf "  %-25s %s\n" "lint" "Run basic 'go vet' checks (quick pre-check)"
	@printf "  %-25s %s\n" "golangci-lint" "Run comprehensive linters with golangci-lint"
	@printf "  %-25s %s\n" "lint-yaml" "Lint YAML files (*.yaml, *.yml) using yamllint"
	@printf "  %-25s %s\n" "test" "Run tests using gotestsum"
	@printf "  %-25s %s\n" "test-debug" "Run tests with verbose output (go test -v)"
	@printf "  %-25s %s\n" "check-line-length" "Check Go file line count (W:$(WARN_LINES), F:$(FAIL_LINES))"
	@printf "\n\033[1mSecurity:\033[0m\n"
	@printf "  %-25s %s\n" "check-vulns" "Scan for known vulnerabilities (fails build if found)"
	@printf "\n\033[1mDeployment:\033[0m\n"
	@printf "  %-25s %s\n" "deploy" "Run all checks, build, then deploy to GCP, show URL & health check"
	@printf "  %-25s %s\n" "health-check" "Perform post-deployment health check (requires params)"
	@printf "\n\033[1mOther:\033[0m\n"
	@printf "  %-25s %s\n" "tree" "Generate project directory tree view in ./docs/"
	@printf "  %-25s %s\n" "help" "Display this help message"
	@printf "\033[0;34m-------------------------------------------------------------\033[0m\n\n"

