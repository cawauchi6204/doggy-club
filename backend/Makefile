# DoggyClub Backend Makefile

# Variables
GO_VERSION := 1.21
APP_NAME := doggyclub-backend
BUILD_DIR := ./bin
DOCKER_COMPOSE := docker-compose

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help build test test-unit test-integration test-coverage run dev clean docker-up docker-down migrate lint fmt vet security deps

# Default target
help: ## Show this help message
	@echo "$(BLUE)DoggyClub Backend - Available Commands$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

# Build commands
build: ## Build the application
	@echo "$(BLUE)Building application...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux ./cmd/api
	@echo "$(GREEN)Linux build completed: $(BUILD_DIR)/$(APP_NAME)-linux$(RESET)"

# Test commands
test: test-unit test-integration ## Run all tests

test-all: test-unit test-integration test-e2e test-load ## Run all tests including E2E and load tests

test-unit: ## Run unit tests
	@echo "$(BLUE)Running unit tests...$(RESET)"
	@go test -v -race -timeout=30s ./pkg/services/... ./pkg/utils/...
	@echo "$(GREEN)Unit tests completed$(RESET)"

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(RESET)"
	@echo "$(YELLOW)Make sure test database is running$(RESET)"
	@go test -v -race -timeout=60s ./pkg/handlers/... -tags=integration
	@echo "$(GREEN)Integration tests completed$(RESET)"

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

test-benchmark: ## Run benchmark tests
	@echo "$(BLUE)Running benchmark tests...$(RESET)"
	@go test -bench=. -benchmem ./pkg/services/...
	@echo "$(GREEN)Benchmark tests completed$(RESET)"

test-load: ## Run load tests (requires test environment)
	@echo "$(BLUE)Running load tests...$(RESET)"
	@go test -v -timeout=600s ./tests/performance/...
	@echo "$(GREEN)Load tests completed$(RESET)"

test-e2e: ## Run end-to-end tests
	@echo "$(BLUE)Running E2E tests...$(RESET)"
	@echo "$(YELLOW)Make sure application is running on localhost:8080$(RESET)"
	@go test -v -timeout=300s ./tests/e2e/...
	@echo "$(GREEN)E2E tests completed$(RESET)"

# Development commands
run: ## Run the application
	@echo "$(BLUE)Starting application...$(RESET)"
	@go run ./cmd/api

dev: ## Run with hot reload (requires air)
	@echo "$(BLUE)Starting development server with hot reload...$(RESET)"
	@air -c .air.toml

# Database commands
migrate: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(RESET)"
	@go run ./cmd/migrate
	@echo "$(GREEN)Migrations completed$(RESET)"

migrate-down: ## Rollback database migrations
	@echo "$(BLUE)Rolling back database migrations...$(RESET)"
	@go run ./cmd/migrate -down
	@echo "$(GREEN)Rollback completed$(RESET)"

seed: ## Seed database with test data
	@echo "$(BLUE)Seeding database...$(RESET)"
	@go run ./cmd/seed
	@echo "$(GREEN)Database seeded$(RESET)"

# Docker commands
docker-up: ## Start Docker containers
	@echo "$(BLUE)Starting Docker containers...$(RESET)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Docker containers started$(RESET)"

docker-down: ## Stop Docker containers
	@echo "$(BLUE)Stopping Docker containers...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Docker containers stopped$(RESET)"

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(RESET)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)Docker image built$(RESET)"

docker-test: ## Run tests in Docker
	@echo "$(BLUE)Running tests in Docker...$(RESET)"
	@$(DOCKER_COMPOSE) -f docker-compose.test.yml up --build --abort-on-container-exit
	@$(DOCKER_COMPOSE) -f docker-compose.test.yml down
	@echo "$(GREEN)Docker tests completed$(RESET)"

# Code quality commands
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(RESET)"
	@golangci-lint run ./...
	@echo "$(GREEN)Linting completed$(RESET)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted$(RESET)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)Vet completed$(RESET)"

security: ## Run security scanner
	@echo "$(BLUE)Running security scanner...$(RESET)"
	@gosec ./...
	@echo "$(GREEN)Security scan completed$(RESET)"

# Dependency commands
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(RESET)"

deps-update: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(RESET)"

deps-verify: ## Verify dependencies
	@echo "$(BLUE)Verifying dependencies...$(RESET)"
	@go mod verify
	@echo "$(GREEN)Dependencies verified$(RESET)"

# Cleanup commands
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleanup completed$(RESET)"

clean-cache: ## Clean Go cache
	@echo "$(BLUE)Cleaning Go cache...$(RESET)"
	@go clean -cache -modcache -testcache
	@echo "$(GREEN)Cache cleaned$(RESET)"

# CI/CD commands
ci: deps fmt vet lint test-coverage security ## Run CI pipeline
	@echo "$(GREEN)CI pipeline completed successfully$(RESET)"

ci-quick: deps fmt vet test-unit ## Run quick CI checks
	@echo "$(GREEN)Quick CI checks completed$(RESET)"

# Performance commands
profile-cpu: ## Run CPU profiling
	@echo "$(BLUE)Running CPU profiling...$(RESET)"
	@go test -cpuprofile=cpu.prof -bench=. ./pkg/services/...
	@echo "$(GREEN)CPU profiling completed: cpu.prof$(RESET)"

profile-mem: ## Run memory profiling
	@echo "$(BLUE)Running memory profiling...$(RESET)"
	@go test -memprofile=mem.prof -bench=. ./pkg/services/...
	@echo "$(GREEN)Memory profiling completed: mem.prof$(RESET)"

# Documentation commands
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(RESET)"
	@godoc -http=:6060
	@echo "$(GREEN)Documentation server started at http://localhost:6060$(RESET)"

# Environment setup
setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@echo "Installing dependencies..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "$(GREEN)Development environment setup completed$(RESET)"

# Database specific commands for testing
test-db-up: ## Start test database
	@echo "$(BLUE)Starting test database...$(RESET)"
	@docker run -d --name doggyclub-test-db \
		-e POSTGRES_DB=doggyclub_test \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5433:5432 \
		postgres:15-alpine
	@echo "$(GREEN)Test database started on port 5433$(RESET)"

test-db-down: ## Stop test database
	@echo "$(BLUE)Stopping test database...$(RESET)"
	@docker stop doggyclub-test-db || true
	@docker rm doggyclub-test-db || true
	@echo "$(GREEN)Test database stopped$(RESET)"

test-redis-up: ## Start test Redis
	@echo "$(BLUE)Starting test Redis...$(RESET)"
	@docker run -d --name doggyclub-test-redis \
		-p 6380:6379 \
		redis:7-alpine
	@echo "$(GREEN)Test Redis started on port 6380$(RESET)"

test-redis-down: ## Stop test Redis
	@echo "$(BLUE)Stopping test Redis...$(RESET)"
	@docker stop doggyclub-test-redis || true
	@docker rm doggyclub-test-redis || true
	@echo "$(GREEN)Test Redis stopped$(RESET)"

test-env-up: test-db-up test-redis-up ## Start test environment
	@echo "$(GREEN)Test environment is ready$(RESET)"

test-env-down: test-db-down test-redis-down ## Stop test environment
	@echo "$(GREEN)Test environment stopped$(RESET)"

# Full test suite with environment
test-full: test-env-up ## Run full test suite with test environment
	@echo "$(BLUE)Running full test suite...$(RESET)"
	@sleep 5  # Wait for databases to be ready
	@TEST_DB_HOST=localhost \
	 TEST_DB_PORT=5433 \
	 TEST_DB_NAME=doggyclub_test \
	 TEST_DB_USER=postgres \
	 TEST_DB_PASSWORD=password \
	 TEST_REDIS_HOST=localhost \
	 TEST_REDIS_PORT=6380 \
	 go test -v -race -timeout=300s ./...
	@$(MAKE) test-env-down
	@echo "$(GREEN)Full test suite completed$(RESET)"