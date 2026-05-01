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

.PHONY: all help run build test clean fmt vet docker-start-pg start-pg test-coverage test-verbose docker-up docker-down docker-logs docker-ps docker-stop-pg docker-reload

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
	@echo "  make docker-up    Start all services via docker-compose"
	@echo "  make docker-down  Stop all services via docker-compose"
	@echo "  make docker-logs  Show logs from all services"
	@echo "  make docker-ps    List running containers"
	@echo "  make docker-stop-pg Stop PostgreSQL container"
	@echo "  make docker-reload  Drop and recreate database from docs/database.sql"
	@echo ""
	@echo "$(GREEN)Testing Commands:$(NC)"
	@echo "  make test              Run all tests"
	@echo "  make test-verbose      Run tests with verbose output"
	@echo "  make test-coverage     Run tests and generate coverage report"
	@echo "  make test-clean        Clean up orphaned test databases"
	@echo "  make compare-responses Run Rails API comparison tests"
	@echo ""
	@echo "$(GREEN)Benchmark Commands:$(NC)"
	@echo "  make benchmark-parallel  Run parallel performance benchmarks"
	@echo "  make benchmark-large-scale Run large-scale benchmarks (10,000+ logs)"
	@echo ""
	@echo "$(GREEN)Examples:$(NC)"
	@echo "  make run              # Start the server on :3000"
	@echo "  make build            # Build binary to bin/$(BINARY_NAME)"
	@echo "  make test             # Run all 121 tests"
	@echo "  make test-coverage    # Generate coverage.out report"
	@echo ""

run: build
	@PORT=$${SERVER_PORT:-3000}; \
	if command -v lsof >/dev/null; then \
		PIDS=$$(lsof -t -i :$$PORT); \
		if [ -n "$$PIDS" ]; then \
			echo "$(YELLOW)Killing existing processes: $$PIDS$(NC)"; \
			kill -9 $$PIDS; \
		fi; \
	else \
		echo "$(YELLOW)Warning: lsof not installed. Skipping port check.$(NC)"; \
	fi; \
 	nohup bin/$(BINARY_NAME) > server.log 2>&1 &

stop:
	@PORT=$${SERVER_PORT:-3000}; \
	if command -v lsof >/dev/null; then \
		PIDS=$$(lsof -t -i :$$PORT); \
		if [ -n "$$PIDS" ]; then \
			echo "$(YELLOW)Killing existing processes: $$PIDS$(NC)"; \
			kill -9 $$PIDS; \
		fi; \
	else \
		echo "$(YELLOW)Warning: lsof not installed. Skipping port check.$(NC)"; \
	fi; \

# Build the binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GO_BUILD) -o bin/$(BINARY_NAME) $(SERVER_CMD)
	@echo "$(GREEN)Build complete: bin/$(BINARY_NAME)$(NC)"

# Run all tests
test:
	@echo "$(BLUE)Running all tests...$(NC)"
	@echo "$(BLUE)Loading test configuration from .env.test...$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && $(GO_TEST) -timeout=5m $(TEST_PKG)
	@echo "$(GREEN)All tests passed!$(NC)"

# Run tests with verbose output
test-verbose:
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@echo "$(BLUE)Loading test configuration from .env.test...$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && $(GO_TEST) -v -timeout=5m $(TEST_PKG)

# Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@echo "$(BLUE)Loading test configuration from .env.test...$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && $(GO_TEST) -coverprofile=$(COVERAGE_FILE) -timeout=5m $(TEST_PKG)
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_FILE)$(NC)"
	$(GO) tool cover -func=$(COVERAGE_FILE)

# Clean up orphaned test databases
test-clean:
	@echo "$(BLUE)Cleaning up orphaned test databases...$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && \
		$(GO) run ./test/cleanup_orphaned_databases.go 2>/dev/null || \
		echo "$(YELLOW)No orphaned databases found or cleanup skipped$(NC)"
	@echo "$(GREEN)Cleanup complete$(NC)"

# Run parallel performance benchmarks
benchmark-parallel:
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)  Running Parallel Performance Benchmarks$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && \
		$(GO) test -bench=BenchmarkParallel -benchmem -count=3 $(TEST_PKG)/performance
	@echo "$(GREEN)Benchmark complete$(NC)"
	@echo "$(YELLOW)Run 'go tool pprof -http=:8080 profile.out' to analyze results$(NC)"

# Run large-scale performance benchmarks (10,000+ logs)
benchmark-large-scale:
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)  Running Large-Scale Performance Benchmarks$(NC)"
	@echo "$(BLUE)  Dataset: 100 projects, 10,000+ logs$(NC)"
	@echo "$(BLUE)  Threshold: P95 < 500ms$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && \
		$(GO) test -bench=BenchmarkLargeScale -benchmem -count=3 $(TEST_PKG)/performance
	@echo "$(GREEN)Large-scale benchmark complete$(NC)"
	@echo "$(YELLOW)Results documented in: docs/performance/large-scale-benchmarks.md$(NC)"

# Run Rails API comparison tests
compare-responses:
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)  Running Rails API Comparison Tests$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(YELLOW)Note: RAILS_API_URL must be set (e.g., http://localhost:3001)$(NC)"
	@echo "$(YELLOW)Make sure Rails API is running on port 3001$(NC)"
	@echo ""
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && \
		$(GO) test -v ./test/integration/... -run ".*Comparison.*"

# Alias for test-clean (convenience)
test-cleanup: test-clean

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

# Docker Compose Commands
docker-up:
	@echo "$(BLUE)Starting services with docker-compose...$(NC)"
	@if ! command -v docker-compose &> /dev/null && ! command -v docker &> /dev/null; then \
		echo "$(RED)Error: Docker or docker-compose not installed or not in PATH$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Starting PostgreSQL, Go API, and Rails API...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)Services started$(NC)"
	@echo "$(YELLOW)Go API: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Rails API: http://localhost:3001$(NC)"
	@echo "$(YELLOW)Logs: make docker-logs$(NC)"

docker-down:
	@echo "$(BLUE)Stopping services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"

docker-logs:
	@echo "$(BLUE)Showing logs...$(NC)"
	docker-compose logs -f

docker-ps:
	@echo "$(BLUE)Listing containers...$(NC)"
	docker-compose ps

docker-stop-pg:
	@echo "$(BLUE)Stopping PostgreSQL container...$(NC)"
	@docker ps -a --format "{{.Names}}" | grep -q reading-log-db && \
		docker stop reading-log-db && \
		docker rm reading-log-db || \
		echo "$(YELLOW)No reading-log-db container found$(NC)"

# Reload database from docs/database.sql
reload: docker-reload

docker-reload:
	@echo "$(YELLOW)========================================$(NC)"
	@echo "$(YELLOW)       DATABASE RELOAD WARNING$(NC)"
	@echo "$(YELLOW)========================================$(NC)"
	@echo ""
	@echo "$(RED)This will permanently delete all database data!$(NC)"
	@echo ""
	@echo "$(YELLOW)Database to be reloaded: $(DB_DATABASE)$(NC)"
	@echo "$(YELLOW)SQL file: docs/database.sql$(NC)"
	@echo ""
	@read -p "Are you sure you want to continue? (yes/no): " ans && \
		if ! echo "$$ans" | grep -qE "^[yY](es)?"; then echo "Reload cancelled"; exit 0; fi
	@echo ""
	@echo "$(BLUE)Checking for Docker...$(NC)"
	@if ! command -v docker &> /dev/null; then \
		echo "$(RED)Error: Docker is not installed or not in PATH$(NC)"; \
		echo "$(YELLOW)Please install Docker$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Docker found$(NC)"
	@echo ""
	@echo "$(BLUE)Checking for docs/database.sql...$(NC)"
	@if [ ! -f docs/database.sql ]; then \
		echo "$(RED)Error: docs/database.sql not found$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Database SQL file found$(NC)"
	@echo ""
	@echo "$(BLUE)Stopping services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"
	@echo ""
	@echo "$(BLUE)Removing volumes...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)Volumes removed$(NC)"
	@echo ""
	@echo "$(BLUE)Starting services...$(NC)"
	docker-compose up postgres -d
	@echo "$(GREEN)Services started$(NC)"
	@echo ""
	@echo "$(BLUE)Waiting for PostgreSQL to be ready...$(NC)"
	@for i in $$(seq 1 30); do \
		if docker exec reading-log-db pg_isready -U $${DB_USER:-postgres} > /dev/null 2>&1; then \
			echo "$(GREEN)PostgreSQL is ready$(NC)"; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "$(RED)Error: PostgreSQL did not become ready in time$(NC)"; \
			echo "$(YELLOW)Check logs with: make docker-logs$(NC)"; \
			exit 1; \
		fi; \
		echo "$(YELLOW)Waiting...$$i$(NC)"; \
		sleep 2; \
	done
	@echo ""
	@echo "$(BLUE)Restoring database from docs/database.sql...$(NC)"
	@echo "$(YELLOW)Note: This may take a few moments...$(NC)"
	docker exec -i reading-log-db psql -U $${DB_USER:-postgres} -d $${DB_DATABASE:-reading_log} -f /docker-entrypoint-initdb.d/database.sql > /dev/null 2>&1 || \
		docker exec -i reading-log-db psql -U $${DB_USER:-postgres} -d $${DB_DATABASE:-reading_log} -c '\i /docker-entrypoint-initdb.d/database.sql' > /dev/null 2>&1 || \
		( \
			echo "$(YELLOW)Trying alternative method...$(NC)"; \
			cat docs/database.sql | docker exec -i -e PGHOST=localhost -e PGPORT=$${DB_PORT:-5432} -e PGUSER=$${DB_USER:-postgres} -e PGDATABASE=$${DB_DATABASE:-reading_log} reading-log-db psql -U $${DB_USER:-postgres} -d $${DB_DATABASE:-reading_log} > /dev/null 2>&1 || \
			( \
				echo "$(RED)Error: Database restoration failed$(NC)"; \
				exit 1; \
			) \
		)
	@echo "$(GREEN)Database restored successfully$(NC)"
	@echo ""
	@echo "$(BLUE)Verifying database restoration...$(NC)"
	@if docker exec reading-log-db psql -U $${DB_USER:-postgres} -d $${DB_DATABASE:-reading_log} -c "SELECT 1 FROM projects LIMIT 1" > /dev/null 2>&1; then \
		echo "$(GREEN)Database verification successful$(NC)"; \
		echo ""; \
		echo "$(GREEN)========================================$(NC)"; \
		echo "$(GREEN)       DATABASE RELOAD COMPLETE$(NC)"; \
		echo "$(GREEN)========================================$(NC)"; \
	else \
		echo "$(YELLOW)Warning: Verification query failed, but restoration may still be complete$(NC)"; \
	fi
	@echo ""
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  make docker-logs    # View container logs"
	@echo "  make docker-down    # Stop all services"
