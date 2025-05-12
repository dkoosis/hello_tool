# Tool installation and verification
# makefiles/tools.mk

.PHONY: install-tools check-env

# --- Tool Installation ---
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

# --- Environment Checks ---
check-env: install-tools
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
	@printf $(LABEL_FMT) "yamllint:"
	@if command -v yamllint >/dev/null 2>&1; then \
		printf "$(ICON_OK) $(GREEN)Installed$(NC) ($$(yamllint --version))\n"; \
	else \
		printf "$(ICON_WARN) $(YELLOW)Not installed$(NC) (needed for 'lint-yaml', will attempt install)\n"; \
	fi
	@printf $(LABEL_FMT) "Go Bin Path:"
	@if [ -x "$(SCRIPT_DIR)/check_go_bin_path.sh" ]; then \
		"$(SCRIPT_DIR)/check_go_bin_path.sh" && \
		printf "$(ICON_OK) $(GREEN)GOBIN is in PATH$(NC)\n" || \
		printf "$(ICON_WARN) $(YELLOW)GOBIN path issue detected$(NC)\n"; \
	else \
		printf "$(ICON_WARN) $(YELLOW)check_go_bin_path.sh not found or not executable$(NC)\n"; \
	fi
	@# Check for govulncheck (but don't fail if not found)
	@printf $(LABEL_FMT) "govulncheck:"
	@if command -v govulncheck >/dev/null 2>&1; then \
		printf "$(ICON_OK) $(GREEN)Installed$(NC)\n"; \
	else \
		printf "$(ICON_WARN) $(YELLOW)Not installed$(NC) (will be auto-installed when needed)\n"; \
	fi
	@printf "  $(ICON_OK) $(GREEN)Tool and environment check complete$(NC)\n"
	@printf "\n"

# --- Vulnerability Scanning ---
.PHONY: check-vulns

check-vulns:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Scanning for vulnerabilities with govulncheck...$(NC)\n"
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		printf "  $(ICON_INFO) $(YELLOW)govulncheck not found. Installing...$(NC)\n"; \
		go install golang.org/x/vuln/cmd/govulncheck@latest && \
		printf "  $(ICON_OK) $(GREEN)govulncheck installed successfully.$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)Failed to install govulncheck.$(NC)\n" && exit 1); \
	else \
		printf "  $(ICON_OK) $(GREEN)govulncheck is already installed.$(NC)\n"; \
	fi
	@govulncheck ./... && \
		printf "  $(ICON_OK) $(GREEN)No known vulnerabilities found.$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)$(BOLD)govulncheck found potential vulnerabilities. Deployment will be halted.$(NC)\n" && exit 1)
	@printf "\n"
