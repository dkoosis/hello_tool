# Test targets
# makefiles/test.mk

.PHONY: test test-debug

# --- Standard Testing ---
test: install-tools
	@printf "$(ICON_START) $(BOLD)$(BLUE)Running tests with gotestsum...$(NC)\n"
	@gotestsum --format testdox -- -race -coverprofile=coverage.out -covermode=atomic ./... && \
		printf "  $(ICON_OK) $(GREEN)Tests passed$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
	@printf "\n"

# --- Debug Testing ---
test-debug: install-tools
	@printf "$(ICON_START) $(BOLD)$(YELLOW)Running tests (verbose debug mode)...$(NC)\n"
	@LOG_LEVEL=debug go test -v -race -count=1 -coverprofile=coverage.out ./... && \
		printf "  $(ICON_OK) $(GREEN)Tests finished (check output for failures)$(NC)\n" || \
		(printf "  $(ICON_FAIL) $(RED)$(BOLD)Tests failed$(NC)\n" && exit 1)
	@printf "\n"
