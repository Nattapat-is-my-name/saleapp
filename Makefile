.PHONY: help build test lint docker-build docker-up docker-down clean backend frontend

# Variables
GO_CMD=go
GO_FILES=$(shell find . -name '*.go' -type f)
BACKEND_DIR=backend
FRONTEND_DIR=frontend
DOCKER_BACKEND_IMAGE=saleapp-backend
DOCKER_FRONTEND_IMAGE=saleapp-frontend
DOCKER_COMPOSE_FILE=docker-compose.yml

# Default target
help:
	@echo "SaleApp Development Makefile"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Backend targets:"
	@echo "  backend-deps       - Download Go dependencies"
	@echo "  backend-build      - Build Go backend binary"
	@echo "  backend-test       - Run Go tests"
	@echo "  backend-lint       - Run Go linter (staticcheck)"
	@echo "  backend-run        - Run Go backend locally"
	@echo ""
	@echo "Frontend targets:"
	@echo "  frontend-deps      - Install npm dependencies"
	@echo "  frontend-build     - Build Next.js frontend"
	@echo "  frontend-test      - Run Next.js tests"
	@echo "  frontend-lint      - Run ESLint"
	@echo "  frontend-dev       - Run Next.js dev server"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build       - Build all Docker images"
	@echo "  docker-up           - Start all services"
	@echo "  docker-down         - Stop all services"
	@echo "  docker-logs         - View logs (use BACKEND=1 or FRONTEND=1 or DB=1)"
	@echo "  docker-clean        - Remove containers, volumes, and images"
	@echo ""
	@echo "Full stack targets:"
	@echo "  build              - Build all (backend + frontend)"
	@echo "  test               - Run all tests"
	@echo "  lint               - Run all linters"
	@echo "  clean              - Clean build artifacts"
	@echo ""

# ===================
# Backend targets
# ===================

backend-deps:
	cd $(BACKEND_DIR) && $(GO_CMD) mod download

backend-build:
	cd $(BACKEND_DIR) && $(GO_CMD) build -v -o saleapp-backend ./cmd/server

backend-test:
	cd $(BACKEND_DIR) && $(GO_CMD) test -v -race -coverprofile=coverage.out ./...

backend-lint:
	cd $(BACKEND_DIR) && $(GO_CMD) vet ./...
	staticcheck ./...

backend-run:
	cd $(BACKEND_DIR) && $(GO_CMD) run ./cmd/server

# ===================
# Frontend targets
# ===================

frontend-deps:
	cd $(FRONTEND_DIR) && npm ci

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-test:
	cd $(FRONTEND_DIR) && npm test

frontend-lint:
	cd $(FRONTEND_DIR) && npm run lint

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

# ===================
# Docker targets
# ===================

docker-build:
	docker compose -f $(DOCKER_COMPOSE_FILE) build

docker-up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs:
ifdef DB
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f postgres
else ifdef BACKEND
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f backend
else ifdef FRONTEND
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f frontend
else
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f
endif

docker-clean:
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi local

# ===================
# Full stack targets
# ===================

build: backend-build frontend-build

test: backend-test frontend-test

lint: backend-lint frontend-lint

clean:
	rm -f $(BACKEND_DIR)/saleapp-backend
	rm -rf $(FRONTEND_DIR)/.next
	rm -f $(BACKEND_DIR)/coverage.out

# ===================
# Development helpers
# ===================

dev: docker-up
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

# Initialize database (run migrations)
db-migrate:
	cd $(BACKEND_DIR) && migrate -path migrations -database "$$DATABASE_URL" up

db-migrate-down:
	cd $(BACKEND_DIR) && migrate -path migrations -database "$$DATABASE_URL" down 1

# Generate Go mocks (requires mockgen)
mocks:
	cd $(BACKEND_DIR) && mockgen -source=internal/repository/product_repo.go -destination=internal/repository/mocks/product_repo_mock.go -package=mocks
	cd $(BACKEND_DIR) && mockgen -source=internal/repository/order_repo.go -destination=internal/repository/mocks/order_repo_mock.go -package=mocks
	cd $(BACKEND_DIR) && mockgen -source=internal/repository/customer_repo.go -destination=internal/repository/mocks/customer_repo_mock.go -package=mocks
