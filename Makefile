.PHONY: help generate build test lint clean install run validate setup install-tools tidy

# Variables
BINARY_NAME=pulumicost-mcp
GO_VERSION=1.24
BUILD_DIR=bin
DESIGN_DIR=design
GEN_DIR=gen
CMD_DIR=cmd/pulumicost-mcp

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

##@ Help

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(COLOR_BOLD)Usage:$(COLOR_RESET)\n  make $(COLOR_BLUE)<target>$(COLOR_RESET)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(COLOR_BLUE)%-15s$(COLOR_RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(COLOR_BOLD)%s$(COLOR_RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

generate: ## Generate Goa code from design
	@echo "$(COLOR_GREEN)Generating Goa code...$(COLOR_RESET)"
	@goa gen github.com/rshade/pulumicost-mcp/design
	@echo "$(COLOR_GREEN)✓ Code generation complete$(COLOR_RESET)"

build: generate ## Build the server binary
	@echo "$(COLOR_GREEN)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "$(COLOR_GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

run: build ## Run the server locally
	@echo "$(COLOR_GREEN)Starting $(BINARY_NAME)...$(COLOR_RESET)"
	@$(BUILD_DIR)/$(BINARY_NAME) --config config.yaml

install: build ## Install binary to $GOPATH/bin
	@echo "$(COLOR_GREEN)Installing $(BINARY_NAME)...$(COLOR_RESET)"
	@go install $(CMD_DIR)/main.go
	@echo "$(COLOR_GREEN)✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(COLOR_RESET)"

##@ Testing

test: ## Run all tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	@go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report: coverage.html$(COLOR_RESET)"

test-unit: ## Run only unit tests
	@echo "$(COLOR_GREEN)Running unit tests...$(COLOR_RESET)"
	@go test -v -short ./...

test-integration: ## Run only integration tests
	@echo "$(COLOR_GREEN)Running integration tests...$(COLOR_RESET)"
	@go test -v -run Integration ./...

bench: ## Run benchmarks
	@echo "$(COLOR_GREEN)Running benchmarks...$(COLOR_RESET)"
	@go test -bench=. -benchmem ./...

##@ Code Quality

lint: ## Run linters
	@echo "$(COLOR_GREEN)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed. Run 'make install-tools'$(COLOR_RESET)"; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	@gofumpt -w .
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

vet: ## Run go vet
	@echo "$(COLOR_GREEN)Running go vet...$(COLOR_RESET)"
	@go vet ./...

validate: lint test ## Run all validation (lint + test)
	@echo "$(COLOR_GREEN)✓ All validation passed$(COLOR_RESET)"

##@ Cleanup

clean: ## Clean generated files and build artifacts
	@echo "$(COLOR_GREEN)Cleaning...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(GEN_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ Cleaned$(COLOR_RESET)"

tidy: ## Tidy go modules
	@echo "$(COLOR_GREEN)Tidying go modules...$(COLOR_RESET)"
	@go mod tidy
	@echo "$(COLOR_GREEN)✓ Modules tidied$(COLOR_RESET)"

##@ Setup

setup: install-tools tidy generate ## Setup development environment
	@echo "$(COLOR_GREEN)✓ Development environment ready$(COLOR_RESET)"

install-tools: ## Install development tools
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	@go install goa.design/goa/v3/cmd/goa@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install mvdan.cc/gofumpt@latest
	@go install gotest.tools/gotestsum@latest
	@echo "$(COLOR_GREEN)✓ Tools installed$(COLOR_RESET)"

check-go-version: ## Check Go version
	@echo "Checking Go version..."
	@GO_CURRENT=$$(go version | awk '{print $$3}' | sed 's/go//'); \
	if [ "$$(printf '%s\n' "$(GO_VERSION)" "$$GO_CURRENT" | sort -V | head -n1)" != "$(GO_VERSION)" ]; then \
		echo "$(COLOR_YELLOW)Warning: Go $(GO_VERSION) or higher required. Current: $$GO_CURRENT$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_GREEN)✓ Go version OK: $$GO_CURRENT$(COLOR_RESET)"; \
	fi

##@ Docker

docker-build: ## Build Docker image
	@echo "$(COLOR_GREEN)Building Docker image...$(COLOR_RESET)"
	@docker build -t $(BINARY_NAME):latest .
	@echo "$(COLOR_GREEN)✓ Docker image built$(COLOR_RESET)"

docker-run: docker-build ## Run Docker container
	@echo "$(COLOR_GREEN)Running Docker container...$(COLOR_RESET)"
	@docker run -p 8080:8080 -v $(PWD)/config.yaml:/app/config.yaml $(BINARY_NAME):latest

##@ Deployment

deploy-local: build ## Deploy locally with systemd
	@echo "$(COLOR_GREEN)Deploying locally...$(COLOR_RESET)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /opt/$(BINARY_NAME)/bin/
	@sudo systemctl restart $(BINARY_NAME)
	@echo "$(COLOR_GREEN)✓ Deployed and restarted$(COLOR_RESET)"

deploy-k8s: ## Deploy to Kubernetes
	@echo "$(COLOR_GREEN)Deploying to Kubernetes...$(COLOR_RESET)"
	@kubectl apply -f deployments/k8s/
	@echo "$(COLOR_GREEN)✓ Deployed to Kubernetes$(COLOR_RESET)"

##@ Documentation

docs: ## Generate documentation
	@echo "$(COLOR_GREEN)Generating documentation...$(COLOR_RESET)"
	@godoc -http=:6060 &
	@echo "$(COLOR_GREEN)✓ Documentation server started at http://localhost:6060$(COLOR_RESET)"

##@ Release

release: validate build ## Create release artifacts
	@echo "$(COLOR_GREEN)Creating release artifacts...$(COLOR_RESET)"
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)/main.go
	@GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)/main.go
	@GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)/main.go
	@GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)/main.go
	@GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)/main.go
	@echo "$(COLOR_GREEN)✓ Release artifacts created in dist/$(COLOR_RESET)"

.DEFAULT_GOAL := help
