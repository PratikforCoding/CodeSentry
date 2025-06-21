.PHONY: all up down build run dev deps test test-coverage fmt lint mongo-shell clean logs

# App and binary names
APP_NAME=CodeSentry
BINARY=bin/server
DOCKER_COMPOSE=docker-compose.yml

# Default target
all: deps build up

# Bring up the entire system (API + MongoDB)
up:
	@echo "Starting services..."
	docker-compose -f $(DOCKER_COMPOSE) up --build -d
	@echo "Services started. API available at http://localhost:8080"

# Show logs
logs:
	docker-compose -f $(DOCKER_COMPOSE) logs -f

# Stop containers but preserve data
down:
	@echo "Stopping services..."
	docker-compose -f $(DOCKER_COMPOSE) down

# Stop containers and delete volumes (CAUTION: deletes DB)
down-hard:
	@echo "Stopping services and removing volumes..."
	docker-compose -f $(DOCKER_COMPOSE) down -v

# Build Go app locally
build:
	@echo "Building $(APP_NAME)..."
	mkdir -p bin
	go build -o $(BINARY) cmd/server/main.go

# Run Go app outside docker (for local dev)
run: build
	@echo "Running $(APP_NAME) locally..."
	./$(BINARY)

# Dev mode (requires `air`)
dev:
	@echo "Starting development mode..."
	air -c .air.toml

# Dependency management
deps:
	@echo "Managing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint installed)
lint:
	@echo "Linting code..."
	golangci-lint run

# Shell into MongoDB container
mongo-shell:
	@echo "Connecting to MongoDB..."
	docker exec -it codesentry-mongo-1 mongosh -u root -p example --authenticationDatabase admin

# Check service status
status:
	docker-compose -f $(DOCKER_COMPOSE) ps

# Restart services
restart: down up

# Clean artifacts
clean:
	@echo "Cleaning artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker system prune -f

# Help target
help:
	@echo "Available targets:"
	@echo "  up          - Start all services"
	@echo "  down        - Stop services"
	@echo "  down-hard   - Stop services and remove volumes"
	@echo "  logs        - Show service logs"
	@echo "  build       - Build Go application"
	@echo "  run         - Run application locally"
	@echo "  dev         - Start development mode"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  status      - Show service status"
	@echo "  mongo-shell - Connect to MongoDB"