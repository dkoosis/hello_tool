# Code quality validation checks
# makefiles/checks.mk

.PHONY: check-gomod fmt lint golangci-lint lint-yaml check-line-length

# --- Go Module Path Check ---
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

# --- Code Formatting ---
fmt: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Formatting code using golangci-lint fmt...$(NC)\n"
	@golangci-lint fmt ./... && \
		printf "  $(ICON_OK) $(GREEN)Code formatted$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)$(BOLD)Formatting failed (see errors above)$(NC)\n" && exit 1)
	@# Ensure go.mod and go.sum are tidy after formatting, as some formatters might change imports.
	@go mod tidy -v > /dev/null
	@printf "\n"

# --- Basic Linting ---
lint:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running basic linter (go vet)...$(NC)\n"
	@go vet ./... && \
		printf "  $(ICON_OK) $(GREEN)go vet passed$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)go vet found issues$(NC)\n" && exit 1)
	@printf "\n"

# --- Comprehensive Linting ---
golangci-lint: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running golangci-lint...$(NC)\n"
	@golangci-lint run ./... && \
		printf "  $(ICON_OK) $(GREEN)golangci-lint passed$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)$(BOLD)golangci-lint failed (see errors above)$(NC)\n" && exit 1)
	@printf "\n"

# --- YAML Linting ---
lint-yaml:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Linting YAML files...$(NC)\n"
	@# Check if yamllint is installed
	@if ! command -v yamllint >/dev/null 2>&1; then \
		printf "  $(ICON_WARN) $(YELLOW)yamllint not found.$(NC)\n"; \
		if command -v pip >/dev/null 2>&1; then \
			printf "  $(ICON_INFO) $(YELLOW)Attempting to install yamllint using pip...$(NC)\n"; \
			pip install --user yamllint && \
			printf "  $(ICON_OK) $(GREEN)yamllint installed successfully via pip.$(NC)\n" || \
			(printf "  $(ICON_FAIL) $(RED)Failed to install yamllint via pip. Please install it manually.$(NC)\n" && exit 1); \
		elif command -v pip3 >/dev/null 2>&1; then \
			printf "  $(ICON_INFO) $(YELLOW)Attempting to install yamllint using pip3...$(NC)\n"; \
			pip3 install --user yamllint && \
			printf "  $(ICON_OK) $(GREEN)yamllint installed successfully via pip3.$(NC)\n" || \
			(printf "  $(ICON_FAIL) $(RED)Failed to install yamllint via pip3. Please install it manually.$(NC)\n" && exit 1); \
		else \
			printf "  $(ICON_FAIL) $(RED)pip/pip3 not found. Please install yamllint manually (e.g., 'pip install yamllint' or using your system package manager).$(NC)\n"; \
			exit 1; \
		fi; \
		echo "# Ensure yamllint installed via pip is in PATH (common issue)"; \
		if ! command -v yamllint >/dev/null 2>&1; then \
			printf "  $(ICON_WARN) $(YELLOW)yamllint installed but might not be in your PATH. Please check your environment.$(NC)\n"; \
			printf "  $(ICON_INFO) $(YELLOW)Common paths include: ~/.local/bin - ensure this is in your PATH environment variable.$(NC)\n"; \
			exit 1; \
		fi; \
	else \
		printf "  $(ICON_OK) $(GREEN)yamllint found: $$(command -v yamllint)$(NC)\n"; \
	fi
	@# Check if there are any YAML files to lint
	@if [ -z "$(YAML_FILES)" ]; then \
		printf "  $(ICON_INFO) $(YELLOW)No YAML files found to lint.$(NC)\n"; \
	else \
		printf "  $(ICON_INFO) Found $(words $(YAML_FILES)) YAML files to check\n"; \
		yamllint $(YAML_FILES); \
		EXIT_CODE=$$?; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			printf "  $(ICON_OK) $(GREEN)YAML linting passed.$(NC)\n"; \
		else \
			printf "  $(ICON_FAIL) $(RED)$(BOLD)YAML linting failed (exit code: $$EXIT_CODE). See errors above.$(NC)\n"; \
			exit $$EXIT_CODE; \
		fi; \
	fi
	@printf "\n"

# --- File Length Checks ---
check-line-length:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Checking file lengths (warn > $(WARN_LINES), fail > $(FAIL_LINES))...$(NC)\n"
	@if [ ! -x "$(SCRIPT_DIR)/check_file_length.sh" ]; then \
		printf "   $(ICON_FAIL) $(RED)Error: Script '$(SCRIPT_DIR)/check_file_length.sh' not found or not executable.$(NC)\n"; \
		exit 1; \
	fi
	@$(SCRIPT_DIR)/check_file_length.sh $(WARN_LINES) $(FAIL_LINES) $(GO_FILES) || \
		(printf "   $(ICON_WARN) $(YELLOW)Line length check reported issues (see script output above)$(NC)\n" && exit 0)
	@printf "\n"
