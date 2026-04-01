# Go Reading Log API - Makefile
# Provides common development commands for building, testing, and running the project

# Variables
BINARY_NAME := server
MODULE_NAME := go-reading-log-api-next
SERVER_CMD := ./cmd/server.go
TEST_PKG := ./...
COVERAGE_FILE := coverage.out

# Go commands
GO := go
GO_RUN := $(GO) run
GO_BUILD := $(GO) build
GO_TEST := $(GO) test
GO_FMT := $(GO) fmt
GO_VET := $(GO) vet

# Colors for output (using tput for portability)
GREEN := $(shell tput -Txterm setaf 2 2>/dev/null || echo "")
YELLOW := $(shell tput -Txterm setaf 3 2>/dev/null || echo "")
BLUE := $(shell tput -Txterm setaf 4 2>/dev/null || echo "")
RED := $(shell tput -Txterm setaf 1 2>/dev/null || echo "")
NC := $(shell tput -Txterm sgr0 2>/dev/null || echo "")

.PHONY: all help run build test clean fmt vet docker-start-pg start-pg test-coverage test-verbose

# Default target
all: help

# Display help information
help:
	@echo "$(BLUE)======================================$(NC)"
	@echo "$(BLUE)  Go Reading Log API - Makefile Help$(NC)"
	@echo "$(BLUE)======================================$(NC)"
	@echo ""
	@echo "$(GREEN)Main Commands:$(NC)"
	@echo "  make run       Build and run the server (development mode)"
	@echo "  make build     Build the binary for production"
	@echo "  make test      Run all tests"
	@echo "  make help      Display this help message"
	@echo ""
	@echo "$(GREEN)Code Quality Commands:$(NC)"
	@echo "  make fmt       Format code with go fmt"
	@echo "  make vet       Run go vet for static analysis"
	@echo "  make clean     Clean up binaries and build artifacts"
	@echo ""
	@echo "$(GREEN)Database Commands:$(NC)"
	@echo "  make start-pg     Start PostgreSQL via Docker (if available)"
	@echo "  make docker-start-pg  Start PostgreSQL via Docker (explicit)"
	@echo ""
	@echo "$(GREEN)Testing Commands:$(NC)"
	@echo "  make test              Run all tests"
	@echo "  make test-verbose      Run tests with verbose output"
	@echo "  make test-coverage     Run tests and generate coverage report"
	@echo ""
	@echo "$(GREEN)Examples:$(NC)"
	@echo "  make run              # Start the server on :3000"
	@echo "  make build            # Build binary to bin/$(BINARY_NAME)"
	@echo "  make test             # Run all 121 tests"
	@echo "  make test-coverage    # Generate coverage.out report"
	@echo ""

# Run the server in development mode
run: build
	@echo "$(YELLOW)Starting server...$(NC)"
	$(GO_RUN) $(SERVER_CMD)

# Build the binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GO_BUILD) -o bin/$(BINARY_NAME) $(SERVER_CMD)
	@echo "$(GREEN)Build complete: bin/$(BINARY_NAME)$(NC)"

# Run all tests
test:
	@echo "$(BLUE)Running all tests...$(NC)"
	$(GO_TEST) $(TEST_PKG)
	@echo "$(GREEN)All tests passed!$(NC)"

# Run tests with verbose output
test-verbose:
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	$(GO_TEST) -v $(TEST_PKG)

# Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO_TEST) -coverprofile=$(COVERAGE_FILE) $(TEST_PKG)
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_FILE)$(NC)"
	$(GO) tool cover -func=$(COVERAGE_FILE)

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO_FMT) $(TEST_PKG)
	@echo "$(GREEN)Code formatted successfully$(NC)"

# Run go vet for static analysis
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO_VET) $(TEST_PKG)
	@echo "$(GREEN)go vet completed$(NC)"

# Clean up build artifacts
clean:
	@echo "$(YELLOW)Cleaning up...$(NC)"
	$(GO) clean -cache -testcache
	rm -f $(COVERAGE_FILE)
	rm -f bin/$(BINARY_NAME)
	@echo "$(GREEN)Clean complete$(NC)"

# Start PostgreSQL via Docker
start-pg: docker-start-pg

# Start PostgreSQL via Docker (explicit)
docker-start-pg:
	@echo "$(BLUE)Checking for Docker...$(NC)"
	@if ! command -v docker &> /dev/null; then \
		echo "$(RED)Error: Docker is not installed or not in PATH$(NC)"; \
		echo "$(YELLOW)Please install Docker or start PostgreSQL manually$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Docker found$(NC)"
	@echo "$(BLUE)Starting PostgreSQL container...$(NC)"
	@docker ps -a --format "{{.Names}}" | grep -q reading-log-db && \
		echo "$(YELLOW)Container already exists. Starting it...$(NC)" && \
		docker start reading-log-db || \
		docker run -d \
			--name reading-log-db \
			-p 5432:5432 \
			-e POSTGRES_USER=$$DB_USER 2>/dev/null || \
		docker run -d \
			--name reading-log-db \
			-p 5432:5432 \
			-e POSTGRES_USER=postgres \
			-e POSTGRES_PASSWORD=postgres \
			-e POSTGRES_DB=reading_log \
			-e PGDATA=/var/lib/postgresql/data/pgdata \
			--health-cmd="pg_isready -U postgres" \
			--health-interval=10s \
			--health-timeout=5s \
			--health-retries=5 \
			postgres:15
	@echo "$(GREEN)PostgreSQL container started$(NC)"
	@echo "$(YELLOW)To connect to the database:$(NC)"
	@echo "  docker exec -it reading-log-db psql -U postgres -d reading_log"
	@echo "$(YELLOW)To stop the container:$(NC)"
	@echo "  docker stop reading-log-db"
	@echo "$(YELLOW)To remove the container:$(NC)"
	@echo "  docker rm reading-log-db"
