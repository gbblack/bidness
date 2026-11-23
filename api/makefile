# Makefile for api (name pending)
#
# Usage: make [target]
# Run 'make help' to see all available targets

BINARY_NAME := api

.PHONY: all build sqlc test test-verbose test-coverage run deps tidy fmt vet lint clean help

# Colors for terminal output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target

all: help 

## Build targets

build: ## Build the api binary
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	go build -o bin/$(BINARY_NAME) ./cmd/api
	@echo "$(GREEN)✓ Built $(BINARY_NAME)$(NC)"

sqlc: ## Generate sqlc code
	@echo "$(GREEN)Genrating $(BINARY_NAME)...$(NC)"
	sqlc generate
	@echo "$(GREEN)✓ Code generated"

## Test targets

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test ./...

test-verbose: ## Run all tests with verbose output
	@echo "$(GREEN)Running tests (verbose)...$(NC)"
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -cover ./...

## Run targets

run: build ## Build and run api
	@echo "$(GREEN)Starting api server...$(NC)"
	./bin/$(BINARY_NAME)

## Setup targets
# TODO: create setup scripts when project is more mature

## Dependency targets

deps: ## Download Go module dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

tidy: ## Tidy Go modules
	@echo "$(GREEN)Tidying modules...$(NC)"
	go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

## Code quality targets

fmt: ## Format Go code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

lint: fmt vet ## Run all linters (fmt + vet)
	@echo "$(GREEN)✓ All linting passed$(NC)"

# install-hooks: ## Install git pre-commit hooks
# 	@echo "$(GREEN)Installing pre-commit hook...$(NC)"
# 	@cp scripts/pre-commit .git/hooks/pre-commit
# 	@chmod +x .git/hooks/pre-commit
# 	@echo "$(GREEN)✓ Pre-commit hook installed$(NC)"
# TODO: keep this in case of pre commit hooks in future

## Cleanup targets

clean: ## Remove build artifacts and test data
	@echo "$(YELLOW)Cleaning up...$(NC)"
	@rm -f $(BINARY_NAME)
	@echo "$(GREEN)✓ Cleaned$(NC)"
# TODO: add this back in when ready -> @rm -rf data/ test_data/

# clean-all: clean ## Remove all artifacts including /tmp/data
# 	@echo "$(YELLOW)Removing $(DATA_DIR)...$(NC)"
# 	@rm -rf $(DATA_DIR)
# 	@echo "$(GREEN)✓ All cleaned$(NC)"
# TODO: add back in later

## Help

help: ## Show this help message
	@echo "bidness - name pending"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

# TODO: add this in when ready
# 	@echo "Quick start:"
# 	@echo "  make setup    # First-time setup"
# 	@echo "  make run      # Build and run server"
# 	@echo ""