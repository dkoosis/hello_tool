# Configuration variables for the hello-tool-base project
# makefiles/config.mk

# --- Application specific variables ---
SERVICE_NAME := hello-tool-base
BINARY_NAME  := $(SERVICE_NAME)

# --- Go module and path configuration ---
MODULE_PATH  := github.com/dkoosis/hello-tool-base
CMD_PATH     := ./cmd/$(SERVICE_NAME)
SCRIPT_DIR   := ./scripts

# --- Build-time versioning ---
# These are determined locally from git for the `deploy` target and local `build`
# For builds within Docker triggered by this Makefile, VERSION_ARG and COMMIT_ARG can be passed.
LOCAL_VERSION := $(shell git describe --tags --always --dirty --match=v* 2>/dev/null || echo "dev")
LOCAL_COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Effective version and commit to use in LDFLAGS.
# Prioritize args passed, then fallback to locally determined git values.
VERSION      := $(or $(VERSION_ARG),$(LOCAL_VERSION))
COMMIT_HASH  := $(or $(COMMIT_ARG),$(LOCAL_COMMIT_HASH))
BUILD_DATE   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# --- LDFLAGS for injecting build information ---
LDFLAGS      := -ldflags "-s -w \
                -X $(MODULE_PATH)/internal/buildinfo.Version=$(VERSION) \
                -X $(MODULE_PATH)/internal/buildinfo.CommitHash=$(COMMIT_HASH) \
                -X $(MODULE_PATH)/internal/buildinfo.BuildDate=$(BUILD_DATE)"

# --- Google Cloud Platform variables ---
PROJECT_ID   := $(shell gcloud config get-value project 2>/dev/null)
GCP_REGION   := us-central1
ARTIFACT_REGISTRY_REPO := my-go-apps
CLOUDBUILD_CONFIG_PATH := ./build/cloudbuild/cloudbuild.yaml

# --- Tool versions ---
GOLANGCILINT_VERSION := latest # Or a specific version e.g., v1.58.0
GOTESTSUM_VERSION    := latest # Or a specific version e.g., v1.11.0

# --- Line length check configuration ---
WARN_LINES   := 350  # Warn if lines exceed this
FAIL_LINES   := 1500 # Fail if lines exceed this

# --- File selectors ---
GO_FILES     := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*")
YAML_FILES   := $(shell find . -type f \( -name "*.yaml" -o -name "*.yml" \) -not -path "./vendor/*" -not -path "./.git/*")

# --- Cloud build substitutions ---
# Use a function to trim whitespace from each variable
define trim
$(strip $(1))
endef

GCLOUD_BUILD_SUBSTITUTIONS := _SERVICE_NAME=$(call trim,$(SERVICE_NAME)),_REGION=$(call trim,$(GCP_REGION)),_ARTIFACT_REGISTRY_REPO=$(call trim,$(ARTIFACT_REGISTRY_REPO)),_MODULE_PATH=$(call trim,$(MODULE_PATH)),_MAKEFILE_VERSION=$(call trim,$(LOCAL_VERSION)),_MAKEFILE_COMMIT=$(call trim,$(LOCAL_COMMIT_HASH))
