# Makefile for Multi-language Application
.PHONY: help build build-dev run test clean deploy logs status stop restart \
         setup build-go build-node build-python install-deps lint \
         docker-build docker-run docker-push release

# Configuration
APP_NAME := my-app
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/your-username/$(APP_NAME)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")
GO_VERSION := 1.21
NODE_VERSION := 20
PYTHON_VERSION := 3.11
JULIA_VERSION := 1.10.0

# Docker configuration
DOCKERFILE := Dockerfile
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_BUILD_ARGS := --build-arg GOLANG_VERSION=$(GO_VERSION) \
                     --build-arg NODE_VERSION=$(NODE_VERSION) \
                     --build-arg JULIA_VERSION=$(JULIA_VERSION)

# Directories
GO_DIR := ./cmd
NODE_DIR := ./frontend
PYTHON_DIR := ./scripts
BUILD_DIR := ./build
DIST_DIR := ./dist

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

# Default target
.DEFAULT_GOAL := help

##@ Help
help: ## Display this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@echo
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo

##@ Development
setup: ## Install all dependencies for development
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	@go mod download
	@echo "$(YELLOW)Installing Node.js dependencies...$(NC)"
	@if [ -f package.json ]; then npm install; fi
	@echo "$(YELLOW)Installing Python dependencies...$(NC)"
	@if [ -f requirements.txt ]; then pip3 install -r requirements.txt; fi
	@echo "$(GREEN)All dependencies installed!$(NC)"

build-go: ## Build Go application
	@echo "$(YELLOW)Building Go application...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/app $(GO_DIR)/main.go
	@echo "$(GREEN)Go application built: $(BUILD_DIR)/app$(NC)"

build-node: ## Build Node.js/TypeScript application
	@echo "$(YELLOW)Building Node.js application...$(NC)"
	@if [ -f package.json ]; then \
		if [ -f tsconfig.json ]; then \
			npm run build || npx tsc; \
		else \
			echo "$(YELLOW)No TypeScript configuration found, skipping build$(NC)"; \
		fi; \
	else \
		echo "$(YELLOW)No Node.js project found$(NC)"; \
	fi

build-python: ## Install Python dependencies
	@echo "$(YELLOW)Installing Python dependencies...$(NC)"
	@if [ -f requirements.txt ]; then \
		python3 -m pip install -r requirements.txt; \
	else \
		echo "$(YELLOW)No requirements.txt found$(NC)"; \
	fi

build: clean build-go build-node build-python ## Build all components
	@echo "$(GREEN)All components built successfully!$(NC)"

run: build-go ## Run the application locally
	@echo "$(YELLOW)Starting application...$(NC)"
	@$(BUILD_DIR)/app

run-dev: ## Run in development mode with hot reload
	@echo "$(YELLOW)Starting development server...$(NC)"
	@if [ -f docker-compose.dev.yml ]; then \
		docker-compose -f docker-compose.dev.yml up --build; \
	else \
		docker-compose up --build; \
	fi

##@ Testing
test: ## Run all tests
	@echo "$(YELLOW)Running Go tests...$(NC)"
	@go test -v ./... -cover
	@echo "$(YELLOW)Running Node.js tests...$(NC)"
	@if [ -f package.json ]; then npm test; fi
	@echo "$(GREEN)All tests completed!$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

lint: ## Run linters for all languages
	@echo "$(YELLOW)Linting Go code...$(NC)"
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed, skipping$(NC)"; \
	fi
	@echo "$(YELLOW)Linting Node.js code...$(NC)"
	@if [ -f package.json ]; then \
		if [ -f node_modules/.bin/eslint ]; then \
			npm run lint; \
		else \
			echo "$(YELLOW)ESLint not configured, skipping$(NC)"; \
		fi; \
	fi
	@echo "$(GREEN)Linting completed!$(NC)"

##@ Docker
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image $(DOCKER_IMAGE):$(VERSION)...$(NC)"
	@docker build $(DOCKER_BUILD_ARGS) -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(VERSION)$(NC)"

docker-run: docker-build ## Run Docker container locally
	@echo "$(YELLOW)Starting Docker container...$(NC)"
	@docker run -d -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE):$(VERSION)
	@echo "$(GREEN)Container started. Access at http://localhost:8080$(NC)"

docker-compose-up: ## Start with Docker Compose
	@echo "$(YELLOW)Starting services with Docker Compose...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started!$(NC)"

docker-compose-down: ## Stop Docker Compose services
	@echo "$(YELLOW)Stopping Docker Compose services...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "$(GREEN)Services stopped!$(NC)"

docker-push: docker-build ## Push Docker image to registry
	@echo "$(YELLOW)Pushing Docker image to registry...$(NC)"
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_IMAGE):latest
	@echo "$(GREEN)Image pushed to registry!$(NC)"

##@ Deployment
deploy: docker-build docker-push ## Build and deploy application
	@echo "$(YELLOW)Deploying application...$(NC)"
	@if [ -f deploy.sh ]; then \
		chmod +x deploy.sh; \
		./deploy.sh; \
	else \
		echo "$(RED)deploy.sh not found$(NC)"; \
		exit 1; \
	fi

deploy-staging: ## Deploy to staging environment
	@echo "$(YELLOW)Deploying to staging...$(NC)"
	@docker-compose -f docker-compose.staging.yml up -d

deploy-production: ## Deploy to production environment
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@docker-compose -f docker-compose.production.yml up -d

##@ Utility
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) coverage.out coverage.html
	@go clean
	@if [ -f package.json ]; then rm -rf node_modules; fi
	@echo "$(GREEN)Clean completed!$(NC)"

logs: ## Show application logs
	@docker logs -f $(APP_NAME) 2>/dev/null || \
	 echo "$(RED)Container $(APP_NAME) not running. Try: make docker-run$(NC)"

status: ## Show container status
	@echo "$(YELLOW)Container status:$(NC)"
	@docker ps -f name=$(APP_NAME)

stop: ## Stop running container
	@echo "$(YELLOW)Stopping container...$(NC)"
	@docker stop $(APP_NAME) 2>/dev/null || true
	@docker rm $(APP_NAME) 2>/dev/null || true
	@echo "$(GREEN)Container stopped!$(NC)"

restart: stop docker-run ## Restart container

##@ Release
release: test lint docker-build ## Create a new release
	@echo "$(YELLOW)Creating release $(VERSION)...$(NC)"
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
	@echo "$(GREEN)Release v$(VERSION) created and pushed!$(NC)"

version: ## Show current version
	@echo "$(GREEN)Version: $(VERSION)$(NC)"
	@echo "$(GREEN)Go: $(GO_VERSION)$(NC)"
	@echo "$(GREEN)Node.js: $(NODE_VERSION)$(NC)"
	@echo "$(GREEN)Python: $(PYTHON_VERSION)$(NC)"
	@echo "$(GREEN)Julia: $(JULIA_VERSION)$(NC)"

##@ Database (if needed)
db-migrate: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	@if [ -f scripts/migrate.py ]; then \
		python3 scripts/migrate.py; \
	elif [ -f cmd/migrate/main.go ]; then \
		go run cmd/migrate/main.go; \
	else \
		echo "$(YELLOW)No migration scripts found$(NC)"; \
	fi

db-seed: ## Seed database with sample data
	@echo "$(YELLOW)Seeding database...$(NC)"
	@if [ -f scripts/seed.py ]; then \
		python3 scripts/seed.py; \
	elif [ -f cmd/seed/main.go ]; then \
		go run cmd/seed/main.go; \
	else \
		echo "$(YELLOW)No seed scripts found$(NC)"; \
	fi

##@ Monitoring
monitor: ## Show system resources
	@echo "$(YELLOW)System monitoring:$(NC)"
	@docker stats $(APP_NAME) 2>/dev/null || echo "$(RED)Container not running$(NC)"

health: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	@curl -f http://localhost:8080/health || \
	 echo "$(RED)Application is not healthy or not running$(NC)"
