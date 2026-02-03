.PHONY: help install dev build test clean docker-up docker-down migrate seed

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing Go dependencies..."
	@cd services/auth-service && go mod download
	@cd services/org-service && go mod download
	@cd services/procurement-service && go mod download
	@cd services/api-gateway && go mod download
	@echo "Dependencies installed!"

dev: ## Run all services in development mode
	@echo "Starting all services..."
	docker-compose up

build: ## Build all services
	@echo "Building all services..."
	@cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go
	@cd services/org-service && go build -o ../../bin/org-service ./cmd/main.go
	@cd services/procurement-service && go build -o ../../bin/procurement-service ./cmd/main.go
	@cd services/api-gateway && go build -o ../../bin/api-gateway ./cmd/main.go
	@echo "Build complete!"

test: ## Run tests for all services
	@echo "Running tests..."
	@cd services/auth-service && go test ./...
	@cd services/org-service && go test ./...
	@cd services/procurement-service && go test ./...
	@echo "Tests complete!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf services/*/bin/
	@echo "Clean complete!"

docker-up: ## Start all Docker services
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-clean: ## Remove all Docker containers and volumes
	docker-compose down -v

seed-units: ## Seed default units and conversions
	@echo "Seeding units..."
	@docker exec -i erp-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin < scripts/seeds/001_insert_units.js
	@echo "Units seeding complete!"

seed-all: ## Seed all data (init + units + more)
	@echo "Seeding all data..."
	@docker exec -i erp-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin < scripts/mongo-init.js
	@for file in scripts/seeds/*.js; do \
		echo "Seeding $$file..."; \
		docker exec -i erp-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin < $$file; \
	done
	@echo "All seeding complete!"

seed: ## Seed database with initial data
	@echo "Seeding database..."
	@docker exec -i erp-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin < scripts/mongo-init.js
	@echo "Seeding complete!"

# proto: ## Generate protobuf files
# 	@echo "Generating protobuf files..."
# 	@protoc --go_out=. --go_opt=paths=source_relative \
# 		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
# 		shared/proto/*.proto
# 	@echo "Protobuf generation complete!"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs/swagger --pd
	@echo "Swagger docs generated!"

swagger-serve: ## Serve Swagger UI documentation
	@echo "Starting Swagger UI on http://localhost:8000"
	@echo "Press Ctrl+C to stop"
	@python3 -m http.server 8000 || python -m SimpleHTTPServer 8000

api-docs: ## Open API documentation in browser
	@echo "Opening API documentation..."
	@open swagger-ui.html || xdg-open swagger-ui.html || start swagger-ui.html

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "Linting complete!"

format: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy
	@echo "Formatting complete!"
