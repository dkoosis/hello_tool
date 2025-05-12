# Build targets
# makefiles/build.mk

.PHONY: build clean deps tree

# --- Building ---
build:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Building $(BINARY_NAME) for linux/amd64...$(NC)\n"
	@printf "  $(ICON_INFO) Using Version: $(VERSION), Commit: $(COMMIT_HASH), BuildDate: $(BUILD_DATE)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH) && \
		printf "  $(ICON_OK) $(GREEN)Build successful: $(PWD)/$(BINARY_NAME)$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)Build failed$(NC)\n" && exit 1)
	@printf "\n"

# --- Cleaning ---
clean:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Cleaning build artifacts...$(NC)\n"
	@rm -f $(BINARY_NAME) coverage.out
	@go clean -cache -testcache
	@printf "  $(ICON_OK) $(GREEN)Cleaned$(NC)\n"
	@printf "\n"

# --- Dependency Management ---
deps:
	@printf "$(ICON_START) $(BOLD)$(BLUE)Synchronizing dependencies...$(NC)\n"
	@printf "  $(ICON_INFO) Running go mod tidy...\n"
	@go mod tidy -v && printf "  $(ICON_OK) $(GREEN)Dependencies tidied successfully$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)Failed to tidy dependencies$(NC)\n" && exit 1)
	@printf "  $(ICON_INFO) Running go mod download...\n"
	@go mod download && printf "  $(ICON_OK) $(GREEN)Dependencies downloaded successfully$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)Failed to download dependencies$(NC)\n" && exit 1)
	@printf "\n"

# --- Project Tree Generation ---
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
