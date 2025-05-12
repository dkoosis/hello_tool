# Color and formatting definitions
# makefiles/colors.mk

# --- Colors for output formatting ---
RESET   := \033[0m
BOLD    := \033[1m
GREEN   := \033[0;32m
YELLOW  := \033[0;33m
RED     := \033[0;31m
BLUE    := \033[0;34m
NC      := $(RESET) # No Color Alias

# --- Icons for visually distinct output ---
ICON_START := $(BLUE)▶$(NC)
ICON_OK    := $(GREEN)✓$(NC)
ICON_WARN  := $(YELLOW)⚠$(NC)
ICON_FAIL  := $(RED)✗$(NC)
ICON_INFO  := $(BLUE)ℹ$(NC)

# --- Formatting Strings for Alignment ---
LABEL_FMT  := "   %-25s" # Indent 3, Pad label to 25 chars, left-aligned
